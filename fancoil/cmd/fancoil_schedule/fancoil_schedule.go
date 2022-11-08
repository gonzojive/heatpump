package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/proto/fancoil"
	"github.com/martinlindhe/unit"
	rrule "github.com/teambition/rrule-go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	fancoilServiceAddress = flag.String("fancoil-service", "localhost:8083", "Remote address of fancoil service.")
)

var localTimezone = func() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	return loc
}()

type recurringCommands struct {
	name           string
	recurrenceRule *rrule.RRule
	commands       []*fancoil.SetStateRequest
}

var (
	daytimeTemp  = &fancoil.Temperature{DegreesCelcius: float32(unit.FromFahrenheit(68).Celsius())}
	nightimeTemp = &fancoil.Temperature{DegreesCelcius: float32(unit.FromFahrenheit(61).Celsius())}
)

const (
	startYear       = 2022
	startMonth      = time.January
	startDayOfMonth = 1
)

var allRecurringCommands = []*recurringCommands{
	{
		name: "Set fan speed to low in the evenings for ease of conversation",
		recurrenceRule: mustNewRRule(rrule.ROption{
			Freq:    rrule.DAILY,
			Count:   4000,
			Dtstart: time.Date(startYear, startMonth, startDayOfMonth, 17, 15, 0, 0, localTimezone),
		}),
		commands: []*fancoil.SetStateRequest{
			{
				PreferenceFanSetting: fancoil.FanSetting_FAN_SETTING_LOW,
			},
		},
	},

	{
		name: "At bedtime, set target temp to 61F and fan speed to low",
		recurrenceRule: mustNewRRule(rrule.ROption{
			Freq:    rrule.DAILY,
			Count:   4000,
			Dtstart: time.Date(startYear, startMonth, startDayOfMonth, 19, 30, 0, 0, localTimezone),
		}),
		commands: []*fancoil.SetStateRequest{
			{
				PreferenceFanSetting:     fancoil.FanSetting_FAN_SETTING_LOW,
				HeatingTargetTemperature: nightimeTemp,
			},
		},
	},

	{
		name: "At 5am, start heating up the room at low fan speed",
		recurrenceRule: mustNewRRule(rrule.ROption{
			Freq:    rrule.DAILY,
			Count:   4000,
			Dtstart: time.Date(startYear, startMonth, startDayOfMonth, 5, 0, 0, 0, localTimezone),
		}),
		commands: []*fancoil.SetStateRequest{
			{
				HeatingTargetTemperature: daytimeTemp,
			},
		},
	},

	{
		name: "At 7:15am, increase the fan speed",
		recurrenceRule: mustNewRRule(rrule.ROption{
			Freq:    rrule.DAILY,
			Count:   4000,
			Dtstart: time.Date(startYear, startMonth, startDayOfMonth, 7, 15, 0, 0, localTimezone),
		}),
		commands: []*fancoil.SetStateRequest{
			{
				PreferenceFanSetting: fancoil.FanSetting_FAN_SETTING_MEDIUM,
			},
		},
	},
}

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	client, err := dial(ctx)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	{
		executionStart := time.Now()
		sort.Slice(allRecurringCommands, func(i, j int) bool {
			a := allRecurringCommands[i].recurrenceRule.After(executionStart, true)
			b := allRecurringCommands[j].recurrenceRule.After(executionStart, true)
			return a.Before(b)
		})
	}

	for _, scheduleItem := range allRecurringCommands {
		scheduleItem := scheduleItem
		eg.Go(func() error {
			for {
				now := time.Now()
				next := scheduleItem.recurrenceRule.After(now, true)
				if next.IsZero() {
					return nil
				}
				glog.Infof("executing next %q command in %s", scheduleItem.name, next.Sub(now))
				time.Sleep(next.Sub(now))
				for _, cmd := range scheduleItem.commands {
					_, err := client.SetState(ctx, cmd)
					if err != nil {
						glog.Errorf("error setting fan coil state: %v", err)
					}
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func dial(ctx context.Context) (fancoil.FanCoilServiceClient, error) {
	conn, err := grpc.DialContext(ctx, *fancoilServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FanCoilService: %w", err)
	}
	return fancoil.NewFanCoilServiceClient(conn), nil
}

func mustNewRRule(arg rrule.ROption) *rrule.RRule {
	value, err := rrule.NewRRule(arg)
	if err != nil {
		panic(err)
	}
	return value
}

type RRuleTimer struct {
	rrule *rrule.RRule
	timer *time.Timer
}
