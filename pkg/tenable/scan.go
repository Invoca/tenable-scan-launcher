package tenable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// TODO: Find a Better Name
type Tenable struct {
	accessKey 	string
	secretKey 	string
	Targets		[]string
	scanID 		string
	fileId		string
	scanUuid 	string
	tenableURL 	string
	status		*scanStatus
}

func SetupClient(accessKey string, secretKey string, scanID string) *Tenable {
	t := Tenable{
		accessKey:  accessKey,
		secretKey:  secretKey,
		Targets:    nil,
		scanID:     scanID,
		tenableURL: "https://cloud.tenable.com",
		status: &scanStatus{
			Pending:   false,
			Running:   false,
		},
	}
	return &t
}

type launchScanBody struct {
	altTargets	[]string `json:"alt_targets"`
}

type scanStatus struct {
	Pending		bool
	Running 	bool
}

func (t *Tenable) createScanRequestBody(targets []string) ([]byte, error) {
	launchScanBody := &launchScanBody{altTargets: targets}
	res, err := json.Marshal(launchScanBody)
	if err != nil {
		return nil, fmt.Errorf("createRequest: Error in json.Marshal(): %s", err)
	}
	return res, nil
}

func (t *Tenable) tenableRequest(url string, method string, headers map[string]string, requestBody io.Reader) ([]byte, error) {
	successCode := 200

	if url == "" {
		return nil, fmt.Errorf("tenableRequest: url cannot be nil")
	} else if method == "" {
		return nil, fmt.Errorf("tenableRequest: method cannot be nil")
	}

	if requestBody == nil {
		log.Debug("requestBody is nil")
	}

	if headers == nil {
		log.Debug("headers nil. Creating a new map[string]string")
		headers = make(map[string]string)
	}

	req, err := http.NewRequest(method, url, requestBody)

	if err != nil {
		return nil, fmt.Errorf("tenableRequest: Error creating export request. %s", err)
	}

	apikeyString := "accessKey=" + t.accessKey + "; secretKey=" + t.secretKey + ";"

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	req.Header.Add("X-ApiKeys", apikeyString)

	log.WithFields(log.Fields{
		"url": url,
		"body": requestBody,
		"headers": headers,
		"method": method,
	}).Debug("HTTP Request created")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("checkScanProgess: Error making request %s", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("checkScanProgess: Error reading body of request %s", err)
	}

	// We do this so we do not accidentally log any binary information
	if res.Header.Get("Content-Type") == "application/octet-stream" {
		log.WithFields(log.Fields{
			"code": res.StatusCode,
			"body": "Binary Length: " + res.Header.Get("Content-Length"),
		}).Debug("Retrieved Octet Stream Response")
	} else {
		log.WithFields(log.Fields{
			"code": res.StatusCode,
			"body": string(body),
		}).Debug("HTTP Response Received")
	}

	if res.StatusCode != successCode {
		return nil, fmt.Errorf("checkScanProgess: Recieved a response code that is not 200. Recieved: %d", res.StatusCode)
	}

	if err != nil {
		return nil, fmt.Errorf("checkScanProgess: Error reading body. %s", err)
	}

	return body, nil
}


