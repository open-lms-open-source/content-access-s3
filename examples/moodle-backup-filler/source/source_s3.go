package source

import (
	"fmt"
	"io"
	"net/http"

	// AWS
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"moodle-backup-filler/config"
	"moodle-backup-filler/source/s3"
)

var (
	s3Client *S3Client
)

// S3Client provides a persistent S3 session across multiple S3ContentReader
// objects.
type S3Client struct {
	credentials *credentials.Credentials
	stsCreds    *credentials.Credentials
	awsConfig   *aws.Config
	session     client.ConfigProvider
	s3Client    *s3wrapper.S3
}

func newS3Client() (*S3Client, error) {
	c := &S3Client{}

	c.credentials = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&ec2rolecreds.EC2RoleProvider{Client: ec2metadata.New(session.New())},
	})

	if _, err := c.credentials.Get(); err != nil {
		return nil, err
	}

	c.awsConfig = &aws.Config{
		Credentials: c.credentials,
		HTTPClient:  http.DefaultClient,
		Region:      aws.String(config.Config.S3Region),
	}

	session, err := session.NewSession(c.awsConfig)
	if err != nil {
		return nil, err
	}
	c.session = session

	if config.Config.S3AssumeRoleARN != "" {
		c.stsCreds = stscreds.NewCredentials(c.session, config.Config.S3AssumeRoleARN)
	}

	if c.stsCreds != nil {
		c.s3Client = s3wrapper.New(c.session, &aws.Config{Credentials: c.stsCreds})
	} else {
		c.s3Client = s3wrapper.New(c.session, nil)
	}

	return c, nil
}

// S3ContentReader implements the ContentReader interface for files
// contained in an S3 bucket.
type S3ContentReader struct {
	reader io.ReadCloser
	size   int64
}

// NewS3ContentReader returns a ContentReader for the given contentHash,
// which reads the file from S3.
func NewS3ContentReader(contentHash string) (*S3ContentReader, error) {
	if s3Client == nil {
		c, err := newS3Client()
		if err != nil {
			return nil, err
		}
		s3Client = c
	}

	paddedHash := fmt.Sprintf("%s____", contentHash) // ensure slices below don't fail if contentHash is invalid
	filePath := fmt.Sprintf("/%s/%s/%s", paddedHash[:2], paddedHash[2:4], contentHash)

	response, err := s3Client.s3Client.GetObjectWithRetry(&s3.GetObjectInput{
		Bucket: &config.Config.S3Bucket,
		Key:    &filePath,
	}, 2000, 4) // timeout=2s, retries=4 (30s total since the timeout is doubled each retry)
	if err != nil {
		return nil, err
	}

	return &S3ContentReader{
		reader: response.Body,
		size:   *response.ContentLength,
	}, nil
}

// Size returns the size of the currently open file.
func (cr *S3ContentReader) Size() int64 {
	return cr.size
}

// Read reads bytes from the currently open file.
func (cr *S3ContentReader) Read(b []byte) (int, error) {
	return cr.reader.Read(b)
}

// Close closes the currently open file.
func (cr *S3ContentReader) Close() error {
	return cr.reader.Close()
}

// vim: nolist expandtab ts=4 sw=4
