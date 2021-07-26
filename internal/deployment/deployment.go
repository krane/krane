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

// Deployment represent a Krane deployment and its configuration, current container resources, and job history
type Deployment struct {
	Config     Config           `json:"config"`
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

// GetDeployment returns a single deployment
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
// deployment creating or re-creating container resources
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

			e.emit(DeploymentSetup, "Preparing deployment environment")

			// ensure secrets collections
			if err := CreateSecretsCollection(deploymentName); err != nil {
				logger.Errorf("[ERROR CODE 1] unable to create secrets collection %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error setting up the deployment environment for %s: ERROR CODE 1", deploymentName))
				return err
			}

			// ensure jobs collections
			if err := CreateJobsCollection(deploymentName); err != nil {
				logger.Errorf("[ERROR CODE 2] unable to create jobs collection %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error setting up the deployment environment for %s: ERROR CODE 2", deploymentName))
				return err
			}

			// get containers (if any) currently part of this deployment
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("[ERROR CODE 3] unable to get containers %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error setting up the deployment environment for %s: ERROR CODE 3", deploymentName))
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
			e.emit(DeploymentPullImage, fmt.Sprintf("Pulling %s:%s", config.Image, config.Tag))
			pullImageReader, err := docker.GetClient().PullImage(config.Registry, config.Image, config.Tag)
			if err != nil {
				logger.Errorf("unable to pull image %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Unable to pull %s:%s: %v", config.Image, config.Tag, err))
				return err
			}
			e.emitStream(DeploymentPullImage, pullImageReader)

			// create containers
			e.emit(DeploymentContainerCreate, fmt.Sprintf("Creating %d container(s)", config.Scale))
			containersCreated := make([]KraneContainer, 0)
			for i := 0; i < config.Scale; i++ {
				c, err := ContainerCreate(config)
				if err != nil {
					logger.Errorf("unable to create container %v", err)
					e.emit(
						DeploymentError,
						fmt.Sprintf("Error creating container resources: %v", err))
					return err
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("%d/%d container(s) for deployment %s created", config.Scale, len(containersCreated), config.Name)

			// start containers
			e.emit(DeploymentContainerStart, "Starting deployment resources")
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					logger.Errorf("unable to start container %v", err)
					e.emit(
						DeploymentError,
						fmt.Sprintf("Error starting container resources: %v", err))
					return err
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("%d/%d container(s) for deployment %s started", len(containersStarted), len(containersCreated), config.Name)

			// health check
			e.emit(DeploymentHealthCheck, "Checking deployment health")
			retries := 10
			if err := RetriableContainersHealthCheck(containersStarted, retries); err != nil {
				logger.Errorf("containers did not pass health check %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Deployment healtcheck failed: %v", err))
				return err
			}
			logger.Debugf("Deployment %s health check complete", config.Name)
			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RunDeploymentJobArgs)

			e.emit(DeploymentCleanup, "Cleaning up unused resources")
			for _, c := range jobArgs.ContainersToRemove {
				err := c.Remove()
				if err != nil {
					logger.Errorf("unable to remove container %v", err)
					e.emit(
						DeploymentError,
						fmt.Sprintf("Error cleaning up unused resources: %v", err))
					return err
				}
			}

			e.emit(DeploymentDone, "Deployment complete")
			return nil
		},
	})

	return nil
}

