package gcloud

import (
	"fmt"
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


func TestGetRegionsForProject(t *testing.T) {

	serviceMock := mocks.GgCloudServiceMock{}
	gcloud := GCloud{}
	gcloud.SetupGCloud(&serviceMock)

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
				serviceMock.On("Zones", mock.Anything).Return(resp, nil)
			},
			shouldError: false,
		},
		{
			desc: "Error returned by region retrieval",
			setup: func() {
				serviceMock.Reset()
				serviceMock.On("Zones", mock.Anything).Return(resp, fmt.Errorf("error"))
			},
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Logf("TestGetRegionsForProject: %s", testCase.desc)
		testCase.setup()

		err := gcloud.getAllRegionsForProject()


		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			fmt.Print(gcloud.regions)
			assert.NoError(t, err)
		}
	}
}

func TestGetInstancesInRegion(t *testing.T) {

	serviceMock := mocks.GgCloudServiceMock{}
	gcloud := GCloud{}
	gcloud.SetupGCloud(&serviceMock)

	resp := []string{
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
		"4.4.4.4",
	}

	testCases := []getRegionsTestCast{
		{
			desc: "successful region retrieval",
			setup: func() {
				serviceMock.Reset()
				serviceMock.On("InstancesIPsInRegion", mock.Anything).Return(resp, nil)
			},
			shouldError: false,
		},
		{
			desc: "Error returned by region retrieval",
			setup: func() {
				serviceMock.Reset()
				serviceMock.On("InstancesIPsInRegion", mock.Anything).Return(resp, fmt.Errorf("error"))
			},
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Logf("TestGetInstancesInRegion: %s", testCase.desc)
		testCase.setup()

		err := gcloud.getInstancesInRegion("test")


		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			fmt.Print(gcloud.regions)
			assert.NoError(t, err)
		}
	}
}