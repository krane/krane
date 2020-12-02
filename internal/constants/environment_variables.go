package constants

const (
	// Environment variables used by Krane
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
	EnvProxyDashboardSecure    = "PROXY_DASHBOARD_SECURE"
	EnvProxyDashboardAlias     = "PROXY_DASHBOARD_ALIAS"
)
