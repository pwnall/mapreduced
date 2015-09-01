package mrd

import (
  "bytes"
  "io/ioutil"
  "testing"
)

func TestFileFetcher(t *testing.T) {
  zipBytes, err := ioutil.ReadFile("../test_tmp/hello.zip")
  if err != nil {
    t.Fatal(err)
  }
  fetcher := NewFileFetcher("../test_tmp/hello.zip")

  if fetcher.Opened() || fetcher.Closed() {
    t.Errorf("Incorrect fetcher state after NewFileFetcher: %v", fetcher)
  }

  if err := fetcher.Close(); err == nil {
    t.Error("Close() without Open() did not error")
  }

  if err := fetcher.Open(); err != nil {
    t.Fatal(err)
  }

  if err := fetcher.Open(); err == nil {
    t.Error("Double Open() did not error")
  }

  readerAt, size := fetcher.Reader()
  if size != len(zipBytes) {
    t.Errorf("Incorrect size from Reader(): %v", size)
  }

  readBuffer := make([]byte, size)
  readSize, err := readerAt.ReadAt(readBuffer, 0)
  if err != nil {
    t.Fatal(err)
  }
  if readSize != size {
    t.Errorf("Incorrect size from ReadAt(): %v", readSize)
  }
  if !bytes.Equal(zipBytes, readBuffer) {
    t.Error("Incorrect data read from ReadAt()")
  }

  if err := fetcher.Close(); err != nil {
    t.Fatal(err)
  }
  if err := fetcher.Open(); err == nil {
    t.Error("Double Close() did not error")
  }
}

func TestFileFetcher_ReaderBeforeOpen(t *testing.T) {
  fetcher := NewFileFetcher("../test_tmp/hello.zip")

  defer func() {
    recovered := recover()
    err, ok := recovered.(error)
    if !ok {
      t.Error("Wrong panic error type: %v", err)
    }
    if err.Error() != "FileFetcher.Open() not called before Reader()" {
      t.Error("Wrong panic error: %v", err)
    }
  }()

  reader, size := fetcher.Reader()
  t.Error("Reader() returned: %v, %v", reader, size)
}

func TestFileFetcher_ReaderAfterClose(t *testing.T) {
  fetcher := NewFileFetcher("../test_tmp/hello.zip")

  if err := fetcher.Open(); err != nil {
    t.Fatal(err)
  }
  if err := fetcher.Close(); err != nil {
    t.Fatal(err)
  }

  defer func() {
    recovered := recover()
    err, ok := recovered.(error)
    if !ok {
      t.Error("Wrong panic error type: %v", err)
    }
    if err.Error() != "FileFetcher.Close() called before Reader()" {
      t.Error("Wrong panic error: %v", err)
    }
  }()

  reader, size := fetcher.Reader()
  t.Error("Reader() returned: %v, %v", reader, size)
}
