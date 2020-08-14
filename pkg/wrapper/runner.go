package wrapper

import "github.com/Invoca/tenable-scan-launcher/pkg/config"

type Runner interface {
	SetupRunner(config *config.BaseConfig) error
	Run() error
}
