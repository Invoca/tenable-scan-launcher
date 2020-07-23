package cloud

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/ec2"
)


type EC2Ips struct {
	ips		[]string
	api		ec2iface.EC2API

}

/*
type clientFactory interface {
	build(d *Driver) Ec2Client
}


func (d *Driver) buildClient() Ec2Client {
	return ec2.New(nil)
}
*/


// GetAWSIPs
func (m *EC2Ips) GetAWSIPs() (error) {
	log.Debug("Getting AWS IPs")

	if m.api == nil {
		fmt.Errorf("GetAWSIPs: api object is nil")
	}

	instances, err := m.getInstances(m.api)
	if err != nil {
		return fmt.Errorf("GetAWSIPs: Could not get list of instances %s", err)
	}
	err = m.parseInstances(instances)
	if err != nil {
		return fmt.Errorf("GetAWSIPs: Could not parse instances given %s", err)
	}
	return nil
}

func (m *EC2Ips) getInstances(ec2Svc ec2iface.EC2API) ([]*ec2.Reservation, error) {

	log.Debug("0")

	// Describe instances. We can add filter flags later if needed
	resp, err := ec2Svc.DescribeInstances(nil)

	log.Debug("1")
	if err != nil {
		return nil, fmt.Errorf("Error listing instances %s", err)
	}
	log.Debug("2")
	log.Debug(resp.Reservations)
	return resp.Reservations, nil
}

func (m *EC2Ips) parseInstances(reservations []*ec2.Reservation) (error) {
	var privateIps []string

	if reservations == nil {
		return fmt.Errorf("parseInstances: Passed empty reservations object")
	}

	for idx, res := range reservations {
		log.Debug("Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range reservations[idx].Instances {

			// Status code 16 is Runnning state
			if *inst.State.Code == 16 {
				log.Debug("Instance private ip: ", *inst.PrivateIpAddress)
				privateIps = append(privateIps, *inst.PrivateIpAddress)
			}
		}
	}

	log.Debug(privateIps)

	m.ips = privateIps
	return nil
}
