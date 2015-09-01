package mrd

import (
  "io/ioutil"
  "reflect"
  "testing"
)

func TestTemplateDescriptor_Init(t *testing.T) {
  templateYaml, err := ioutil.ReadFile("../testdata/hello/mapreduced.yml")
  if err != nil {
    t.Fatal(err)
  }
  var descriptor TemplateDescriptor
  if err := descriptor.Init(templateYaml); err != nil {
    t.Fatal(err)
  }

  if descriptor.ItemCount != 3 {
    t.Errorf("Incorrect item count: %v", descriptor.ItemCount)
  }

  if descriptor.Mapper.WorkDir != "/usr/mrd" {
    t.Errorf("Incorrect map work directory: %v", descriptor.Mapper.WorkDir)
  }
  if descriptor.Mapper.InputPath != "/usr/mrd/map-input" {
    t.Errorf("Incorrect map input path: %v", descriptor.Mapper.InputPath)
  }
  if descriptor.Mapper.OutputPath != "/usr/mrd/map-output" {
    t.Errorf("Incorrect map output path: %v", descriptor.Mapper.OutputPath)
  }
  if !reflect.DeepEqual(descriptor.Mapper.EntryPoint,
                        []string{ "/bin/sh", "/usr/mrd/mapper.sh" }) {
    t.Errorf("Incorrect map entry point: %v", descriptor.Mapper.EntryPoint)
  }

  if descriptor.Reducer.WorkDir != "/" {
    t.Errorf("Incorrect reduce work directory: %v", descriptor.Reducer.WorkDir)
  }
  if descriptor.Reducer.InputPath != "/usr/mrd/" {
    t.Errorf("Incorrect reduce input path: %v", descriptor.Reducer.InputPath)
  }
  if descriptor.Reducer.OutputPath != "/usr/mrd/reduce-output" {
    t.Errorf("Incorrect reduce output path: %v", descriptor.Reducer.OutputPath)
  }
  if !reflect.DeepEqual(descriptor.Reducer.EntryPoint,
                        []string{ "/bin/sh", "/usr/mrd/reducer.sh" }) {
    t.Errorf("Incorrect reduce entry point: %v", descriptor.Reducer.EntryPoint)
  }
}

func TestDockerDescriptor_GetDockerFile(t *testing.T) {
  templateYaml, err := ioutil.ReadFile("../testdata/hello/mapreduced.yml")
  if err != nil {
    t.Fatal(err)
  }
  var descriptor TemplateDescriptor
  if err := descriptor.Init(templateYaml); err != nil {
    t.Fatal(err)
  }

  mapperGolden, err := ioutil.ReadFile("../testdata/Dockerfile.hello.mapper")
  if err != nil {
    t.Fatal(err)
  }
  mapperFile := descriptor.Mapper.GetDockerFile("mapreduced/hello_base",
                                                "input")
  if string(mapperGolden) != mapperFile {
    t.Errorf("Incorrect mapper Dockerfile: %v", mapperFile)
  }

  reducerGolden, err := ioutil.ReadFile("../testdata/Dockerfile.hello.reducer")
  if err != nil {
    t.Fatal(err)
  }
  reducerFile := descriptor.Reducer.GetDockerFile("mapreduced/hello_base", ".")
  if string(reducerGolden) != reducerFile {
    t.Errorf("Incorrect reducer Dockerfile: %v", reducerFile)
  }
}
