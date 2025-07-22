package testResult

type TestResult string

const (
	// ALL represents all ports are online
	ALL TestResult = "all"

	// PART represents some ports are online
	PART TestResult = "part"

	// NONE represents no ports are online
	NONE TestResult = "none"

	// UNKNOWN represents an unknown test result
	UNKNOWN TestResult = "unknown"
)

// String returns the string representation of the TestResult
func (tr TestResult) String() string {
	switch tr {
	case ALL:
		return "all"
	case PART:
		return "part"
	case NONE:
		return "none"
	default:
		return "unknown"
	}
}

// IsValid checks if the TestResult is valid
func (tr TestResult) IsValid() bool {
	return tr == ALL || tr == PART || tr == NONE
}

// ParseTestResult parses a string into a TestResult
func ParseTestResult(s string) TestResult {
	switch s {
	case "all":
		return ALL
	case "part":
		return PART
	case "none":
		return NONE
	default:
		return UNKNOWN
	}
}

// ParseTestResults parses a slice of strings into a slice of TestResult
func ParseTestResults(results []string) []TestResult {
	var parsedResults []TestResult
	for _, result := range results {
		parsedResults = append(parsedResults, ParseTestResult(result))
	}
	return parsedResults
}

// CheckTestResult checks if the given TestResult is valid
func (tr TestResult) CheckTestResult() bool {
	return tr.IsValid()
}
