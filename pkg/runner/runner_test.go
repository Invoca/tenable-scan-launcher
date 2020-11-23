package runner

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
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

	ec2Mock := &mocks.MockCloudAPI{}

	gcloudMock := &mocks.MockCloudAPI{}

	tenableMock := &mocks.MockTenableAPI{}

	runner := Runner{
		awsSvc:         ec2Mock,
		gcloud:         gcloudMock,
		tenable:        tenableMock,
		includeGCloud:  true,
		includeAWS:     true,
		generateReport: true,
		osFs:           afero.NewMemMapFs(),
		fileLocation:   "/tmp/report",
	}

	file, _ := runner.osFs.Create(runner.fileLocation)
	file.Write([]byte(""))
	file.Close()

	testCases := []testCase{
		{
			desc: "AWS and GCloud are able to Gather IPs without issue",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
				tenableMock.On("StartExport", mock.Anything).Return(nil)
				tenableMock.On("WaitForExport", mock.Anything).Return(nil)
				tenableMock.On("DownloadExport", mock.Anything).Return(nil)
				ec2Mock.On("UploadFile", mock.Anything).Return(nil)
			},
			shouldError: false,
		},
		{
			desc: "AWS is not able to GatherIPs",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()

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
			desc: "Tenable is not able to StartExport",
			setup: func() {
				ec2Mock.Reset()
				gcloudMock.Reset()
				tenableMock.Reset()

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
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

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
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

				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances

				ec2Mock.On("GatherIPs", mock.Anything).Return(awsInstances, nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(gcloudInstances, nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan", mock.Anything).Return(nil)
				tenableMock.On("WaitForScanToComplete", mock.Anything).Return(nil)
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
