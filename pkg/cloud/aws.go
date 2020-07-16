package cloud

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetAWSIPs
func GetAWSIPs() ([]string, error) {
	log.Debug("Getting AWS IPs")
	instances, err := getInstances()
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could  not get list of instances %s", err)
	}
	ips, err := parseInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could  not parse instances given %s", err)
	}
	return ips, nil
}

func getInstances() ([]*ec2.Reservation, error) {

	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	ec2Svc := ec2.New(sess)

	// Describe instances. We can add filter flags later if needed
	resp, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return nil, fmt.Errorf("Error listing instances %s", err)
	}
	log.Debug(resp.Reservations)
	return resp.Reservations, nil
}

func parseInstances(reservations []*ec2.Reservation) ([]string, error) {
	var privateIps []string

	if reservations == nil {
		return nil, fmt.Errorf("parseInstances: Passed empty reservations object")
	}

	for idx, res := range reservations {
		log.Debug("Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range reservations[idx].Instances {

			// Status code 16 is Runnning state
			if *inst.State.Code == 16 {
				log.Debug("    - Instance private ip: ", *inst.PrivateIpAddress)
				privateIps = append(privateIps, *inst.PrivateIpAddress)
			}
		}
	}
	log.Debug(privateIps)
	return privateIps, nil
}
