package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
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
 */


var(
	logLevel = []string{"s"}
)

func TestSetupBaseConfig(t *testing.T) {
	newCmd := new(cobra.Command)
	initCmd(newCmd)
	cmd := NewRootCmd()
	newCmd.SetArgs([]string{
		"--tenable-access-key", "tak",
		"--tenable-secret-key", "tsk",
		"--tenable-scan-id", "tsi",
//		"--include-gcloud",
		"--gcloud-service-account-path", "/gsap",
		"--gcloud-project", "gp",
//		"--include-aws",
		"--generate-report",
		"--low-severity",
		"--medium-severity",
		"--high-severity",
		"--critical-severity",
		"--filter-search-type", "and",
		"--report-format", "pdf",
//		"--full-report",
		"--report-file-location", "./rfl",
	})
	fmt.Println("Here?")
	baseConfig, err := setupBaseConfig(cmd)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, baseConfig.IncludeAWS, true)
	assert.Equal(t, baseConfig.IncludeGCloud, true)
	assert.Equal(t, baseConfig.GCloudConfig.ProjectName, "gp")
	assert.Equal(t, baseConfig.GCloudConfig.ServiceAccountPath, "/gsap")
	assert.Equal(t, baseConfig.TenableConfig.GenerateReport, true)
	assert.Equal(t, baseConfig.TenableConfig.FilePath, "./rfl")
	assert.Equal(t, baseConfig.TenableConfig.SearchType, "and")
	assert.Equal(t, baseConfig.TenableConfig.ScanID, "tsi")
	assert.Equal(t, baseConfig.TenableConfig.SecretKey, "tsk")
	assert.Equal(t, baseConfig.TenableConfig.AccessKey, "tak")
	assert.Equal(t, baseConfig.TenableConfig.Format, "pdf")
}
