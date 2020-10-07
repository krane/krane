package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logging"
	"github.com/biensupernice/krane/internal/scheduler"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func init() {
	fmt.Println("Starting Krane...")

	utils.RequireEnv("KRANE_PRIVATE_KEY")
	utils.EnvOrDefault("LISTEN_ADDRESS", "127.0.0.1:8500")
	utils.EnvOrDefault("LOG_LEVEL", logging.INFO)
	utils.EnvOrDefault("DB_PATH", "/tmp/krane.db")
	utils.EnvOrDefault("WORKERPOOL_SIZE", "1")
	utils.EnvOrDefault("JOB_QUEUE_SIZE", "1")
	utils.EnvOrDefault("JOB_MAX_RETRY_POLICY", "5")
	utils.EnvOrDefault("DEPLOYMENT_RETRY_POLICY", "1")
	utils.EnvOrDefault("SCHEDULER_INTERVAL_MS", "30000")
	utils.EnvOrDefault("WATCH_MODE", "false")
	logging.ConfigureLogrus()

	docker.ClientFromEnv()
	store.New(os.Getenv("DB_PATH"))
}

func main() {
	// api
	go api.Run()

	// embedded key store for storing deployment configurations
	db := store.Instance()
	defer db.Shutdown()

	// shared job queue used for deployment jobs
	qsize := utils.GetUIntEnv("JOB_QUEUE_SIZE")
	queue := job.NewJobQueue(qsize)

	// if watch mode is enabled, the scheduler will run
	// in a separate thread polling and queuing jobs to
	// maintain the containers state in parity with desired deployment state
	if utils.GetBoolEnv("WATCH_MODE") {
		dockerClient := docker.GetClient()
		enqueuer := job.NewEnqueuer(db, queue)
		interval := os.Getenv("SCHEDULER_INTERVAL_MS")

		jobScheduler := scheduler.New(db, dockerClient, enqueuer, interval)
		go jobScheduler.Run()
	}

	// workers for executing jobs. If no workers are initiated
	// the queued jobs will stay blocked until a worker frees up or is initiated
	wpSize := utils.GetUIntEnv("WORKERPOOL_SIZE")
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
