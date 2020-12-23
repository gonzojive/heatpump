// Program waterpi reads temperature sensors.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/omron"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	rawOutputLog   = flag.String("rs485-raw-log", "/tmp/serial-log.bin", "Path to USB-to-RS485 raw output log.")
	markdown       = goldmark.New(goldmark.WithExtensions(extension.NewTable()))
	logAllCommands = flag.Bool("log-commands", false, "Whether to print all commands")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	logContents, err := ioutil.ReadFile(*rawOutputLog)
	if err != nil {
		return err
	}
	lines := strings.Split(string(logContents), "\r")
	glog.Infof("got %d commands", len(lines))
	var commands []*omron.CommandFrame
	for i, line := range lines {
		if line == "" {
			continue
		}
		cmd, err := omron.ParseCommandFrame(line)
		if err != nil {
			glog.Errorf("%v", fmt.Errorf("error parsing line %d: %w", i+1, err))
			continue
		}
		if cmd.ParsedCommand != nil && *logAllCommands {
			glog.Infof("cmd[%d]: %s", i, cmd.ParsedCommand)
		} else if *logAllCommands {
			glog.Infof("cmd[%d]: %+v", i, cmd)
		}
		commands = append(commands, cmd)
	}

	counts := map[rwCountKey][]uint16{}
	valueCounts := map[summaryTableKey]int{}
	for i := 0; i < len(commands); i++ {
		c := commands[i]
		switch cmd := c.ParsedCommand.(type) {
		case *omron.DMAreaReadCommand:
			if i+1 >= len(commands) {
				break
			}
			next := commands[i+1]
			resp, ok := next.ParsedCommand.(*omron.DMAreaReadResponse)
			if !ok {
				return fmt.Errorf("command %d should be a read response, got %+v", i, next)
			}
			for j, value := range resp.Data {
				key := rwCountKey{Code: omron.RD, Word: cmd.BeginningWord + uint16(j)}
				counts[key] = append(counts[key], value)
				valueCounts[summaryTableKey{key.Code, key.Word, value}]++
			}
			i++
		case *omron.DMAreaWriteCommand:
			if i+1 >= len(commands) {
				break
			}
			next := commands[i+1]
			resp, ok := next.ParsedCommand.(*omron.DMAreaWriteResponse)
			if !ok {
				return fmt.Errorf("command %d should be a write response, got %+v", i, next)
			} else if resp.EndCode != omron.OK {
				return fmt.Errorf("bad WD response %s", cmd)
			}
			i++
			for j, value := range cmd.Data {
				key := rwCountKey{Code: omron.WD, Word: cmd.BeginningWord + uint16(j)}
				counts[key] = append(counts[key], value)
				valueCounts[summaryTableKey{key.Code, key.Word, value}]++
			}
		}
	}
	var entries []summaryTableEntry
	for key, count := range valueCounts {
		//entries = append(entries, rwCountEntry{rwCountKey: key, Values: values})
		entries = append(entries, summaryTableEntry{Key: key, Count: count})
	}
	sort.Slice(entries, func(i, j int) bool {
		if ci, cj := entries[i].Key.Code, entries[j].Key.Code; ci != cj {
			return ci < cj
		}
		wi, wj := entries[i].Key.Word, entries[j].Key.Word
		if wi == wj {
			return entries[i].Key.Value < entries[j].Key.Value
		}
		return wi < wj
	})
	for _, e := range entries {
		glog.Infof("%s | %d | %d | %d", e.Key.Code, e.Key.Word, e.Key.Value, e.Count)
	}
	return nil
}

type rwCountKey struct {
	Code omron.HeaderCode
	Word uint16
}

type summaryTableKey struct {
	Code  omron.HeaderCode
	Word  uint16
	Value uint16
}

type summaryTableEntry struct {
	Key   summaryTableKey
	Count int
}
