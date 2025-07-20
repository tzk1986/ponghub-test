package main

import (
	"fmt"

	ponghub "github.com/wcy-dt/ponghub/internal"
)

func main() {
	cfg, err := ponghub.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	results := ponghub.CheckServices(cfg)
	if err := ponghub.OutputResults(results, cfg.MaxLogDays); err != nil {
		fmt.Println("Error outputting results:", err)
	}
	if err := ponghub.GenerateReport("data/ponghub_log.json", "data/index.html"); err != nil {
		fmt.Println("Failed to generate report:", err)
	} else {
		fmt.Println("Report generated: data/index.html")
	}
}
