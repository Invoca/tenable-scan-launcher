package cloud

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	"sync"
)


type GCloud struct {
	IPs		[]string
	computeService wrapper.GCloudWrapper
	project string
	regions []string
	mux sync.Mutex
}


func (g *GCloud) SetupGCloud(computeService wrapper.GCloudWrapper, project string) {
	if &computeService == nil {
		fmt.Errorf("SetupGCloud: computeService cannot be nil")
	}

	if &project == nil {
		fmt.Errorf("SetupGCloud: project cannot be nil")
	}

	g.computeService = computeService
	g.project = project
}


func (g *GCloud) getAllRegionsForProject() error {
	regions, err := g.computeService.Zones()

	if err != nil {
		fmt.Errorf("getAllRegionsForProject: Error Getting Zones")
	}

	g.regions = regions
	return nil
}

func (g *GCloud) getInstancesInRegion(region string) {
	if &region == nil {
		fmt.Errorf("getInstancesInRegion: region cannot be nil")
	}

	fmt.Println(region)

	for _, region := range g.regions {
		privateIps, err := g.computeService.InstancesIPsInRegion(region)

		if err != nil {
			fmt.Errorf("getInstancesInRegion: Error Instances in zone")
		}

		g.mux.Lock()
		g.IPs = append(g.IPs, privateIps...)
		g.mux.Unlock()

	}
}

// GetGCloudIPs
func (g *GCloud) GetGCloudIPs() {
	fmt.Println("Getting IPs from Google Cloud")

	if &g.computeService == nil {
		fmt.Errorf("getAllRegionsForProject: computeService cannot be nil")
	}

	if &g.project == nil {
		fmt.Errorf("getAllRegionsForProject: project cannot be nil")
	}

	g.getAllRegionsForProject()

	if &g.regions == nil {
		fmt.Errorf("getAllRegionsForProject: regions cannot be nil")
	}

	for _, region := range g.regions {
		go g.getInstancesInRegion(region)
	}
}
