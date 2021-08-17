package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/krane/krane/internal/api"
	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/job"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/scheduler"
	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

func init() {
	utils.RequireEnv(constants.EnvKranePrivateKey)
	utils.EnvOrDefault(constants.EnvLogLevel, "info")
	utils.EnvOrDefault(constants.EnvListenAddress, "0.0.0.0:8500")
	utils.EnvOrDefault(constants.EnvDatabasePath, "/tmp/krane.db")
	utils.EnvOrDefault(constants.EnvWorkerPoolSize, "1")
	utils.EnvOrDefault(constants.EnvJobQueueSize, "1")
	utils.EnvOrDefault(constants.EnvJobMaxRetryPolicy, "5")
	utils.EnvOrDefault(constants.EnvDeploymentRetryPolicy, "1")
	utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, "30000")
	utils.EnvOrDefault(constants.EnvWatchMode, "false")
	utils.EnvOrDefault(constants.EnvProxyEnabled, "true")
	utils.EnvOrDefault(constants.EnvProxyDashboardSecure, "false")
	utils.EnvOrDefault(constants.EnvProxyDashboardAlias, "")
	utils.EnvOrDefault(constants.EnvLetsEncryptEmail, "")

	logger.Configure()
	logger.Info("Setting up Krane")

	docker.Connect()
	store.Connect(os.Getenv(constants.EnvDatabasePath))
}

func main() {
	logger.Info("Starting Krane")

	// rest api
	go api.Run()

	// embedded database
	db := store.Client()
	defer db.Disconnect()

	// deployment job queue
	qsize := utils.UIntEnv(constants.EnvJobQueueSize)
	queue := job.NewBufferedQueue(qsize)

	// if watch mode is enabled, the scheduler will run in a separate routine polling
	// and queuing jobs to maintain the deployment state in parity with the desired state
	if utils.BoolEnv(constants.EnvWatchMode) {
		logger.Warn("Watch mode is an experimental feature. Krane will maintain container state in parity with your deployment configuration.")
		enqueuer := job.NewEnqueuer(queue)
		interval := utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, utils.TwoMinMs)

		jobScheduler := scheduler.New(db, docker.GetClient(), enqueuer, interval)
		go jobScheduler.Run()
	}

	// workers for executing deployment jobs; when no workers are instantiated,
	// queued jobs will block until a worker is added to the worker pool.
	wpSize := utils.UIntEnv(constants.EnvWorkerPoolSize)
	workers := job.NewWorkerPool(wpSize, queue, store.Client())
	workers.Start()

	// ensure internal services are running
	// the network proxy is currently the only dependant 
	// service for krane (others include the ui)
	EnsureNetworkProxy()

	// block until an exit signal is received
	wait()

	// when an exit signal is received, workers are stopped and cleaned up.
	workers.Stop()

	logger.Info("Shutdown complete")
}

func wait() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan // blocking
}
