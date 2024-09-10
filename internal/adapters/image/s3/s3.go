package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io"
	"net/http"
)

type Adapter struct {
	Session  s3iface.S3API
	s3Bucket string
}

func NewAdapter(awsRegion, s3Bucket string) (*Adapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		Session:  s3.New(sess),
		s3Bucket: s3Bucket,
	}, nil
}

func (a Adapter) Upload(imageName, imageUrl string) (err error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = a.Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(a.s3Bucket),
		Key:    aws.String(imageName),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a Adapter) CheckImageAlreadyUploaded(imageName string) bool {
	_, err := a.Session.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(a.s3Bucket),
		Key:    aws.String(imageName),
	})
	if err == nil {
		return true
	}

	return false
}
