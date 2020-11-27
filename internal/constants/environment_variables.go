package constants

// List of environment variables used by Krane
const (
	EnvKranePrivateKey         = "KRANE_PRIVATE_KEY"
	EnvLogLevel                = "LOG_LEVEL"
	EnvListenAddress           = "LISTEN_ADDRESS"
	EnvWatchMode               = "WATCH_MODE"
	EnvDatabasePath            = "DB_PATH"
	EnvWorkerPoolSize          = "WORKERPOOL_SIZE"
	EnvJobQueueSize            = "JOB_QUEUE_SIZE"
	EnvJobMaxRetryPolicy       = "JOB_MAX_RETRY_POLICY"
	EnvDeploymentRetryPolicy   = "DEPLOYMENT_RETRY_POLICY"
	EnvSchedulerIntervalMs     = "SCHEDULER_INTERVAL_MS"
	EnvDockerBasicAuthUsername = "DOCKER_BASIC_AUTH_USERNAME"
	EnvDockerBasicAuthPassword = "DOCKER_BASIC_AUTH_PASSWORD"
)
