package deployment

type JobArgs struct {
	Config Config

	// OldContainers represent containers created for a previous deployment execution
	OldContainers []KraneContainer

	// NewContainers represent containers which have been created as
	// part of the most recent or current deployment execution
	NewContainers []KraneContainer
}
