package job

type Status struct {
	ExecutionCount uint    `json:"execution_count"`
	FailureCount   uint    `json:"failure_count"`
	Failures       []Error `json:"failures"`
}
