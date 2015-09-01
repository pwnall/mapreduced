package mrd

import (
  "encoding/json"

  yaml "gopkg.in/yaml.v2"
)

// The template information provided in the Yaml file.
type TemplateDescriptor struct {
  // Number of times the mapper will run.
  ItemCount int `yaml:"items"`

  // Information about the map phase of the job.
  Mapper MapperDescriptor `yaml:"mapper"`

  // Information about the reduce phase of the job.
  Reducer ReducerDescriptor `yaml:"reducer"`
}

type DockerDescriptor struct {
  // Executable + paramters used to run the Docker image.
  EntryPoint []string `yaml:"cmd"`

  // Current directory when the Docker image starts.
  WorkDir string `yaml:"chdir"`

  // Path where the job's input will be copied to.
  InputPath string `yaml:"input"`

  // Path where the job's output will be copied from.
  OutputPath string `yaml:"output"`
}

// Information about the mapper in a Map-Reduce job.
type MapperDescriptor struct {
  DockerDescriptor `yaml:",inline"`

  // The name of the environment variable holding the current item's index.
  ItemEnvVar string `yaml:"env"`
}

// Information about the reducer in a Map-Reduce job.
type ReducerDescriptor struct {
  DockerDescriptor `yaml:",inline"`

  // The name of the environment variable holding the number of items.
  ItemCountVar string `yaml:"env"`
}

// Init reads a YAML template definition into the descriptor's fields.
func (d *TemplateDescriptor) Init(templateYaml []byte) error {
  if err := yaml.Unmarshal(templateYaml, &d); err != nil {
    return err
  }
  return nil
}

// GetDockerFile builds the Dockerfile that creates a phase container's image.
func (s *DockerDescriptor) GetDockerFile(sourceImage string,
                                         copyFrom string) string {
  fromLine := "FROM " + sourceImage + "\n"
  copyLine := "COPY " + copyFrom + " " + s.InputPath + "\n"
  workdirLine := "WORKDIR " + s.WorkDir + "\n"

  entryPointBytes, _ := json.Marshal(s.EntryPoint)
  entryLine := "ENTRYPOINT " + string(entryPointBytes) + "\n"

  return fromLine + copyLine + workdirLine + entryLine
}
