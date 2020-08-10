package wrapper

import "github.com/Invoca/tenable-scan-launcher/pkg/config"

type CloudWrapper interface {
	Setup(config *config.BaseConfig) error
	RetrieveIPs() error
	FetchIPs() []string
}

