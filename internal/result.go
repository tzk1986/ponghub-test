package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wcy-dt/ponghub/protos/defaultConfig"
	"github.com/wcy-dt/ponghub/protos/testResult"
)

// MergeOnlineStatus merges a list of online statuses into a single status
func MergeOnlineStatus(statusList []testResult.TestResult) testResult.TestResult {
	if len(statusList) == 0 {
		return testResult.NONE
	}

	hasNone, hasAll := false, false
	for _, s := range statusList {
		switch s {
		case testResult.NONE:
			hasNone = true
		case testResult.ALL:
			hasAll = true
		}
	}

	switch {
	case hasNone && !hasAll:
		return testResult.NONE
	case !hasNone && hasAll:
		return testResult.ALL
	default:
		return testResult.PART
	}
}

// OutputResults writes the check results to a JSON file and updates the log file
func OutputResults(results []CheckResult, maxLogDays int) error {
	// Open result file
	f, err := os.Create(defaultConfig.GetResultPath())
	if err != nil {
		log.Fatalln("Failed to create result file:", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("Error closing result file:", err)
		}
	}(f)

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		log.Fatalln("Failed to write results to file:", err)
	}

	// Get existing log data or create a new map
	var logData = make(map[string]map[string]any)
	if b, err := os.ReadFile(defaultConfig.GetLogPath()); err == nil {
		if err := json.Unmarshal(b, &logData); err != nil {
			log.Fatalln("Failed to read existing log file:", err)
		}
	}

	now := time.Now()
	for _, svc := range results {
		// Service history
		if _, ok := logData[svc.Name]; !ok {
			logData[svc.Name] = map[string]any{
				"service_history": []any{},
				"ports":           map[string][]any{},
			}
		}
		// Handle service_history type
		svcHistoryRaw := logData[svc.Name]["service_history"]
		var svcHistory []map[string]string
		switch v := svcHistoryRaw.(type) {
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					entry := map[string]string{}
					for k, val := range m {
						entry[k] = fmt.Sprintf("%v", val)
					}
					svcHistory = append(svcHistory, entry)
				}
			}
		case []map[string]string:
			svcHistory = v
		}
		svcHistory = append(svcHistory, map[string]string{
			"time":   svc.StartTime,
			"online": svc.Online.String(),
		})
		// Clean up timeout records
		var filteredSvcHistory []map[string]string
		for _, entry := range svcHistory {
			t, err := time.Parse(time.RFC3339, entry["time"])
			if err == nil && now.Sub(t).Hours() <= float64(maxLogDays*24) {
				filteredSvcHistory = append(filteredSvcHistory, entry)
			}
		}
		logData[svc.Name]["service_history"] = filteredSvcHistory

		// Handle ports type
		portsRaw := logData[svc.Name]["ports"]
		portsMap := map[string][]map[string]string{}
		switch v := portsRaw.(type) {
		case map[string]any:
			for url, arr := range v {
				var portHistory []map[string]string
				if arrList, ok := arr.([]any); ok {
					for _, item := range arrList {
						if m, ok := item.(map[string]any); ok {
							entry := map[string]string{}
							for k, val := range m {
								entry[k] = fmt.Sprintf("%v", val)
							}
							portHistory = append(portHistory, entry)
						}
					}
				}
				portsMap[url] = portHistory
			}
		case map[string][]map[string]string:
			portsMap = v
		}
		// Only record one port entry for each unique URL per complete run
		urlStatusMap := map[string][]string{}
		urlTimeMap := map[string]string{}
		for _, pr := range svc.Health {
			urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online.String())
			if urlTimeMap[pr.URL] == "" {
				urlTimeMap[pr.URL] = pr.StartTime
			}
		}
		for _, pr := range svc.API {
			urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online.String())
			if urlTimeMap[pr.URL] == "" {
				urlTimeMap[pr.URL] = pr.StartTime
			}
		}
		for url, statusList := range urlStatusMap {
			mergedStatus := MergeOnlineStatus(testResult.ParseTestResults(statusList))
			entry := map[string]string{
				"time":   urlTimeMap[url],
				"online": mergedStatus.String(),
			}
			portsMap[url] = append(portsMap[url], entry)
		}
		// Clean up expired port records
		for url, history := range portsMap {
			var filteredPortHistory []map[string]string
			for _, entry := range history {
				t, err := time.Parse(time.RFC3339, entry["time"])
				if err == nil && now.Sub(t).Hours() <= float64(maxLogDays*24) {
					filteredPortHistory = append(filteredPortHistory, entry)
				}
			}
			portsMap[url] = filteredPortHistory
		}
		logData[svc.Name]["ports"] = portsMap
	}
	logBytes, _ := json.MarshalIndent(logData, "", "  ")
	_ = os.WriteFile(defaultConfig.GetLogPath(), logBytes, 0644)
	return nil
}
