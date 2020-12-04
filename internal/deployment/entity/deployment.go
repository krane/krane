package entity

import "github.com/biensupernice/krane/internal/deployment/kconfig"

type Deployment struct {
	Config kconfig.Kconfig `json:"config"`
}
