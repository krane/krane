package service

import "github.com/biensupernice/krane/internal/deployment/config"

func GetDeploymentByName(name string) (config.Config, error) { return config.Get(name) }

func GetAllDeployments() ([]config.Config, error) { return config.GetAll() }
