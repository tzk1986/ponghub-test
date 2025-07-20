package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func OutputResults(results []CheckResult, maxLogDays int) error {
	// Write main result
	f, err := os.Create("data/ponghub_result.json")
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return err
	}

	// Write log
	logPath := "data/ponghub_log.json"
	var logData = make(map[string]map[string]interface{})
	if b, err := os.ReadFile(logPath); err == nil {
		_ = json.Unmarshal(b, &logData)
	}
	now := time.Now()
	for _, svc := range results {
		// Service history
		if _, ok := logData[svc.Name]; !ok {
			logData[svc.Name] = map[string]interface{}{
				"service_history": []interface{}{},
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
			"online": svc.Online,
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
		uniquePorts := map[string]map[string]string{}
		for _, pr := range svc.Health {
			uniquePorts[pr.URL] = map[string]string{
				"time":   pr.StartTime,
				"online": pr.Online,
			}
		}
		for _, pr := range svc.API {
			uniquePorts[pr.URL] = map[string]string{
				"time":   pr.StartTime,
				"online": pr.Online,
			}
		}
		for url, entry := range uniquePorts {
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
	_ = os.WriteFile(logPath, logBytes, 0644)
	return nil
}
