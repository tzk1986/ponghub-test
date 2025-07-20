package internal

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type CheckResult struct {
	Name          string       `json:"name"`
	Online        string       `json:"online"`
	Health        []PortResult `json:"health,omitempty"`
	API           []PortResult `json:"api,omitempty"`
	StartTime     string       `json:"start_time"`
	EndTime       string       `json:"end_time"`
	TotalAttempts int          `json:"total_attempts"`
	SuccessCount  int          `json:"success_count"`
}
type PortResult struct {
	URL           string   `json:"url"`
	Method        string   `json:"method"`
	Body          string   `json:"body,omitempty"`
	Online        string   `json:"online"`
	StatusCode    int      `json:"status_code,omitempty"`
	RegexMatch    bool     `json:"regex_match,omitempty"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
	TotalAttempts int      `json:"total_attempts"`
	SuccessCount  int      `json:"success_count"`
	Failures      []string `json:"failures,omitempty"`
	ResponseBody  string   `json:"response_body,omitempty"`
}

func CheckPort(cfg *PortConfig, timeout int, retry int, svcName string, portType string) PortResult {
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
		fmt.Printf("[%s][%s] Checking %s (attempt %d/%d)\n", svcName, portType, cfg.URL, attempt, retry)
		req, err := http.NewRequest(method, cfg.URL, nil)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if cfg.Body != "" {
			req.Body = io.NopCloser(strings.NewReader(cfg.Body))
		}
		resp, err := client.Do(req)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		defer resp.Body.Close()
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
	var onlineStatus string
	switch successCount {
	case actualAttempts:
		onlineStatus = "all"
	case 0:
		onlineStatus = "none"
	default:
		onlineStatus = "part"
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
			pr := CheckPort(&h, timeout, retry, svc.Name, "health")
			res.Health = append(res.Health, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == "all" {
				onlinePorts++
			}
		}
		for _, a := range svc.API {
			pr := CheckPort(&a, timeout, retry, svc.Name, "api")
			res.API = append(res.API, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == "all" {
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
			res.Online = "all"
		case 0:
			res.Online = "none"
		default:
			res.Online = "part"
		}
		results = append(results, res)
	}
	return results
}
