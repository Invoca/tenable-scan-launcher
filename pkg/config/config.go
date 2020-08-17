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
	SummaryReport bool
	FullReport bool

}

type GCloudConfig struct {
	ServiceAccountPath string
	ProjectName string
}

type BaseConfig struct {
	IncludeAWS bool
	IncludeGCloud bool
	HighSeverity bool
	TenableConfig *TenableConfig
	GCloudConfig *GCloudConfig
}
