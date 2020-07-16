package tenable

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SetupClient() {
	fmt.Println("SetupClient")
}

func LaunchScan() {
	fmt.Println("LaunchScan")
}

func CheckScanProgess() {
	fmt.Println("CheckScanProgess")
}

func StartExport() {
	fmt.Println("StartExport")
}

func CheckExport() {
	fmt.Println("CheckExport")
}

func DownloadExport() {
	fmt.Println("DownloadExport")
}

// Ignore functions below. They are copy and pasted from the Tenable docs.

func launchScan() {
	url := "https://cloud.tenable.com/scans/scan_id/launch"

	payload := strings.NewReader("{\"alt_targets\":[\"127.0.0.1\"]}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	// return scan_uuid of body
	//{
	//"scan_uuid":"e7f6c3f2-1718-4451-b459-1e8aa2ec6cdf"
	//}
}

func checkScanStatus() {
	url := "https://cloud.tenable.com/scans/scan_id/latest-status"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func startExportReport() {
	url := "https://cloud.tenable.com/scans/scan_id/export"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func downloadScanReport() {
	url := "https://cloud.tenable.com/scans/scan_id/export/file_id/download"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/octet-stream")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
