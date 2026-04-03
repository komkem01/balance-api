package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel/trace"
)

const (
	transactionSlipMaxSize = int64(10 * 1024 * 1024)
	profileImageMaxSize    = int64(10 * 1024 * 1024)
)

var (
	ErrStorageNotConfigured = errors.New("storage is not configured")
	ErrImageRequired        = errors.New("image file is required")
	ErrImageEmpty           = errors.New("image file is empty")
	ErrImageTooLarge        = errors.New("image file exceeds 10MB limit")
	ErrProfileImageTooLarge = errors.New("profile image file exceeds 10MB limit")
	ErrImageInvalidType     = errors.New("invalid file type")
	ErrImageInvalidContent  = errors.New("invalid file content type")
)

type Client interface {
	UploadTransactionSlip(ctx context.Context, walletID string, fileHeader *multipart.FileHeader) (string, error)
	UploadProfileImage(ctx context.Context, memberID string, fileHeader *multipart.FileHeader) (string, error)
	DisplayImageURL(ctx context.Context, rawURL string) string
	Enabled() bool
}

type storageClient struct {
	client  *minio.Client
	bucket  string
	baseURL string
	enabled bool
}

type Service struct {
	tracer trace.Tracer
	cli    Client
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, cli: opt.client}
}

func (s *Service) Enabled() bool {
	if s == nil || s.cli == nil {
		return false
	}

	return s.cli.Enabled()
}

func (s *Service) UploadTransactionSlip(ctx context.Context, walletID string, fileHeader *multipart.FileHeader) (string, error) {
	if s == nil || s.cli == nil {
		return "", ErrStorageNotConfigured
	}

	return s.cli.UploadTransactionSlip(ctx, walletID, fileHeader)
}

func (s *Service) UploadProfileImage(ctx context.Context, memberID string, fileHeader *multipart.FileHeader) (string, error) {
	if s == nil || s.cli == nil {
		return "", ErrStorageNotConfigured
	}

	return s.cli.UploadProfileImage(ctx, memberID, fileHeader)
}

func (s *Service) DisplayImageURL(ctx context.Context, rawURL string) string {
	if s == nil || s.cli == nil {
		return rawURL
	}

	return s.cli.DisplayImageURL(ctx, rawURL)
}

func NewFromEnv() Client {
	return newStorageClientFromEnv()
}

func newStorageClientFromEnv() *storageClient {
	endpointRaw := firstNonEmpty(
		strings.TrimSpace(os.Getenv("STORAGE_RAILWAY__ENDPOINT")),
		strings.TrimSpace(os.Getenv("RAILWAY_STORAGE__ENDPOINT")),
	)
	bucket := firstNonEmpty(
		strings.TrimSpace(os.Getenv("STORAGE_RAILWAY__BUCKET")),
		strings.TrimSpace(os.Getenv("RAILWAY_STORAGE__BUCKET")),
	)
	accessKey := firstNonEmpty(
		strings.TrimSpace(os.Getenv("STORAGE_RAILWAY__ACCESS_KEY")),
		strings.TrimSpace(os.Getenv("RAILWAY_STORAGE__ACCESS_KEY")),
	)
	secretKey := firstNonEmpty(
		strings.TrimSpace(os.Getenv("STORAGE_RAILWAY__SECRET_KEY")),
		strings.TrimSpace(os.Getenv("RAILWAY_STORAGE__SECRET_KEY")),
	)

	if endpointRaw == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return &storageClient{enabled: false}
	}

	endpoint, secure, baseURL, err := parseStorageEndpoint(endpointRaw)
	if err != nil {
		return &storageClient{enabled: false}
	}

	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return &storageClient{enabled: false}
	}

	return &storageClient{client: cli, bucket: bucket, baseURL: baseURL, enabled: true}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func parseStorageEndpoint(raw string) (string, bool, string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", false, "", err
	}
	if u.Scheme == "" {
		return raw, false, strings.TrimRight("https://"+raw, "/"), nil
	}

	endpoint := u.Host
	secure := strings.EqualFold(u.Scheme, "https")
	base := strings.TrimRight(u.Scheme+"://"+u.Host, "/")
	return endpoint, secure, base, nil
}

