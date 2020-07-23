package runner

import (
	"fmt"
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type Runner struct {
	ec2Svc *ec2.EC2
	gcloud cloud.GCloud
	awsInterface cloud.EC2Ips
}

func (r *Runner) Run() {
	fmt.Println("Run")
	r.setup()
	tenable.LaunchScan()
	tenable.CheckScanProgess()
	tenable.StartExport()
	tenable.CheckExport()
	tenable.DownloadExport()
	r.getIPs()
	fmt.Println("Run Finished")
}

func (r *Runner) setup() {
	setupBasedOnFlags()
	//TODO: Use Environment variables instead of $HOME/.aws/config file.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	r.ec2Svc = ec2.New(sess)
	r.awsInterface = cloud.EC2Ips{}

	r.gcloud = cloud.GCloud{}
	r.gcloud.SetupGCloud(context.Background(), "development-156617")
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

	r.gcloud.GetGCloudIPs()

	awsStrct := cloud.EC2Ips{}
	err := awsStrct.GetAWSIPs(r.ec2Svc)
	if err != nil {
		fmt.Errorf("getIPs: %s", err)
	}
	fmt.Println("AWS: %s", awsStrct.IPs)
}
