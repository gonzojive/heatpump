package loglib

import (
	"bytes"
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
