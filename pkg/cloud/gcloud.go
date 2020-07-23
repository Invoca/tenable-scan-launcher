package cloud

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
)


type GCloudIps struct {
	IPs		[]string
	api		compute.Instance

}

// GetGCloudIPs
func GetGCloudIPs() {
	fmt.Println("Getting IPs from Google Cloud")

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)

	if err != nil {
		fmt.Println(err)
	}

	//instances := computeService.Instances.Get("development-156617", "*", "*")


	listRegions := computeService.Zones.List("development-156617")
	regions, err := listRegions.Do()

	if regions == nil {
		fmt.Println("No regions available")
	}

	fmt.Println(regions)


	for _, region := range regions.Items {

		fmt.Println(region)

		regionName := region.Name
		fmt.Println(regionName)

		listInstances := computeService.Instances.List("development-156617", regionName)

		resList, err := listInstances.Do()

		if err != nil {
			fmt.Println(err)
		}

		for _, resItem := range resList.Items {
			fmt.Println(resItem.Name)
			for _, device := range resItem.NetworkInterfaces {
				fmt.Println(device.NetworkIP)
			}
		}
	}



	// Do executes the "compute.instances.get" call.
	// Exactly one of *Instance or error will be non-nil. Any non-2xx status
	// code is an error. Response headers are in either
	// *Instance.ServerResponse.Header or (if a response was returned at
	// all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified
	// to check whether the returned error was because
	// http.StatusNotModified was returned

	//res, err := instances.Do()

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(res.MachineType)
}
