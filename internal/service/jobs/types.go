package jobs

var (
	// A job name should be self descriptive. The name itself should follow this naming convention
	// whenever possible `{action}_{subject}_{context}`. The action should be CRUD like, the subject should
	// be the entity you want to apply the action to, and the context provides extra details on which the action should be taken on.
	StartDeploymentJobName       = "START_DEPLOYMENT"
	DeleteDeploymentJobName      = "DELETE_DEPLOYMENT"
	StopDeploymentJobName        = "STOP_DEPLOYMENT"
	UpdateDeploymentAliasJobName = "UPDATE_DEPLOYMENT_ALIAS"
	DeleteDeploymentAliasJobName = "DELETE_DEPLOYMENT_ALIAS"
)
