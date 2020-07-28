package cloud

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	"google.golang.org/api/compute/v1"
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
		return fmt.Errorf("getAllRegionsForProject: Error Getting Zones")
	}

	g.regions = regions
	return nil
}

func (g *GCloud) getInstancesInRegion(region string) error {
	if &region == nil {
		fmt.Errorf("getInstancesInRegion: region cannot be nil")
	}

	fmt.Println(region)

	privateIps, err := g.computeService.InstancesIPsInRegion(region)

	if err != nil {
		return fmt.Errorf("getInstancesInRegion: Error Instances in zone")
	}

	fmt.Print("getInstancesInRegion: ", privateIps)

	g.mux.Lock()
	g.IPs = append(g.IPs, privateIps...)
	g.mux.Unlock()
	return nil
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
		fmt.Errorf("InstancesIPsInRegion: region cannot be nil")
	}

	var privateIps []string
	fmt.Println(region)

	listInstances := g.computeService.Instances.List(g.project, region)

	resList, err := listInstances.Do()

	if err != nil {
		fmt.Errorf("InstancesIPsInRegion: Error getting instances %s", err)
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
