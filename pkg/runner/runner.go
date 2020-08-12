package runner

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/aws"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/gcloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	ec2Svc wrapper.CloudWrapper
	gcloud wrapper.CloudWrapper
	tenable wrapper.Tenable
	includeGCloud bool
	includeAWS bool
	generateReport bool
}

func (r *Runner) SetupRunner(config *config.BaseConfig) error {

	r.includeAWS = config.IncludeAWS
	r.includeGCloud = config.IncludeGCloud

	if r.includeAWS {
		ec2Svc := &aws.AWSEc2{}
		err := ec2Svc.Setup(config)
		if err != nil {
			return fmt.Errorf("SetupRunner: Error setting up AWS")
		}
		r.ec2Svc = ec2Svc
	}


	if r.includeGCloud {
		r.gcloud = &gcloud.GCloud{}
		err := r.gcloud.Setup(config)
		if err != nil {
			return fmt.Errorf("SetupRunner: Error setting up GCloud")
		}
	}

	tenableClient, err := tenable.SetupTenable(config.TenableConfig)
	if err != nil {
		return fmt.Errorf("SetupRunner: Error creating tenable client")
	}
	r.tenable = tenableClient

	r.generateReport = config.TenableConfig.GenerateReport
	return nil
}

func (r *Runner) Run() error {
	log.Debug("Run")
	err := r.getIPs()
	if err != nil {
		return fmt.Errorf("Run: Error getting ips %s", err)
	}

	log.Debug("IPs have been gathered from sources. Launching scan.")

	err = r.tenable.LaunchScan()
	if err != nil {
		return fmt.Errorf("Run: Error launching scan %s", err)
	}

	log.Debug("Scan launched. Waiting for scan to complete.")

	err = r.tenable.WaitForScanToComplete()
	if err != nil {
		return fmt.Errorf("Run: Error Waiting For Scan To Complete %s", err)
	}

	log.Debug("Scan complete.")

	if r.generateReport {
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
	}

	log.Debug("Run Finished")
	return nil
}

func (r *Runner) getIPs() error {
	log.Debug("getIPs")
	var  ips []string

	if r.includeGCloud {
		log.Debug("Gathering Google Cloud IPs")
		gcloudIPs, err := r.gcloud.GatherIPs()
		log.Debug(err)
		if err != nil {
			return fmt.Errorf("getIPs: Error retrieving GCloud IPs %s", err)
		}
		if len(gcloudIPs) == 0 {
			log.Debug("No GCloud IPs found")
		}
		ips = append(ips, gcloudIPs...)
	}

	if r.includeAWS {
		log.Debug("Gathering AWS IPs")
		awsIPs, err := r.ec2Svc.GatherIPs()
		if err != nil {
			return fmt.Errorf("getIPs: Error retrieving AWS IPs %s", err)
		}
		if len(awsIPs) == 0 {
			log.Debug("No AWS IPs found")
		}

		ips = append(ips, awsIPs...)
	}

	// targets is just for testing to make the scan go quicker
	/*
	var targets []string

	target1 := "127.0.0.1"
	targets = append(targets, target1)

	r.tenable.SetTargets(targets)
	 */

	err := r.tenable.SetTargets(ips)

	if err != nil {
		return fmt.Errorf("getIPs: Error setting Tenable targets")
	}

	log.Debug("\n\nALL:", ips)
	return nil
}
