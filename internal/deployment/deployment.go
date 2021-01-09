package deployment

import (
	"fmt"

	"github.com/docker/distribution/uuid"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/job"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/utils"
)

type Deployment struct {
	Config
	Containers []KraneContainer `json:"containers"`
	Jobs       []job.Job        `json:"jobs"`
}

// Exist returns true if a deployment exist, false otherwise
func Exist(deployment string) bool {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return false
	}

	if config.Empty() {
		return false
	}

	return true
}

// GetAllDeployments returns a single deployment
func GetDeployment(deployment string) (Deployment, error) {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return Deployment{}, err
	}

	containers, err := GetContainersByDeployment(deployment)
	if err != nil {
		return Deployment{}, err
	}

	jobs, err := GetJobsByDeployment(deployment, uint(utils.OneYear))
	if err != nil {
		return Deployment{}, err
	}

	return Deployment{
		Config:     config,
		Containers: containers,
		Jobs:       jobs,
	}, nil
}

// GetAllDeployments returns a list of all deployments
func GetAllDeployments() ([]Deployment, error) {
	configs, err := GetAllDeploymentConfigs()
	if err != nil {
		return []Deployment{}, err
	}

	deployments := make([]Deployment, 0)
	for _, config := range configs {
		d, err := GetDeployment(config.Name)
		if err != nil {
			return []Deployment{}, err
		}

		deployments = append(deployments, d)
	}

	return deployments, nil
}

// Run a deployment runs the current configuration for a
// deployment creating or re-creating container resources.
// Returns the job id to track deployment events
func Run(deployment string) error {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return err
	}

	type RunDeploymentJobArgs struct {
		Config             Config
		ContainersToRemove []KraneContainer
	}

	jobID := uuid.Generate().String()
	e := createEventEmitter(config.Name, jobID)
	go enqueue(job.Job{
		ID:          jobID,
		Deployment:  config.Name,
		Type:        string(RunDeploymentJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: &RunDeploymentJobArgs{
			Config:             config,
			ContainersToRemove: []KraneContainer{},
		},
		Setup: func(args interface{}) error {
			jobArgs := args.(*RunDeploymentJobArgs)
			deploymentName := jobArgs.Config.Name

			// ensure secrets collections
			if err := CreateSecretsCollection(deploymentName); err != nil {
				logger.Errorf("unable to create secrets collection %v", err)
				return err
			}

			// ensure jobs collections
			if err := CreateJobsCollection(deploymentName); err != nil {
				logger.Errorf("unable to create jobs collection %v", err)
				return err
			}

			// get containers (if any) currently part of this deployment
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable to get containers %v", err)
				return err
			}

			// update job arguments to process them for deletion later on
			jobArgs.ContainersToRemove = containers

			return nil
		},
		Run: func(args interface{}) error {
			jobArgs := args.(*RunDeploymentJobArgs)
			config := jobArgs.Config

			// pull image
			logger.Debugf("Pulling image for deployment %s", config.Name)
			pullImageReader, err := docker.GetClient().PullImage(config.Registry, config.Image, config.Tag)
			if err != nil {
				logger.Errorf("unable to pull image %v", err)
				return err
			}
			e.emitS(pullImageReader)

			// create containers
			containersCreated := make([]KraneContainer, 0)
			for i := 0; i < config.Scale; i++ {
				c, err := ContainerCreate(config)
				if err != nil {
					logger.Errorf("unable to create container %v", err)
					return err
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("%d/%d container(s) for deployment %s created", config.Scale, len(containersCreated), config.Name)

			// start containers
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					logger.Errorf("unable to start container %v", err)
					return err
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("%d/%d container(s) for deployment %s started", len(containersStarted), len(containersCreated), config.Name)

			// health check
			retries := 10
			if err := RetriableContainerHealthCheck(containersStarted, retries); err != nil {
				logger.Errorf("containers did not pass health check %v", err)
				return err
			}
			logger.Debugf("Deployment %s health check complete", config.Name)
			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RunDeploymentJobArgs)

			for _, c := range jobArgs.ContainersToRemove {
				logger.Debugf("Purging container %s", c.Name)
				err := c.Remove()
				if err != nil {
					logger.Errorf("unable to remove container %v", err)
					return err
				}
			}
			return nil
		},
	})

	return nil
}

