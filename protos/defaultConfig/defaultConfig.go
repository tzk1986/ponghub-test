package defaultConfig

const (
	// timeout is the default timeout for service checks in seconds
	timeout = 5

	// retry is the default retry count for service checks
	retry = 2

	// maxLogDays is the default maximum number of days to keep logs
	maxLogDays = 30
)

// GetDefaultTimeout returns the default timeout for service checks
func GetDefaultTimeout() int {
	return timeout
}

// GetDefaultRetry returns the default retry count for service checks
func GetDefaultRetry() int {
	return retry
}

// GetDefaultMaxLogDays returns the default maximum number of days to keep logs
func GetDefaultMaxLogDays() int {
	return maxLogDays
}

// SetDefaultTimeout sets the default timeout for a given configuration pointer
func SetDefaultTimeout(cfg *int) {
	if cfg == nil || *cfg <= 0 {
		*cfg = GetDefaultTimeout()
	}
}

// SetDefaultRetry sets the default retry count for a given configuration pointer
func SetDefaultRetry(cfg *int) {
	if cfg == nil || *cfg <= 0 {
		*cfg = GetDefaultRetry()
	}
}

// SetDefaultMaxLogDays sets the default maximum number of days to keep logs for a given configuration pointer
func SetDefaultMaxLogDays(cfg *int) {
	if cfg == nil || *cfg <= 0 {
		*cfg = GetDefaultMaxLogDays()
	}
}

const (
	// configPath is the default path to the configuration file
	configPath = "config.yaml"

	// logPath is the default path to the data file where logs are stored
	logPath = "data/ponghub_log.json"

	// reportPath is the default path to the HTML report file
	reportPath = "data/index.html"
)

// GetConfigPath returns the default path to the configuration file
func GetConfigPath() string {
	return configPath
}

// GetLogPath returns the default path to the data file where logs are stored
func GetLogPath() string {
	return logPath
}

// GetReportPath returns the default path to the HTML report file
func GetReportPath() string {
	return reportPath
}
