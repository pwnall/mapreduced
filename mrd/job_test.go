package mrd

import (
  "testing"
)

func TestJob_Init(t *testing.T) {
  var job Job
  if err := job.Init("42", "mrd1337"); err != nil {
    t.Fatal(err)
  }

  if job.MapperImageName != "mrd1337/map_42" {
    t.Error("Incorrect mapper image name: %v", job.MapperImageName)
  }
  if job.ReducerImageName != "mrd1337/reduce_42" {
    t.Error("Incorrect reducer image name: %v", job.ReducerImageName)
  }
}
