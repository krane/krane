package jobs

import (
	"github.com/biensupernice/krane/internal/service/deployment"
	"github.com/biensupernice/krane/pkg/bbq"
)

func StartDeployment(deployment deployment.Deployment, imageTag string) {
	props := map[string]string{"tag": imageTag}
	bbq.Queue(bbq.Job{
		Body:    deployment,
		Props:   props,
		JobType: StartDeploymentJobName,
		Process: deployment.Start,
		Done:    onJobDone,
		OnError: onJobError,
	})
}

func StopDeployment(deployment deployment.Deployment) {
	bbq.Queue(bbq.Job{
		Body:    deployment,
		JobType: StopDeploymentJobName,
		Process: deployment.Stop,
		Done:    onJobDone,
		OnError: onJobError,
	})
}

func DeleteDeployment(deployment deployment.Deployment) {
	bbq.Queue(bbq.Job{
		Body:    deployment,
		JobType: DeleteDeploymentJobName,
		Process: deployment.Delete,
		Done:    onJobDone,
		OnError: onJobError,
	})
}

func DeleteDeploymentAlias(deployment deployment.Deployment) {
	bbq.Queue(bbq.Job{
		Body:    deployment,
		JobType: DeleteDeploymentJobName,
		Process: deployment.DeleteAlias,
		Done:    onJobDone,
		OnError: onJobError,
	})
}

func UpdateDeploymentAlias(deployment deployment.Deployment, alias string) {
	props := map[string]string{"alias": alias}
	bbq.Queue(bbq.Job{
		Body:    deployment,
		Props:   props,
		JobType: UpdateDeploymentAliasJobName,
		Process: deployment.UpdateAlias,
		Done:    onJobDone,
		OnError: onJobError,
	})
}
