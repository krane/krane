package service

//
// func getCurrentContainers(args job.Args) error {
// 	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
//
// 	containers, err := container.GetKraneContainersByDeployment(cfg.Name)
// 	if err != nil {
// 		return err
// 	}
//
// 	args[CurrentContainersJobArgName] = &containers
// 	return nil
// }
