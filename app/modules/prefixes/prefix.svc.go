package prefixes

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.PrefixEntity // Database interface for prefix entities
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.PrefixEntity // Database interface for prefix entities
}

func newService(opt *Options) *Service {
	return &Service{
		tracer: opt.tracer,
		db:     opt.db,
	}
}
