package runner

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCase struct {
	desc        string
	setup       func()
	shouldError bool
}

func TestRun(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	gcloudInstances := []string{"1.1.1.1"}
	awsInstances := []string{"2.2.2.2"}

	//instanceList := []string{"1.1.1.1", "2.2.2.2"}

	//instanceList = append(gcloudInstances, awsInstances...)

	ec2Mock := &mocks.MockCloudAPI{}

	gcloudMock := &mocks.MockCloudAPI{}

	tenableMock := &mocks.MockTenableAPI{}

	slackMock := mocks.SlackInterfaceMock{}

	alerts := tenable.Alerts{

		Vulnerabilities:         []tenable.Vulnerabilities{},
		TotalVulnerabilityCount: 1,
		TotalAssetCount:         1,
	}

	runner := Runner{
		ec2Svc:         ec2Mock,
		gcloud:         gcloudMock,
		tenable:        tenableMock,
		slackSvc:       &slackMock,
		includeGCloud:  true,
		includeAWS:     true,
		generateReport: true,
	}

	testCases := []testCase{
		{
			desc: "AWS and GCloud are able to Gather IPs without issue",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, nil)
				slackMock.On("PrintAlerts", mock.Anything).Return(nil)
				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("StartExport", mock.Anything).Return(nil)
				tenableMock.On("WaitForExport", mock.Anything).Return(nil)
				tenableMock.On("DownloadExport", mock.Anything).Return(nil)
			},
			shouldError: false,
		},
		{
			desc: "AWS is not able to GatherIPs",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = nil
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil, fmt.Errorf("Error!"))
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)
			},
			shouldError: true,
		},
		{
			desc: "Gcloud API is not able to GatherIPs",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil, fmt.Errorf("Error!"))
			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to SetTargets",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(fmt.Errorf("Tenable Error"))
			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to LaunchScan",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(fmt.Errorf("Tenable Error"))
			},
			shouldError: true,
		},

		{
			desc: "Tenable is not able to WaitForScanToComplete",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(fmt.Errorf("Tenable Error"))
			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to get list of vulnerabilites from Tenable dashboard",
			setup: func() {

				//TODO: Create Test cases with Max
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, fmt.Errorf("Error getting Vulnerabilities from Dashboard"))

			},
			shouldError: true,
		},
		{
			desc: "Alerts are not able to be posted to slack",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, nil)
				slackMock.On("PrintAlerts", mock.Anything).Return(fmt.Errorf("Error"))

			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to StartExport",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, nil)
				slackMock.On("PrintAlerts", mock.Anything).Return(nil)
				tenableMock.On("StartExport", mock.Anything).Return(fmt.Errorf("Error!"))
			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to WaitForExport",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, nil)
				slackMock.On("PrintAlerts", mock.Anything).Return(nil)
				tenableMock.On("StartExport", mock.Anything).Return(nil)
				tenableMock.On("WaitForExport", mock.Anything).Return(fmt.Errorf("Error!"))
			},
			shouldError: true,
		},
		{
			desc: "Tenable is not able to DownloadExport",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()
				slackMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("GetVulnerabilities", mock.Anything).Return(&alerts, nil)
				slackMock.On("PrintAlerts", mock.Anything).Return(nil)
				tenableMock.On("StartExport", mock.Anything).Return(nil)
				tenableMock.On("WaitForExport", mock.Anything).Return(nil)
				tenableMock.On("DownloadExport", mock.Anything).Return(fmt.Errorf("Error!"))
			},
			shouldError: true,
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc":        testCase.desc,
			"shouldError": testCase.shouldError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		if testCase.setup != nil {
			log.Debug("Setting up testCase")
			testCase.setup()
		}

		err := runner.Run()

		log.WithFields(log.Fields{
			"shouldError": testCase.shouldError,
			"Error":       err,
		}).Debug("Run() complete")

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
