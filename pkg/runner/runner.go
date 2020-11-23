package runner

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/aws"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/gcloud"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"time"
)

type Runner struct {
	awsSvc         wrapper.CloudWrapper
	gcloud         wrapper.CloudWrapper
	tenable        wrapper.Tenable
	includeGCloud  bool
	includeAWS     bool
	generateReport bool
	fileLocation   string
	storeName      string
	fileName       string
	osFs           afero.Fs
}

func (r *Runner) SetupRunner(config *config.BaseConfig) error {

	r.includeAWS = config.IncludeAWS
	r.includeGCloud = config.IncludeGCloud
	r.osFs = afero.NewOsFs()

	if config.TenableConfig == nil {
		return fmt.Errorf("SetupRunner: TenableConfig in config cannot be nil")
	}

	if r.includeAWS {
		ec2Svc := &aws.AwsSvc{}
		err := ec2Svc.Setup(config)
		if err != nil {
			return fmt.Errorf("SetupRunner: Error setting up AWS")
		}
		r.awsSvc = ec2Svc
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
	r.fileLocation = config.TenableConfig.FilePath
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

		data, err := r.osFs.Open(r.fileLocation)
		if err != nil {
			return fmt.Errorf("Run: Error opening file %s	", err)
		}

		var fileData []byte

		_, err = data.Read(fileData)
		if err != nil {
			return fmt.Errorf("Run: Error reading file %s", err)
		}

		err = r.awsSvc.UploadFile(r.storeName, r.fileName+"-"+time.Now().String()+".pdf", fileData)
		if err != nil {
			return fmt.Errorf("Run: Error uploading object %s", err)
		}
	}

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
		awsIPs, err := r.awsSvc.GatherIPs()
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
