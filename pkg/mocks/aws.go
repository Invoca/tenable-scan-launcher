package mocks

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (m *MockEC2API) UploadFile(s3Key string, bucketName string, data []byte) error {
	args := m.Called()
	return args.Error(0)
}

type MockS3API struct {
	s3iface.S3API
	ResettableMock
}

func (s *MockS3API) CreateMultipartUploadWithContext(aws.Context, *s3.CreateMultipartUploadInput, ...request.Option) (*s3.CreateMultipartUploadOutput, error) {
	args := s.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.CreateMultipartUploadOutput), args.Get(1).(error)
}

func (s *MockS3API) PutObjectRequest(*s3.PutObjectInput) (*request.Request, *s3.PutObjectOutput) {
	args := s.Called()
	return args.Get(0).(*request.Request), args.Get(1).(*s3.PutObjectOutput)
}

func (s *MockS3API) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := s.Called()
	return args.Get(0).(*s3manager.UploadOutput), args.Error(1)
}
