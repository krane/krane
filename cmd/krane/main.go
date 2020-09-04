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
	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/storage/boltdb"
)

func init() {
	log.Println("Starting Krane...")

	requireEnv("KRANE_PRIVATE_KEY")
	envOrDefault("LOG_LEVEL", logging.INFO)
	envOrDefault("LISTEN_ADDRESS", "127.0.0.1:8500")
	envOrDefault("STORE_PATH", "/tmp/krane.db")
	envOrDefault("WORKERPOOL_SIZE", "1")
	envOrDefault("JOBQUEUE_SIZE", "1")

	logging.ConfigureLogrus()
	docker.Init()
	boltdb.Init(os.Getenv("STORE_PATH"))
}

func main() {
	// Docker Client
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		logrus.Fatalf("Unable to connect to Docker client %s", err.Error())
	}

	// Storage
	store := storage.GetInstance()
	defer store.Shutdown()

	// Job Queue
	jobQueueSize, _ := strconv.ParseUint(os.Getenv("JOBQUEUE_SIZE"), 10, 8)
	jobQueue := job.NewJobQueue(uint(jobQueueSize))

	// Job Enqueuer
	jobEnqueuer := job.NewEnqueuer(&store, jobQueue)

	// Job Worker Pool
	wpSize, _ := strconv.ParseUint(os.Getenv("WORKERPOOL_SIZE"), 10, 8)
	jobWorkers := job.NewWorkerPool(uint(wpSize), jobQueue, &store)
	jobWorkers.Start()

	// Scheduler
	jobScheduler := scheduler.New(store, dockerClient, jobEnqueuer)
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