// Delete removes a deployments containers and configuration. This will also remove any existing
// collections created for the deployment (Secrets, Jobs, etc...)
func Delete(deployment string) error {
	type DeleteDeploymentJobArgs struct {
		Deployment string
	}

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  deployment,
		Type:        string(DeleteDeploymentJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: DeleteDeploymentJobArgs{
			Deployment: deployment,
		},
		Run: func(args interface{}) error {
			jobArgs := args.(DeleteDeploymentJobArgs)
			deploymentName := jobArgs.Deployment

			// get current containers
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable get containers %v", err)
				return err
			}

			// remove containers
			for _, c := range containers {
				if err := c.Remove(); err != nil {
					logger.Errorf("unable to remove container %v", err)
					return err
				}
			}

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(DeleteDeploymentJobArgs)
			deploymentName := jobArgs.Deployment

			// delete secrets collection
			logger.Debugf("removing secrets collection for deployment %s", deploymentName)
			if err := DeleteSecretsCollection(deploymentName); err != nil {
				logger.Errorf("unable to remove secrets collection %v", err)
				return err
			}

			// delete jobs collection
			logger.Debugf("removing jobs collection for deployment %s", deploymentName)
			if err := DeleteJobsCollection(deploymentName); err != nil {
				logger.Errorf("unable to remove jobs collection %v", err)
				return err
			}

			// delete deployment configuration
			logger.Debugf("removing config for deployment %s", deploymentName)
			if err := DeleteConfig(deploymentName); err != nil {
				logger.Errorf("unable to remove deployment configuration %v", err)
				return err
			}

			return nil
		},
	})

	return nil
}

func StartContainers(deployment string) error {
	type StartContainersJobArgs struct {
		Deployment string
	}

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  deployment,
		Type:        string(StartContainersJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: StartContainersJobArgs{
			Deployment: deployment,
		},
		Run: func(args interface{}) error {
			jobArgs := args.(StartContainersJobArgs)
			deploymentName := jobArgs.Deployment

			// get current containers
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable to get containers %v", err)
				return err
			}

			if len(containers) == 0 {
				return fmt.Errorf("deployment %s has 0 containers to start", deploymentName)
			}

			// start containers
			for _, c := range containers {
				if err := c.Start(); err != nil {
					logger.Errorf("unable to start container %v", err)
					return err
				}
			}

			return nil
		},
	})
	return nil
}

func StopContainers(deployment string) error {
	type StopContainersJobArgs struct {
		Deployment string
	}

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  deployment,
		Type:        string(StopContainersJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: StopContainersJobArgs{
			Deployment: deployment,
		},
		Run: func(args interface{}) error {
			jobArgs := args.(StopContainersJobArgs)
			deploymentName := jobArgs.Deployment

			// get current containers
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable to get containers %v", err)
				return err
			}

			if len(containers) == 0 {
				return fmt.Errorf("deployment %s has 0 containers to stop", deploymentName)
			}

			// stop containers
			for _, c := range containers {
				if err := c.Stop(); err != nil {
					logger.Errorf("unable to stop container %v", err)
					return err
				}
			}

			return nil
		},
	})
	return nil
}

func RestartContainers(deployment string) error {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return fmt.Errorf("unable to get configuration for deployment %s", deployment)
	}

	type RestartContainersJobArgs struct {
		Config             Config
		ContainersToRemove []KraneContainer
	}

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  deployment,
		Type:        string(RestartContainersJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: &RestartContainersJobArgs{
			ContainersToRemove: []KraneContainer{},
			Config:             config,
		},
		Setup: func(args interface{}) error {
			jobArgs := args.(*RestartContainersJobArgs)
			deploymentName := jobArgs.Config.Name

			// get current containers (if any) which will be removed after new containers are created
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable to get containers %v", err)
				return err
			}

			jobArgs.ContainersToRemove = containers
			return nil
		},
		Run: func(args interface{}) error {
			jobArgs := args.(*RestartContainersJobArgs)
			config := jobArgs.Config

			// create containers
			containersCreated := make([]KraneContainer, 0)
			for i := 0; i < config.Scale; i++ {
				c, err := ContainerCreate(config)
				if err != nil {
					logger.Errorf("unable to create container %v", err)
					return err
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("Created %d container(s) for deployment %s", len(containersCreated), config.Name)

			// start containers
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					logger.Errorf("unable to start container %v", err)
					return err
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("Started %d container(s) for deployment %s", len(containersStarted), config.Name)

			retries := 10
			if err := RetriableContainerHealthCheck(containersStarted, retries); err != nil {
				logger.Errorf("containers did not pass health check %v", err)
				return err
			}

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RestartContainersJobArgs)
			for _, c := range jobArgs.ContainersToRemove {
				if err := c.Remove(); err != nil {
					logger.Errorf("unable to remove container %v", err)
					return err
				}
			}

			return nil
		},
	})
	return nil
}
