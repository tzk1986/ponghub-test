package internal

import (
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
	RegexMatch    bool                  `json:"regex_match,omitempty"`
	StartTime     string                `json:"start_time"`
	EndTime       string                `json:"end_time"`
	TotalAttempts int                   `json:"total_attempts"`
	SuccessCount  int                   `json:"success_count"`
	Failures      []string              `json:"failures,omitempty"`
	ResponseBody  string                `json:"response_body,omitempty"`
}

// CheckPort checks a single port based on the provided configuration
func CheckPort(cfg *PortConfig, timeout int, retry int, svcName string, portType portType.PortType) PortResult {
	failures := []string{}
	successCount := 0
	actualAttempts := 0
	start := time.Now()
	var statusCode int
	var regexMatch bool
	var responseBody string
	method := cfg.Method
	if method == "" {
		method = "GET"
	}
	for attempt := 1; attempt <= retry; attempt++ {
		actualAttempts++
		client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
		log.Printf("[%s] Checking %s (attempt %d/%d)\n", svcName, cfg.URL, attempt, retry)
		req, err := http.NewRequest(method, cfg.URL, nil)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			continue
		}
		if cfg.Body != "" {
			req.Body = io.NopCloser(strings.NewReader(cfg.Body))
		}
		resp, err := client.Do(req)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			continue
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.Printf("Error closing response body for %s: %v", cfg.URL, err)
			}
		}(resp.Body)
		body, _ := io.ReadAll(resp.Body)
		responseBody = string(body)
		statusCode = resp.StatusCode
		var online bool
		if cfg.StatusCode == 0 && cfg.ResponseRegex == "" {
			online = true
		} else {
			online = true
			if cfg.StatusCode != 0 {
				online = online && resp.StatusCode == cfg.StatusCode
			}
			if cfg.ResponseRegex != "" {
				matched, _ := regexp.Match(cfg.ResponseRegex, body)
				regexMatch = matched
				online = online && matched
			}
		}
		if online {
			successCount++
			responseBody = ""
			break
		} else {
			failures = append(failures, fmt.Sprintf("StatusCode: %d, RegexMatch: %v", resp.StatusCode, regexMatch))
		}
	}
	end := time.Now()
	var onlineStatus testResult.TestResult
	switch successCount {
	case actualAttempts:
		onlineStatus = testResult.ALL
	case 0:
		onlineStatus = testResult.NONE
	default:
		onlineStatus = testResult.PART
	}
	return PortResult{
		URL:           cfg.URL,
		Method:        method,
		Body:          cfg.Body,
		Online:        onlineStatus,
		StatusCode:    statusCode,
		RegexMatch:    regexMatch,
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
		svcStart := time.Now()
		res := CheckResult{Name: svc.Name}
		timeout := svc.Timeout
		retry := svc.Retry
		totalAttempts := 0
		successCount := 0
		totalPorts := 0
		onlinePorts := 0
		for _, h := range svc.Health {
			pr := CheckPort(&h, timeout, retry, svc.Name, portType.HEALTH)
			res.Health = append(res.Health, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == testResult.ALL {
				onlinePorts++
			}
		}
		for _, a := range svc.API {
			pr := CheckPort(&a, timeout, retry, svc.Name, portType.API)
			res.API = append(res.API, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == testResult.ALL {
				onlinePorts++
			}
		}
		svcEnd := time.Now()
		res.StartTime = svcStart.Format(time.RFC3339)
		res.EndTime = svcEnd.Format(time.RFC3339)
		res.TotalAttempts = totalAttempts
		res.SuccessCount = successCount
		switch onlinePorts {
		case totalPorts:
			res.Online = testResult.ALL
		case 0:
			res.Online = testResult.NONE
		default:
			res.Online = testResult.PART
		}
		results = append(results, res)
	}
	return results
}
