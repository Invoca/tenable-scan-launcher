package cloud

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
)


type GCloud struct {
	IPs		[]string
	ctx 	context.Context
	computeService compute.Service
	project string
	regions []*compute.Zone
}


func (g *GCloud) SetupGCloud(ctx context.Context, project string) {
	//TODO: Setup GCloud SDK to use json from Service Account
	computeService, err := compute.NewService(ctx)

	if err != nil {
		fmt.Errorf("SetupGCloud: Error getting compute.Service object %s", err)
	}

	g.computeService = *computeService
	g.project = project
	g.getAllRegionsForProject()
}


func (g *GCloud) getAllRegionsForProject() {
	if g.ctx == nil {
		fmt.Println("err")
	}

	//TODO: FIX!
	//if g.computeService != nil {
	//	fmt.Println("err")
	//}

	listRegions := g.computeService.Zones.List(g.project)
	regions, err := listRegions.Do()

	if regions == nil {
		fmt.Errorf("GetRegions: No regions available")
	}

	if err != nil {
		fmt.Errorf("GetRegions: Error Getting Regions %s", err)
	}
	fmt.Println(regions)

	g.regions = regions.Items


}

func (g *GCloud) getInstancesInRegion(region compute.Zone) {
	regionName := region.Name
	fmt.Println(regionName)

	listInstances := g.computeService.Instances.List(g.project, regionName)

	resList, err := listInstances.Do()

	if err != nil {
		fmt.Errorf("getInstancesInRegion: Error getting instances %s", err)
	}

	for _, resItem := range resList.Items {
		fmt.Println(resItem.Name)
		for _, device := range resItem.NetworkInterfaces {
			fmt.Println(device.NetworkIP)
		}
	}
}

// GetGCloudIPs
func (g *GCloud) GetGCloudIPs() {
	fmt.Println("Getting IPs from Google Cloud")

	if g.ctx == nil {
		fmt.Println("err")
	}

	//TODO: FIX!
	//if g.computeService != nil {
	//	fmt.Println("err")
	//}


	for _, region := range g.regions {
		g.getInstancesInRegion(*region)
	}
}
