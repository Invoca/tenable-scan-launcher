package mocks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type MockRunner struct {
	ResettableMock
}

func (m *MockRunner) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	fmt.Println("DescribeInstances Mock")
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
	}
}
