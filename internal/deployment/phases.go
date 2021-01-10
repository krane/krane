package deployment

// Phases represents a particular step a deployment could be going through during its deployment cycle.
// They are attached to jobs allowing clients to react or filter for particular phases of the deployments cycle.
type Phase string

const (
	SetupPhase           Phase = "DEPLOYMENT_SETUP"
	HealthCheckPhase     Phase = "DEPLOYMENT_HEALTCHECK"
	TeardownPhase        Phase = "DEPLOYMENT_TEARDOWN"
	DonePhase            Phase = "DEPLOYMENT_DONE"
	PullImagePhase       Phase = "PULL_IMAGE"
	CreateContainerPhase Phase = "CREATE_CONTAINER"
	StartContainerPhase  Phase = "START_CONTAINER"
)
