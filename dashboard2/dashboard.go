// Package dashboard runs a web dashboard that displays information about a CX34 heat pump.
package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"google.golang.org/grpc"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/mdtable"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/gonzojive/heatpump/units"
)

var markdown = goldmark.New(goldmark.WithExtensions(extension.NewTable()))

const (
	reportInterval = time.Minute

	dashboardWindow              = time.Hour * 24 * 60
	dashboardPointSpacingDefault = time.Minute * 5
	registerChangeSampleCount    = 5
	//dashboardWindow = time.Hour * 6

	timeLayout        = "2006-01-02T15:04:05"
	machineTimeLayout = "2006-01-02 15:04:05.999999" // based on https://plotly.com/chart-studio-help/date-format-and-time-series/
)

// Run runs a dashboard that displays information about the heat pump.
func Run(ctx context.Context, historianAddr string, httpPort int) error {
	glog.Infof("dialing %s...", historianAddr)
	// Set up a connection to the server.
	conn, err := grpc.Dial(historianAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %w", err)
	}
	defer conn.Close()
	c := chiltrix.NewHistorianClient(conn)

	glog.Infof("starting cache...")
	cache, closer, err := newCache(ctx, c)
	if err != nil {
		return err
	}
	defer closer.Close()

	server := &server{c, cache}
	server.registerHandlers()

	glog.Infof("starting up http server...")
	return (&http.Server{Addr: fmt.Sprintf(":%d", httpPort)}).ListenAndServe()
}

type server struct {
	c     chiltrix.HistorianClient
	cache *cache
}

func (s *server) registerHandlers() {
	http.HandleFunc("/index.md", s.handleReport)
	http.HandleFunc("/index.csv", s.handleReport)
	http.HandleFunc("/", s.handleReport)

	http.Handle("/index.js", staticHandler(mainScript))
}

func (s *server) handleSetTemp(w http.ResponseWriter, r *http.Request) {
	writeErr := func(err error) {
		w.Header().Set("Content-Type", textContent.headerValue())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
	}
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		writeErr(err)
		return
	}
	t, err := strconv.ParseFloat(u.Query().Get("target-heat"), 64)
	if err != nil {
		writeErr(err)
		return
	}
	fmt.Fprintf(w, "set temp to %f", units.FromCelsius(t).Celsius())
}

