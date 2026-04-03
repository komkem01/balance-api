package transactions

import (
	"balance/app/modules/entities/ent"
	"balance/app/modules/storage"
	"balance/internal/config"
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     TransactionStore
	sto    storage.Client
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     TransactionStore
	sto    storage.Client
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db, sto: opt.sto}
}

func parseTransactionType(value string) (ent.TransactionType, bool) {
	switch value {
	case string(ent.TransactionTypeIncome):
		return ent.TransactionTypeIncome, true
	case string(ent.TransactionTypeExpense):
		return ent.TransactionTypeExpense, true
	default:
		return "", false
	}
}

func (s *Service) resolveImageURL(ctx context.Context, rawURL string) string {
	if s.sto == nil {
		return rawURL
	}

	return s.sto.DisplayImageURL(ctx, rawURL)
}
