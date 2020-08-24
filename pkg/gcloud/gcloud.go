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
	regions        []string
	mux            sync.Mutex
	concurrency    int
}

func (g *GCloud) Setup(config *config.BaseConfig) error {
	GCloudwrapper, err := createGCloudInterface(config)
	if err != nil {
		return fmt.Errorf("Setup: Error Creating GCloud Interface")
	}

	g.computeService = GCloudwrapper
	g.concurrency = config.GCloudConfig.Concurrency
	return nil
}

func createGCloudInterface(baseConfig *config.BaseConfig) (*GCloudWrapper, error) {
	options := option.WithCredentialsFile(baseConfig.GCloudConfig.ServiceAccountPath)

	computeService, err := compute.NewService(context.Background(), options)
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
	if computeService == nil {
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

	g.regions = regions
	return nil
}

func (g *GCloud) storeIPs(ips []string) {
	g.mux.Lock()
	g.IPs = append(g.IPs, ips...)
	g.mux.Unlock()
}

func (g *GCloud) getInstancesInRegion(region string, wg *sync.WaitGroup) error {
	defer wg.Done()

	if region == "" {
		return fmt.Errorf("getInstancesInRegion: region cannot be empty")
	}

	log.Debug(region)

	privateIps, err := g.computeService.InstancesIPsInRegion(region)

	if err != nil {
		return fmt.Errorf("getInstancesInRegion: Error Instances in zone")
	}

	log.Debug("getInstancesInRegion: ", privateIps)

	g.storeIPs(privateIps)
	return nil
}

// GetGCloudIPs
func (g *GCloud) GatherIPs() ([]string, error) {
	var wg sync.WaitGroup
	log.Debug("Getting IPs from Google Cloud")

	if g.computeService == nil {
		return nil, fmt.Errorf("getAllRegionsForProject: computeService cannot be nil")
	}

	err := g.getAllRegionsForProject()
	if err != nil {
		return nil, fmt.Errorf("GatherIPs: Error getting regions. %s", err)
	}

	if g.regions == nil {
		return nil, fmt.Errorf("getAllRegionsForProject: regions cannot be nil")
	}

	for index, region := range g.regions {
		wg.Add(1)
		go g.getInstancesInRegion(region, &wg)
		if index > 0 && (index % g.concurrency) == 0 {
			log.WithFields(log.Fields{
				"index":     index,
				"concurrency":    g.concurrency,
			}).Debug("Waiting")
			wg.Wait()
		}
	}

	wg.Wait()

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

	if project == "" {
		return nil, fmt.Errorf("NewCloudWrapper: project cannot be empty")
	}

	return &GCloudWrapper{computeService: computeService, project: project}, nil
}

func (g *GCloudWrapper) Zones() ([]string, error) {
	var regionNames []string
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
	if region == "" {
		return nil, fmt.Errorf("InstancesIPsInRegion: region cannot be nil")
	}

	var privateIps []string

	listInstances := g.computeService.Instances.List(g.project, region)

	resList, err := listInstances.Do()

	if err != nil {
		return nil, fmt.Errorf("InstancesIPsInRegion: Error getting instances %s", err)
	}

	for _, resItem := range resList.Items {
		fmt.Println(resItem.Name)
		for _, device := range resItem.NetworkInterfaces {
			privateIps = append(privateIps, device.NetworkIP)
		}
	}
	return privateIps, nil
}
