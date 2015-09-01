package mrd

import (
  "bytes"
  "strings"
  "strconv"

  dockerclient "github.com/fsouza/go-dockerclient"
)

type Controller struct {
  // The prefix added to Docker container and image names.
  // Used to prevent clashes with other images.
  NamePrefix string

  // The value for the mapreduced.ctl label added to containers and images.
  // Used to clean up objects left over from previous runs of the controller.
  LabelValue string

  // Available job templates, indexed by ID.
  templates map[string]*Template

  // The number assigned to the next created Job.
  NextJobId int

  // Used to talk to a Docker API server (Docker or Swarm).
  docker *dockerclient.Client
}

// Init sets up the controller's Docker connection.
func (c *Controller) Init(namePrefix string) error {
  var err error
  if c.docker, err = dockerclient.NewClientFromEnv(); err != nil {
    return err
  }

  c.NamePrefix = namePrefix
  c.NextJobId = 1
  c.templates = make(map[string]*Template)
  return nil
}

// FindTemplate locates or creates the job template with the given ID.
func (c *Controller) FindTemplate(id string) (*Template, error) {
  if template := c.templates[id]; template != nil {
    return template, nil
  }
  template := new(Template)
  if err := template.Init(id); err != nil {
    return nil, err
  }
  c.templates[id] = template
  return template, nil
}

// BuildTemplate gets a Template from Allocated to Ready.
func (c *Controller) BuildTemplate(template *Template, fetcher Fetcher) error {
  var reader DefinitionReader
  if err := reader.Init(fetcher, template); err != nil {
    return err
  }

  var tarBuffer bytes.Buffer
  if err := reader.WriteImageTar(&tarBuffer); err != nil {
    return err
  }

  var outputBuffer bytes.Buffer
  imageName := c.TemplateImageName(template)
  err := c.docker.BuildImage(dockerclient.BuildImageOptions{
    Name: imageName, InputStream: &tarBuffer, OutputStream: &outputBuffer,
  })
  if err != nil {
    return err
  }
  if err := template.SetBuildOutput(outputBuffer.Bytes()); err != nil {
    return err
  }

  image, err := c.docker.InspectImage(imageName)
  if err != nil {
    return err
  }

  if err := template.SetBaseImage(image.ID); err != nil {
    return err
  }

  return nil
}

// NewJob creates a job to be run by this controller, from a YAML description.
func (c *Controller) NewJob() (*Job, error) {
  var job Job
  if err := job.Init(strconv.Itoa(c.NextJobId), c.NamePrefix); err != nil {
    return nil, err
  }

  c.NextJobId += 1
  return &job, nil
}

func (c *Controller) RunJob(job *Job) error {

  return nil
}

// CleanOldState removes all the Docker objects left over from a previous run.
func (c *Controller) CleanOldState() error {
  if err := c.cleanOldContainers(); err != nil {
    return err
  }
  return c.cleanOldImages()
}

// cleanOldContainers removes Docker containers left over from a previous run.
func (c *Controller) cleanOldContainers() error {
  containers, err := c.docker.ListContainers(
      dockerclient.ListContainersOptions{
      All: false, Filters: c.DockerFilters(),
      })
  if err != nil {
    return err
  }

  for _, container := range containers {
    err = c.docker.RemoveContainer(dockerclient.RemoveContainerOptions{
      ID: container.ID, RemoveVolumes: true, Force: false,
    })

    if err != nil {
      return err
    }
  }

  return nil
}

// cleanOldImages removes Docker images left over from a previous run.
func (c *Controller) cleanOldImages() error {
  images, err := c.docker.ListImages(dockerclient.ListImagesOptions{
    All: false,
  })
  if err != nil {
    return err
  }

  imageNamePrefix := c.NamePrefix + "/"
  for _, image := range images {
    repoTag := ""
    for _, tag := range image.RepoTags {
      if strings.HasPrefix(tag, imageNamePrefix) {
        repoTag = tag
        break
      }
    }
    if repoTag == "" {
      continue
    }

    err = c.docker.RemoveImageExtended(image.ID,
        dockerclient.RemoveImageOptions{NoPrune: false, Force: false,})

    if err != nil {
      return err
    }
  }

  return nil
}

// TemplateImageName computes the name of a template's base Docker image.
func (c *Controller) TemplateImageName(template *Template) string {
  return c.NamePrefix + "/base_" + template.ID
}

// MapperImageName computes the name of a job's mapper image.
func (c *Controller) MapperImageName(job *Job) string {
  return c.NamePrefix + "/map_" + job.ID
}

// MapperImageName computes the name of a job's mapper image.
func (c *Controller) ReducerImageName(job *Job) string {
  return c.NamePrefix + "/reduce_" + job.ID
}

// DockerFilters returns the filters used to find this controller's objects.
func (c *Controller) DockerFilters() map[string][]string {
  labelFilter := "mapreduced.ctl=\"" + c.LabelValue + "\""

  return map[string][]string{"label": []string{labelFilter}}
}
