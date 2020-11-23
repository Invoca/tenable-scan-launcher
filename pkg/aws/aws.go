package aws

import (
	"bytes"
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

type AwsSvc struct {
	IPs    []string
	Ec2svc ec2iface.EC2API
	s3Manager s3Wrapper
}

func (m *AwsSvc) Setup(config *config.BaseConfig) error {
	if config.IncludeAWS == false {
		return fmt.Errorf("Setup: AWS is not supposed to be included")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	m.Ec2svc = ec2.New(sess)
	m.s3Manager = s3manager.NewUploader(sess)

	return nil
}

func (m *AwsSvc) GatherIPs() ([]string, error) {
	log.Debug("Getting AWS IPs")

	if m.Ec2svc == nil {
		return nil, fmt.Errorf("GetAWSIPs: Ec2svc object is nil")
	}

	instances, err := m.getInstances()
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could not get list of instances %s", err)
	}

	err = m.parseInstances(instances)
	if err != nil {
		return nil, fmt.Errorf("GetAWSIPs: Could not parse instances given %s", err)
	}

	return m.IPs, nil
}

func (m *AwsSvc) getInstances() ([]*ec2.Reservation, error) {

	if m.Ec2svc == nil {
		return nil, fmt.Errorf("getInstances: Passed empty ec2iface object")
	}

	// Describe instances. We can add filter flags later if needed
	resp, err := m.Ec2svc.DescribeInstances(nil)

	if err != nil {
		return nil, fmt.Errorf("Error listing instances %s", err)
	}

	log.Debug(resp.Reservations)
	return resp.Reservations, nil
}

func (m *AwsSvc) parseInstances(reservations []*ec2.Reservation) error {
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

func (a *AwsSvc) UploadFile(s3Key string, bucketName string, data []byte) error {
	_, err := a.s3Manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("GetFileFromS3: Error getting resp from s3")
	}

	return nil
}

type s3Wrapper interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}
