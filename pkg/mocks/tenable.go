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
	if args.Get(0) == nil {
		return fmt.Errorf("input nil")
	}
	return nil
}

func (m *MockTenableAPI) LaunchScan() error {
	fmt.Println("LaunchScan Mock")
	m.Called()
	return nil
}

func (m *MockTenableAPI) StartExport() error {
	fmt.Println("StartExport Mock")
	m.Called()
	return nil
}

func (m *MockTenableAPI) WaitForExport() error {
	fmt.Println("WaitForExport Mock")
	m.Called()
	return nil
}

func (m *MockTenableAPI) DownloadExport() error {
	fmt.Println("DownloadExport Mock")
	m.Called()
	return nil
}

func (m *MockTenableAPI) WaitForScanToComplete() error {
	fmt.Println("WaitForScanToComplete Mock")
	m.Called()
	return nil
}

