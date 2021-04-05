package wrapper

import (
	t "github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type SlackSvc interface {
	PrintAlerts(t.Alerts) error
}
