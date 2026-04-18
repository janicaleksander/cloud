package aws

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client *s3.Client
}

func NewAWSStorage(client *s3.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) StoreFile(ctx context.Context, bucket string, fileID string, contentType string, reader io.Reader) error {
	fmt.Println(reader)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileID),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return err
	}
	slog.Info("File stored in S3", "bucket", bucket, "fileID", fileID)
	return nil
}

func (s *Storage) GetFile(ctx context.Context, bucket string, fileID string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}

	slog.Info("File retrieved from S3", "bucket", bucket, "fileID", fileID)
	return resp.Body, nil
}

func (s *Storage) RemoveFile(ctx context.Context, bucket string, fileID string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return err
	}

	slog.Info("File removed from S3", "bucket", bucket, "fileID", fileID)
	return nil
}
