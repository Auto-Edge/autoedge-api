package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	client     *s3.Client
	bucketName string
}

func NewStorageService(client *s3.Client, bucketName string) *StorageService {
	return &StorageService{
		client:     client,
		bucketName: bucketName,
	}
}

// UploadFile uploads a file stream to S3 and returns the key
func (s *StorageService) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error) {
	key := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), fileName)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

// GetPresignedURL generates a temporary URL for downloading
func (s *StorageService) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return req.URL, nil
}
