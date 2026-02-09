package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/yurisasc/algafood-go/internal/config"
)

// StorageService interface for file storage operations
type StorageService interface {
	Store(filename string, contentType string, content io.Reader) (string, error)
	Retrieve(filename string) (io.ReadCloser, error)
	Delete(filename string) error
	GetURL(filename string) string
}

// NewStorageService creates a new storage service based on configuration
func NewStorageService(cfg *config.StorageConfig) (StorageService, error) {
	switch cfg.Type {
	case "s3":
		return NewS3StorageService(cfg)
	default:
		return NewLocalStorageService(cfg)
	}
}

// LocalStorageService stores files locally
type LocalStorageService struct {
	directory string
}

// NewLocalStorageService creates a new local storage service
func NewLocalStorageService(cfg *config.StorageConfig) (*LocalStorageService, error) {
	// Create directory if not exists
	if err := os.MkdirAll(cfg.Local.Directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorageService{directory: cfg.Local.Directory}, nil
}

func (s *LocalStorageService) Store(filename string, contentType string, content io.Reader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext

	fullPath := filepath.Join(s.directory, newFilename)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return newFilename, nil
}

func (s *LocalStorageService) Retrieve(filename string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.directory, filename)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

func (s *LocalStorageService) Delete(filename string) error {
	fullPath := filepath.Join(s.directory, filename)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorageService) GetURL(filename string) string {
	return filepath.Join(s.directory, filename)
}

// S3StorageService stores files in AWS S3
type S3StorageService struct {
	client    *s3.Client
	bucket    string
	directory string
}

// NewS3StorageService creates a new S3 storage service
func NewS3StorageService(cfg *config.StorageConfig) (*S3StorageService, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.S3.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3StorageService{
		client:    client,
		bucket:    cfg.S3.Bucket,
		directory: cfg.S3.Directory,
	}, nil
}

func (s *S3StorageService) Store(filename string, contentType string, content io.Reader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	key := s.getKey(newFilename)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        content,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return newFilename, nil
}

func (s *S3StorageService) Retrieve(filename string) (io.ReadCloser, error) {
	key := s.getKey(filename)
	result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	return result.Body, nil
}

func (s *S3StorageService) Delete(filename string) error {
	key := s.getKey(filename)
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

func (s *S3StorageService) GetURL(filename string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, s.getKey(filename))
}

func (s *S3StorageService) getKey(filename string) string {
	if s.directory != "" {
		return s.directory + "/" + filename
	}
	return filename
}
