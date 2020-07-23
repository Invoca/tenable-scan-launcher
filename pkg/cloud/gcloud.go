package cloud

import (
	"fmt"
	"google.golang.org/api/compute/v1"
)

// GetGCloudIPs
func GetGCloudIPs() {
	fmt.Println("Getting IPs from Google Cloud")
	compute.NewService(nil)
}
