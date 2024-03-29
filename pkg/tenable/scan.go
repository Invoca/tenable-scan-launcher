package tenable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Filter struct {
	filter  string
	quality string
	value   string
}

func CreateFilter(filter string, quality string, value string) (*Filter, error) {
	if filter == "" {
		return nil, fmt.Errorf("CreateSeverityFilter: filter cannot be nil")
	}
	if quality == "" {
		return nil, fmt.Errorf("CreateSeverityFilter: quality cannot be nil")
	}
	if value == "" {
		return nil, fmt.Errorf("CreateSeverityFilter: value cannot be nil")
	}

	return &Filter{
		filter:  filter,
		quality: quality,
		value:   value,
	}, nil

}

func setupSeverityFilter(low bool, medium bool, high bool, critical bool) ([]*Filter, error) {
	if (low || medium || high || critical) == false {
		return nil, fmt.Errorf("setupSeverityFilter: Cannot generate a report without specifiying severity filters")
	}

	var filters []*Filter
	filterName := "severity"
	filterQuality := "eq"

	if low {
		newFilter, err := CreateFilter(filterName, filterQuality, "Low")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Low Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if medium {
		newFilter, err := CreateFilter(filterName, filterQuality, "Medium")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Medium Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if high {
		newFilter, err := CreateFilter(filterName, filterQuality, "High")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating High Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	if critical {
		newFilter, err := CreateFilter(filterName, filterQuality, "Critical")
		if err != nil {
			return nil, fmt.Errorf("setupSeverityFilter: Error Creating Critical Severity Filter")
		}

		filters = append(filters, newFilter)
	}
	return filters, nil
}

type ExportSettings struct {
	filter     []*Filter
	chapters   string
	searchType string
	format     string
	filePath   string
}

// TODO: Find a Better Name
type Tenable struct {
	accessKey      string
	secretKey      string
	Targets        []string
	scanID         string
	fileId         string
	scanUuid       string
	tenableURL     string
	status         *scanStatus
	export         *ExportSettings
	generateReport bool
	osFs           afero.Fs
}

func (t *Tenable) SetTargets(targets []string) error {
	t.Targets = targets
	return nil
}

func SetupTenable(tenableConfig *config.TenableConfig) (*Tenable, error) {
	var filters []*Filter
	var err error

	if tenableConfig.SecretKey == "" || tenableConfig.AccessKey == "" {
		return nil, fmt.Errorf("SetupTenable: Cannot have empty secret or access keys")
	}

	es := &ExportSettings{}

	format := tenableConfig.Format
	if tenableConfig.GenerateReport {
		// supported formats  are Nessus, HTML, PDF, CSV, or DB
		if format != "nessus" && format != "html" && format != "pdf" && format != "csv" && format == "db" {
			return nil, fmt.Errorf("SetupTenable: Invalid format %s", format)
		}

		filters, err = setupSeverityFilter(tenableConfig.LowSeverity, tenableConfig.MediumSeverity, tenableConfig.HighSeverity, tenableConfig.CriticalSeverity)
		if err != nil {
			return nil, fmt.Errorf("SetupTenable: Error creating severityFilter")
		}
		chapters := tenableConfig.Chapters
		if tenableConfig.SummaryReport == true {
			chapters = "vuln_hosts_summary"
		}
		if tenableConfig.FullReport == true {
			chapters = "vuln_hosts_summary; vuln_by_host; compliance_exec; remediations; vuln_by_plugin; compliance"
		}
		es = &ExportSettings{
			filter:     filters,
			chapters:   chapters,
			searchType: tenableConfig.SearchType,
			format:     tenableConfig.Format,
			filePath:   tenableConfig.FilePath,
		}
	}
	t := &Tenable{
		accessKey:  tenableConfig.AccessKey,
		secretKey:  tenableConfig.SecretKey,
		Targets:    nil,
		scanID:     tenableConfig.ScanID,
		fileId:     "",
		scanUuid:   "",
		tenableURL: "https://cloud.tenable.com",
		status: &scanStatus{
			Pending: false,
			Running: false,
		},
		export:         es,
		generateReport: tenableConfig.GenerateReport,
		osFs:           afero.NewOsFs(),
	}
	return t, nil
}

type launchScanBody struct {
	AltTargets []string `json:"alt_targets"`
}

type scanStatus struct {
	Pending bool
	Running bool
}

func (t *Tenable) createScanRequestBody(targets []string) ([]byte, error) {
	launchScanBody := &launchScanBody{AltTargets: targets}
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
		"url":     url,
		"body":    requestBody,
		"headers": headers,
		"method":  method,
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
	log.Debug("Launching scan")

	if t.scanID == "" {
		return fmt.Errorf("LaunchScan: scanID cannot be nil")
	}

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

func (t *Tenable) WaitForScanToComplete() error {
	fmt.Println("Waiting for scan to complete")

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
# https://cloud.tenable.com/scans/111/latest-status
# {"status":"pending"}
# {"status":"running"}
# {"status":"completed"}
*/

func (t *Tenable) checkScanProgess() (string, error) {
	fmt.Println("Checking progress of the scan")

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
	if t.generateReport == false {
		return fmt.Errorf("StartExport: generateReport has been set to false. Method should not be called")
	}

	fmt.Println("Starting Export")

	if t.scanID == "" {
		return fmt.Errorf("StartExport: scanID cannot be nil")
	}

	if t.export == nil {
		return fmt.Errorf("StartExport: export cannot be nil")
	}

	headers := make(map[string]string)
	url := t.tenableURL + "/scans/" + t.scanID + "/export"

	bodyMap := make(map[string]interface{})

	for index, filter := range t.export.filter {
		bodyMap["filter."+strconv.Itoa(index)+".filter"] = filter.filter
		bodyMap["filter."+strconv.Itoa(index)+".quality"] = filter.quality
		bodyMap["filter."+strconv.Itoa(index)+".value"] = filter.value
	}

	bodyMap["filter.search_type"] = t.export.searchType
	bodyMap["format"] = t.export.format
	bodyMap["chapters"] = t.export.chapters

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

	t.fileId = strconv.FormatFloat(data["file"].(float64), 'f', -1, 64)
	return nil
}

func (t *Tenable) WaitForExport() error {
	log.Debug("Waiting for export to complete")
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
	log.Debug("checking status of export")

	if t.fileId == "" {
		return "", fmt.Errorf("checkExport: fileId cannot be nil")
	}

	if t.scanID == "" {
		return "", fmt.Errorf("checkExport: scanID cannot be nil")
	}

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
	if t.generateReport == false {
		return fmt.Errorf("DownloadExport: generateReport has been set to false. Method should not be called")
	}

	log.Debug("Downloading export")

	if t.fileId == "" {
		return fmt.Errorf("DownloadExport: fileId cannot be nil")
	}

	if t.scanID == "" {
		return fmt.Errorf("DownloadExport: scanID cannot be nil")
	}

	headers := make(map[string]string)
	url := t.tenableURL + "/scans/" + t.scanID + "/export/" + t.fileId + "/download"

	headers["accept"] = "application/octet-stream"

	body, err := t.tenableRequest(url, "GET", headers, nil)
	if err != nil {
		return fmt.Errorf("DownloadExport: Error making request %s", err)
	}

	f, err := t.osFs.Create(t.export.filePath)
	if err != nil {
		return fmt.Errorf("DownloadExport: Error creating file %s", err)
	}

	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return fmt.Errorf("DownloadExport: Error writing to file %s", err)
	}

	log.Debug("Report written to file")

	return nil
}

type Vulnerabilities struct {
	Count              int                `json:"count"`
	PluginFamily       string             `json:"plugin_family"`
	PluginID           int                `json:"plugin_id"`
	PluginName         string             `json:"plugin_name"`
	VulnerabilityState string             `json:"vulnerability_state"`
	VprState           string             `json:"vpr_state"`
	VprScore           float32            `json:"vpr_score"`
	AcceptedCount      int                `json:"accepted_count"`
	RecastedCount      int                `json:"recasted_count"`
	CountsBySeverity   []CountsBySeverity `json:"counts_by_severity"`
}
type CountsBySeverity struct {
	Count int `json:"count"`
	Value int `json:"value"`
}
type Alerts struct {
	Vulnerabilities         []Vulnerabilities `json:"vulnerabilities"`
	TotalVulnerabilityCount int               `json:"total_vulnerability_count"`
	TotalAssetCount         int               `json:"total_asset_count"`
}

func (t *Tenable) GetVulnerabilities() (*Alerts, error) {
	url := "https://cloud.tenable.com/workbenches/vulnerabilities"

	// TODO This should just use tenableRequest()
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKeys", "accessKey="+t.accessKey+"; secretKey="+t.secretKey+";")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetVulnerabilities: Error performing request. %s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error pulling vulnerabilities from Tenable %s", err)
	}
	alerts := Alerts{}
	err = json.Unmarshal([]byte(body), &alerts)

	if err != nil {
		return nil, fmt.Errorf("Error Unmarshalling Json to Alerts object %s", err)
	}
	fmt.Println(alerts.TotalAssetCount)
	fmt.Println(alerts.TotalVulnerabilityCount)
	return &alerts, nil
}
