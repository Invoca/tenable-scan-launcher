package mocks

import "github.com/Invoca/tenable-scan-launcher/pkg/tenable"

type SlackInterfaceMock struct {
	ResettableMock
}

func (s *SlackInterfaceMock) PrintAlerts(alerts tenable.Alerts) error {
	args := s.Called(nil)
	if args.Get(0) == nil {
		return args.Error(0)
	} else {
		return args.Error(0)
	}
}
