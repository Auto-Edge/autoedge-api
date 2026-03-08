package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// StorageService defines the contract for our cloud storage operations
type StorageService interface {
	GeneratePresignedUploadURL(ctx context.Context, bucket string, key string, expireMinutes int) (string, error)
}

// s3Service holds the specific AWS S3 presign client
type s3Service struct {
	presignClient *s3.PresignClient
}

// NewStorageService initializes the service with a base AWS S3 client
func NewStorageService(client *s3.Client) StorageService {
	return &s3Service{
		presignClient: s3.NewPresignClient(client),
	}
}

// GeneratePresignedUploadURL creates the secure ticket for the React frontend
func (s *s3Service) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expireMinutes int) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireMinutes) * time.Minute
	})

	if err != nil {
		return "", err
	}

	return req.URL, nil
}