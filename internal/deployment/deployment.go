package deployment

import (
	"fmt"

	"github.com/docker/distribution/uuid"
	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type DeploymentAction string

const (
	RunDeploymentAction     DeploymentAction = "RUN_DEPLOYMENT"
	DeleteDeploymentAction  DeploymentAction = "DELETE_DEPLOYMENT"
	StopContainersAction    DeploymentAction = "STOP_CONTAINERS"
	StartContainersAction   DeploymentAction = "START_CONTAINERS"
	RestartContainersAction DeploymentAction = "RESTART_CONTAINERS"
)

// enqueue queue's up deployment job for processing
func enqueue(j job.Job) {
	enqueuer := job.NewEnqueuer(job.Queue())
	queuedJob, err := enqueuer.Enqueue(j)
	if err != nil {
		logger.Errorf("Error enqueuing deployment job %v", err)
		return
	}
	logger.Debugf("Deployment job %s queued for processing", queuedJob.Deployment)
	return
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

// Save a deployment configuration into the db
func Save(config Config) error {
	config.applyDefaults()

	if err := config.isValid(); err != nil {
		return err
	}

	bytes, _ := config.Serialize()
	return store.Client().Put(constants.DeploymentsCollectionName, config.Name, bytes)
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

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  config.Name,
		Type:        string(RunDeploymentAction),
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error creating secrets collection for deployment %s during job execution", deploymentName))
			}

			// ensure jobs collections
			if err := CreateJobsCollection(deploymentName); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error creating jobs collection for deployment %s during job execution", deploymentName))
			}

			// get containers (if any) currently part of this deployment
			containers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting containers for deployment %s during job execution", deploymentName))
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
			err := docker.GetClient().PullImage(config.Registry, config.Image, config.Tag)
			if err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error pulling image %s for deployment %s during job execution", config.Image, config.Name))
			}

			// create containers
			containersCreated := make([]KraneContainer, 0)
			for i := 0; i < config.Scale; i++ {
				c, err := ContainerCreate(config)
				if err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error creating container for deployment %s during job execution", config.Name))
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("Created %d container(s) for deployment %s", len(containersCreated), config.Name)

			// start containers
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error starting container %s for deployment %s during job execution", c.ID, config.Name))
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("Started %d container(s) for deployment %s", len(containersStarted), config.Name)

			retries := 10
			if err := RetriableContainerHealthCheck(containersStarted, retries); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error checking containers health for deployment %s during job execution", config.Name))
			}

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RunDeploymentJobArgs)

			for _, c := range jobArgs.ContainersToRemove {
				err := c.Remove()
				if err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error removing container %s for deployment %s during job execution", c.ID, config.Name))
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
		Type:        string(DeleteDeploymentAction),
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting containers for deployment %s", deploymentName))
			}

			// remove containers
			for _, c := range containers {
				if err := c.Remove(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error removing container %s for deployment %s", c.ID, deploymentName))
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing secrets collection for deployment %s", deploymentName))
			}

			// delete jobs collection
			logger.Debugf("removing jobs collection for deployment %s", deploymentName)
			if err := DeleteJobsCollection(deploymentName); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing jobs collection for deployment %s", deploymentName))
			}

			// delete deployment configuration
			logger.Debugf("removing config for deployment %s", deploymentName)
			if err := DeleteDeploymentConfig(deploymentName); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing configuration for deployment %s", deploymentName))
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
		Type:        string(StartContainersAction),
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting containers for deployment %s", deploymentName))
			}

			if len(containers) == 0 {
				return fmt.Errorf("deployment %s has 0 containers to start", deploymentName)
			}

			// start containers
			for _, c := range containers {
				if err := c.Start(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error starting container %s for deployment %s", c.ID, deploymentName))
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
		Type:        string(StopContainersAction),
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting containers for deployment %s", deploymentName))
			}

			if len(containers) == 0 {
				return fmt.Errorf("deployment %s has 0 containers to stop", deploymentName)
			}

			// stop containers
			for _, c := range containers {
				if err := c.Stop(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error stopping container %s for deployment %s", c.ID, deploymentName))
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
		Type:        string(RestartContainersAction),
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
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting oldContainers for %s during job execution", deploymentName))
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
					return errors.Wrap(
						err,
						fmt.Sprintf("error creating container %s for deployment %s during job execution", c.ID, config.Name))
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("Created %d container(s) for deployment %s", len(containersCreated), config.Name)

			// start containers
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error starting container %s for deployment %s during job execution", c.ID, config.Name))
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("Started %d container(s) for deployment %s", len(containersStarted), config.Name)

			retries := 10
			if err := RetriableContainerHealthCheck(containersStarted, retries); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error checking containers health for deployment %s during job execution", config.Name))
			}

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*RestartContainersJobArgs)
			config := jobArgs.Config

			for _, c := range jobArgs.ContainersToRemove {
				if err := c.Remove(); err != nil {
					return errors.Wrap(
						err,
						fmt.Sprintf("error removing container %s for deployment %s during job execution", c.ID, config.Name))
				}
			}

			return nil
		},
	})
	return nil
}
