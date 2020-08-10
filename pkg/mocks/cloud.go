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
	 m.Called(input)
	return nil
}

func (m *MockCloudAPI) GatherIPs() error {
	fmt.Println("GatherIPs Mock")
	m.Called()
	return nil
}

func (m *MockCloudAPI) FetchIPs() []string {
	fmt.Println("GatherIPs Mock")
	m.Called()
	return m.IPs
}
