package runner

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

type testCase struct {
	desc        string
	setup       func()
	post 		func()
	shouldError bool
	returnError bool
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

	runner := Runner{
		ec2Svc:         ec2Mock,
		gcloud:         gcloudMock,
		tenable:        tenableMock,
		includeGCloud:  true,
		includeAWS:     true,
		generateReport: false,
	}
	
	testCases := []testCase{
		{
			desc: "AWS and GCloud are able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			post: func() {
				ec2Mock.IPs = *new([]string)
				gcloudMock.IPs = *new([]string)
			},
			shouldError: false,
		},
		{
			desc: "AWS is not able to GatherIPs",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = gcloudInstances
				ec2Mock.On("FetchIPs", mock.Anything).Return(fmt.Errorf("EC2 Error"))
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "Gcloud API is not able to GatherIPs",
			setup: func() {
				ec2Mock.IPs = nil
				gcloudMock.IPs = gcloudInstances
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(fmt.Errorf("GCloud Error"))

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(fmt.Errorf("AWS Error"))
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(fmt.Errorf("GCloud Error"))

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(fmt.Errorf("Tenable Error"))
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(fmt.Errorf("Tenable Error"))
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(fmt.Errorf("Tenable Error"))
				tenableMock.On("WaitForScanToComplete").Return(nil)
			},
			shouldError: true,
		},
		{
			desc: "AWS and GCloud are both not able to Gather IPs without issue",
			setup: func() {
				ec2Mock.IPs = awsInstances
				gcloudMock.IPs = nil
				ec2Mock.On("FetchIPs", mock.Anything).Return(nil)
				gcloudMock.On("FetchIPs", mock.Anything).Return(nil)

				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)

				tenableMock.On("SetTargets", mock.Anything).Return(nil)
				tenableMock.On("LaunchScan").Return(nil)
				tenableMock.On("SetTargets").Return(nil)
				tenableMock.On("WaitForScanToComplete").Return(fmt.Errorf("Tenable Error"))
			},
			shouldError: true,
		},

	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"returnError": testCase.returnError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		if testCase.setup != nil {
			testCase.setup()
		}

		err := runner.Run()

		if testCase.post != nil {
			testCase.post()
		}

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}