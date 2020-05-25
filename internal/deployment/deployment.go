package deployment

import (
	"github.com/docker/docker/api/types"
)

// Deployment : representation
type Deployment struct {
	Template   Template          `json:"template"`
	Containers []types.Container `json:"containers"`
}
