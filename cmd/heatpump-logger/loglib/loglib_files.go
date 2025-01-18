package loglib

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// LogFile represents a single log file written by MultiFileTFRecordWriter.
type LogFile struct {
	// Path of the TFRecord log file.
	Path string
	// The timestamp of file creation.
	Time time.Time
}

// newLogFile creates a new LogFile, parsing the time from the filename.
func newLogFile(path string) (*LogFile, error) {
	base := filepath.Base(path)
	parts := logFileRegexp.FindStringSubmatch(base)
	if len(parts) < 2 {
		return nil, fmt.Errorf("unable to parse timestamp from filename %q", path)
	}
	timestamp := parts[1]
	t, err := time.Parse(TimestampedFileTimeFormat, timestamp)
	if err != nil {
		return nil, fmt.Errorf("error parsing time from filename %q: %w", path, err)
	}
	return &LogFile{Path: path, Time: t}, nil
}

// ListAllLogFiles lists all of the LogFiles in a directory.
func ListAllLogFiles(dir string) ([]*LogFile, error) {
	var logFiles []*LogFile
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isLogFile(path) {
			logFile, err := newLogFile(path)
			if err != nil {
				return err
			}
			logFiles = append(logFiles, logFile)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory %q: %w", dir, err)
	}
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Time.After(logFiles[j].Time)
	})
	return logFiles, nil
}

var (
	logFileRegexp = regexp.MustCompile(`^.*\.(\d{8}-\d{4})\.tfrecord$`)
)

// isLogFile checks if a given path is likely a log file based on its name.
// This is a heuristic and might need adjustment based on the actual naming scheme.
func isLogFile(path string) bool {
	// Example pattern: "prefix.YYYYMMDD-HHMM.suffix"
	return logFileRegexp.MatchString(filepath.Base(path))
}

// ReadLogFile reads all of the records in a log file using ReadAllRecords
// and calls a provided parsing function on them.
func ReadLogFile[T any](logFile *LogFile, parse func([]byte) (T, error)) iter.Seq[RecordOrErr[T]] {
	return func(yield func(RecordOrErr[T]) bool) {
		file, err := os.Open(logFile.Path)
		if err != nil {
			yield(RecordOrErr[T]{Error: fmt.Errorf("failed to open log file %q: %w", logFile.Path, err)})
			return
		}
		defer file.Close()

		for rec := range ReadAllRecords(file) {
			if rec.Error != nil {
				yield(RecordOrErr[T]{Error: rec.Error})
				return
			}
			parsed, err := parse(rec.Record)
			if !yield(RecordOrErr[T]{Record: parsed, Error: err}) {
				return
			}
		}
	}
}

// ReadLogFiles reads all of the records in a set of log files and calls
// a provided parsing function on them. The records from all files are
// interleaved into a single sequence.
func ReadLogFiles[T any](logFiles []*LogFile, parse func([]byte) (T, error)) iter.Seq[RecordOrErr[T]] {
	return func(yield func(RecordOrErr[T]) bool) {
		for _, logFile := range logFiles {
			for result := range ReadLogFile(logFile, parse) {
				if !yield(result) || result.Error != nil {
					return
				}
			}
		}
	}
}
