package defaultConfig

const (
	// TIMEOUT is the default timeout for service checks in seconds
	TIMEOUT = 5

	// RETRY is the default retry count for service checks
	RETRY = 2

	// MAX_LOG_DAYS is the default maximum number of days to keep logs
	MAX_LOG_DAYS = 30
)

// GetDefaultTimeout returns the default timeout for service checks
func GetDefaultTimeout() int {
	return TIMEOUT
}

// GetDefaultRetry returns the default retry count for service checks
func GetDefaultRetry() int {
	return RETRY
}

// GetDefaultMaxLogDays returns the default maximum number of days to keep logs
func GetDefaultMaxLogDays() int {
	return MAX_LOG_DAYS
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
	// CONFIG_PATH is the default path to the configuration file
	CONFIG_PATH = "config.yaml"

	// RESULT_PATH is the default path to the data file where logs are stored
	RESULT_PATH = "data/ponghub_result.json"

	// LOG_PATH is the default path to the data file where logs are stored
	LOG_PATH = "data/ponghub_log.json"

	// REPORT_PATH is the default path to the HTML report file
	REPORT_PATH = "data/index.html"
)

// GetConfigPath returns the default path to the configuration file
func GetConfigPath() string {
	return CONFIG_PATH
}

// GetResultPath returns the default path to the data file where logs are stored
func GetResultPath() string {
	return RESULT_PATH
}

// GetLogPath returns the default path to the data file where logs are stored
func GetLogPath() string {
	return LOG_PATH
}

// GetReportPath returns the default path to the HTML report file
func GetReportPath() string {
	return REPORT_PATH
}
