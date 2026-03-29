package members

import (
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

type Service struct {
	tracer trace.Tracer
	db     MemberStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     MemberStore
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}
