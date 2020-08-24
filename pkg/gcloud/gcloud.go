package gcloud

import (
	"context"
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/wrapper"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"sync"
)

type GCloud struct {
	IPs            []string
	computeService wrapper.GCloudWrapper
	regions        *[]string
	mux            sync.Mutex
}

func (g *GCloud) Setup(config *config.BaseConfig) error {
	wrapper, err := createGCloudInterface(config)
	if err != nil {
		return fmt.Errorf("Setup: Error Creating GCloud Interface")
	}

	g.computeService = wrapper
	g.IPs = *new([]string)

	return nil
}

func createGCloudInterface(baseConfig *config.BaseConfig) (*GCloudWrapper, error) {
	option := option.WithCredentialsFile(baseConfig.GCloudConfig.ServiceAccountPath)

	computeService, err := compute.NewService(context.Background(), option)
	if err != nil {
		return nil, fmt.Errorf("SetupRunner: Error getting compute.Service object %s", err)
	}

	gCloudInterface, err := newCloudWrapper(computeService, baseConfig.GCloudConfig.ProjectName)
	if err != nil {
		return nil, fmt.Errorf("SetupRunner: Error creating GCloud wrapper %s", err)
	}
	return gCloudInterface, nil
}

func (g *GCloud) SetupGCloud(computeService wrapper.GCloudWrapper) error {
	if &computeService == nil {
		return fmt.Errorf("SetupGCloud: computeService cannot be nil")
	}

	g.computeService = computeService
	return nil
}

func (g *GCloud) getAllRegionsForProject() error {
	regions, err := g.computeService.Zones()

	if err != nil {
		return fmt.Errorf("getAllRegionsForProject: Error Getting Zones %s", err)
	}

	g.regions = &regions
	return nil
}

func (g *GCloud) addIPsToStruct(ips []string) {
	g.mux.Lock()
	g.IPs = append(g.IPs, ips...)
	g.mux.Unlock()
}

func (g *GCloud) getInstancesInRegion(region string) error {
	if &region == nil {
		fmt.Errorf("getInstancesInRegion: region cannot be nil")
	}

	log.Debug(region)

	privateIps, err := g.computeService.InstancesIPsInRegion(region)

	if err != nil {
		return fmt.Errorf("getInstancesInRegion: Error Instances in zone")
	}

	log.Debug("getInstancesInRegion: ", privateIps)

	g.addIPsToStruct(privateIps)
	return nil
}

// GetGCloudIPs
func (g *GCloud) GatherIPs() ([]string, error) {
	log.Debug("Getting IPs from Google Cloud")

	if &g.computeService == nil {
		return nil, fmt.Errorf("getAllRegionsForProject: computeService cannot be nil")
	}

	err := g.getAllRegionsForProject()
	if err != nil {
		return nil, fmt.Errorf("GatherIPs: Error getting regions. %s", err)
	}

	if &g.regions == nil {
		return nil, fmt.Errorf("getAllRegionsForProject: regions cannot be nil")
	}

	for _, region := range *g.regions {
		go g.getInstancesInRegion(region)
	}

	return g.IPs, nil
}

type GCloudWrapper struct {
	computeService *compute.Service
	project        string
}

func newCloudWrapper(computeService *compute.Service, project string) (*GCloudWrapper, error) {
	if computeService == nil {
		return nil, fmt.Errorf("NewCloudWrapper: computeService cannot be nil")
	}

	if &project == nil {
		return nil, fmt.Errorf("NewCloudWrapper: project cannot be nil")
	}

	return &GCloudWrapper{computeService: computeService, project: project}, nil
}

func (g *GCloudWrapper) Zones() ([]string, error) {
	regionNames := *new([]string)
	listRegions := g.computeService.Zones.List(g.project)
	regions, err := listRegions.Do()

	if regions == nil {
		return nil, fmt.Errorf("Zones: No Zones Available")
	}

	if err != nil {
		return nil, fmt.Errorf("Zones: Error Getting zones %s", err)
	}

	for _, region := range regions.Items {
		regionNames = append(regionNames, region.Name)
	}

	return regionNames, nil
}

func (g *GCloudWrapper) InstancesIPsInRegion(region string) ([]string, error) {
	if &region == nil {
		return nil, fmt.Errorf("InstancesIPsInRegion: region cannot be nil")
	}

	var privateIps []string
	fmt.Println(region)

	listInstances := g.computeService.Instances.List(g.project, region)

	resList, err := listInstances.Do()

	if err != nil {
		return nil, fmt.Errorf("InstancesIPsInRegion: Error getting instances %s", err)
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
