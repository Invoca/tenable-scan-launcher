package mocks

import (
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	log "github.com/sirupsen/logrus"
)

type MockCloudAPI struct {
	ResettableMock
	IPs []string
}

func (m *MockCloudAPI) Setup(input *config.BaseConfig) error {
	log.Debug("Setup Mock")
	args := m.Called(input)
	log.Debug(args.Error(0))
	return args.Error(0)
}

func (m *MockCloudAPI) GatherIPs() ([]string, error) {
	log.Debug("GatherIPs Mock")
	args := m.Called()
	log.Debug(args)
	return m.IPs, args.Error(1)
}