func (s *server) handleReport(w http.ResponseWriter, r *http.Request) {
	wantContentType := htmlContent
	if strings.HasSuffix(r.URL.Path, ".md") {
		wantContentType = markdownContent
	}
	if strings.HasSuffix(r.URL.Path, ".csv") {
		wantContentType = csvContent
	}

	resp, err := s.dashboardContent(r.Context(), wantContentType)

	finalContentType := wantContentType
	if wantContentType == markdownContent && err == nil && resp.contentType != markdownContent {
		err = fmt.Errorf("internal error: requested markdown but got %s", resp.contentType)
	}

	w.Header().Set("Content-Type", finalContentType.headerValue())

	if err != nil {
		glog.Errorf("error generating dashboard response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if resp.contentType == wantContentType {
		w.Write([]byte(resp.text))
		return
	}
	if wantContentType == htmlContent && resp.contentType == markdownContent {
		mdHTML := &strings.Builder{}
		if err := markdown.Convert([]byte(resp.text), mdHTML); err != nil {
			w.Header().Set("Content-Type", "text; charset=UTF-8")
			glog.Errorf("error rendering markdown: %s", err)
			return
		}
		fmt.Fprintf(w, `
		<html>
			<head>
				<title>waterpi report</title>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<style>%s</style>
				<script src='https://cdn.plot.ly/plotly-latest.min.js'></script>
			</head>
			<body>
				<div id='copDiv'><!-- Plotly chart will be drawn inside this DIV --></div>
				<div id='tempDiv'><!-- Plotly chart will be drawn inside this DIV --></div>
				<div id='pumpDiv'><!-- Plotly chart will be drawn inside this DIV --></div>
				<div id='powerQualityDiv'><!-- Plotly chart will be drawn inside this DIV --></div>
				<div class="markdown-body">%s</div>
				<script src='index.js'></script>
			</body>
		</html>
		`, css, mdHTML)
		return
	}
	glog.Errorf("unhandled report request")

}

func (s *server) queryStates(ctx context.Context, span span) ([]*cx34.State, error) {
	states, err := s.cache.queryStates(ctx, span)
	if err != nil {
		return nil, err
	}
	var sampled []*cx34.State
	var mostRecentSample *cx34.State
	now := time.Now()
	for i := len(states) - 1; i >= 0; i-- {
		s := states[i]
		if mostRecentSample == nil || mostRecentSample.CollectionTime().Sub(s.CollectionTime()) >= spacingByTimeAgo(now.Sub(s.CollectionTime())) {
			sampled = append(sampled, s)
			mostRecentSample = s
			continue
		}

	}
	sort.Slice(sampled, func(i, j int) bool {
		return sampled[i].CollectionTime().Before(sampled[j].CollectionTime())
	})
	return sampled, nil
}

func (s *server) dashboardContent(ctx context.Context, wantContentType contentType) (content, error) {
	end := time.Now()
	start := end.Add(-1 * dashboardWindow)

	states, err := s.queryStates(ctx, span{start, end})
	if err != nil {
		return content{
			text:        fmt.Sprintf("## Error\n\n```\n%s\n```\n", err.Error()),
			contentType: markdownContent,
		}, nil
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].CollectionTime().After(states[j].CollectionTime())
	})

	table := &mdtable.Builder{}

	machineReadable := wantContentType == csvContent
	if machineReadable {
		table.SetHeader([]string{
			"Time",
			"Target Temp",
			"Inlet Temp",
			"Outlet Temp",
			"Ambient Temp",
			"Flow Rate",
			"Pump Speed",
			"Approx Power",
			"Voltage",
			"COP",
		})
	} else {
		table.SetHeader([]string{
			"Time",
			"Target Temp",
			"Inlet Temp",
			"Outlet Temp",
			"Ambient Temp",
			"Flow Rate",
			"Pump Speed",
			"Approx Power",
			"COP",
		})
	}

	totalBytes := 0

	formatTemp := func(t units.Temperature) string {
		return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius(), t.Fahrenheit())
	}

	for _, s := range states {
		totalBytes += proto.Size(s.Proto())
		if machineReadable {
			cop := ""
			if copFrac, ok := s.COP(); ok {
				cop = fmt.Sprintf("%f", copFrac.Float64())
			}

			table.AddRow([]string{
				s.CollectionTime().Local().Format(machineTimeLayout),
				fmt.Sprintf("%.1f", s.ACTargetTemp().Celsius()),
				fmt.Sprintf("%.1f", s.ACInletWaterTemp().Celsius()),
				fmt.Sprintf("%.1f", s.ACOutletWaterTemp().Celsius()),
				fmt.Sprintf("%.1f", s.AmbientTemp().Celsius()),
				fmt.Sprintf("%.1f", s.FlowRate().LitersPerMinute()),
				fmt.Sprintf("%d", s.InternalPumpSpeed()),
				fmt.Sprintf("%.4f", s.ApparentPower().Kilowatts()),
				fmt.Sprintf("%.4f", s.ACVoltage().Volts()),
				cop,
			})
		} else {
			cop := "n/a"
			if copFrac, ok := s.COP(); ok {
				cop = fmt.Sprintf("%.1f%%", copFrac.Float64()*100)
			}
			cop = fmt.Sprintf("%s (%.1fkW/%.1fkW)", cop, s.UsefulHeatRate().Kilowatts(), s.ApparentPower().Kilowatts())
			table.AddRow([]string{
				s.CollectionTime().Local().Format(timeLayout),
				formatTemp(s.ACTargetTemp()),
				formatTemp(s.ACInletWaterTemp()),
				formatTemp(s.ACOutletWaterTemp()),
				formatTemp(s.AmbientTemp()),
				fmt.Sprintf("%.1fL/m", s.FlowRate().LitersPerMinute()),
				s.InternalPumpSpeed().String(),
				fmt.Sprintf("%.1fV⋅%.1fA = %.3fkVA", s.ACVoltage().Volts(), s.ACCurrent().Amperes(), s.ApparentPower().Kilowatts()),
				cop,
			})
		}
	}

	if wantContentType == csvContent {
		return content{
			contentType: csvContent,
			text:        table.BuildCSV(),
		}, nil
	}

	out := &strings.Builder{}

	hpReport := "## Heat Pump report\n"

	hpReport += fmt.Sprintf("\n%d state (%d bytes) snapshots between %s and %s", len(states), totalBytes, start.Local().Format(timeLayout), end.Local().Format(timeLayout))
	out.WriteString(hpReport)
	out.WriteString(fmt.Sprintf("\n\n%s\n", table.Build()))

	if len(states) > 0 {
		fmt.Fprintf(out, "\n%s\n", s.registersTable(states[0]).text)
	}

	fmt.Fprintf(out, "\n%s\n", s.registerChangesTable(states).text)
	return content{
		text:        out.String(),
		contentType: markdownContent,
	}, nil

	// TODO: Use https://plotly.com/javascript/time-series/ to render
	// time series plots.

}

func copyStates(states []*cx34.State) []*cx34.State {
	var out []*cx34.State
	for _, s := range states {
		out = append(out, s)
	}
	return out
}

