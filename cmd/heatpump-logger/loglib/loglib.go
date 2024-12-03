package loglib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"sync"
	"time"

	"github.com/ryszard/tfutils/go/tfrecord"
)

// overheadPerRecord is the number of bytes written for each TFRecord other than
// the data for the record.
//
// https://www.tensorflow.org/tutorials/load_data/tfrecord#tfrecords_format_details
// states each record is stored in
//
// uint64 length
// uint32 masked_crc32_of_length
// byte   data[length]
// uint32 masked_crc32_of_data
const overheadPerRecord = 8 + 4 + 4

// SingleFileTFRecordWriter is used to write a set of records to a single TFRecord file.
type SingleFileTFRecordWriter struct {
	ioWriter             io.Writer
	lock                 sync.Mutex
	successfulWriteCount int
	bytesWritten         int
}

// Write writs a single TFRecord record and returns any error that might occur.
func (w *SingleFileTFRecordWriter) Write(record []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	err := tfrecord.Write(w.ioWriter, record)
	if err != nil {
		return fmt.Errorf("error writing TFRecord from SIngleifleTFRecordWriter: %w", err)
	}
	w.bytesWritten += len(record) + overheadPerRecord
	w.successfulWriteCount++
	return nil
}

// MultiFileTFRecordWriter is used to write a set of records to a single TFRecord file.
type MultiFileTFRecordWriter struct {
	lock               sync.Mutex
	writerName         string
	singleFileWriter   *SingleFileTFRecordWriter
	closed             bool
	shouldFinalizeFile func(name string, recordCount, byteCount int) bool
	finalizeFile       func(name string, writer io.Writer) error
	newWriter          func() (string, io.Writer, error)
}

// TimestampedNewFileCreator returns a function that creates a new file with the
// given prefix and a -YYYYMMDD-HHMM suffix based on the current time.
//
// The returned function takes no arguments. It generates a filename by
// appending the current date and time in YYYYMMDD-HHMM format to the prefix and
// then opens a file with that name. It returns the generated filename, an
// io.Writer for the opened file, and an error if there was any problem opening
// the file.
//
// When the function is called and no error is returned, the caller must call
// Close() on the returned writer.
func TimestampedNewFileCreator(ctx context.Context, prefix, suffix string, now func() time.Time) func() (string, io.WriteCloser, error) {
	return func() (string, io.WriteCloser, error) {
		name := fmt.Sprintf("%s.%s%s", prefix, now().Format("20060102-1504"), suffix)
		f, err := os.Create(name)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create file: %w", err)
		}
		return name, f, nil
	}
}

// NewMultiFileTFRecordWriter returns an object for writing records to a set of files.
func NewMultiFileTFRecordWriter(
	newFile func() (string, io.Writer, error),
	shouldFinalizeFile func(name string, recordCount, byteCount int) bool,
	finalizeFile func(name string, writer io.Writer) error,
) *MultiFileTFRecordWriter {
	return &MultiFileTFRecordWriter{
		singleFileWriter:   nil,
		newWriter:          newFile,
		shouldFinalizeFile: shouldFinalizeFile,
		finalizeFile:       finalizeFile,
	}
}

// Write writs a single TFRecord record and returns any error that might occur.
func (w *MultiFileTFRecordWriter) Write(record []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.closed {
		return fmt.Errorf("invalid operation: attempted to write to closed writer")
	}
	if w.singleFileWriter == nil {
		name, ioWriter, err := w.newWriter()
		if err != nil {
			return fmt.Errorf("failed to create new single-file writer: %w", err)
		}
		w.writerName = name
		w.singleFileWriter = &SingleFileTFRecordWriter{ioWriter: ioWriter}
	}
	if err := w.singleFileWriter.Write(record); err != nil {
		return err
	}
	if !w.shouldFinalizeFile(w.writerName, w.singleFileWriter.successfulWriteCount, w.singleFileWriter.bytesWritten) {
		return nil
	}
	if err := w.finalizeFile(w.writerName, w.singleFileWriter.ioWriter); err != nil {
		return fmt.Errorf("error finalizing file %q: %w", w.writerName, err)
	}
	w.writerName = ""
	w.singleFileWriter = nil
	return nil
}

// Close finalizes the active active file and disables further writing.
func (w *MultiFileTFRecordWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.closed = true

	if w.singleFileWriter == nil {
		return nil
	}
	if err := w.finalizeFile(w.writerName, w.singleFileWriter.ioWriter); err != nil {
		return fmt.Errorf("error finalizing file %q: %w", w.writerName, err)
	}
	w.writerName = ""
	w.singleFileWriter = nil
	return nil
}

// NewPeriodicMultiFileTFRecordWriter returns a [MultiFileTFRecordWriter] that creates a new file
// periodically, based on the specified singleFileInterval.
//
// The returned MultiFileTFRecordWriter will create a new file when the current time is
// greater than or equal to the most recent file creation time plus the singleFileInterval.
// Each file will be named using the filePrefix and a timestamp in the format YYYYMMDD-HHMM.
//
// The now parameter allows for injecting a custom time source for testing or other purposes.
func NewPeriodicMultiFileTFRecordWriter(
	ctx context.Context,
	now func() time.Time,
	filePrefix, fileSuffix string,
	singleFileInterval time.Duration,
) *MultiFileTFRecordWriter {
	mostRecentFileCreateTime := now()
	var mostRecentFile io.WriteCloser
	createFile := TimestampedNewFileCreator(ctx, filePrefix, fileSuffix, now)

	return NewMultiFileTFRecordWriter(
		func() (string, io.Writer, error) {
			mostRecentFileCreateTime = now()
			name, writer, err := createFile()
			mostRecentFile = writer
			return name, writer, err
		},
		func(_ string, _, _ int) bool {
			return now().Sub(mostRecentFileCreateTime) >= singleFileInterval
		},
		func(name string, writer io.Writer) error {
			if writer != mostRecentFile {
				return fmt.Errorf("programing assumption error: expected writer %v to equal mostRecent file %v", writer, mostRecentFile)
			}
			return mostRecentFile.Close()
		},
	)
}

type RecordOrErr[T any] struct {
	Record T
	Error  error
}

func ReadAllRecords(reader io.Reader) iter.Seq[RecordOrErr[[]byte]] {
	return func(yield func(RecordOrErr[[]byte]) bool) {
		for {
			record, err := tfrecord.Read(reader)
			if errors.Is(err, io.EOF) {
				return
			}
			yield(RecordOrErr[[]byte]{record, err})
			if err != nil {
				break
			}
		}
	}
}
