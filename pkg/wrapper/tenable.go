package wrapper

type Tenable interface {
	SetTargets([]string) error
	LaunchScan() error
	WaitForScanToComplete() error
	StartExport() error
	WaitForExport() error
	DownloadExport() error
}
