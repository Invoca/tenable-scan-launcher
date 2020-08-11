package mocks

import (
	"fmt"
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
	fmt.Println("DescribeInstances Mock")
	args := m.Called()
	fmt.Println(args)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}

func (m *MockCloudAPI) FetchIPs() []string {
	log.Debug("FetchIPs Mock")
	m.Called()
	return m.IPs
}
