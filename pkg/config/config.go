package config

type TenableConfig struct {
	LowSeverity bool
	MediumSeverity bool
	HighSeverity bool
	CriticalSeverity bool
	GenerateReport bool
	SearchType string
	Format string
	Chapters string
	FilePath string
	AccessKey string
	SecretKey string
	ScanID string
}

type GCloudConfig struct {
	ServiceAccountPath string
	ProjectName string
}

type RunnerConfig struct {
	IncludeAWS bool
	IncludeGCloud bool
	HighSeverity bool
	TenableConfig *TenableConfig
	GCloudConfig *GCloudConfig
}
