package runner

import (
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


	ec2Mock := mock.Mock{}
	runner := Runner{
		ec2Svc:         ec2Mock,
		gcloud:         nil,
		tenable:        nil,
		includeGCloud:  false,
		includeAWS:     false,
		generateReport: false,
	}
	
	testCases := []testCase{
		{
			desc: "create export returns successfully",
			setup: func() {
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

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		/*
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.DownloadExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
		*/
	}
}