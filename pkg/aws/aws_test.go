package aws

import (
	"fmt"
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type getInstanceIpsTestCast struct {
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
						PrivateIpAddress: 	   aws.String("1.1.1.1"),
						State: &runningState,
					},
					{
						PrivateIpAddress: 	   aws.String("2.2.2.2"),
						State: &runningState,
					},
					{
						PrivateIpAddress: 	   aws.String("3.3.3.3"),
						State: &runningState,
					},
				},

			},
		},
	}

	testCases := []getInstanceIpsTestCast{
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

	for _, testCase := range testCases {
		t.Logf("TestGetAWSInstances: %s", testCase.desc)
		testCase.setup()


		ec2api := AWSEc2{}

		_, err := ec2api.getInstances(mockEc2)

		mockEc2.AssertExpectations(t)


		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	t.Logf("TestGetAWSInstances: pass nil object to getInstances")

	ec2api := AWSEc2{}
	_, err := ec2api.getInstances(nil)
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
						PrivateIpAddress: 	   aws.String("1.1.1.1"),
						State: &runningState,
					},
					{
						PrivateIpAddress: 	   aws.String("2.2.2.2"),
						State: &runningState,
					},
					{
						PrivateIpAddress: 	   aws.String("3.3.3.3"),
						State: &nonRunningState,
					},
				},

			},
		}
	ec2api := AWSEc2{}
	log.Debug("TestParseInstances: Instances Are passed to parseInstance")
	err := ec2api.parseInstances(resp)
	assert.NoError(t, err)
	assert.Equal(t, ec2api.IPs, runningInstances)

	log.Debug("TestParseInstances: Nothing is passed to parseInstance")
	err = ec2api.parseInstances(nil)
	assert.Error(t, err)
}