func (s *storageClient) Enabled() bool {
	return s != nil && s.enabled
}

func (s *storageClient) UploadTransactionSlip(ctx context.Context, walletID string, fileHeader *multipart.FileHeader) (string, error) {
	if s == nil || !s.enabled {
		return "", ErrStorageNotConfigured
	}
	if fileHeader == nil {
		return "", ErrImageRequired
	}
	if fileHeader.Size <= 0 {
		return "", ErrImageEmpty
	}
	if fileHeader.Size > transactionSlipMaxSize {
		return "", ErrImageTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", ErrImageInvalidType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, transactionSlipMaxSize+1))
	if err != nil {
		return "", err
	}
	if int64(len(data)) > transactionSlipMaxSize {
		return "", ErrImageTooLarge
	}

	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", ErrImageInvalidContent
	}

	objectExt := path.Ext(strings.ToLower(fileHeader.Filename))
	walletPath := strings.TrimSpace(walletID)
	if walletPath == "" {
		walletPath = "unknown-wallet"
	}
	objectKey := fmt.Sprintf(
		"transaction-slips/%s/%d-%s%s",
		walletPath,
		time.Now().Unix(),
		uuid.NewString(),
		objectExt,
	)

	_, err = s.client.PutObject(
		ctx,
		s.bucket,
		objectKey,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.bucket, objectKey), nil
}

func (s *storageClient) UploadProfileImage(ctx context.Context, memberID string, fileHeader *multipart.FileHeader) (string, error) {
	if s == nil || !s.enabled {
		return "", ErrStorageNotConfigured
	}
	if fileHeader == nil {
		return "", ErrImageRequired
	}
	if fileHeader.Size <= 0 {
		return "", ErrImageEmpty
	}
	if fileHeader.Size > profileImageMaxSize {
		return "", ErrProfileImageTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", ErrImageInvalidType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, profileImageMaxSize+1))
	if err != nil {
		return "", err
	}
	if int64(len(data)) > profileImageMaxSize {
		return "", ErrProfileImageTooLarge
	}

	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", ErrImageInvalidContent
	}

	objectExt := path.Ext(strings.ToLower(fileHeader.Filename))
	memberPath := strings.TrimSpace(memberID)
	if memberPath == "" {
		memberPath = "unknown-member"
	}
	objectKey := fmt.Sprintf(
		"profile-images/%s/%d-%s%s",
		memberPath,
		time.Now().Unix(),
		uuid.NewString(),
		objectExt,
	)

	_, err = s.client.PutObject(
		ctx,
		s.bucket,
		objectKey,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.bucket, objectKey), nil
}

func (s *storageClient) DisplayImageURL(ctx context.Context, rawURL string) string {
	if s == nil || !s.enabled {
		return rawURL
	}

	objectKey, ok := s.extractObjectKey(rawURL)
	if !ok {
		return rawURL
	}

	presigned, err := s.client.PresignedGetObject(ctx, s.bucket, objectKey, 24*time.Hour, nil)
	if err != nil {
		return rawURL
	}

	return presigned.String()
}

func (s *storageClient) extractObjectKey(rawURL string) (string, bool) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", false
	}

	if !strings.Contains(trimmed, "://") {
		return strings.TrimPrefix(trimmed, "/"), true
	}

	u, err := url.Parse(trimmed)
	if err != nil {
		return "", false
	}

	pathValue := strings.TrimPrefix(u.Path, "/")
	if strings.HasPrefix(pathValue, s.bucket+"/") {
		return strings.TrimPrefix(pathValue, s.bucket+"/"), true
	}

	return "", false
}
