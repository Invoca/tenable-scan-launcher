package runner

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type Runner struct {
	ec2Svc *ec2.EC2
}

func (r *Runner) Run() {
	fmt.Println("Run")
	r.setup()
	tenable.LaunchScan()
	tenable.CheckScanProgess()
	tenable.StartExport()
	tenable.CheckExport()
	tenable.DownloadExport()
	fmt.Println("Run Finished")
}

func (r *Runner) setup() {
	setupBasedOnFlags()
	r.ec2Svc = ec2.New(nil)
	r.getIPs()
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

func (r *Runner) getIPs() {
	//ips := []string
	cloud.GetGCloudIPs()
	// ips = append(ips, ^)
	awsStrct := cloud.EC2Ips{}
	err := awsStrct.GetAWSIPs()
	if err != nil {
		fmt.Errorf("getIPs: %s", err)
	}
}
