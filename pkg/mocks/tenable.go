package mocks

import (
	"fmt"
)

type MockTenableAPI struct {
	ResettableMock
}

func (m *MockTenableAPI) SetTargets(input []string) error {
	fmt.Println("Setup Mock")
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockTenableAPI) LaunchScan() error {
	fmt.Println("LaunchScan Mock")
	args := m.Called()
	return args.Error(0)
}

func (m *MockTenableAPI) StartExport() error {
	fmt.Println("StartExport Mock")
	args := m.Called()
	return args.Error(0)
}

func (m *MockTenableAPI) WaitForExport() error {
	fmt.Println("WaitForExport Mock")
	args := m.Called()
	return args.Error(0)
}

func (m *MockTenableAPI) DownloadExport() error {
	fmt.Println("DownloadExport Mock")
	args := m.Called()
	return args.Error(0)
}

func (m *MockTenableAPI) WaitForScanToComplete() error {
	fmt.Println("WaitForScanToComplete Mock")
	args := m.Called()
	return args.Error(0)
}

