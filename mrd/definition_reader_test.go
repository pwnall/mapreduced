package mrd

import (
  "archive/tar"
  "bytes"
  "io"
  "testing"
)

func TestDefinitionReader_Init(t *testing.T) {
  var template Template
  if err := template.Init("template-id"); err != nil {
    t.Fatal(err)
  }

  fetcher := NewFileFetcher("../test_tmp/hello.zip")
  var reader DefinitionReader
  if err := reader.Init(fetcher, &template); err != nil {
    t.Fatal(err)
  }

  if template.State != TemplateDefined {
    t.Errorf("Incorrect template state: %v\n", template.State)
  }
  if template.ItemCount != 3 {
    t.Errorf("Incorrect template YAML parsing: %v\n", template)
  }

  if fetcher.Closed() {
    t.Error("Init() closed the Fetcher")
  }
}

func TestDefinitionReader_WriteImageTar(t *testing.T) {
  var template Template
  if err := template.Init("template-id"); err != nil {
    t.Fatal(err)
  }

  fetcher := NewFileFetcher("../test_tmp/hello.zip")
  var reader DefinitionReader
  if err := reader.Init(fetcher, &template); err != nil {
    t.Fatal(err)
  }

  var tarBuffer bytes.Buffer
  if err := reader.WriteImageTar(&tarBuffer); err != nil {
    t.Fatal(err)
  }
  if !fetcher.Closed() {
    t.Error("WriteImageTar() did not close the Fetcher")
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
