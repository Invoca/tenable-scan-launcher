package main

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/runner"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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

var (
	rootCmd = &cobra.Command{
		Use:   "tenable-scanner",
		Short: "Gets IPs and launches scans",
		Long: `tenable-scanner collects ip address from Google Cloud and AWS and launches a scan on the ips of the 
instances given based on the scanner id. It is also able to export the scans and downloads them`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return setupLogging(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			runner, err := setupRunner(cmd)
			if err != nil {
				return fmt.Errorf("RunE: Error seting up runner %s", err)
			}
			err = runner.Run()
			if err != nil {
				return fmt.Errorf("RunE: Error running runner %s", err)
			}
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "", "", "Log level (trace,info,fatal,panic,warn, debug) default is debug")
	rootCmd.PersistentFlags().StringP("log-type", "", "", "Log type (text,json)")

	rootCmd.PersistentFlags().StringP("tenable-access-key", "a", "", "tenable access key")
	rootCmd.PersistentFlags().StringP("tenable-secret-key", "s", "", "tenable secret key")
	rootCmd.PersistentFlags().StringP("tenable-scan-id", "i", "", "tenable scanID")


	rootCmd.PersistentFlags().BoolP("include-gcloud", "g", false, "Include Google Cloud Instances In Report")
	rootCmd.PersistentFlags().StringP("gcloud-service-account-path", "", "", "Path of service account token. Uses default if not specified")
	rootCmd.PersistentFlags().StringP("gcloud-project", "p", "", "GCloud project to list instances from")

	rootCmd.PersistentFlags().BoolP("include-aws", "A", false, "Include AWS Instances In Report")


	rootCmd.PersistentFlags().BoolP("generate-report", "R", false, "Generate A report after the scan is complete")
	rootCmd.PersistentFlags().BoolP("low-severity", "L", false, "Add Low Severity To Report")
	rootCmd.PersistentFlags().BoolP("medium-severity", "M", false, "Add Medium Severity To Report")
	rootCmd.PersistentFlags().BoolP("high-severity", "H", false, "Add High Severity To Report")
	rootCmd.PersistentFlags().BoolP("critical-severity", "C", false, "Add Critical Severity To Report")

	rootCmd.PersistentFlags().StringP("filter-search-type", "", "", "Search type to use in report. Only (and, or) are supported")
	rootCmd.PersistentFlags().StringP("report-format", "", "", "Report Format of the scan. Support formats are Nessus, HTML, PDF, CSV, or DB")
	rootCmd.PersistentFlags().StringP("report-chapters", "", "", "Chapters to include in the report")
	rootCmd.PersistentFlags().BoolP("summary-report", "S", false, "Generate A report in summary format")
	rootCmd.PersistentFlags().BoolP("full-report", "F", false, "Generate A report with all chapters")
	rootCmd.PersistentFlags().StringP("report-file-location", "", "", "File Location of the report")
}


func setupLogging(cmd *cobra.Command) error {
	logLevel, err  := cmd.Flags().GetString("log-level")
	if err != nil {
		return fmt.Errorf("setupLogging: error getting flag log-level")
	}

	logType, err  := cmd.Flags().GetString("log-type")
	if err != nil {
		return fmt.Errorf("setupLogging: error getting flag log-type")
	}

	if logType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "panic" {
		log.SetLevel(log.PanicLevel)
	} else if logLevel == "fatal" {
		log.SetLevel(log.FatalLevel)
	} else if logLevel == "trace" {
		log.SetLevel(log.TraceLevel)
	} else  {
		log.SetLevel(log.WarnLevel)
	}
	return nil
}

func setupSeverityFilter(low bool, medium bool, high bool, critical bool) ([]*tenable.Filter, error) {
	if (low || medium || high || critical) == false {
		return nil, fmt.Errorf("setupSeverityFilter: Cannot generate a report without specifiying severity filters")
	}

	var filters []*tenable.Filter
	filterName := "severity"
	filterQuality := "eq"

	if low {
		newFilter, err := tenable.CreateFilter(filterName, filterQuality, "Low")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Low Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if medium {
		newFilter, err := tenable.CreateFilter(filterName, filterQuality, "Medium")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Medium Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if high {
		newFilter, err := tenable.CreateFilter(filterName, filterQuality, "High")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating High Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if critical {
		newFilter, err := tenable.CreateFilter(filterName, filterQuality, "Critical")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Critical Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	return filters, nil
}

func setupTenableExport(cmd *cobra.Command) (*tenable.ExportSettings, error) {
	lowSeverity, err := cmd.Flags().GetBool("low-severity")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag low-severity")
	}

	mediumSeverity, err := cmd.Flags().GetBool("medium-severity")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag medium-severity")
	}

	highSeverity, err := cmd.Flags().GetBool("high-severity")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag high-severity")
	}

	criticalSeverity, err := cmd.Flags().GetBool("critical-severity")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag critical-severity")
	}

	fullReport, err := cmd.Flags().GetBool("full-report")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag full-report")
	}

	summaryReport, err  := cmd.Flags().GetBool("summary-report")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag summary-report")
	}

	searchType, err := cmd.Flags().GetString("filter-search-type")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag filter-search-type")
	}

	format, err := cmd.Flags().GetString("report-format")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag report-format")
	}

	chapters, err  := cmd.Flags().GetString("report-chapters")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag report-chapters")
	}

	filePath, err  := cmd.Flags().GetString("report-file-location")
	if err != nil {
		return nil, fmt.Errorf("setupTenableExport: error getting flag report-file-location")
	}

	if summaryReport {
		chapters = "vuln_hosts_summary"
	}

	if fullReport {
		chapters = "vuln_hosts_summary; vuln_by_host; compliance_exec; remediations; vuln_by_plugin; compliance"
	}

	if searchType == "" {
		return nil, fmt.Errorf("setupTenable: filter-search-type cannot be nil")
	}
	if format == "" {
		return nil, fmt.Errorf("setupTenable: format cannot be nil")
	}
	if chapters == "" {
		return nil, fmt.Errorf("setupTenable: chapters cannot be nil")
	}
	if filePath == "" {
		return nil, fmt.Errorf("setupTenable: filePath cannot be nil")
	}

	filters, err := setupSeverityFilter(lowSeverity, mediumSeverity, highSeverity, criticalSeverity)
	if err != nil {
		return nil, fmt.Errorf("setup: Error setting up export filterst %s", err)
	}

	return tenable.SetupExportSettings(filters, searchType, format, chapters, filePath)

}

