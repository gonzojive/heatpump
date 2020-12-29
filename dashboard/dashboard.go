// Package dashboard runs a web dashboard that displays information about a CX34 heat pump.
package dashboard

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/mdtable"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/gonzojive/heatpump/units"
)

var markdown = goldmark.New(goldmark.WithExtensions(extension.NewTable()))

const (
	reportInterval = time.Minute

	//dashboardWindow = time.Hour * 24 * 14
	dashboardWindow = time.Hour * 6

	timeLayout        = "2006-01-02T15:04:05"
	machineTimeLayout = "2006-01-02 15:04:05.999999" // based on https://plotly.com/chart-studio-help/date-format-and-time-series/
)

// Run runs a dashboard that displays information about the heat pump.
func Run(ctx context.Context, historianAddr string, httpPort int) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(historianAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %w", err)
	}
	defer conn.Close()
	c := chiltrix.NewHistorianClient(conn)

	server := &server{c}
	server.registerHandlers()

	return (&http.Server{Addr: fmt.Sprintf(":%d", httpPort)}).ListenAndServe()
}

type server struct {
	c chiltrix.HistorianClient
}

func (s *server) registerHandlers() {
	http.HandleFunc("/index.md", s.handleReport)
	http.HandleFunc("/index.csv", s.handleReport)
	http.HandleFunc("/", s.handleReport)

	http.Handle("/index.js", staticHandler(mainScript))
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
				<div class="markdown-body">%s</div>
				<script src='index.js'></script>
			</body>
		</html>
		`, css, mdHTML)
		return
	}
	glog.Errorf("unhandled report request")

}

func (s *server) dashboardContent(ctx context.Context, wantContentType contentType) (content, error) {
	end := time.Now()
	start := end.Add(-1 * dashboardWindow)

	queryClient, err := s.c.QueryStream(ctx, &chiltrix.QueryStreamRequest{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	})
	if err != nil {
		return content{
			text:        fmt.Sprintf("## Error\n\n```\n%s\n```\n", err.Error()),
			contentType: markdownContent,
		}, nil
	}

	table := (&mdtable.Builder{}).SetHeader([]string{
		"Time", "Target Temp", "Inlet Temp", "Outlet Temp", "Ambient Temp", "Flow Rate", "Pump Speed", "Approx Power",
		"COP",
	})

	var states []*cx34.State
	totalBytes := 0
	for {
		resp, err := queryClient.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return content{
				text:        fmt.Sprintf("## Error\n\n```\n%s\n```\n", err.Error()),
				contentType: markdownContent,
			}, err
		}
		s, err := cx34.StateFromProto(resp.GetState())
		if err != nil {
			return content{
				text:        fmt.Sprintf("## Error\n\n```\n%s\n```\n", err.Error()),
				contentType: markdownContent,
			}, err
		}

		states = append(states, s)
		totalBytes += proto.Size(resp)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Proto().GetCollectionTime().AsTime().After(states[j].Proto().CollectionTime.AsTime())
	})
	glog.Infof("got states: %s", cx34.DebugSequenceInfo(states))

	machineReadable := wantContentType == csvContent

	formatTemp := func(t units.Temperature) string {
		return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius(), t.Fahrenheit())
	}

	for _, s := range states {
		if machineReadable {
			cop := ""
			if copFrac, ok := s.COP(); ok {
				cop = fmt.Sprintf("%f", copFrac.Float64())
			}

			table.AddRow([]string{
				s.Proto().GetCollectionTime().AsTime().Local().Format(machineTimeLayout),
				fmt.Sprintf("%.1f", s.ACHeatingTargetTemp().Celsius()),
				fmt.Sprintf("%.1f", s.ACInletWaterTemp().Celsius()),
				fmt.Sprintf("%.1f", s.ACOutletWaterTemp().Celsius()),
				fmt.Sprintf("%.1f", s.AmbientTemp().Celsius()),
				fmt.Sprintf("%.1f", s.FlowRate().LitersPerMinute()),
				fmt.Sprintf("%d", s.InternalPumpSpeed()),
				fmt.Sprintf("%.4f", s.ApparentPower().Kilowatts()),
				cop,
			})
		} else {
			cop := "n/a"
			if copFrac, ok := s.COP(); ok {
				cop = fmt.Sprintf("%.1f%%", copFrac.Float64()*100)
			}
			cop = fmt.Sprintf("%s (%.1fkW/%.1fkW)", cop, s.UsefulHeatRate().Kilowatts(), s.ApparentPower().Kilowatts())
			table.AddRow([]string{
				s.Proto().GetCollectionTime().AsTime().Local().Format(timeLayout),
				formatTemp(s.ACHeatingTargetTemp()),
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
	return content{
		text:        out.String(),
		contentType: markdownContent,
	}, nil

	// TODO: Use https://plotly.com/javascript/time-series/ to render
	// time series plots.

}

type content struct {
	text        string
	contentType contentType
}

type contentType string

const (
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
