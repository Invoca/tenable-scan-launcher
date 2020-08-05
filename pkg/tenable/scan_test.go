package tenable

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type getRegionsTestCast struct {
	desc        string
	setup       func()
	shouldError bool
	expectedPath string
	returnError bool
	requestBodies [][]byte
}

func setupTenable(t *testing.T, tenable *Tenable, requestBodies [][]byte, returnError bool, expectedPath string) *Tenable {
	counter := 0
	apiKeyFormat := "accessKey=" + tenable.accessKey + "; secretKey=" +  tenable.secretKey + ";"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if returnError {
			rw.WriteHeader(500)
		} else {
			generatedUrl := req.URL.Hostname() + req.URL.Port() + expectedPath
			apiKeyHeader := req.Header.Get("X-ApiKeys")
			assert.Equal(t, apiKeyHeader, apiKeyFormat)
			assert.Equal(t, generatedUrl, req.URL.String())

			if counter < len(requestBodies) {
				rw.Write(requestBodies[counter])
				counter += 1
			}
		}
	}))

	// defer server.Close()

	tenable.tenableURL = server.URL
	return tenable
}


func TestLaunchScan(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		Targets:    nil,
		scanID:     scanID,
		status: &scanStatus{
			Pending:   false,
			Running:   false,
		},
		scanUuid: scanUuid,
	}

	statusPath := "/scans/" + tenable.scanID + "/launch"

	testCases := []getRegionsTestCast{
		{
			desc: "launching a scan returns successfully",
			setup: func() {
			},
			shouldError: false,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"scan_uuid":"ABC"}`),
			},
		},
		{
			desc: "launching a scan returns invalid json",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{Not valid JSON}`),
			},
		},
		{
			desc: "launching a scan returns empty body",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "launching a scan with an empty scanId is not set",
			setup: func() {
				tenable.scanUuid = ""
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "launching a scan returns a non-200 code",
			setup: func() {
				tenable.scanUuid = scanUuid
			},
			shouldError: true,
			returnError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"scan_uuid":"ABC"}`),
			},
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"expectedPath": testCase.expectedPath,
			"returnError": testCase.returnError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()
		tenable := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.LaunchScan()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestWaitForScanToComplete(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		Targets:    nil,
		scanID:     scanID,
		status: &scanStatus{
			Pending:   false,
			Running:   false,
		},
		scanUuid: scanUuid,
	}

	statusPath := "/scans/" + scanID + "/latest-status"

	testCases := []getRegionsTestCast{
		{
			desc: "scan finishes successfully",
			setup: func() {
			},
			shouldError: false,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"status":"pending"}`),
				[]byte(`{"status":"running"}`),
				[]byte(`{"status":"completed"}`),
			},
		},
		{
			desc: "scan status returns invalid json",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{Not valid JSON}`),
			},
		},
		{
			desc: "scan status returns empty body",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "scanUuid is not set",
			setup: func() {
				tenable.scanUuid = ""
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "scan status returns a non-200 code",
			setup: func() {
				tenable.scanUuid = scanUuid
			},
			shouldError: true,
			returnError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"status":"completed"}`),
			},
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"expectedPath": testCase.expectedPath,
			"returnError": testCase.returnError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()
		tenable := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.WaitForScanToComplete()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStartExport(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		Targets:    nil,
		scanID:     scanID,
		status: &scanStatus{
			Pending:   false,
			Running:   false,
		},
		scanUuid: scanUuid,
	}

	statusPath := "/scans/" + scanID + "/export"

	testCases := []getRegionsTestCast{
		{
			desc: "create export returns successfully",
			setup: func() {
			},
			shouldError: false,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"file": 123}`),
			},
		},
		{
			desc: "create export returns invalid json",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{Not valid JSON}`),
			},
		},
		{
			desc: "create export returns empty body",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "create export scanID is not set",
			setup: func() {
				tenable.scanID = ""
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "scan status returns a non-200 code",
			setup: func() {
				tenable.scanID = scanID
			},
			shouldError: true,
			returnError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"file": 123}`),
			},
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"expectedPath": testCase.expectedPath,
			"returnError": testCase.returnError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()
		tenable := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.StartExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestWaitForExport(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"
	fileId    := "111"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		scanID:     scanID,
		scanUuid: scanUuid,
		fileId: fileId,
	}

	statusPath := "/scans/" + scanID + "/export/" + fileId + "/status"

	testCases := []getRegionsTestCast{
		{
			desc: "create export returns successfully",
			setup: func() {
			},
			shouldError: false,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"status": "pending"}`),
				[]byte(`{"status": "pending"}`),
				[]byte(`{"status": "ready"}`),
			},
		},
		{
			desc: "create export returns invalid json",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{Not valid JSON}`),
			},
		},
		{
			desc: "create export returns empty body",
			setup: func() {
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(``),
			},
		},
		{
			desc: "create export scanID is not set",
			setup: func() {
				tenable.scanID = ""
				tenable.fileId = ""
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"status": "ready"}`),
			},
		},
		{
			desc: "scan status returns a non-200 code",
			setup: func() {
				tenable.scanID = scanID
				tenable.fileId = fileId
			},
			shouldError: true,
			returnError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`{"status": "ready"}`),
			},
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"expectedPath": testCase.expectedPath,
			"returnError": testCase.returnError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()
		tenable := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.WaitForExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}