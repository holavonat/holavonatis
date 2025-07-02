package r2

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Cloudflare struct {
	BucketName        string
	ObjectPath        string
	AccessKeyID       string
	SecretAccessKey   string
	EndpointURL       string
	Client            *s3.Client
	PublicEndpointURL string
}

func NewClient(cloudflare Cloudflare) (Cloudflare, error) {
	if cloudflare.AccessKeyID == "" || cloudflare.SecretAccessKey == "" || cloudflare.EndpointURL == "" {
		return Cloudflare{}, errors.New("missing required credentials")
	}

	if cloudflare.SecretAccessKey != "" && len(cloudflare.SecretAccessKey) <= 10 {
		return Cloudflare{}, errors.New("invalid secret access key length")
	}

	if cloudflare.BucketName == "" {
		return Cloudflare{}, errors.New("missing bucket name")
	}

	if cloudflare.EndpointURL == "" {
		return Cloudflare{}, errors.New("missing endpoint url")
	}

	if cloudflare.PublicEndpointURL == "" {
		return Cloudflare{}, errors.New("missing public endpoint")
	}

	client := s3.NewFromConfig(aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			cloudflare.AccessKeyID, cloudflare.SecretAccessKey, ""),
		Region: "auto",
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           cloudflare.EndpointURL,
				SigningRegion: "auto",
			}, nil
		}),
	})

	cloudflare.Client = client

	return cloudflare, nil
}

type UploadedFile struct {
	ETag       string `json:"etag"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mime_type"`
	PublicLink string `json:"public_link"`
}

func (c *Cloudflare) UploadFile(filename string, file []byte, contentType string, contentEncoding string) (UploadedFile, error) {
	if c.Client == nil {
		return UploadedFile{}, errors.New("client not initialized")
	}

	var objectPath string
	if c.ObjectPath == "" {
		objectPath = filename
	} else {
		objectPath = strings.TrimSuffix(c.ObjectPath, "/") + "/" + filename
	}

	output, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:          aws.String(c.BucketName),
		Key:             aws.String(objectPath),
		Body:            bytes.NewReader(file),
		ContentType:     aws.String(contentType),
		ContentEncoding: aws.String(contentEncoding),
	})
	if err != nil {
		return UploadedFile{}, err
	}

	return UploadedFile{
		ETag:       *output.ETag,
		Filename:   filename,
		MimeType:   contentType,
		PublicLink: c.PublicEndpointURL + "/" + objectPath,
	}, nil
}
