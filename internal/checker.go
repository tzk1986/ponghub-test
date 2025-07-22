package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/wcy-dt/ponghub/protos/portType"
	"github.com/wcy-dt/ponghub/protos/testResult"
)

// CheckResult defines the structure for the result of checking a service
type CheckResult struct {
	Name          string                `json:"name"`
	Online        testResult.TestResult `json:"online"`
	Health        []PortResult          `json:"health,omitempty"`
	API           []PortResult          `json:"api,omitempty"`
	StartTime     string                `json:"start_time"`
	EndTime       string                `json:"end_time"`
	TotalAttempts int                   `json:"total_attempts"`
	SuccessCount  int                   `json:"success_count"`
}

// PortResult defines the structure for the result of checking a port
type PortResult struct {
	URL           string                `json:"url"`
	Method        string                `json:"method"`
	Body          string                `json:"body,omitempty"`
	Online        testResult.TestResult `json:"online"`
	StatusCode    int                   `json:"status_code,omitempty"`
	StartTime     string                `json:"start_time"`
	EndTime       string                `json:"end_time"`
	TotalAttempts int                   `json:"total_attempts"`
	SuccessCount  int                   `json:"success_count"`
	Failures      []string              `json:"failures,omitempty"`
	ResponseBody  string                `json:"response_body,omitempty"`
}

// getHttpMethod converts a string method to an HTTP method constant
func getHttpMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return http.MethodGet
	case "POST":
		return http.MethodPost
	case "PUT":
		return http.MethodPut
	case "DELETE":
		log.Fatalln(errors.New("method not supported"))
	case "HEAD":
		log.Fatalln(errors.New("method not supported"))
	case "PATCH":
		log.Fatalln(errors.New("method not supported"))
	case "OPTIONS":
		log.Fatalln(errors.New("method not supported"))
	case "TRACE":
		log.Fatalln(errors.New("method not supported"))
	case "CONNECT":
		log.Fatalln(errors.New("method not supported"))
	default:
		return http.MethodGet // Default to GET if method is unknown
	}
	return http.MethodGet
}

// getTestResult determines the test result based on the success count and actual attempts
func getTestResult(successCount, actualAttempts int) testResult.TestResult {
	switch successCount {
	case actualAttempts:
		return testResult.ALL
	case 0:
		return testResult.NONE
	default:
		return testResult.PART
	}
}

// isSuccessfulResponse checks if the response from the server is successful based on the configuration
func isSuccessfulResponse(cfg *PortConfig, resp *http.Response, body []byte) bool {
	// responseRegex is set, and the response body does not match the regex
	if cfg.ResponseRegex != "" {
		matched, err := regexp.Match(cfg.ResponseRegex, body)
		if err != nil {
			log.Fatalln("Error parsing regexp:", err)
		}
		if !matched {
			return false
		}
	}

	// statusCode and responseRegex are not set, and the response is OK
	if cfg.StatusCode == 0 && cfg.ResponseRegex == "" && resp.StatusCode == http.StatusOK {
		return true
	}

	// statusCode is not set, and the responseRegex matches
	if cfg.StatusCode == 0 && cfg.ResponseRegex != "" {
		return true
	}

	// statusCode is set, and the response matches the expected status code
	if cfg.StatusCode != 0 && resp.StatusCode == cfg.StatusCode {
		return true
	}

	return false
}

// CheckPort checks a single port based on the provided configuration
func CheckPort(cfg *PortConfig, timeout int, retryTimes int, svcName string, portType portType.PortType) PortResult {
	failures := []string{}
	successCount := 0
	actualAttempts := 0

	var statusCode int
	var responseBody string

	httpMethod := getHttpMethod(cfg.Method)

	// start timer
	start := time.Now()

	for attemptTimes := range retryTimes {
		actualAttempts++
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
		log.Printf("[%s] %s %s (attempt %d/%d)\n",
			svcName, httpMethod, cfg.URL, attemptTimes+1, retryTimes)

		// build the request
		req, err := http.NewRequest(httpMethod, cfg.URL, nil)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			log.Printf("FAILED - Error: %s", err.Error())
			continue
		}
		if cfg.Body != "" {
			req.Body = io.NopCloser(strings.NewReader(cfg.Body))
		}

		// get the response
		resp, err := client.Do(req)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			log.Printf("FAILED - Error: %s", err.Error())
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: %d, Error: %s", resp.StatusCode, err.Error()))
			log.Printf("FAILED - StatusCode: %d, Error: %s", resp.StatusCode, err.Error())
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body for %s: %v", cfg.URL, err)
			}
			continue
		}
		responseBody = string(body)
		statusCode = resp.StatusCode

		// check the response
		isOnline := isSuccessfulResponse(cfg, resp, body)
		if isOnline {
			successCount++
			responseBody = ""
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body for %s: %v", cfg.URL, err)
			}
			break
		}
		failures = append(failures, fmt.Sprintf("StatusCode or ResponseRegex mismatch: %d", resp.StatusCode))
		log.Printf("FAILED - StatusCode or ResponseRegex mismatch: %d", resp.StatusCode)
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body for %s: %v", cfg.URL, err)
		}
	}

	// end timer
	end := time.Now()

	return PortResult{
		URL:           cfg.URL,
		Method:        httpMethod,
		Body:          cfg.Body,
		Online:        getTestResult(successCount, actualAttempts),
		StatusCode:    statusCode,
		StartTime:     start.Format(time.RFC3339),
		EndTime:       end.Format(time.RFC3339),
		TotalAttempts: actualAttempts,
		SuccessCount:  successCount,
		Failures:      failures,
		ResponseBody:  responseBody,
	}
}

// CheckServices checks all services defined in the configuration
func CheckServices(cfg *Config) []CheckResult {
	results := []CheckResult{}
	for _, svc := range cfg.Services {
		// start timer
		svcStart := time.Now()

		totalAttempts := 0
		successCount := 0
		totalPorts := 0
		onlinePorts := 0

		// check health ports
		healthResults := []PortResult{}
		for _, h := range svc.Health {
			pr := CheckPort(&h, svc.Timeout, svc.Retry, svc.Name, portType.HEALTH)
			healthResults = append(healthResults, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == testResult.ALL {
				onlinePorts++
			}
		}

		// check API ports
		apiResults := []PortResult{}
		for _, a := range svc.API {
			pr := CheckPort(&a, svc.Timeout, svc.Retry, svc.Name, portType.API)
			apiResults = append(apiResults, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == testResult.ALL {
				onlinePorts++
			}
		}

		// end timer
		svcEnd := time.Now()

		res := CheckResult{
			Name:          svc.Name,
			Online:        getTestResult(onlinePorts, totalPorts),
			Health:        healthResults,
			API:           apiResults,
			StartTime:     svcStart.Format(time.RFC3339),
			EndTime:       svcEnd.Format(time.RFC3339),
			TotalAttempts: totalAttempts,
			SuccessCount:  successCount,
		}
		results = append(results, res)
	}
	return results
}
