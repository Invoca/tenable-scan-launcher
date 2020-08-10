package runner

import (
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
	shouldError bool
	returnError bool
}

func TestRun(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	/*
	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		scanID:     scanID,
		scanUuid: scanUuid,
		fileId: fileId,
	}

	statusPath := "/scans/" + scanID + "/export/" + fileId + "/download"
	 */


	ec2Mock := &mocks.MockCloudAPI{}

	gcloudMock := &mocks.MockCloudAPI{}

	tenableMock := &mocks.MockTenableAPI{}

	runner := Runner{
		ec2Svc:         ec2Mock,
		gcloud:         gcloudMock,
		tenable:        tenableMock,
		includeGCloud:  false,
		includeAWS:     false,
		generateReport: false,
	}
	
	testCases := []testCase{
		{
			desc: "create export returns successfully",
			setup: func() {
				ec2Mock.On("GatherIPs", mock.Anything).Return(nil)
				gcloudMock.On("GatherIPs", mock.Anything).Return(nil)
			},
			shouldError: false,
		},
		{
			desc: "create export returns empty body",
			setup: func() {
			},
			shouldError: false,
		},
		{
			desc: "create export scanID is not set",
			setup: func() {
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

		testCase.setup()
		err := runner.Run()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}