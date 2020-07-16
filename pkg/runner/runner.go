package runner

import (
	"fmt"

	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type Runner struct {
	test string
}

func Run() {
	fmt.Println("Run")
	setup()
	tenable.LaunchScan()
	tenable.CheckScanProgess()
	tenable.StartExport()
	tenable.CheckExport()
	tenable.DownloadExport()
	fmt.Println("Run Finished")
}

func setup() {
	setupBasedOnFlags()
	getIPs()
	tenable.SetupClient()
}

func launchScan() {
	fmt.Println("launchScan")
}

func waitForScanToComplete() {
	fmt.Println("waitForScanToComplete")
}

func setupBasedOnFlags() {
	fmt.Println("setupBasedOnFlags")
}

func getIPs() {
	// ips := []string
	cloud.GetGCloudIPs()
	// ips = append(ips, ^)
	cloud.GetAWSIPs()
	// ips = append(ips, ^)
}
