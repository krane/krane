package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/biensupernice/krane/internal/api"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/scheduler"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func init() {
	utils.RequireEnv(constants.EnvKranePrivateKey)
	utils.EnvOrDefault(constants.EnvLogLevel, "info")
	utils.EnvOrDefault(constants.EnvListenAddress, "127.0.0.1:8500")
	utils.EnvOrDefault(constants.EnvDatabasePath, "/tmp/krane.db")
	utils.EnvOrDefault(constants.EnvWorkerPoolSize, "1")
	utils.EnvOrDefault(constants.EnvJobQueueSize, "1")
	utils.EnvOrDefault(constants.EnvJobMaxRetryPolicy, "5")
	utils.EnvOrDefault(constants.EnvDeploymentRetryPolicy, "1")
	utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, "30000")
	utils.EnvOrDefault(constants.EnvWatchMode, "false")
	utils.EnvOrDefault(constants.EnvDockerBasicAuthUsername, "")
	utils.EnvOrDefault(constants.EnvDockerBasicAuthPassword, "")
	utils.EnvOrDefault(constants.EnvProxyEnabled, "false")
	utils.EnvOrDefault(constants.EnvProxyDashboardSecure, "false")
	utils.EnvOrDefault(constants.EnvProxyDashboardAlias, "")

	logger.Configure()
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

	// shared job queue for deployment jobs
	qsize := utils.UIntEnv(constants.EnvJobQueueSize)
	queue := job.NewBufferedQueue(qsize)

	// if watch mode is enabled, the scheduler will run
	// in a separate thread polling and queuing jobs to
	// maintain the containers state in parity with desired deployment state
	if utils.BoolEnv(constants.EnvWatchMode) {
		logger.Warn("This feature is experimental, dont use unless ")
		enqueuer := job.NewEnqueuer(queue)
		interval := utils.EnvOrDefault(constants.EnvSchedulerIntervalMs, "30000")

		jobScheduler := scheduler.New(db, docker.GetClient(), enqueuer, interval)
		go jobScheduler.Run()
	}

	// workers for executing deployment jobs; when no workers are instantiated,
	// queued jobs will block until a worker is added to the worker pool.
	wpSize := utils.UIntEnv(constants.EnvWorkerPoolSize)
	workers := job.NewWorkerPool(wpSize, queue, store.Client())
	workers.Start()

	// if enabled, ensure network proxy is running
	EnsureNetworkProxy()

	// wait will block
	// until an exit signal is received by the program.
	// the exit signal can be ctrl+c, an IDE stop, or any system termination signal.
	// once the signal is received it will shutdown all workers.
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
