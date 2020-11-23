package aws

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

type testCase struct {
	desc        string
	setup       func()
	shouldError bool
}

func TestGetAWSInstances(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	mockEc2 := &mocks.MockEC2API{}

	runningCode := int64(16)
	runningState := ec2.InstanceState{Code: &runningCode}

	resp := ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				ReservationId: aws.String("123ABC"),
				Instances: []*ec2.Instance{
					{
						PrivateIpAddress: aws.String("1.1.1.1"),
						State:            &runningState,
					},
					{
						PrivateIpAddress: aws.String("2.2.2.2"),
						State:            &runningState,
					},
					{
						PrivateIpAddress: aws.String("3.3.3.3"),
						State:            &runningState,
					},
				},
			},
		},
	}

	testCases := []testCase{
		{
			desc: "successful ip retrieval",
			setup: func() {
				mockEc2.Reset()
				mockEc2.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&resp, nil)
			},
			shouldError: false,
		},
		{
			desc: "error returned by ip retrieval",
			setup: func() {
				mockEc2.Reset()
				mockEc2.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&resp, fmt.Errorf("error"))
			},
			shouldError: true,
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc":        testCase.desc,
			"shouldError": testCase.shouldError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()

		ec2api := AwsSvc{}
		ec2api.Ec2svc = mockEc2
		_, err := ec2api.getInstances()

		mockEc2.AssertExpectations(t)

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	t.Logf("TestGetAWSInstances: pass nil object to getInstances")

	ec2api := AwsSvc{}
	_, err := ec2api.getInstances()
	assert.Error(t, err)

}

func TestParseInstances(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	runningCode := int64(16)
	runningState := ec2.InstanceState{Code: &runningCode}

	nonRunningCode := int64(0)
	nonRunningState := ec2.InstanceState{Code: &nonRunningCode}

	runningInstances := []string{
		"1.1.1.1",
		"2.2.2.2",
	}

	resp := []*ec2.Reservation{
		{
			ReservationId: aws.String("123ABC"),
			Instances: []*ec2.Instance{
				{
					PrivateIpAddress: aws.String("1.1.1.1"),
					State:            &runningState,
				},
				{
					PrivateIpAddress: aws.String("2.2.2.2"),
					State:            &runningState,
				},
				{
					PrivateIpAddress: aws.String("3.3.3.3"),
					State:            &nonRunningState,
				},
			},
		},
	}
	ec2api := AwsSvc{}
	log.Debug("TestParseInstances: Instances Are passed to parseInstance")
	err := ec2api.parseInstances(resp)
	assert.NoError(t, err)
	assert.Equal(t, ec2api.IPs, runningInstances)

	log.Debug("TestParseInstances: Nothing is passed to parseInstance")
	err = ec2api.parseInstances(nil)
	assert.Error(t, err)
}

func TestUploadFile(t *testing.T) {
	mockS3 := &mocks.MockS3API{}

	uploadOutput := s3manager.UploadOutput{}

	testCases := []testCase{
		{
			desc: "Uploads file to S3 without error",
			setup: func() {
				mockS3.Reset()
				mockS3.On("Upload", mock.Anything).Return(&uploadOutput, nil)
			},
			shouldError: false,
		},
		{
			desc: "Uploads file to S3 withs error",
			setup: func() {
				mockS3.Reset()
				mockS3.On("Upload", mock.Anything).Return(&uploadOutput, fmt.Errorf("Error"))
			},
			shouldError: true,
		},
	}

	for index, testCase := range testCases {
		log.WithFields(log.Fields{
			"desc":        testCase.desc,
			"shouldError": testCase.shouldError,
		}).Debug("Starting testCase " + strconv.Itoa(index))

		testCase.setup()

		ec2api := AwsSvc{}
		ec2api.s3Manager = mockS3

		err := ec2api.UploadFile("s3Key", "bucketName", []byte("data"))

		mockS3.AssertExpectations(t)

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