func (t *Tenable) LaunchScan() error {
	log.Debug("LaunchScan")

	url := t.tenableURL + "/scans/" + t.scanID + "/launch"
	headers := make(map[string]string)

	marshalledPayload, err := t.createScanRequestBody(t.Targets)
	log.Debug(marshalledPayload)

	if err != nil {
		return fmt.Errorf("LaunchScan(): Error creating json from targets. %s", err)
	}

	payloadBuffer := bytes.NewBuffer(marshalledPayload)

	headers["accept"] = "application/json"
	headers["content-type"] = "application/json"

	body, err := t.tenableRequest(url, "POST", headers, payloadBuffer)

	if err != nil {
		return fmt.Errorf("LaunchScan: Error making request %s", err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("checkScanProgess(): Error unmarshalling json body: %s", string(body))
	}

	uuid := data["scan_uuid"].(string)


	t.scanUuid = uuid
	return nil
}

func (t *Tenable)  WaitForScanToComplete() error {
	fmt.Println("WaitForScanToComplete")

	if t.scanID == "" {
		return fmt.Errorf("waitForScanToComplete: scanID cannot be nil")
	}

	for {
		currentStatus, err := t.checkScanProgess()
		if err != nil {
			return fmt.Errorf("waitForScanToComplete: Error getting scan progress %s", err)
		}

		if currentStatus == "pending" && t.status.Pending == false {
			t.status.Pending = true
			log.Debug("Entered Pending State")
		} else if currentStatus == "running" && t.status.Running == false {
			t.status.Running = true
			log.Debug("Entered Running State")
		} else if currentStatus == "completed" {
			log.Debug("Scan has been completed")
			return nil
		}
		log.Debug("Sleeping For 5 Seconds and trying again")
		time.Sleep(5)
	}
}

/*
# https://cloud.tenable.com/scans/***REMOVED***/latest-status
# {"status":"pending"}
# {"status":"running"}
# {"status":"completed"}
 */

func (t *Tenable) checkScanProgess() (string, error) {
	fmt.Println("checkScanProgess")

	if t.scanID == "" {
		return "", fmt.Errorf("checkScanProgess: scanID cannot be nil")
	}

	url := t.tenableURL + "/scans/" + t.scanID + "/latest-status"

	body, err := t.tenableRequest(url, "GET", nil, nil)
	if err != nil {
		return "", fmt.Errorf("checkScanProgess: Error making request %s", err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("checkScanProgess: Error unmarshalling json body: %s", string(body))
	}

	status := data["status"].(string)
	return status, nil
}

func (t *Tenable) StartExport() error {
	fmt.Println("StartExport")

	headers := make(map[string]string)
	url := t.tenableURL + "/scans/" + t.scanID + "/export"

	//TODO: Create Better map than this.
	bodyMap := make(map[string]interface{})

	bodyMap["filter.0.filter"] 		= "severity"
	bodyMap["filter.0.quality"] 	= "eq"
	bodyMap["filter.0.value"] 		= "Critical"
	bodyMap["filter.1.filter"] 		= "severity"
	bodyMap["filter.1.quality"] 	= "eq"
	bodyMap["filter.1.value"] 		= "High"
	bodyMap["filter.2.filter"] 		= "severity"
	bodyMap["filter.2.quality"] 	= "eq"
	bodyMap["filter.2.value"] 		= "Medium"
	bodyMap["filter.3.filter"] 		= "severity"
	bodyMap["filter.3.quality"] 	= "eq"
	bodyMap["filter.3.value"] 		= "Low"
	bodyMap["filter.search_type"] 	= "or"

	bodyMap["format"] = "pdf"
	bodyMap["chapters"] = "vuln_hosts_summary; vuln_by_host; compliance_exec; remediations; vuln_by_plugin; compliance"

	reqBody, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("StartExport: Error Marshalling Body")
	}
	log.Debug(string(reqBody))

	bodyBuffer := bytes.NewBuffer(reqBody)


	headers["accept"] = "application/json"
	headers["content-type"] = "application/json"

	body, err := t.tenableRequest(url, "POST", headers, bodyBuffer)

	if err != nil {
		return fmt.Errorf("StartExport: Error making request %s", err)
	}


	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("StartExport: Error unmarshalling json body: %s", string(body))
	}

	t.fileId =  strconv.FormatFloat(data["file"].(float64), 'f', -1, 64)
	return nil
}

func (t *Tenable) WaitForExport() error {
	log.Debug("WaitForExport")
	if t.fileId == "" {
		return fmt.Errorf("WaitForExport: fileId cannot be nil")
	}

	if t.scanID == "" {
		return fmt.Errorf("WaitForExport: scanID cannot be nil")
	}

	for {
		status, err := t.checkExport()

		if err != nil {
			return fmt.Errorf("WaitForExport(): Error waiting for report %s", err)
		}

		if status == "ready" {
			return nil
		}
		log.Debug("File not ready for export. Sleeping 5 seconds and checking again")
		time.Sleep(5)
	}
}

func (t *Tenable) checkExport() (string, error) {
	log.Debug("checkExport")

	headers := make(map[string]string)
	url := t.tenableURL + "/scans/" + t.scanID + "/export/" + t.fileId + "/status"


	headers["accept"] = "application/json"
	headers["content-type"] = "application/json"


	body, err := t.tenableRequest(url, "GET", headers, nil)

	if err != nil {
		return "", fmt.Errorf("checkExport: Error making request %s", err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("checkScanProgess: Error unmarshalling json body: %s", string(body))
	}

	status := data["status"].(string)
	log.Debug(status)

	return status, nil
}

func (t *Tenable) DownloadExport() error {
	log.Debug("DownloadExport")

	headers := make(map[string]string)
	url := t.tenableURL + "/scans/" + t.scanID + "/export/" + t.fileId + "/download"

	headers["accept"] = "application/octet-stream"


	body, err := t.tenableRequest(url, "GET", headers, nil)

	err = ioutil.WriteFile("./temp_result.pdf", body, 0777)
	if err != nil {
		return fmt.Errorf("checkScanProgess: Writing to file %s", err)
	}
	log.Debug("Completed Writing to file")


	return nil

}
