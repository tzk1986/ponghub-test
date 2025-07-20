package internal

import (
	"encoding/json"
	"html/template"
	"os"
)

func GenerateReport(logPath, outPath string) error {
	b, err := os.ReadFile(logPath)
	if err != nil {
		return err
	}
	var logData map[string]map[string]any
	if err := json.Unmarshal(b, &logData); err != nil {
		return err
	}

	// Build template data structure
	type ServiceHistory struct {
		Status string
		Time   string
	}
	type PortHistory struct {
		URL    string
		Time   string
		Status string
	}
	type ServiceResult struct {
		Name    string
		History []ServiceHistory
		Ports   map[string][]PortHistory
	}
	var results []ServiceResult
	var latestTime string
	for svcName, svcData := range logData {
		serviceHistory := []ServiceHistory{}
		if sh, ok := svcData["service_history"].([]any); ok {
			for _, entry := range sh {
				m, _ := entry.(map[string]any)
				status, _ := m["online"].(string)
				time, _ := m["time"].(string)
				serviceHistory = append(serviceHistory, ServiceHistory{Status: status, Time: time})
				if time > latestTime {
					latestTime = time
				}
			}
		}
		ports := map[string][]PortHistory{}
		if portMap, ok := svcData["ports"].(map[string]any); ok {
			for url, historyRaw := range portMap {
				if historyArr, ok := historyRaw.([]any); ok {
					for _, entry := range historyArr {
						m, _ := entry.(map[string]any)
						status, _ := m["online"].(string)
						time, _ := m["time"].(string)
						ports[url] = append(ports[url], PortHistory{URL: url, Time: time, Status: status})
						if time > latestTime {
							latestTime = time
						}
					}
				}
			}
		}
		results = append(results, ServiceResult{Name: svcName, History: serviceHistory, Ports: ports})
	}

	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"until": func(n int) []int {
			arr := make([]int, n)
			for i := range n {
				arr[i] = i
			}
			return arr
		},
	}
	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFiles("templates/report.html")
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, map[string]interface{}{"Results": results, "UpdateTime": latestTime})
}
