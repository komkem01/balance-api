package wallets

import (
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     WalletStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     WalletStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}
