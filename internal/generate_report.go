package internal

import (
	"encoding/json"
	"github.com/wcy-dt/ponghub/protos/testResult"
	"html/template"
	"log"
	"os"
)

// GenerateReport generates an HTML report from the log data at logPath and writes it to outPath
func GenerateReport(logPath, outPath string) error {
	b, err := os.ReadFile(logPath)
	if err != nil {
		log.Fatalln("Failed to read log file:", err)
	}
	var logData map[string]map[string]any
	if err := json.Unmarshal(b, &logData); err != nil {
		log.Fatalln("Failed to parse log data:", err)
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
		Name         string
		History      []ServiceHistory
		Ports        map[string][]PortHistory
		Availability float64
	}
	var results []ServiceResult
	var latestTime string
	for svcName, svcData := range logData {
		serviceHistory := []ServiceHistory{}
		allCount := 0   // Count of "ALL" status in service history
		totalCount := 0 // Total count of service history entries
		if sh, ok := svcData["service_history"].([]any); ok {
			for _, entry := range sh {
				m, _ := entry.(map[string]any)
				status, _ := m["online"].(string)
				time, _ := m["time"].(string)
				serviceHistory = append(serviceHistory, ServiceHistory{Status: status, Time: time})
				totalCount++
				if status == testResult.ALL.String() {
					allCount++
				}
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

		// Calculate availability
		availability := float64(0)
		if totalCount > 0 {
			availability = float64(allCount) / float64(totalCount)
		}

		results = append(results, ServiceResult{
			Name:         svcName,
			History:      serviceHistory,
			Ports:        ports,
			Availability: availability,
		})
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
		"mul": func(a, b float64) float64 { return a * b },
	}
	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFiles("templates/report.html")
	if err != nil {
		log.Fatal("Failed to parse report template:", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatal("Failed to create report file:", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("Error closing report file:", err)
		}
	}(f)
	return tmpl.Execute(f, map[string]interface{}{
		"Results":    results,
		"UpdateTime": latestTime,
	})
}
