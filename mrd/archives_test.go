package mrd

import (
  "archive/tar"
  "archive/zip"
  "bytes"
  "io"
  "testing"
)

func TestTarDirectory(t *testing.T) {
  var tarBuffer bytes.Buffer
  if err := TarDirectory("../testdata/hello", &tarBuffer); err != nil {
    t.Fatal(err)
  }

  tarReader := tar.NewReader(bytes.NewBuffer(tarBuffer.Bytes()))

  foundFiles := map[string]bool { }
  expectedFiles := map[string]string {
    "data/hello.txt": "Hello world!\n",
    "data/goodbye.txt": "Goodbye, cruel world!\n",
  }
  for {
    header, err := tarReader.Next()
    if err == io.EOF {
      break
    }
    if err != nil {
      t.Fatal(err)
    }
    t.Logf("Tar entry: %v\n", header)
    if goldenText, ok := expectedFiles[header.Name]; ok {
      foundFiles[header.Name] = true
      textBuffer := make([]byte, header.Size)
      if _, err := tarReader.Read(textBuffer); err != nil {
        t.Fatal(err)
      }
      if string(textBuffer) != goldenText {
        t.Errorf("Incorrect %s text: %v", header.Name, string(textBuffer))
      }
    }
  }

  if len(foundFiles) != len(expectedFiles) {
    t.Errorf("Did not find all expected files. Found %v\n", foundFiles)
  }
}

func TestZipDirectory(t *testing.T) {
  var zipBuffer bytes.Buffer
  if err := ZipDirectory("../testdata/hello", &zipBuffer); err != nil {
    t.Fatal(err)
  }

  zipReader, err := zip.NewReader(bytes.NewReader(zipBuffer.Bytes()),
                                  int64(zipBuffer.Len()))
  if err != nil {
    t.Fatal(err)
  }

  foundFiles := map[string]bool { }
  expectedFiles := map[string]string {
    "data/hello.txt": "Hello world!\n",
    "data/goodbye.txt": "Goodbye, cruel world!\n",
  }
  for _, zipFile := range zipReader.File {
    t.Logf("Zip entry: %v\n", zipFile)
    if goldenText, ok := expectedFiles[zipFile.Name]; ok {
      foundFiles[zipFile.Name] = true
      textBuffer := make([]byte, zipFile.UncompressedSize)
      file, err := zipFile.Open()
      if err != nil {
        t.Fatal(err)
      }
      if _, err := file.Read(textBuffer); err != nil {
        t.Fatal(err)
      }
      if string(textBuffer) != goldenText {
        t.Errorf("Incorrect %s text: %v", zipFile.Name, string(textBuffer))
      }
    }
  }

  if len(foundFiles) != len(expectedFiles) {
    t.Errorf("Did not find all expected files. Found %v\n", foundFiles)
  }
}
