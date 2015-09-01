package mrd

import (
  "testing"
)

func TestController_BuildTemplate(t *testing.T) {
  var controller Controller
  controller.Init("mapreduced0")

  template, err := controller.FindTemplate("hello")
  if err != nil {
    t.Fatal(err)
  }

  fetcher := NewFileFetcher("../test_tmp/hello.zip")

  if err := controller.BuildTemplate(template, fetcher); err != nil {
    t.Fatal(err)
  }

  if template.State != TemplateReady {
    t.Errorf("Wrong template state: %v", template.State)
  }
  if !fetcher.Closed() {
    t.Error("BuildTemplate did not close the Fetcher")
  }

  image, err := controller.docker.InspectImage("mapreduced0/base_hello")
  if err != nil {
    t.Errorf("Built image not found: %v", err)
  }
  if image.ID != template.BaseImageID {
    t.Errorf("Incorrect BaseImageID: %v, expected %v", image.ID,
             template.BaseImageID)
  }

  if err := controller.CleanOldState(); err != nil {
    t.Fatal(err)
  }
  _, err = controller.docker.InspectImage("mapreduced0_hello/base")
  if err == nil {
    t.Error(
        "mapreduced0_base_hello was not deleted by the controller cleanup")
  }
}

func TestController_NewJob(t *testing.T) {
  var controller Controller
  controller.Init("mapreduced")

  job1, err := controller.NewJob()
  if err != nil {
    t.Fatal(err)
  }
  job2, err := controller.NewJob()
  if err != nil {
    t.Fatal(err)
  }

  if job1.ID != "1" || job2.ID != "2" {
    t.Errorf("Incorrect job IDs: %v, %v", job1.ID, job2.ID)
  }
  if job1.MapperImageName != "mapreduced/map_1" {
    t.Errorf("Incorrect job namePrefix. Mapper image: %v",
             job1.MapperImageName)
  }
}

func TestController_TemplateImageName(t *testing.T) {
  var controller Controller
  controller.Init("mapreduced")

  var template Template
  template.ID = "fib"

  imageName := controller.TemplateImageName(&template)
  if imageName != "mapreduced/base_fib" {
    t.Errorf("Incorrect image name: %v", imageName)
  }
}
