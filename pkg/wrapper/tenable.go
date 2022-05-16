package wrapper

import (
	t "github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type Tenable interface {
	SetTargets([]string) error
	LaunchScan() error
//	WaitForScanToComplete() error
//	StartExport() error
//	WaitForExport() error
//	DownloadExport() error
	GetVulnerabilities() (*t.Alerts, error)
}
