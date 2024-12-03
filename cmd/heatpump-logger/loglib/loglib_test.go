package loglib

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSingleFileTFRecordWriter(t *testing.T) {
	var b bytes.Buffer
	writer := &SingleFileTFRecordWriter{ioWriter: &b}
	records := []string{"r1", "r2x"}
	wantBytesWritten := []int{16 + 2, 16 + 2 + 16 + 3}
	for i, rec := range records {
		if err := writer.Write([]byte(rec)); err != nil {
			t.Fatalf("Write() failed: %v", err)
		}
		if got, want := writer.successfulWriteCount, i+1; got != want {
			t.Errorf("got successfulWriteCount = %d, want %d", got, want)
		}
		if got, want := writer.bytesWritten, wantBytesWritten[i]; got != want {
			t.Errorf("got bytesWritten = %d, want %d", got, want)
		}
	}

	reader := bytes.NewReader(b.Bytes())
	recordSeq := ReadAllRecords(reader)
	var recordStrings []string
	for rec := range recordSeq {
		if rec.Error != nil {
			t.Fatalf("error reading records: %v", rec.Error)
		}
		recordStrings = append(recordStrings, string(rec.Record))
	}
	if diff := cmp.Diff(records, recordStrings); diff != "" {
		t.Errorf("unxpected diff reading back records (-want, +got):\n%s", diff)
	}
}

func TestNewMultiFileTFRecordWriter(t *testing.T) {
	for _, tc := range []struct {
		name               string
		records            []string
		shouldFinalizeFile func(name string, recordCount, byteCount int) bool
		wantRecordFiles    map[string][]string
	}{
		{
			records: []string{"a", "b", "c", "d", "e"},
			shouldFinalizeFile: func(name string, recordCount, byteCount int) bool {
				return recordCount == 2
			},
			wantRecordFiles: map[string][]string{
				"1": {"a", "b"},
				"2": {"c", "d"},
				"3": {"e"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fileNum := 0
			files := map[string]*bytes.Buffer{}
			finalizedFiles := map[string][]byte{}
			mfWriter := NewMultiFileTFRecordWriter(
				func() (string, io.Writer, error) {
					buf := &bytes.Buffer{}
					fileNum++
					fileName := fmt.Sprintf("%d", fileNum)
					files[fileName] = buf
					return fileName, buf, nil
				},
				tc.shouldFinalizeFile,
				func(name string, writer io.Writer) error {
					finalizedFiles[name] = files[name].Bytes()
					return nil
				},
			)
			for _, rec := range tc.records {
				if err := mfWriter.Write([]byte(rec)); err != nil {
					t.Fatalf("error writing record")
				}
			}
			if err := mfWriter.Close(); err != nil {
				t.Fatalf("error closing writer: %v", err)
			}
			gotRecordFiles := map[string][]string{}
			for fileName, fileContents := range finalizedFiles {
				gotRecordFiles[fileName] = mustReadAllRecords(t, bytes.NewReader(fileContents))
			}
			if diff := cmp.Diff(tc.wantRecordFiles, gotRecordFiles); diff != "" {
				t.Errorf("unxpected diff reading back records (-want, +got):\n%s", diff)
			}
		})
	}
}

func mustReadAllRecords(t *testing.T, reader io.Reader) []string {
	recordSeq := ReadAllRecords(reader)
	var recordStrings []string
	for rec := range recordSeq {
		if rec.Error != nil {
			t.Fatalf("error reading records: %v", rec.Error)
		}
		recordStrings = append(recordStrings, string(rec.Record))
	}
	return recordStrings
}
