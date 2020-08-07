package runner

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/cloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	ec2Svc *ec2.EC2
	gcloud cloud.GCloud
	awsInterface cloud.EC2Ips
	tenable *tenable.Tenable
	includeGCloud bool
	includeAWS bool
}

func SetupRunner (tenableClient *tenable.Tenable, gCloudInterface *cloud.GCloudWrapper, ec2Interface *ec2.EC2, includeGCloud bool, includeAWS bool) (*Runner, error) {
	r := &Runner{}

	r.includeAWS = includeAWS
	r.includeGCloud = includeGCloud

	if includeAWS {
		r.ec2Svc = ec2Interface
		r.awsInterface = cloud.EC2Ips{}
	}


	if includeGCloud {
		r.gcloud = cloud.GCloud{}
		r.gcloud.SetupGCloud(gCloudInterface)
	}

	r.tenable = tenableClient
	return r, nil
}

func (r *Runner) Run() error {
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

func (r *Runner) getIPs() error {
	log.Debug("getIPs")
	var ips []string
	var err error

	if r.includeGCloud {
		err = r.gcloud.GetGCloudIPs()
		if err != nil {
			return fmt.Errorf("getIPs: Error fetching GCloud IPs %s", err)
		}
		if len(r.gcloud.IPs) == 0 {
			log.Debug("No GCloud IPs found")
		}
		ips = append(ips, r.gcloud.IPs...)
	}

	if r.includeAWS {
		awsSvc := cloud.EC2Ips{}
		err = awsSvc.GetAWSIPs(r.ec2Svc)
		if err != nil {
			return fmt.Errorf("getIPs: Error fetching AWS IPs %s", err)
		}

		ips = append(ips, awsSvc.IPs...)
	}

	r.tenable.Targets = ips

	log.Debug("\n\nALL:", ips)
	return nil
}
