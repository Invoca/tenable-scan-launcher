package wrapper

import (
	"fmt"
	"google.golang.org/api/compute/v1"
)

type GCloudWrapper interface {
	Zones() ([]string, error)
	InstancesIPsInRegion(string) ([]string, error)
}

type gCloudWrapper struct{
	computeService 	*compute.Service
	project			string
}


func NewCloudWrapper(computeService *compute.Service, project string) (*gCloudWrapper, error) {
	if computeService == nil {
		fmt.Errorf("NewgCloudWrapper: computeService cannot be nil")
	}

	if &project == nil {
		fmt.Errorf("NewgCloudWrapper: project cannot be nil")
	}

	return &gCloudWrapper{computeService: computeService, project: project}, nil
}

func (g *gCloudWrapper) Zones() ([]string, error) {
	regionNames := *new([]string)
	listRegions := g.computeService.Zones.List(g.project)
	regions, err := listRegions.Do()

	if regions == nil {
		fmt.Errorf("GetRegions: No regions available")
	}

	if err != nil {
		fmt.Errorf("GetRegions: Error Getting Regions %s", err)
	}

	for _, region := range regions.Items {
		regionNames = append(regionNames,region.Name)
	}

	return regionNames, nil
}

func (g *gCloudWrapper) InstancesIPsInRegion(region string) ([]string, error) {
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
	return privateIps, nil
}