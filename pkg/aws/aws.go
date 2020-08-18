package aws

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/ec2"
)


type AWSEc2 struct {
	IPs		[]string
	Ec2svc ec2iface.EC2API
}

func (m *AWSEc2)  Setup(config *config.BaseConfig) error {
	if config.IncludeAWS == false {
		return fmt.Errorf("Setup: AWS is not supposed to be included")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	if sess == nil {
		return fmt.Errorf("setupAWS: Error creating session object")
	}
	m.IPs = *new([]string)
	m.Ec2svc = ec2.New(sess)
	return nil
}

func (m *AWSEc2) GatherIPs() ([]string, error) {
	log.Debug("Getting AWS IPs")

	if m.Ec2svc == nil {
		fmt.Errorf("GetAWSIPs: Ec2svc object is nil")
	}

	instances, err := m.getInstances(m.Ec2svc)
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could not get list of instances %s", err)
	}

	err = m.parseInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could not parse instances given %s", err)
	}

	return m.IPs, nil
}

func (m *AWSEc2) getInstances(ec2Svc ec2iface.EC2API) ([]*ec2.Reservation, error) {

	if ec2Svc == nil {
		return nil, fmt.Errorf("getInstances: Passed empty ec2iface object")
	}

	// Describe instances. We can add filter flags later if needed
	resp, err := ec2Svc.DescribeInstances(nil)

	if err != nil {
		return nil, fmt.Errorf("Error listing instances %s", err)
	}

	log.Debug(resp.Reservations)
	return resp.Reservations, nil
}

func (m *AWSEc2) parseInstances(reservations []*ec2.Reservation) error {
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

	m.IPs = privateIps
	return nil
}
