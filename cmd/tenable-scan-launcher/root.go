package main

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/runner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
	Environment Variables needed:
	AWS_ACCESS_KEY
	AWS_SECRET_KEY

	TENABLE_ACCESS_KEY
	TENABLE_SECRET_KEY

	Flags needed:
	--log_level Log Level
	--log_type Log Type


	--include_gcloud Use GCloud
	--gcloud_json GCLoud credentials location
	--gcloud_project GCLoud project to use


	--include_aws

	--scanner_id Which Tenable Scanner to use


	--generate_report boolean, determine
	Report Filters https://developer.tenable.com/docs/scan-export-filters-tio
	filter.n.value (low, medium, high for now)
	filter.n.quality (eq, neq for now)
	filter.n.filter (only severity for now)
	filter.search_type (and,or)

	--low_severity
	--medium_severity
	--high_severity
	--search_type

	--format file format (Nessus, HTML, PDF, CSV, or DB)
	--chapters (vuln_hosts_summary, vuln_by_host, compliance_exec, remediations, vuln_by_plugin, compliance)
	--full-report (vuln_hosts_summary, vuln_by_host, compliance_exec, remediations, vuln_by_plugin, compliance)
    --summary-report(vuln_hosts_summary)
	--report-file-location

*/

type logConfig struct {
	LogLevel string
	LogType string
}

func NewRootCmd() *cobra.Command {
	baseConfig := config.BaseConfig{}
	tenableConfig := config.TenableConfig{}
	gcloudConfig := config.GCloudConfig{}
	baseConfig.TenableConfig = &tenableConfig
	baseConfig.GCloudConfig = &gcloudConfig

	logConfig := logConfig{}
	cmd := &cobra.Command{
		Use:   "tenable-scanner",
		Short: "Gets IPs and launches scans",
		Long: `tenable-scanner collects ip address from Google Cloud and AWS and launches a scan on the ips of the 
instances given based on the scanner id. It is also able to export the scans and downloads them`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setupLogging(&logConfig)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Setting up runner")

			runnerSvc, err := setupRunner(&baseConfig)
			if err != nil {
				return fmt.Errorf("RunE: Error seting up runner %s", err)
			}

			log.Debug("Setup completed. Running command")
			err = runnerSvc.Run()
			if err != nil {
				return fmt.Errorf("RunE: Error running runner %s", err)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&logConfig.LogLevel,"log-level", "", "", "Log level (trace,info,fatal,panic,warn, debug) default is debug")
	f.StringVarP(&logConfig.LogType, "log-type", "", "", "Log type (text,json)")

	f.StringVarP(&baseConfig.TenableConfig.AccessKey, "tenable-access-key", "a", "", "tenable access key")
	f.StringVarP(&baseConfig.TenableConfig.SecretKey, "tenable-secret-key", "s", "", "tenable secret key")
	f.StringVarP(&baseConfig.TenableConfig.ScanID, "tenable-scan-id", "i", "", "tenable scanID")


	f.BoolVarP(&baseConfig.IncludeGCloud, "include-gcloud", "g",false, "Include Google Cloud Instances In Report")
	f.StringVarP(&baseConfig.GCloudConfig.ServiceAccountPath ,"gcloud-service-account-path", "", "", "Path of service account token. Uses default if not specified")
	f.StringVarP(&baseConfig.GCloudConfig.ProjectName, "gcloud-project", "p", "", "GCloud project to list instances from")

	f.BoolVarP(&baseConfig.IncludeAWS, "include-aws", "A", false, "Include AWS Instances In Report")


	f.BoolVarP(&baseConfig.TenableConfig.GenerateReport, "generate-report", "R", false, "Generate A report after the scan is complete")
	f.BoolVarP(&baseConfig.TenableConfig.LowSeverity, "low-severity", "L", false, "Add Low Severity To Report")
	f.BoolVarP(&baseConfig.TenableConfig.MediumSeverity, "medium-severity", "M", false, "Add Medium Severity To Report")
	f.BoolVarP(&baseConfig.TenableConfig.HighSeverity, "high-severity", "H", false, "Add High Severity To Report")
	f.BoolVarP(&baseConfig.TenableConfig.CriticalSeverity, "critical-severity", "C", false, "Add Critical Severity To Report")

	f.StringVarP(&baseConfig.TenableConfig.SearchType, "filter-search-type", "", "", "Search type to use in report. Only (and, or) are supported")
	f.StringVarP(&baseConfig.TenableConfig.Format, "report-format", "", "", "Report Format of the scan. Support formats are Nessus, HTML, PDF, CSV, or DB")
	f.StringVarP(&baseConfig.TenableConfig.Chapters, "report-chapters", "", "", "Chapters to include in the report")
	f.BoolVarP(&baseConfig.TenableConfig.SummaryReport, "summary-report", "S", false, "Generate A report in summary format")
	f.BoolVarP(&baseConfig.TenableConfig.FullReport, "full-report", "F", false, "Generate A report with all chapters")
	f.StringVarP(&baseConfig.TenableConfig.FilePath, "report-file-location", "", "", "File Location of the report")

	return cmd
}

func setupLogging(logConfig *logConfig) error {

	if logConfig.LogType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	if logConfig.LogLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if logConfig.LogLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logConfig.LogLevel == "panic" {
		log.SetLevel(log.PanicLevel)
	} else if logConfig.LogLevel == "fatal" {
		log.SetLevel(log.FatalLevel)
	} else if logConfig.LogLevel == "trace" {
		log.SetLevel(log.TraceLevel)
	} else  {
		log.SetLevel(log.WarnLevel)
	}
	return nil
}

func setupRunner(baseConfig *config.BaseConfig) (*runner.Runner, error) {
	runner := &runner.Runner{}

	err := runner.SetupRunner(baseConfig)
	if err != nil {
		return nil, fmt.Errorf("setupRunner: Error seting up runner %s", err)
	}

	return runner, nil
}