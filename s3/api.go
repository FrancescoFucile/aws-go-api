package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// S3 represents a configurable interface to communicate with S3.
type S3 struct {
	AwsConfig aws.Config
}

// Default generates an S3 interface using default values
func (s *S3) Default() {
	s.AwsConfig = aws.Config{
		Credentials:      credentials.NewStaticCredentials(defId, defEndpoint, defToken),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(defRegion),
		Endpoint:         aws.String(defEndpoint),
	}
}
