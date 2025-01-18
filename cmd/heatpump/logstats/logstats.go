package logstats

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gonzojive/heatpump/cmd/heatpump-logger/loglib"
	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/internal/collections"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/gonzojive/heatpump/proto/logs"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var LogStatsCommand = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stats",
		Short:   "Prints a summary of logs",
		GroupID: "logs",
	}

	dataDir := cmd.PersistentFlags().String("data-dir", "", "Data directory where logs are stored")

	main := func(ctx context.Context) error {
		glog.Infof("generating stats using data dir %s", *dataDir)
		logFiles, err := loglib.ListAllLogFiles(*dataDir)
		if err != nil {
			return fmt.Errorf("error listing all logs files: %w", err)
		}
		glog.Infof("got %d log files", len(logFiles))
		entries, err := collectLogEntries(ctx, logFiles, time.Hour*24)
		if err != nil {
			return err
		}
		glog.Infof("got %d total log entries", len(entries))

		return nil
	}

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		if err := main(cmd.Context()); err != nil {
			glog.Errorf("fatal error: %v", err)
			os.Exit(1)
		}
	}
	return cmd
}()

func collectLogEntries(ctx context.Context, logFiles []*loglib.LogFile, maxDur time.Duration) ([]*cx34.State, error) {
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Time.After(logFiles[j].Time)
	})
	maxCollectionTime := func(i int) time.Time {
		if i == 0 {
			return time.Now()
		}
		return logFiles[i-1].Time.Add(time.Minute)
	}

	var filteredLogEntries []*cx34.State
	var mostRecentLogEntry *cx34.State
	var lock sync.Mutex
	isRecentEnough := func(t time.Time) bool {
		if mostRecentLogEntry == nil || !t.Before(mostRecentLogEntry.CollectionTime().Add(-maxDur)) {
			return true
		}
		return false
	}
	maybeAddLogEntry := func(entry *cx34.State) {
		lock.Lock()
		defer lock.Unlock()
		if mostRecentLogEntry == nil || entry.CollectionTime().After(mostRecentLogEntry.CollectionTime()) {
			mostRecentLogEntry = entry
			filteredLogEntries = append(filteredLogEntries, entry)
			filteredLogEntries = collections.Filter(filteredLogEntries, func(elem *cx34.State) bool {
				return isRecentEnough(elem.CollectionTime())
			})
		} else if isRecentEnough(entry.CollectionTime()) {
			filteredLogEntries = append(filteredLogEntries, entry)
		}
	}

	{
		eg, ctx := errgroup.WithContext(ctx)
		for i, logFile := range logFiles {
			i, logFile := i, logFile
			eg.Go(func() error {
				if !isRecentEnough(maxCollectionTime(i)) {
					return nil
				}
				for entry := range loglib.ReadLogFile(logFile, parseLogEntry) {
					if entry.Error != nil {
						return fmt.Errorf("error reading entry within log file %s: %w", logFile.Path, entry.Error)
					}
					if !isRecentEnough(maxCollectionTime(i)) {
						// The entire file is too old. Return early.
						return nil
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						maybeAddLogEntry(entry.Record)
					}
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}
	collections.SortBy(filteredLogEntries, func(a, b *cx34.State) bool {
		return a.CollectionTime().Before(b.CollectionTime())
	})

	return filteredLogEntries, nil
}

func parseLogEntry(entryBytes []byte) (*cx34.State, error) {
	var entry logs.HeatpumpLogEntry
	if err := proto.Unmarshal(entryBytes, &entry); err != nil {
		return nil, fmt.Errorf("error parsing log entry: %w", err)
	}
	state, err := cx34.StateFromProto(&chiltrix.State{
		CollectionTime: entry.GetCollectionTime(),
		RegisterValues: entry.GetRegisterValues(),
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing heatpump state from proto: %w", err)
	}

	return state, nil
}

var LogsCollectCommand = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collect",
		Short:   "Start a process that polls the heat pump and dumps state to logs files.",
		GroupID: "logs",
	}

	dataDir := cmd.PersistentFlags().String("data-dir", "", "Data directory where logs are stored")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		glog.Infof("generating stats using data dir %s", *dataDir)
	}
	return cmd
}()

var LogCommand = &cobra.Command{
	Use:   "logs",
	Short: "Print the version number of heatpump",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logs-related commands")
		os.Exit(1)
	},
}

func init() {
	LogCommand.AddGroup(&cobra.Group{
		ID: "logs",
	})
	LogCommand.AddCommand(LogStatsCommand)
	LogCommand.AddCommand(LogsCollectCommand)
}
