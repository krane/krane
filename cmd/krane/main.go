package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api"
	"github.com/biensupernice/krane/internal/docker"
	job "github.com/biensupernice/krane/internal/jobs"
	"github.com/biensupernice/krane/internal/logging"
	"github.com/biensupernice/krane/internal/scheduler"
	"github.com/biensupernice/krane/internal/store"
)

func init() {
	log.Println("Starting Krane...")

	requireEnv("KRANE_PRIVATE_KEY")
	envOrDefault("LOG_LEVEL", logging.INFO)
	envOrDefault("LISTEN_ADDRESS", "127.0.0.1:8500")
	envOrDefault("STORE_PATH", "/tmp/krane.db")
	envOrDefault("WORKERPOOL_SIZE", "1")
	envOrDefault("JOBQUEUE_SIZE", "1")
	envOrDefault("SCHEDULER_INTERVAL_MS", "5000")

	logging.ConfigureLogrus()
	docker.Init()
	store.New(os.Getenv("STORE_PATH"))
}

func main() {
	defer store.Instance().Shutdown()

	// Docker Client
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		logrus.Fatalf("Unable to connect to Docker client %s", err.Error())
	}

	// Job Queue
	jobQueueSize, _ := strconv.ParseUint(os.Getenv("JOBQUEUE_SIZE"), 10, 8)
	jobQueue := job.NewJobQueue(uint(jobQueueSize))

	// Job Enqueuer
	jobEnqueuer := job.NewEnqueuer(store.Instance(), jobQueue)

	// Job Worker Pool
	wpSize, _ := strconv.ParseUint(os.Getenv("WORKERPOOL_SIZE"), 10, 8)
	jobWorkers := job.NewWorkerPool(uint(wpSize), jobQueue, store.Instance())
	jobWorkers.Start()

	// Scheduler
	interval := os.Getenv("SCHEDULER_INTERVAL_MS")
	jobScheduler := scheduler.New(store.Instance(), dockerClient, jobEnqueuer, interval)
	jobScheduler.Run()

	go api.Run()

	wait()

	logrus.Info("Shutdown signal received")
	jobWorkers.Stop()
	logrus.Info("All workers done shutting down")
}

// wait : for a signal to quit
func wait() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan // Blocks here until interrupted
}

func requireEnv(key string) {
	_, hasEnv := os.LookupEnv(key)
	if !hasEnv {
		log.Fatalf("Missing required env %s", key)
	}
}

func envOrDefault(key string, fallback string) string {
	val, hasEnv := os.LookupEnv(key)
	if !hasEnv {
		log.Printf("%s not set, defaulting to %s", key, fallback)
		os.Setenv(key, fallback)
		return fallback
	}

	logrus.Infof("%s already set with value %s", key, val)
	return os.Getenv(key)
}
