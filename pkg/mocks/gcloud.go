package mocks
import (
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	"fmt"
)

type GCloudServiceMock interface {
	wrapper.GCloudWrapper
}

type GgCloudServiceMock struct {
	ResettableMock
}

func (g *GgCloudServiceMock) Zones() ([]string, error) {
	fmt.Println("Zones() Mock")
	args := g.Called(nil)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}

func (g *GgCloudServiceMock) InstancesIPsInRegion(region string) ([]string, error) {
	fmt.Println("InstancesIPsInRegion() Mock")
	args := g.Called(region)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}
