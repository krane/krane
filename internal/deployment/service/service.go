package service

//
// type Deployment interface {
// 	StartContainers(deploymentName string) (job.Job, error)
// 	StopContainers(deploymentName string) (job.Job, error)
// 	DeleteContainers(deploymentName string) (job.Job, error)
// 	RestartContainers(deploymentName string) (job.Job, error)
// }
//
// type KraneDeployment struct {
// 	Config     config.DeploymentConfig
// 	containers []container.KraneContainer
// }
//
// func (d KraneDeployment) StartContainers() (job.Job, error) {
// 	deploymentJob, err := d.job(StartContainers)
// 	if err != nil {
// 		return job.Job{}, err
// 	}
// 	d.enqueue(deploymentJob)
// 	return deploymentJob, nil
// }
//
// func (d KraneDeployment) StopContainers() (job.Job, error) {
// 	return job.Job{}, nil
// }
//
// func (d KraneDeployment) DeleteContainers() (job.Job, error) {
// 	return job.Job{}, nil
// }
//
// func (d KraneDeployment) RestartContainers() (job.Job, error) {
// 	return job.Job{}, nil
// }
//
// func (d KraneDeployment) job(action DeploymentAction) (job.Job, error) {
// 	switch action {
// 	case CreateContainers:
// 		return createContainersJob(d.Config), nil
// 	case DeleteContainers:
// 		return deleteContainersJob(d.Config), nil
// 	case StopContainers:
// 		return stopContainersJob(d.Config), nil
// 	default:
// 		return job.Job{}, fmt.Errorf("unknown deployment action %s", action)
// 	}
// }
//
// func (d KraneDeployment) enqueue(deploymentJob job.Job) {
// 	enqueuer := job.NewEnqueuer(job.Queue())
// 	queuedJob, err := enqueuer.Enqueue(deploymentJob)
// 	if err != nil {
// 		logger.Errorf("Error enqueuing deployment job %v", err)
// 		return
// 	}
// 	logger.Debugf("Queued job for %s", queuedJob.Deployment)
// }
//
// func applyDeployment(cfg config.DeploymentConfig, action DeploymentAction) error {
// 	deploymentJob, err := createDeploymentJob(cfg, action)
// 	if err != nil {
// 		return err
// 	}
// 	go enqueueDeploymentJob(deploymentJob)
// 	return nil
// }
//
// // enqueueDeploymentJob : enqueues a deployment job for processing
// func enqueueDeploymentJob(deploymentJob job.Job) {
// 	enqueuer := job.NewEnqueuer(job.Queue())
// 	queuedJob, err := enqueuer.Enqueue(deploymentJob)
// 	if err != nil {
// 		logger.Errorf("Error enqueuing deployment job %v", err)
// 		return
// 	}
// 	logger.Debugf("Queued job for %s", queuedJob.Deployment)
// }
//
// // StartDeploymentContainers : starts a deployments container resources
// func StartDeploymentContainers(cfg config.DeploymentConfig) error {
// 	return applyDeployment(cfg, CreateContainers)
// }
//
// // DeleteDeployment : deletes a deployment and its container resources
// func DeleteDeployment(cfg config.DeploymentConfig) error {
// 	return applyDeployment(cfg, DeleteContainers)
// }
//
// // StopDeploymentContainers : stops a deployments container resources
// func StopDeploymentContainers(cfg config.DeploymentConfig) error {
// 	return applyDeployment(cfg, StopContainers)
// }
