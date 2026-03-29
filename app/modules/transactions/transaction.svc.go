package transactions

import (
	"balance/app/modules/entities/ent"
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     TransactionStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     TransactionStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
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
