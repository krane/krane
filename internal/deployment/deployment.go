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
	CreateContainers  DeploymentAction = "CREATE_CONTAINERS"
	DeleteContainers  DeploymentAction = "DELETE_CONTAINERS"
	StopContainers    DeploymentAction = "STOP_CONTAINERS"
	StartContainers   DeploymentAction = "START_CONTAINERS"
	ReStartContainers DeploymentAction = "RESTART_CONTAINERS"
)

// enqueue queue's up deployment job for processing
func enqueue(deploymentJob job.Job) {
	enqueuer := job.NewEnqueuer(job.Queue())
	queuedJob, err := enqueuer.Enqueue(deploymentJob)
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

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  config.Name,
		Type:        string(ReStartContainers),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: &JobArgs{
			Config:        config,
			OldContainers: []KraneContainer{},
			NewContainers: []KraneContainer{},
		},
		Setup: func(args interface{}) error {
			jobArgs := args.(*JobArgs)
			deploymentName := jobArgs.Config.Name

			// ensure secrets collections
			if err := CreateSecretsCollection(deploymentName); err != nil {
				return err
			}

			// ensure jobs collections
			if err := job.CreateCollection(deploymentName); err != nil {
				return err
			}

			// get oldContainers (if any) part of a previous deployment run
			oldContainers, err := GetContainersByDeployment(deploymentName)
			if err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error getting oldContainers for %s during job execution", deploymentName))
			}

			// update job arguments to process them for deletion later on
			jobArgs.OldContainers = oldContainers

			return nil
		},
		Run: func(args interface{}) error {
			jobArgs := args.(*JobArgs)
			config := jobArgs.Config

			// pull image
			logger.Debugf("Pulling image for deployment %s", config.Name)
			dockerPullImageErr := docker.GetClient().PullImage(config.Registry, config.Image, config.Tag)
			if dockerPullImageErr != nil {
				return dockerPullImageErr
			}

			// create containers
			containersCreated := make([]KraneContainer, 0)
			for i := 0; i < config.Scale; i++ {
				c, err := ContainerCreate(config)
				if err != nil {
					return err
				}
				containersCreated = append(containersCreated, c)
			}
			logger.Debugf("Created %d container(s) for deployment %s", len(containersCreated), config.Name)

			// start containers
			containersStarted := make([]KraneContainer, 0)
			for _, c := range containersCreated {
				if err := c.Start(); err != nil {
					return err
				}
				containersStarted = append(containersStarted, c)
			}
			logger.Debugf("Started %d container(s) for deployment %s", len(containersStarted), config.Name)

			// update job arguments to reflect recently started containers
			jobArgs.NewContainers = containersStarted

			// check recently started containers are in a running state
			retries := 10
			if err := RetriableContainerHealthCheck(containersStarted, retries); err != nil {
				return err
			}

			return nil
		},
		Finally: func(args interface{}) error {
			jobArgs := args.(*JobArgs)

			for _, c := range jobArgs.OldContainers {
				err := c.Remove()
				if err != nil {
					return err
				}
			}

			return nil
		},
	})

	return nil
}

// Delete removes a deployments containers and configuration. This will also remove any existing
// collections created for the deployment (Secrets, Logs, etc...)
func Delete(deployment string) error {
	config, err := GetDeploymentConfig(deployment)
	if err != nil {
		return err
	}

	go enqueue(job.Job{
		ID:          uuid.Generate().String(),
		Deployment:  config.Name,
		Type:        string(DeleteContainers),
		RetryPolicy: utils.UIntEnv(constants.EnvDeploymentRetryPolicy),
		Args: JobArgs{
			Config: config,
		},
		Run: func(args interface{}) error {
			jobArgs := args.(JobArgs)
			deploymentName := jobArgs.Config.Name

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
			jobArgs := args.(JobArgs)
			deploymentName := jobArgs.Config.Name

			// delete secrets collection
			logger.Debugf("removing secrets collection for deployment %s", deploymentName)
			if err := DeleteSecretsCollection(deploymentName); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing secrets collection for deployment %s", deploymentName))
			}

			// delete jobs collection
			logger.Debugf("removing jobs collection for deployment %s", deploymentName)
			if err := job.DeleteCollection(deploymentName); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing jobs collection for deployment %s", deploymentName))
			}

			// delete deployment configuration
			logger.Debugf("removing config for deployment %s", deploymentName)
			if err := DeleteDeploymentConfig(jobArgs.Config.Name); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("error removing config for deployment %s", deploymentName))
			}

			return nil
		},
	})

	return nil
}
