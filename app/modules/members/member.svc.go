package members

import (
	"context"
	"mime/multipart"

	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type MemberStore interface {
	entitiesinf.MemberEntity
	entitiesinf.GenderEntity
	entitiesinf.PrefixEntity
	entitiesinf.MemberAccountEntity
}

type StorageStore interface {
	UploadProfileImage(ctx context.Context, memberID string, fileHeader *multipart.FileHeader) (string, error)
	DisplayImageURL(ctx context.Context, rawURL string) string
	Enabled() bool
}

type Service struct {
	tracer trace.Tracer
	db     MemberStore
	sto    StorageStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     MemberStore
	sto    StorageStore
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
		sto:    opt.sto,
	}
}

func (s *Service) displayProfileImageURL(ctx context.Context, rawURL string) string {
	if s == nil || s.sto == nil {
		return rawURL
	}

	return s.sto.DisplayImageURL(ctx, rawURL)
}
