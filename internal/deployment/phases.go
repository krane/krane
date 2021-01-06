package deployment

type Phase string

const (
	SetupPhase           = "DEPLOYMENT_SETUP"
	TeardownPhase        = "DEPLOYMENT_TEARDOWN"
	HealthCheckPhase     = "DEPLOYMENT_HEALTCHECK"
	PullImagePhase       = "PULL_IMAGE"
	CreateContainerPhase = "CREATE_CONTAINER"
	StartContainerPhase  = "START_CONTAINER"
	DonePhase            = "DONE"
)
