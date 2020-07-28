package wrapper

type GCloudWrapper interface {
	Zones() ([]string, error)
	InstancesIPsInRegion(string) ([]string, error)
}
