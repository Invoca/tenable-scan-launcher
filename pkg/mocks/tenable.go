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
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockTenableAPI) LaunchScan() error {
	fmt.Println("LaunchScan Mock")
	args := m.Called()
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockTenableAPI) StartExport() error {
	fmt.Println("StartExport Mock")
	args := m.Called()
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockTenableAPI) WaitForExport() error {
	fmt.Println("WaitForExport Mock")
	args := m.Called()
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockTenableAPI) DownloadExport() error {
	fmt.Println("DownloadExport Mock")
	args := m.Called()
	if args.Get(0) == nil {
		return args.Error(1)
	} else {
		return args.Error(1)
	}
}

func (m *MockTenableAPI) WaitForScanToComplete() error {
	fmt.Println("WaitForScanToComplete Mock")
	args := m.Called()
	if args.Get(0) == nil {
	return args.Error(1)
	} else {
		return args.Error(1)
	}
}

