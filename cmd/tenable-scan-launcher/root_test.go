package main

import (
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type setupRunnerTestCast struct {
	desc        string
	baseConfig   *config.BaseConfig
	shouldError bool
}

func TestSetupBaseConfig(t *testing.T) {

	boolFlags := []string{
		"include-gcloud",
		"include-aws",
		"generate-report",
		"low-severity",
		"medium-severity",
		"high-severity",
		"critical-severity",
		"full-report",
		"summary-report",
	}

	stringFlags := []string{
	"tenable-access-key",
	"tenable-secret-key",
	"tenable-scan-id",
	"gcloud-service-account-path",
	"gcloud-project",
	"filter-search-type",
	"report-format",
	"report-file-location"}
	newCmd := NewRootCmd()

	for _, f := range stringFlags {
		if newCmd.Flags().Lookup(f) == nil {
			t.Fatalf("generate command should have flag %s", f)
		}
		_, err := newCmd.Flags().GetString(f)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	for _, f := range boolFlags {
		if newCmd.Flags().Lookup(f) == nil {
			t.Fatalf("generate command should have flag %s", f)
		}
		_, err := newCmd.Flags().GetBool(f)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}
}

type loggingPair struct {
	loglevelFromFlag string
	expectedLoglevel log.Level
}

//TODO: Add method of testing logging type
func TestSetupLogging(t *testing.T) {

	lp := []loggingPair{
		{
			loglevelFromFlag: "trace",
			expectedLoglevel: log.TraceLevel,
		},
		{
			loglevelFromFlag: "debug",
			expectedLoglevel: log.DebugLevel,
		},
		{
			loglevelFromFlag: "info",
			expectedLoglevel: log.InfoLevel,
		},
		{
			loglevelFromFlag: "panic",
			expectedLoglevel: log.PanicLevel,
		},
		{
			loglevelFromFlag: "fatal",
			expectedLoglevel: log.FatalLevel,
		},
	}

	lc := &logConfig{}

	for _, logPair := range lp {
		lc.LogLevel = logPair.loglevelFromFlag
		err := setupLogging(lc)
		if err != nil {
			t.Fatalf("Error! %s", err)
		}

		if log.GetLevel() != logPair.expectedLoglevel {
			t.Fatalf("Error! Log level not expected")
		}
	}
}

func TestSetupRunner(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	testCases := []*setupRunnerTestCast {
		{
			desc: "It should fail when the config struct is empty.",
			baseConfig: &config.BaseConfig{},
			shouldError: true,
		},
		{
			desc: "It should fail when the TenableConfig struct is nil.",
			baseConfig: &config.BaseConfig{
				GCloudConfig: &config.GCloudConfig{},
			},
			shouldError: true,
		},
		{
			desc: "It should fail when the GCloudConfig struct is nil.",
			baseConfig: &config.BaseConfig{
				TenableConfig: &config.TenableConfig{},
			},
			shouldError: true,
		},
		{
			desc: "It should fail when no Tenable credentials are set",
			baseConfig: &config.BaseConfig{
				GCloudConfig: &config.GCloudConfig{},
				TenableConfig: &config.TenableConfig{},
			},
			shouldError: true,
		},
		{
			desc: "It should fail when export report is set and no severity is set",
			baseConfig: &config.BaseConfig{
				GCloudConfig: &config.GCloudConfig{},
				TenableConfig: &config.TenableConfig{
					AccessKey: "ak",
					SecretKey: "sk",
					GenerateReport: true,
				},
			},
			shouldError: true,
		},
		{
			desc: "It should not fail when all of the required fields are set",
			baseConfig: &config.BaseConfig{
				GCloudConfig: &config.GCloudConfig{},
				TenableConfig: &config.TenableConfig{
					AccessKey: "ak",
					SecretKey: "sk",
					GenerateReport: true,
					LowSeverity: true,
				},
			},
			shouldError: false,
		},
		{
			desc: "It should not fail when exporting a report is not set and is missing fields for exporting.",
			baseConfig: &config.BaseConfig{
				GCloudConfig: &config.GCloudConfig{},
				TenableConfig: &config.TenableConfig{
					AccessKey: "ak",
					SecretKey: "sk",
				},
			},
			shouldError: false,
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		_, err := setupRunner(testCase.baseConfig)

		log.WithFields(log.Fields{
			"shouldError": testCase.shouldError,
			"err": err,
		}).Debug("Finished running testCase " + strconv.Itoa(index))

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
// func setupRunner(baseConfig *config.BaseConfig) (*runner.Runner, error) {

