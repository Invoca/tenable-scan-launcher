package cloud

import (
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	log "github.com/sirupsen/logrus"
)

type getInstanceIpsTestCast struct {
	desc        string
	setup       func()
	shouldError bool
}

func TestGetAWSIPs(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	mockEc2 := new(mocks.MockEC2API)
	//	mockInstances := new(mocks.MockEC2API)

	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//		SharedConfigState: session.SharedConfigEnable,
	//	}))

	/*
	mockOptions := mocks.MockOptions{
		ResettableMock:    mocks.ResettableMock{},
		SharedConfigState: false,

	}

	mockSession := new(mocks.MockSessionAPI)
	*/
	resp := ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				ReservationId: aws.String("123ABC"),
				Instances: []*ec2.Instance{
					{
						PrivateIpAddress: 	   aws.String("1.1.1.1"),
					},
					{
						PrivateIpAddress: 	   aws.String("2.2.2.2"),
					},
					{
						PrivateIpAddress: 	   aws.String("3.3.3.3"),
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
				mockEc2.On("New", mock.Anything).Return(mockEc2)
				mockEc2.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&mockEc2)
				mockEc2.On("Reservations", mock.AnythingOfType("[]*Reservation")).Return(&resp)
			},
			shouldError: false,
		},
		{
			desc: "error returned by snapshot delete",
			setup: func() {
				mockEc2.Reset()
				//mockEc2.On("New", mock.AnythingOfType("*EC2")).Return(mockEc2)
				mockEc2.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(resp)
			},
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		//t.Logf("TEST: %s", testCase.desc)
		testCase.setup()

		_, err := getInstances()

		mockEc2.AssertExpectations(t)


		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

/*
type snapshotTestCast struct {
	desc        string
	setup       func()
	shouldError bool
}

func TestSnapshotDelete(t *testing.T) {
	mock := new(mocks.MockEC2API)

	snapshotId := "snap-1"
	snapshot := aws.NewSnapshot(mock, snapshotId, nil)

	deleteSnapshotInput := &ec2.DeleteSnapshotInput{SnapshotId: &snapshotId}
	deleteSnapshotOutput := &ec2.DeleteSnapshotOutput{}

	testCases := []snapshotTestCast{
		{
			desc: "successful snapshot delete",
			setup: func() {
				mock.Reset()
				mock.On("DeleteSnapshot", deleteSnapshotInput).Return(deleteSnapshotOutput, nil)
			},
			shouldError: false,
		},
		{
			desc: "error returned by snapshot delete",
			setup: func() {
				mock.Reset()
				mock.On("DeleteSnapshot", deleteSnapshotInput).Return(nil, errors.New("snapshot delete error"))
			},
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Logf("TEST: %s", testCase.desc)
		testCase.setup()

		err := snapshot.Delete()

		mock.AssertExpectations(t)

		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
*/
