package mocks

import (
	"fmt"

	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"
)

type SlackInterfaceMock struct {
	ResettableMock
}

func (s *SlackInterfaceMock) PrintAlerts(alerts tenable.Alerts) error {
	fmt.Println("Slack Mock")
	args := s.Called()
	err := args.Error(0)
	fmt.Println(err)
	return err

}
