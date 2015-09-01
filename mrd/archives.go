package mrd

import (
  "archive/tar"
  "archive/zip"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
)

// TarDirectory compresses a directory into a tar archive.
func TarDirectory(dirPath string, writer io.Writer) error {
  if dirPath[len(dirPath) - 1] != '/' {
    dirPath += "/"
  }
  dirPathLen := len(dirPath)
  tarWriter := tar.NewWriter(writer)

  walkFn := func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    header, err := tar.FileInfoHeader(info, "")
    if err != nil {
      return err
    }
    if strings.HasPrefix(path, dirPath) {
      header.Name = path[dirPathLen:]
    } else {
      return nil
    }
    if err := tarWriter.WriteHeader(header); err != nil {
      return err
    }

    if info.IsDir() {
      return nil
    }
    fileBytes, err := ioutil.ReadFile(path)
    if err != nil {
      return err
    }
    if _, err := tarWriter.Write(fileBytes); err != nil {
      return err
    }
    return nil
  }
  if err := filepath.Walk(dirPath, walkFn); err != nil {
    return err
  }

  return tarWriter.Close()
}

// ZipDirectory compresses a directory into a zip archive.
func ZipDirectory(dirPath string, writer io.Writer) error {
  if dirPath[len(dirPath) - 1] != '/' {
    dirPath += "/"
  }
  dirPathLen := len(dirPath)
  zipWriter := zip.NewWriter(writer)

  walkFn := func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if info.IsDir() {
      return nil
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
      return err
    }
    if strings.HasPrefix(path, dirPath) {
      header.Name = path[dirPathLen:]
    } else {
      return nil
    }
    zipFileWriter, err := zipWriter.CreateHeader(header)
    if err != nil {
      return err
    }

    file, err := os.Open(path)
    if err != nil {
      return err
    }
    defer file.Close()
    if _, err := io.Copy(zipFileWriter, file); err != nil {
      return err
    }
    return nil
  }
  if err := filepath.Walk(dirPath, walkFn); err != nil {
    return err
  }

  return zipWriter.Close()
}
