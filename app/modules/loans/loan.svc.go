package loans

import (
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     LoanStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     LoanStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}
