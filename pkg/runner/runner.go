package runner

import (
	"fmt"

	"github.com/Invoca/tenable-scan-launcher/pkg/aws"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/gcloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/slack"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"

	log "github.com/sirupsen/logrus"
)

type Runner struct {
	ec2Svc         wrapper.CloudWrapper
	gcloud         wrapper.CloudWrapper
	slackSvc       wrapper.SlackSvc
	tenable        wrapper.Tenable
	includeGCloud  bool
	includeAWS     bool
	generateReport bool
}

func (r *Runner) SetupRunner(config *config.BaseConfig) error {

	var err error
	r.includeAWS = config.IncludeAWS
	r.includeGCloud = config.IncludeGCloud

	if config.TenableConfig == nil {
		return fmt.Errorf("SetupRunner: TenableConfig in config cannot be nil")
	}

	if config.SlackConfig == nil {
		return fmt.Errorf("SetupRunner: SlackConfig in config cannot be nil ")
	}
	r.slackSvc, err = slack.New(*config)

	if r.includeAWS {
		ec2Svc := &aws.AWSEc2{}
		err := ec2Svc.Setup(config)
		if err != nil {
			return fmt.Errorf("SetupRunner: Error setting up AWS")
		}
		r.ec2Svc = ec2Svc
	}

	if r.includeGCloud {
		if config.GCloudConfig == nil {
			return fmt.Errorf("SetupRunner: GCloudConfig in config cannot be nil")
		}
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

	/**
	err = r.tenable.WaitForScanToComplete()
	if err != nil {
		return fmt.Errorf("Run: Error Waiting For Scan To Complete %s", err)
	}

	log.Debug("Scan complete.")

	log.Debug("Fetching Critical alerts from Tenable Dashboard")
	alerts, err := r.tenable.GetVulnerabilities()

	if err != nil {
		return fmt.Errorf("Run: Error Fetching alerts from Tenable Dashboard %s", err)
	}
	if alerts.TotalVulnerabilityCount > 0 {
		log.Debug("Posting to Slack")
		err = r.slackSvc.PrintAlerts(*alerts)
		if err != nil {
			return fmt.Errorf("Run: Error posting to slack %s", err)
		}
	}

	if r.generateReport {
		err = r.tenable.StartExport()
		if err != nil {
			return fmt.Errorf("Run: Error Starting Scan %s", err)
		}

		log.Debug("Export Started. Waiting for file to be ready.")

		err = r.tenable.WaitForExport()
		if err != nil {
			return fmt.Errorf("Run: Error Waiting For Export %s", err)
		}

		log.Debug("Starting file download")

		err = r.tenable.DownloadExport()
		if err != nil {
			return fmt.Errorf("Run: Error Downloading Export %s", err)
		}
		log.Debug("File successfully downloaded")
	}
	**/
	log.Debug("Run Finished")

	return nil
}

func (r *Runner) getIPs() error {
	log.Debug("getIPs")
	var ips []string

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

	err := r.tenable.SetTargets(ips)

	if err != nil {
		return fmt.Errorf("getIPs: Error setting Tenable targets")
	}

	log.Debug("\n\nALL:", ips)
	return nil
}
