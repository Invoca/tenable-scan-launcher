package mocks

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
)

type GCloudServiceMock interface {
	wrapper.GCloudWrapper
}

type GCloudMock struct {
	ResettableMock
}

func (g *GCloudMock) Zones() ([]string, error) {
	fmt.Println("Zones() Mock")
	args := g.Called(nil)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}

func (g *GCloudMock) InstancesIPsInRegion(region string) ([]string, error) {
	fmt.Println("InstancesIPsInRegion() Mock")
	args := g.Called(region)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}
