package slack

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type slackTestCase struct {
	desc        string
	setup       func() *httptest.Server
	shouldError bool
}

func PrintAlerts(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	slackInterface := slack{}
	slackInterface.rateLimit = &rateLimitedHTTPClient{
		client:   http.DefaultClient,
		rlClient: rate.NewLimiter(rate.Every(10*time.Second), 10),
	}

	testCases := []slackTestCase{
		{
			desc: "Able to post to slack",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
			shouldError: false,
		},
		{
			desc: "Error posting to slack",
			setup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
			},
			shouldError: true,
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc":        testCase.desc,
			"shouldError": testCase.shouldError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testServer := testCase.setup()
		slackInterface.slackUrl = testServer.URL

		err := slackInterface.PrintAlerts(tenable.Alerts{})

		testServer.Close()

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
