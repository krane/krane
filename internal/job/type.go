package job

type Type string

const (
	ContainerCreate = "CREATE_CONTAINER"
	ContainerDelete = "DELETE_CONTAINER"
)

var AllowedJobs = []string{ContainerCreate, ContainerDelete}

func isAllowedJobType(jobType Type) bool {
	for _, t := range AllowedJobs {
		if t == string(jobType) {
			return true
		}
	}
	return false
}
