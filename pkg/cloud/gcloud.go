package cloud

import (
	"fmt"
	"google.golang.org/api/compute/v1"
	"sync"
)


type GCloud struct {
	IPs		[]string
	computeService compute.Service
	project string
	regions []*compute.Zone
	mux sync.Mutex
}


func (g *GCloud) SetupGCloud(computeService compute.Service, project string) {
	if &computeService == nil {
		fmt.Errorf("SetupGCloud: computeService cannot be nil")
	}

	if &project == nil {
		fmt.Errorf("SetupGCloud: project cannot be nil")
	}

	g.computeService = computeService
	g.project = project
}


func (g *GCloud) getAllRegionsForProject() {
	if &g.computeService == nil {
		fmt.Errorf("getAllRegionsForProject: computeService cannot be nil")
	}

	listRegions := g.computeService.Zones.List(g.project)
	regions, err := listRegions.Do()

	if regions == nil {
		fmt.Errorf("GetRegions: No regions available")
	}

	if err != nil {
		fmt.Errorf("GetRegions: Error Getting Regions %s", err)
	}

	g.regions = regions.Items
}

func (g *GCloud) getInstancesInRegion(region string) {
	if &region == nil {
		fmt.Errorf("getInstancesInRegion: region cannot be nil")
	}

	var privateIps []string
	fmt.Println(region)

	listInstances := g.computeService.Instances.List(g.project, region)

	resList, err := listInstances.Do()

	if err != nil {
		fmt.Errorf("getInstancesInRegion: Error getting instances %s", err)
	}

	for _, resItem := range resList.Items {
		fmt.Println(resItem.Name)
		for _, device := range resItem.NetworkInterfaces {
			deviceIP := device.NetworkIP
			fmt.Println(device.NetworkIP)
			privateIps = append(privateIps, deviceIP)
		}
	}

	g.mux.Lock()
	g.IPs = append(g.IPs, privateIps...)
	g.mux.Unlock()
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
		go g.getInstancesInRegion(region.Name)
	}
}
