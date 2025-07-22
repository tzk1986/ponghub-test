package main

import (
	"log"

	ponghub "github.com/wcy-dt/ponghub/internal"
	"github.com/wcy-dt/ponghub/protos/defaultConfig"
)

func main() {
	// load the default configuration
	cfg, err := ponghub.LoadConfig(defaultConfig.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", defaultConfig.GetConfigPath(), ":", err)
	}

	// check services based on the configuration
	results := ponghub.CheckServices(cfg)
	if err := ponghub.OutputResults(results, cfg.MaxLogDays); err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := ponghub.GenerateReport(defaultConfig.GetLogPath(), defaultConfig.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", defaultConfig.GetReportPath())
	}
}
