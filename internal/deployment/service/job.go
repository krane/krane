package service

//
// type DeploymentAction string
//
// const (
// 	CreateContainers DeploymentAction = "CREATE_CONTAINERS"
// 	DeleteContainers DeploymentAction = "DELETE_CONTAINERS"
// 	StopContainers   DeploymentAction = "STOP_CONTAINERS"
// 	StartContainers  DeploymentAction = "START_CONTAINERS"
// )
//
// const (
// 	DeploymentConfigJobArgName  = "deployment_config"
// 	CurrentContainersJobArgName = "current_containers"
// 	NewContainersJobArgName     = "new_containers"
// )
//
// func createDeploymentJob(config config.DeploymentConfig, action DeploymentAction) (job.Job, error) {
// 	switch action {
// 	case CreateContainers:
// 		return createContainersJob(config), nil
// 	case DeleteContainers:
// 		return deleteContainersJob(config), nil
// 	case StopContainers:
// 		return stopContainersJob(config), nil
// 	default:
// 		return job.Job{}, fmt.Errorf("unknown deployment DeploymentAction %s", action)
// 	}
// }
//
// func createContainersJob(config config.DeploymentConfig) job.Job {
// 	currContainers := make([]container.KraneContainer, 0)
// 	newContainers := make([]container.KraneContainer, 0)
// 	retryPolicy := utils.UIntEnv("DEPLOYMENT_RETRY_POLICY")
//
// 	jobsArgs := job.Args{
// 		DeploymentConfigJobArgName:  config,
// 		CurrentContainersJobArgName: &currContainers,
// 		NewContainersJobArgName:     &newContainers,
// 	}
//
// 	return job.Job{
// 		ID:          uuid.Generate().String(),
// 		Deployment:  config.Name,
// 		Type:        string(CreateContainers),
// 		Args:        jobsArgs,
// 		RetryPolicy: retryPolicy,
// 		// Run:         createContainerResources,
// 	}
// }
//
// func deleteContainersJob(config config.DeploymentConfig) job.Job {
// 	containers := make([]container.KraneContainer, 0)
// 	retryPolicy := utils.UIntEnv("DEPLOYMENT_RETRY_POLICY")
//
// 	jobsArgs := job.Args{
// 		DeploymentConfigJobArgName:  config,
// 		CurrentContainersJobArgName: &containers,
// 	}
//
// 	return job.Job{
// 		ID:          uuid.Generate().String(),
// 		Deployment:  config.Name,
// 		Type:        string(DeleteContainers),
// 		Args:        jobsArgs,
// 		RetryPolicy: retryPolicy,
// 		// Run:         deleteContainerResources,
// 	}
// }

// func stopContainersJob(config config.DeploymentConfig) job.Job {
// 	containers := make([]container.KraneContainer, 0)
// 	retryPolicy := utils.UIntEnv("DEPLOYMENT_RETRY_POLICY")
//
// 	jobsArgs := job.Args{
// 		DeploymentConfigJobArgName:  config,
// 		CurrentContainersJobArgName: &containers,
// 	}
//
// 	return job.Job{
// 		ID:          uuid.Generate().String(),
// 		Deployment:  config.Name,
// 		Type:        string(StopContainers),
// 		Args:        jobsArgs,
// 		RetryPolicy: retryPolicy,
// 		// Run:         stopContainerResources,
// 	}
// }
