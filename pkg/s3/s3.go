package s3

import (
	"bytes"
	"vacabulary/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Manager struct {
	awsSession *session.Session
	config     config.AWSConfig
}

func (s *S3Manager) StorePdfFile(pathToFile string, file *bytes.Reader) (string, error) {
	sess := s._getAwsSession()

	s3 := s3manager.NewUploader(sess)

	res, err := s3.Upload(&s3manager.UploadInput{
		Bucket: aws.String("collections-words"),
		Key:    aws.String(pathToFile),
		ACL:    aws.String("public-read"),
		Body:   file,
	})
	if err != nil {
		return "", nil
	}

	return res.Location, nil
}

func (s *S3Manager) _getAwsSession() *session.Session {
	if s.awsSession == nil {
		s._setAwsSession()
	}

	return s.awsSession
}

func (s *S3Manager) _setAwsSession() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &s.config.Region,
		Credentials: credentials.NewStaticCredentials(s.config.AccessId, s.config.Secret, ""),
	}))

	s.awsSession = sess
}

func NewS3Manager(config config.AWSConfig) S3Manager {
	return S3Manager{
		config: config,
	}
}
