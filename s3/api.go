package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
)

// Bucket represents a configurable interface to an S3 bucket.
type Bucket struct {
	Name string
}

// LocalSession returns a session configured to work with LocalStack.
func LocalSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(defId, defEndpoint, defToken),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(defRegion),
		Endpoint:         aws.String(defEndpoint),
	}))
}

// Get executes a download from S3.
func (s *Bucket) Get(key string) (out []byte, err error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Name),
		Key:    aws.String(key),
	}
	svc := s3.New(LocalSession())
	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, fmt.Errorf(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			return nil, fmt.Errorf(err.Error())
		}
		return
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return body, nil
}

// Put uploads a file to S3.
func (s *Bucket) Put(key string, data []byte) (err error) {
	_, err = s3.New(LocalSession()).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s.Name),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(data),
		ContentLength:        aws.Int64(int64(len(data))),
		ContentType:          aws.String(http.DetectContentType(data)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return fmt.Errorf("cannot upload to S3: " + err.Error())
	}
	return nil
}
