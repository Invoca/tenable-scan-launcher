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



/*
func (m *MockEC2API) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.DescribeVolumesOutput), args.Error(1)
	}
}

func (m *MockEC2API) CreateSnapshot(input *ec2.CreateSnapshotInput) (*ec2.Snapshot, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.Snapshot), args.Error(1)
	}
}

func (m *MockEC2API) DeleteSnapshot(input *ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.DeleteSnapshotOutput), args.Error(1)
	}
}

func (m *MockEC2API) DescribeSnapshots(input *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*ec2.DescribeSnapshotsOutput), args.Error(1)
	}
}

type MockVolumeRetriever struct {
	ResettableMock
}

func (m *MockVolumeRetriever) Snapshots() ([]*aws.Snapshot, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]*aws.Snapshot), args.Error(1)
	}
}

type MockVolume struct {
	ResettableMock
}

func (m *MockVolume) Id() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockVolume) CreateSnapshot() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockVolume) Snapshots() ([]*aws.Snapshot, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]*aws.Snapshot), args.Error(1)
	}
}

func (m *MockVolume) CollectSnapshotsForDelete(deleter util.CollectionDeleter) error {
	args := m.Called(deleter)
	return args.Error(0)
}

type MockVolumeFactory struct {
	ResettableMock
}

func (m *MockVolumeFactory) VolumesMatchingTags(tags map[string]string) ([]aws.Volume, error) {
	args := m.Called(tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]aws.Volume), args.Error(1)
	}
}

func (m *MockVolumeFactory) Create(id string, tags map[string]string) aws.Volume {
	args := m.Called(id, tags)
	if args.Get(0) == nil {
		return nil
	} else {
		return args.Get(0).(aws.Volume)
	}
}
*/