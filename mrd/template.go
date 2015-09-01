package mrd

import (
  "errors"
)

// Information common to Map-Reduce jobs.
type Template struct {
  // The descriptor fields are only initialzed if State >= TemplateRead.
  TemplateDescriptor

  // Identifier assigned when the template is created
  ID string
  // The current state of the template
  State TemplateState
  // The first error that occurred while processing the template
  Err error
  // The output produced by Docker while building the base Docker image
  BuildOutput []byte
  // The ID of the template's base Docker image
  BaseImageID string
}

type TemplateState int

const (
  // Init was not called on this template.
  InvalidTemplateState = 0
  // The template was assigned an ID, and is blank otherwise.
  TemplateAllocated = 1
  // The template's YAML definition was read into the descriptor.
  TemplateDefined = 1
  // The template's base Docker image was created.
  TemplateBuilt = 2
  // The template's base Docker image is ready for use.
  TemplateReady = 3
  // The template's base Docker image is pending deletion.
  TemplateDeleting = 4
  // Errors were encountered while processing the template.
  TemplateErrored = 5
)

// Init gets the template to the TemplateAllocated state.
func (t *Template) Init(id string) error {
  if t.State != InvalidTemplateState {
    err := errors.New("Incorrect template state for Init")
    t.SetError(err)
    return err
  }
  t.ID = id
  t.State = TemplateAllocated
  return nil
}

// ReadDefinition gets the template to the TemplateRead state.
func (t *Template) ReadDefinition(templateYaml []byte) error {
  if t.State != TemplateAllocated {
    err := errors.New("Incorrect template state for ReadDefinition")
    t.SetError(err)
    return err
  }
  if err := t.TemplateDescriptor.Init(templateYaml); err != nil {
    t.SetError(err)
    return err
  }
  t.State = TemplateDefined
  return nil
}

// SetBuildOutput gets the template to the TemplateBuilt state.
func (t *Template) SetBuildOutput(buildOutput []byte) error {
  if t.State != TemplateDefined {
    err := errors.New("Incorrect template state for SetBuildOutput")
    t.SetError(err)
    return err
  }
  t.BuildOutput = buildOutput
  t.State = TemplateBuilt
  return nil
}

// ReadDefinition gets the template to the TemplateReady state.
func (t *Template) SetBaseImage(baseImageID string) error {
  if t.State != TemplateBuilt {
    err := errors.New("Incorrect template state for SetBaseImage")
    t.SetError(err)
    return err
  }
  t.BaseImageID = baseImageID
  t.State = TemplateReady
  return nil
}


// SetError gets the template to the TemplateErrored state.
// Templates ignore new errors if they're already in the TemplateErrored state.
func (t *Template) SetError(err error) {
  if t.State == TemplateErrored {
    return
  }
  t.State = TemplateErrored
  t.Err = err
}
