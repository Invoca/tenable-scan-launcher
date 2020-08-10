package mocks

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
)

type MockCloudAPI struct {
	ResettableMock
	IPs []string
}

func (m *MockCloudAPI) Setup(input *config.BaseConfig) error {
	fmt.Println("Setup Mock")
	args := m.Called(input)
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockCloudAPI) GatherIPs() error {
	fmt.Println("GatherIPs Mock")
	args := m.Called()
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockCloudAPI) FetchIPs() []string {
	fmt.Println("GatherIPs Mock")
	m.Called()
	return m.IPs
}
