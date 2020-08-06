package runner

import (
	"context"
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

type runner struct {
	ec2Svc *ec2.EC2
	gcloud cloud.GCloud
	awsInterface cloud.EC2Ips
	tenable *tenable.Tenable
}

/*
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
*/

func SetupRunner (tenableClient *tenable.Tenable) (*runner, error) {
	r := &runner{}

	//TODO: Use Environment variables instead of $HOME/.aws/config file.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	r.ec2Svc = ec2.New(sess)
	r.awsInterface = cloud.EC2Ips{}

	//TODO: Setup GCloud SDK to use json from Service Account
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("SetupRunner: Error getting compute.Service object %s", err)
	}


	gCloudInterface, err := cloud.NewCloudWrapper(computeService, "***REMOVED***")
	if err != nil {
		return nil, fmt.Errorf("SetupRunner: Error creating GCloud wrapper %s", err)
	}


	r.gcloud = cloud.GCloud{}
	r.gcloud.SetupGCloud(gCloudInterface)

	r.tenable = tenableClient
	return r, nil
}

func (r *runner) Run() error {
	log.Debug("Run")
	err := r.getIPs()
	if err != nil {
		return fmt.Errorf("Run: Error getting ips %s", err)
	}

	if len(r.tenable.Targets) == 0 {
		return fmt.Errorf("Run: No targets added to scan")
	}

	// targets is just for testing to make the scan go quicker
	var targets []string

	target1 := "127.0.0.1"
	targets = append(targets, target1)
	r.tenable.Targets = targets

	err = r.tenable.LaunchScan()
	if err != nil {
		return fmt.Errorf("Run: Error launching scan %s", err)
	}

	err = r.tenable.WaitForScanToComplete()
	if err != nil {
		return fmt.Errorf("Run: Error Waiting For Scan To Complete %s", err)
	}

	err = r.tenable.StartExport()
	if err != nil {
		return fmt.Errorf("Run: Error Starting Scan %s", err)
	}

	err = r.tenable.WaitForExport()
	if err != nil {
		return fmt.Errorf("Run: Error Waiting For Export %s", err)
	}

	err = r.tenable.DownloadExport()
	if err != nil {
		return fmt.Errorf("Run: Error Downloading Export %s", err)
	}

	log.Debug("Run Finished")
	return nil
}

func (r *runner) getIPs() error {
	log.Debug("getIPs")
	var ips []string

	err := r.gcloud.GetGCloudIPs()
	if err != nil {
		return fmt.Errorf("getIPs: Error fetching GCloud IPs %s", err)
	}

	if len(r.gcloud.IPs) == 0 {
		log.Debug("No GCloud IPs found")
	}

	ips = append(ips, r.gcloud.IPs...)

	awsSvc := cloud.EC2Ips{}
	err = awsSvc.GetAWSIPs(r.ec2Svc)
	if err != nil {
		return fmt.Errorf("getIPs: Error fetching AWS IPs %s", err)
	}

	ips = append(ips, awsSvc.IPs...)
	r.tenable.Targets = ips

	log.Debug("\n\nALL:", ips)
	return nil
}
