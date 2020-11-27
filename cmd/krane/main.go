package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logging"
	"github.com/biensupernice/krane/internal/scheduler"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func init() {
	fmt.Println("Starting Krane...")

	utils.RequireEnv(constants.EnvKranePrivateKey)
	utils.EnvOrDefault(constants.EnvListenAddress, "127.0.0.1:8500")
	utils.EnvOrDefault(constants.EnvLogLevel, logging.INFO)
	utils.EnvOrDefault(constants.EnvDatabasePath, "/tmp/krane.db")
	utils.EnvOrDefault(constants.EnvWorkerPoolSize, "1")
	utils.EnvOrDefault(constants.EnvJobQueueSize, "1")
	utils.EnvOrDefault(constants.EnvJobMaxRetryPolicy, "5")
	utils.EnvOrDefault(constants.EnvDeploymentRetryPolicy, "1")
	utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, "30000")
	utils.EnvOrDefault(constants.EnvWatchMode, "false")
	utils.EnvOrDefault(constants.EnvDockerBasicAuthUsername, "")
	utils.EnvOrDefault(constants.EnvDockerBasicAuthPassword, "")

	logging.ConfigureLogrus()
	docker.ClientFromEnv()
	store.NewInstance(os.Getenv(constants.EnvDatabasePath))
}

func main() {
	// api
	go api.Run()

	// embedded key store for storing deployment configurations
	db := store.Instance()
	defer db.Shutdown()

	// shared job queue used for deployment jobs
	qsize := utils.GetUIntEnv(constants.EnvJobQueueSize)
	queue := job.NewJobQueue(qsize)

	// if watch mode is enabled, the scheduler will run
	// in a separate thread polling and queuing jobs to
	// maintain the containers state in parity with desired deployment state
	if utils.GetBoolEnv(constants.EnvWatchMode) {
		logrus.Warn("This feature is experimental, dont use unless ")
		enqueuer := job.NewEnqueuer(queue)
		interval := utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, "30000")

		jobScheduler := scheduler.New(db, docker.GetClient(), enqueuer, interval)
		go jobScheduler.Run()
	}

	// workers for executing jobs. If no workers are initiated
	// the queued jobs will stay blocked until a worker frees up or is initiated
	wpSize := utils.GetUIntEnv(constants.EnvWorkerPoolSize)
	workers := job.NewWorkerPool(wpSize, queue, store.Instance())
	workers.Start()

	// This wait statement will block until an exit signal is received by the program.
	// The exit signal can be ctrl+c, any IDE stop, or any system termination signal.
	// Once the signal is received it will shutdown all workers as gracefully as it can.
	wait()

	// workers run in background threads.
	// When an exit signal is received, the workers are stopped and cleaned up
	workers.Stop()

	logrus.Info("Shutdown complete")
}

func wait() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan // blocking
}
