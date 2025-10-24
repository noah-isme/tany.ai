package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/tanydotai/tanyai/backend/internal/config"
)

// S3Storage implements uploads backed by an S3-compatible service.
type S3Storage struct {
	client     *s3.Client
	bucket     string
	region     string
	publicBase string
}

// NewS3Storage creates a new S3Storage instance.
func NewS3Storage(cfg config.S3Config) (*S3Storage, error) {
	if cfg.Region == "" || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("storage: incomplete s3 configuration")
	}

	awsCfg := aws.Config{
		Region: cfg.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.EndpointResolver = s3.EndpointResolverFromURL(cfg.Endpoint)
			o.UsePathStyle = true
		} else if cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})

	return &S3Storage{
		client:     client,
		bucket:     cfg.Bucket,
		region:     cfg.Region,
		publicBase: strings.TrimSuffix(cfg.PublicBaseURL, "/"),
	}, nil
}

// Put uploads an object to S3 and returns the public URL.
func (s *S3Storage) Put(ctx context.Context, key string, content []byte, contentType string) (string, error) {
	if len(content) == 0 {
		return "", fmt.Errorf("storage: empty content")
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(content),
		ContentType:   aws.String(contentType),
		ContentLength: int64(len(content)),
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("storage: s3 upload failed: %w", err)
	}

	return s.publicURL(key), nil
}

func (s *S3Storage) publicURL(key string) string {
	cleaned := strings.TrimPrefix(key, "/")
	if s.publicBase != "" {
		return s.publicBase + "/" + cleaned
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, cleaned)
}