// Delete removes a deployments container resources and configuration.
// Note: This will also remove any existing collections created for the deployment (Secrets, Jobs, Config etc...)
func Delete(deployment string) error {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return err
	}

	type DeleteDeploymentJobArgs struct {
		Deployment string
	}

	jobID := uuid.Generate().String()
	e := createEventEmitter(config.Name, jobID)
	go enqueue(job.Job{
		ID:          jobID,
		Deployment:  config.Name,
		Type:        string(DeleteDeploymentJobType),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: DeleteDeploymentJobArgs{
			Deployment: config.Name,
		},
		Run: func(args interface{}) error {
			jobArgs := args.(DeleteDeploymentJobArgs)
			deploymentName := jobArgs.Deployment

			e.emit(DeploymentContainerRemove, "Finding deployment resources to remove")

			// get current containers
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable get containers %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error finding deployment containers to remove: %v", err))
				return err
			}

			// remove containers
			e.emit(
				DeploymentContainerRemove,
				fmt.Sprintf("Removing %d container(s)", len(containers)))
			for _, c := range containers {
				if err := c.Remove(); err != nil {
					logger.Errorf("unable to remove container %v", err)
					e.emit(
						DeploymentError,
						fmt.Sprintf("Error removing deployment containers: %v", err))
					return err
				}
			}
			logger.Debugf("%d container(s) for deployment %s removed", len(containers), deploymentName)

			e.emit(DeploymentCleanup, "Container(s) successfully removed")

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(DeleteDeploymentJobArgs)
			deploymentName := jobArgs.Deployment

			e.emit(DeploymentCleanup, "Cleaning up remaining deployment resources")

			// delete secrets collection
			logger.Debugf("removing secrets collection for deployment %s", deploymentName)
			if err := DeleteSecretsCollection(deploymentName); err != nil {
				logger.Errorf("unable to remove secrets collection %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error cleaning up deployment resources: %v", err))
				return err
			}

			// delete jobs collection
			logger.Debugf("removing jobs collection for deployment %s", deploymentName)
			if err := DeleteJobsCollection(deploymentName); err != nil {
				logger.Errorf("unable to remove jobs collection %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error cleaning up deployment resources: %v", err))
				return err
			}

			// delete deployment configuration
			logger.Debugf("removing config for deployment %s", deploymentName)
			if err := DeleteConfig(deploymentName); err != nil {
				logger.Errorf("unable to remove deployment configuration %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error cleaning up deployment resources: %v", err))
				return err
			}

			e.emit(DeploymentDone, "Deployment deletion complete")

			return nil
		},
	})

	return nil
}

// StartContainers starts current existing containers (if any) for a deployment
// Note: this does not re-create container resources, only start existing ones
func StartContainers(deployment string) error {
	type StartContainersJobArgs struct {
		Deployment string
	}

	jobID := uuid.Generate().String()
	e := createEventEmitter(deployment, jobID)
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

			e.emit(DeploymentContainerStart, "Looking for container resources to start")

			// get current containers
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				logger.Errorf("unable to get containers %v", err)
				e.emit(
					DeploymentError,
					fmt.Sprintf("Error getting container resources: %v", err))
				return err
			}

			if len(containers) == 0 {
				e.emit(DeploymentDone, "Deployment has 0 containers to start up")
				return fmt.Errorf("deployment %s has 0 containers to start", deploymentName)
			}

			e.emit(DeploymentContainerStart,
				fmt.Sprintf("Starting %d container resources", len(containers)))

			// start containers
			for _, c := range containers {
				logger.Debugf("Starting container %s", c.Name)
				if err := c.Start(); err != nil {
					e.emit(
						DeploymentError,
						fmt.Sprintf("Error starting container resources: %v", err))
					logger.Errorf("unable to start container %v", err)
					return err
				}
			}
			logger.Debugf("%d container(s) for deployment %s started", len(containers), deploymentName)

			e.emit(
				DeploymentDone,
				fmt.Sprintf("%s container resources started", len(containers)))

			return nil
		},
	})
	return nil
}

// StopContainers stops current existing containers (if any) for a deployment
// Note: this does not re-create container resources, only stop existing ones
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
				logger.Debugf("Stopping container %s", c.Name)
				if err := c.Stop(); err != nil {
					logger.Errorf("unable to stop container %v", err)
					return err
				}
			}
			logger.Debugf("%d container(s) for deployment %s stopped", len(containers), deploymentName)

			return nil
		},
	})
	return nil
}

// RestartContainers will re-create container resources for a deployment
// Note: this almost the same call as 'Run' since they both re-create container resources based on the current configuration
func RestartContainers(deployment string) error {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return fmt.Errorf("unable to get configuration for deployment %s", deployment)
	}

	type RestartContainersJobArgs struct {
		Config             Config
		ContainersToRemove []KraneContainer
	}

	jobID := uuid.Generate().String()
	e := createEventEmitter(config.Name, jobID)
	go enqueue(job.Job{
		ID:          jobID,
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

			// pull image
			logger.Debugf("Pulling image for deployment %s", config.Name)
			pullImageReader, err := docker.GetClient().PullImage(config.Registry, config.Image, config.Tag)
			if err != nil {
				logger.Errorf("unable to pull image %v", err)
				return err
			}
			e.emitStream(DeploymentPullImage, pullImageReader)

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
			logger.Debugf("%d/%d container(s) for deployment %s created", len(containersCreated), config.Scale, config.Name)

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

			retries := 10
			if err := RetriableContainersHealthCheck(containersStarted, retries); err != nil {
				logger.Errorf("containers did not pass health check %v", err)
				return err
			}
			logger.Debugf("Deployment %s health check complete", config.Name)

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RestartContainersJobArgs)
			for _, c := range jobArgs.ContainersToRemove {
				logger.Debugf("Removing container %s", c.Name)
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