func (s *server) registersTable(state *cx34.State) content {
	table := &mdtable.Builder{}
	table.SetHeader([]string{
		"Register",
		"Register Name",
		"Raw Value",
	})
	type entry struct {
		reg   cx34.Register
		value uint16
	}
	var entries []entry
	for r, v := range state.RegisterValues() {
		entries = append(entries, entry{r, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].reg < entries[j].reg
	})

	for _, e := range entries {
		table.AddRow([]string{fmt.Sprintf("%d", int(e.reg)), e.reg.String(), fmt.Sprintf("%d", e.value)})
	}

	return content{
		text:        "## Latest modbus register values\n\n" + table.Build(),
		contentType: markdownContent,
	}
}

func (s *server) registerChangesTable(states []*cx34.State) content {
	states = copyStates(states)
	sort.Slice(states, func(i, j int) bool {
		return states[i].CollectionTime().Before(states[j].CollectionTime())
	})
	type change struct {
		register cx34.Register
		old, new uint16
	}
	type timestampedChange struct {
		register cx34.Register
		old, new uint16
		t        time.Time
	}
	changeCounts := make(map[change]int)
	distinctValues := make(map[cx34.Register]map[uint16]int)
	changeSamples := make(map[cx34.Register][]timestampedChange)
	var changeSampleRows [][]string

	for _, s := range states {
		for r, v := range s.RegisterValues() {
			m, ok := distinctValues[r]
			if !ok {
				m = make(map[uint16]int)
				distinctValues[r] = m
			}
			m[v]++
		}
	}

	for i := len(states) - 2; i >= 0; i-- {
		ith, next := states[i], states[i+1]
		for r, old := range ith.RegisterValues() {

			newVal := next.RegisterValues()[r]
			if old == newVal {
				continue
			}
			changeCounts[change{r, old, newVal}]++
			if len(changeSamples[r]) < registerChangeSampleCount {
				changeSamples[r] = append(changeSamples[r], timestampedChange{
					t:        next.CollectionTime(),
					register: r,
					old:      old,
					new:      newVal,
				})
				changeSampleRows = append(changeSampleRows, []string{
					r.String(),
					fmt.Sprintf("%d", len(distinctValues[r])),
					next.CollectionTime().Local().Format(machineTimeLayout),
					fmt.Sprintf("%d", old),
					fmt.Sprintf("%d", newVal),
				})
			}
		}
	}

	table := &mdtable.Builder{}
	table.SetHeader([]string{
		"Register",
		"Register distinct value count",
		"Timestamp",
		"Old value",
		"New Value",
	})

	sort.Slice(changeSampleRows, func(i, j int) bool {
		distinctValueCountI, _ := strconv.Atoi(changeSampleRows[i][1])
		distinctValueCountJ, _ := strconv.Atoi(changeSampleRows[j][1])
		if distinctValueCountI != distinctValueCountJ {
			return distinctValueCountI < distinctValueCountJ
		}
		if changeSampleRows[i][0] != changeSampleRows[j][0] {
			return changeSampleRows[i][0] < changeSampleRows[j][0]
		}
		if changeSampleRows[i][2] != changeSampleRows[j][2] {
			return changeSampleRows[i][2] > changeSampleRows[j][2]
		}
		return false
	})
	for _, r := range changeSampleRows {
		table.AddRow(r)
	}

	return content{
		text:        "## Sample of changed modbus register values\n\n" + table.Build(),
		contentType: markdownContent,
	}

}

type content struct {
	text        string
	contentType contentType
}

type contentType string

const (
	textContent       contentType = "text/plain; charset=UTF-8"
	htmlContent       contentType = "text/html; charset=UTF-8"
	markdownContent   contentType = "text/markdown; charset=UTF-8"
	javascriptContent contentType = "text/javascript; charset=UTF-8"
	csvContent        contentType = "text/csv; charset=UTF-8"
)

// headerValue returns the Content-Type HTTP header for this content type.
func (ct contentType) headerValue() string {
	return string(ct)
}

func indent(content, indent string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}

// staticHandler returns an http.Handler that serves static content with a given
// mime type.
func staticHandler(content content) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting the headers so we can set the correct mime type
		headers := w.Header()
		headers["Content-Type"] = []string{content.contentType.headerValue()}
		fmt.Fprint(w, content.text)
	})
}

var dashboardPointSpacingByRecency = []struct {
	ago     time.Duration
	spacing time.Duration
}{
	{time.Minute * 30, time.Second},
	{time.Hour * 5, time.Second * 20},
	{time.Hour * 2, time.Minute},
	{time.Hour * 48, time.Minute * 5},
}

func spacingByTimeAgo(ago time.Duration) time.Duration {
	for _, entry := range dashboardPointSpacingByRecency {
		if ago < entry.ago {
			return entry.spacing
		}
	}
	return dashboardPointSpacingDefault
}
