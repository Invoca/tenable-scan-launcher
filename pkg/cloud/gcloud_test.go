package cloud

import (
	"github.com/Invoca/tenable-scan-launcher/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)


type getRegionsTestCast struct {
	desc        string
	setup       func()
	shouldError bool
}


func TestgetAllRegionsForProject(t *testing.T) {

	serviceMock := mocks.GgCloudServiceMock{}
	gcloud := GCloud{}
	gcloud.SetupGCloud(&serviceMock, "test")

	resp := []string{
		"Never",
		"Eat",
		"Soggy",
		"Waffles",
	}

	testCases := []getRegionsTestCast{
		{
			desc: "successful region retrieval",
			setup: func() {
				serviceMock.Reset()
				serviceMock.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&resp, nil)
			},
			shouldError: false,
		},
		{
			desc: "Error returned by region retrieval",
			setup: func() {
				serviceMock.Reset()
				serviceMock.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&resp, nil)
			},
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Logf("TestGetAWSInstances: %s", testCase.desc)
		testCase.setup()

		err := gcloud.getAllRegionsForProject()


		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

