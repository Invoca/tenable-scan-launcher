package mocks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/mock"
)

type MockEC2API struct {
	ec2iface.EC2API
	ResettableMock
}

func (m *MockEC2API) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	fmt.Println("DescribeInstances Mock")
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
	}
}

// Define a mock struct to be used in your unit tests of myFunc.
type mockEC2Client struct {
	ec2iface.EC2API
	resp   ec2.DescribeInstancesOutput
	result []string
}

func (m *mockEC2Client) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &m.resp, nil
}


type fakeEC2 struct {
	*ec2.EC2
}

type fakeEC2DescribeInstance struct {
	*fakeEC2
	ReturnInstance ec2.Instance
	mock.Mock
}

func (f *fakeEC2DescribeInstance) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	f.Called(input)

	return &ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					&f.ReturnInstance,
				},
			},
		},
	}, nil
}
