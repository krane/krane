package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	utils.EnvOrDefault("LOG_LEVEL", logging.INFO)
	utils.EnvOrDefault("LISTEN_ADDRESS", "127.0.0.1:8500")
	utils.EnvOrDefault("STORE_PATH", "/tmp/krane.db")
	utils.EnvOrDefault("WORKERPOOL_SIZE", "1")
	utils.EnvOrDefault("JOBQUEUE_SIZE", "1")
	utils.EnvOrDefault("SCHEDULER_INTERVAL_MS", "10000")
	logging.ConfigureLogrus()

	docker.NewClient()
	store.New(os.Getenv("STORE_PATH"))
}

func main() {
	// store
	db := store.Instance()
	defer db.Shutdown()

	// queue
	qsize, _ := strconv.ParseUint(os.Getenv("JOBQUEUE_SIZE"), 10, 8)
	queue := job.NewJobQueue(uint(qsize))

	// enqueuer
	enqueuer := job.NewEnqueuer(store.Instance(), queue)

	// scheduler
	interval := os.Getenv("SCHEDULER_INTERVAL_MS")
	scheduler := scheduler.New(db, docker.GetClient(), enqueuer, interval)
	go scheduler.Run()

	// api
	go api.Run()

	// workers
	wpSize := utils.GetUIntEnv("WORKERPOOL_SIZE")
	workers := job.NewWorkerPool(wpSize, queue, store.Instance())
	workers.Start()

	// This wait statement will block until an exit signal is received by the program.
	// The exit signal can be ctrl+c, any IDE stop, or any system termination signal.
	// Once the signal is received it will shutdown all workers as gracefully as it can.
	wait()

	workers.Stop()
	logrus.Info("Shutdown complete")
}

// wait : for a signal to quit
func wait() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
}
