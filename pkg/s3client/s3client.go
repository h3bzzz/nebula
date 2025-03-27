package s3client

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"bytes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client provides methods to interact with AWS S3
type S3Client struct {
	client *s3.Client
	bucket string
}

// New creates a new S3Client instance
func New() (*S3Client, error) {
	// Load AWS configuration from environment variables or AWS credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Client{
		client: client,
		bucket: os.Getenv("AWS_BUCKET"),
	}, nil
}

// ListObjects lists all objects in the bucket with the given prefix
func (s *S3Client) ListObjects(prefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

// GetObject retrieves an object from S3
func (s *S3Client) GetObject(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetObjectURL returns a presigned URL for an object
func (s *S3Client) GetObjectURL(key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object URL: %w", err)
	}

	return presignResult.URL, nil
}

// UploadObject uploads an object to S3
func (s *S3Client) UploadObject(key string, data []byte, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

// DeleteObject deletes an object from S3
func (s *S3Client) DeleteObject(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}
