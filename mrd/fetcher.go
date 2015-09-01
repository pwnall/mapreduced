package mrd

import (
  "bytes"
  "errors"
  "io"
  "io/ioutil"
)

// Fetcher has the functions used by DefinitionReader to read raw .zip bytes.
type Fetcher interface {
  // Activates this fetcher.
  // The fetcher must not do any work until Open() is called.
  // This must not be called multiple times.
  Open() error

  // Returns the information needed to create a zip.Reader.
  // This must not be called multiple times.
  Reader() (io.ReaderAt, int)

  // Deactivates this fetcher. The reader it returns is no longer usable.
  Close() error
}

// FileFetcher demonstrates how a Fetcher should be implemented.
type FileFetcher struct {
  opened, closed bool
  filePath string
  readerAt io.ReaderAt
  size int
}

// Exposes FileFetcher functionality that is only used in testing.
type StatefulFetcher interface {
  Fetcher

  // Returns true if Open() has been called on this fetcher.
  Opened() bool
  // Returns true if Close() has been called on this fetcher.
  Closed() bool
}

// NewFileFetcher creates a Fetcher that reads from a file.
func NewFileFetcher(filePath string) StatefulFetcher {
  fetcher := new(FileFetcher)
  fetcher.filePath = filePath
  fetcher.opened = false
  fetcher.closed = false
  return fetcher
}

func (f *FileFetcher) Open() error {
  // NOTE: Fetchers are not responsible for checking for multiple Open calls.
  //       FileFetcher implements this so our test suite can catch
  //       implementation errors.
  if f.opened {
    return errors.New("FileFetcher.Open() called multiple times")
  }

  fileBytes, err := ioutil.ReadFile(f.filePath)
  if err != nil {
    return err
  }

  f.opened = true
  f.readerAt = bytes.NewReader(fileBytes)
  f.size = len(fileBytes)
  return nil
}

func (f *FileFetcher) Reader() (io.ReaderAt, int) {
  if !f.opened {
    panic(errors.New("FileFetcher.Open() not called before Reader()"))
  }
  if f.closed {
    panic(errors.New("FileFetcher.Close() called before Reader()"))
  }

  return f.readerAt, f.size
}

func (f *FileFetcher) Close() error {
  if !f.opened {
    return errors.New("FileFetcher.Open() not called before Close()")
  }
  if f.closed {
    return errors.New("FileFetcher.Close() called before Close()")
  }
  f.closed = true
  return nil
}

func (f *FileFetcher) Opened() bool {
  return f.opened
}
func (f *FileFetcher) Closed() bool {
  return f.closed
}