func setupTenable(cmd *cobra.Command) (*tenable.Tenable, error) {
	var tenableExportSettings *tenable.ExportSettings
	var err error

	accessKey, err := cmd.Flags().GetString("tenable-access-key")
	if err != nil {
		return nil, fmt.Errorf("setupTenable: error getting flag tenable-access-key")
	}

	secretKey, err := cmd.Flags().GetString("tenable-secret-key")
	if err != nil {
		return nil, fmt.Errorf("setupTenable: error getting flag tenable-secret-key")
	}

	scanID, err := cmd.Flags().GetString("tenable-scan-id")
	if err != nil {
		return nil, fmt.Errorf("setupTenable: error getting flag tenable-scan-id")
	}

	generateReport, err := cmd.Flags().GetBool("generate-report")
	if err != nil {
		return nil, fmt.Errorf("setupTenable: error getting flag generate-report")
	}

	if accessKey == "" {
		return nil, fmt.Errorf("setupTenable: accessKey cannot be nil")
	}
	if secretKey == "" {
		return nil, fmt.Errorf("setupTenable: secretKey cannot be nil")
	}
	if scanID == "" {
		return nil, fmt.Errorf("setupTenable: scanID cannot be nil")
	}

	if generateReport {
		tenableExportSettings, err = setupTenableExport(cmd)
		if err != nil {
			return nil, fmt.Errorf("setupTenable: Error creating Tenable Export Object %s", err)
		}
	} else {
		tenableExportSettings = &tenable.ExportSettings{}
	}

	tenableClient := tenable.SetupClient(accessKey,secretKey,scanID, tenableExportSettings)
	return tenableClient, nil
}

// CreateGCloudInterface

func setupGCloud(cmd *cobra.Command) (*cloud.GCloudWrapper, error) {
	serviceAccountPath, err := cmd.Flags().GetString("gcloud-service-account-path")
	if err != nil {
		return nil, fmt.Errorf("setupGCloud: error getting flag gcloud-service-account-path")
	}

	gcloudProject, err := cmd.Flags().GetString("gcloud-project")
	if err != nil {
		return nil, fmt.Errorf("setupGCloud: error getting flag tenable-secret-key")
	}
	gCloudWrapper, err := cloud.CreateGCloudInterface(gcloudProject, serviceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("setupGCloud: error creating GCloud Interface")
	}
	return gCloudWrapper, nil
}

func setupAWS() (*ec2.EC2, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	if sess == nil {
		return nil, fmt.Errorf("setupAWS: Error creating session object")
	}
	return ec2.New(sess), nil
}


func setupRunner(cmd *cobra.Command) (*runner.Runner, error) {
	var gCloud *cloud.GCloudWrapper
	var aws *ec2.EC2

	includeGCloud, err := cmd.Flags().GetBool("include-gcloud")
	if err != nil {
		return nil, fmt.Errorf("setupRunner: error getting flag include-gcloud")
	}

	includeAWS, err := cmd.Flags().GetBool("include-aws")
	if err != nil {
		return nil, fmt.Errorf("setupRunner: error getting flag include-aws")
	}

	tenableClient, err := setupTenable(cmd)
	if err != nil {
		return nil, fmt.Errorf("setupRunner: Error seting up tenableClient %s", err)
	}

	if includeGCloud {
		gCloud, err = setupGCloud(cmd)
		if err != nil {
			return nil, fmt.Errorf("setupRunner: Error seting up GCloud %s", err)
		}
	}

	if includeAWS {
		aws, err = setupAWS()
		if err != nil {
			return nil, fmt.Errorf("setupRunner: Error seting up GCloud %s", err)
		}
	}


	runner, err := runner.SetupRunner(tenableClient, gCloud, aws, includeGCloud, includeAWS)
	if err != nil {
		return nil, fmt.Errorf("setupRunner: Error seting up runner %s", err)
	}
	return runner, nil
}