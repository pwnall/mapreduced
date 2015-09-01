package mrd

type Job struct {
  // The job's unique ID.
  ID string

  // Names for the Docker images created by this job.
  MapperImageName string
  ReducerImageName string

  // The job description submitted to the Web service.
  Template *Template

  // Prepended to all Docker objects belonging to the job's controller.
  dockerNamePrefix string
}


// Init sets the state to reflect a newly created job.
func (j *Job) Init(id string, dockerNamePrefix string) error {
  j.ID = id
  j.dockerNamePrefix = dockerNamePrefix

  j.MapperImageName = dockerNamePrefix + "/map_" + id
  j.ReducerImageName = dockerNamePrefix + "/reduce_" + id
  return nil
}
