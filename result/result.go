package result

// Result represents successful operation or a failed one
// Ok: If operation was succesful or failure
// Value: If succesful, contains the value of the operation
// Errors: If failure, represents the errors for the operation
type Result struct {
	Ok     bool        `json:"ok"`
	Value  interface{} `json:"value"`
	Errors []string    `json:"errors"`
}
