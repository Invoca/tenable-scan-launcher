package tenable

import (
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type testCase struct {
	desc        string
	setup       func()
	shouldError bool
	expectedPath string
	returnError bool
	requestBodies [][]byte
}

func setupTenable(t *testing.T, tenable *Tenable, requestBodies [][]byte, returnError bool, expectedPath string) (*Tenable, *httptest.Server) {
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

	tenable.tenableURL = server.URL
	return tenable, server
}

func TestSetupTenable(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	tc := &config.TenableConfig{}

	testCases := []testCase{
		{
			desc: "Should fail when Aceess Key and Secret Key are not set",
			setup: func() {
				tc.AccessKey = ""
				tc.SecretKey = ""
			},
			shouldError:  true,
		},
		{
			desc: "Should fail when no severity levels are specified when creating a report",
			setup: func() {
				tc.AccessKey = "ak"
				tc.SecretKey = "sk"
				tc.GenerateReport = true
			},
			shouldError:  true,
		},
		{
			desc: "Should not fail when no required export settings are set and Generate Report is set to false",
			setup: func() {
				tc.AccessKey = "ak"
				tc.SecretKey = "sk"
				tc.GenerateReport = false
			},
			shouldError:  false,
		},
		{
			desc: "Should not fail GenerateReport is set to true and all required fields are passed to it",
			setup: func() {
				tc.AccessKey = "ak"
				tc.SecretKey = "sk"
				tc.GenerateReport = true
				tc.LowSeverity = true
			},
			shouldError:  false,
		},
	}
	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc": testCase.desc,
			"shouldError": testCase.shouldError,
			"expectedPath": testCase.expectedPath,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()

		_,  err := SetupTenable(tc)

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
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

	testCases := []testCase{
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
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)
		err := tenable.LaunchScan()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
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

	testCases := []testCase{
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
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.WaitForScanToComplete()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
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
		export: &ExportSettings{
			filter: []*Filter{
				{
					filter:  "filter",
					quality: "quality",
					value:   "value",
				},
			},
			searchType: "and",
			format: "pdf",
			chapters: "vuln_hosts_summary",
		},
		generateReport: true,
	}

	statusPath := "/scans/" + scanID + "/export"

	testCases := []testCase{
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
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.StartExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
	}
}

func TestWaitForExport(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"
	fileId    := "***REMOVED***"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		scanID:     scanID,
		scanUuid: scanUuid,
		fileId: fileId,
	}

	statusPath := "/scans/" + scanID + "/export/" + fileId + "/status"

	testCases := []testCase{
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
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.WaitForExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
	}
}

func TestDownloadExport(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	accessKey := "access"
	secretKey := "secret"
	scanID	  := "123"
	scanUuid  := "scanUuid"
	fileId    := "***REMOVED***"

	tenable := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		scanID:     scanID,
		scanUuid: scanUuid,
		fileId: fileId,
		generateReport: true,
		export: &ExportSettings{
			filePath: "./blah",
		},
		osFs: afero.NewMemMapFs(),
	}

	statusPath := "/scans/" + scanID + "/export/" + fileId + "/download"

	testCases := []testCase{
		{
			desc: "create export returns successfully",
			setup: func() {
			},
			shouldError: false,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`Binary Blob`),
			},
		},
		{
			desc: "create export returns empty body",
			setup: func() {
			},
			shouldError: false,
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
				[]byte(`Binary Blob`),
			},
		},
		{
			desc: "create export fileId is not set",
			setup: func() {
				tenable.scanID = scanID
				tenable.fileId = ""
			},
			shouldError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`Binary Blob`),
			},
		},
		{
			desc: "create export returns a non-200 code",
			setup: func() {
				tenable.scanID = scanID
				tenable.fileId = fileId
			},
			shouldError: true,
			returnError: true,
			expectedPath: statusPath,
			requestBodies: [][]byte{
				[]byte(`Binary Blob`),
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
		tenable, server := setupTenable(t, &tenable, testCase.requestBodies, testCase.returnError, testCase.expectedPath)

		err := tenable.DownloadExport()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		server.Close()
	}
}