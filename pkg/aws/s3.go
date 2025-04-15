package aws

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	getBucketName() string
	UploadFile(ctx context.Context, key string, file *multipart.File) error
}

type s3Client struct {
	client *s3.Client
}

func NewS3Client(ctx context.Context) (*s3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, err
	}

	s3 := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.Region = "us-west-1"
	})

	return &s3Client{
		client: s3,
	}, nil
}

func (s *s3Client) getBucketName() string {
	return "zchat-bucket"
}

func (s *s3Client) UploadFile(ctx context.Context, key string, file *multipart.File) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.getBucketName()),
		Key:    aws.String(key),
		Body:   *file,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", s.getBucketName())
		} else {
			log.Printf("Couldn't upload file to %v:%v. Here's why: %v\n",
				s.getBucketName(), s.getBucketName(), err)
		}
	} else {
		err = s3.NewObjectExistsWaiter(s.client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(s.getBucketName()), Key: aws.String(key)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist.\n", key)
		}
	}

	return err
}
