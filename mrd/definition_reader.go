package mrd

import (
  "archive/tar"
  "archive/zip"
  "errors"
  "io"
)

// Manages a template definition in a .zip file.
type DefinitionReader struct {
  fetcher Fetcher
  zipReader *zip.Reader
}

// Init sets the reader's source .zip and target Template.
func (r *DefinitionReader) Init(fetcher Fetcher, template *Template) error {
  if err := fetcher.Open(); err != nil {
    return err
  }

  zipReaderAt, zipSize := fetcher.Reader()
  zipReader, err := zip.NewReader(zipReaderAt, int64(zipSize))
  if err != nil {
    return err
  }

  for _, zipFile := range zipReader.File {
    if zipFile.Name != "mapreduced.yml" {
      continue
    }
    zipBytes := make([]byte, zipFile.UncompressedSize)
    file, err := zipFile.Open()
    if err != nil {
      return err
    }
    if _, err := file.Read(zipBytes); err != nil {
      return err
    }

    if err := template.ReadDefinition(zipBytes); err != nil {
      return err
    }
  }

  if template.State != TemplateDefined {
    return errors.New("mapreduced.yml not found in definition zip")
  }

  r.fetcher = fetcher
  r.zipReader = zipReader
  return nil
}

// WriteImageTar saves the .tar context for the template's base Docker image.
func (r *DefinitionReader) WriteImageTar(writer io.Writer) error {
  tarWriter := tar.NewWriter(writer)

  for _, zipFile := range r.zipReader.File {
    zipFileInfo := zipFile.FileInfo()
    header, err := tar.FileInfoHeader(zipFileInfo, "")
    if err != nil {
      return err
    }
    header.Name = zipFile.Name
    if err := tarWriter.WriteHeader(header); err != nil {
      return err
    }

    if zipFileInfo.IsDir() {
      continue
    }
    // TODO(pwnall): Create directory entries, if necessary.

    file, err := zipFile.Open()
    if err != nil {
      return err
    }
    fileBytes := make([]byte, zipFileInfo.Size())
    if _, err := file.Read(fileBytes); err != nil {
      return err
    }
    if _, err := tarWriter.Write(fileBytes); err != nil {
      return err
    }
  }
  if err := tarWriter.Close(); err != nil {
    return nil
  }
  return r.fetcher.Close()
}
