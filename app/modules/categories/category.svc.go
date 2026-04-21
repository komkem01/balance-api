package categories

import (
	"balance/app/modules/entities/ent"
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     CategoryStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     CategoryStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func parseCategoryType(value string) (ent.CategoryType, bool) {
	switch value {
	case string(ent.CategoryTypeIncome):
		return ent.CategoryTypeIncome, true
	case string(ent.CategoryTypeExpense):
		return ent.CategoryTypeExpense, true
	default:
		return "", false
	}
}

func parseCategoryPurpose(value string) (ent.CategoryPurpose, bool) {
	switch value {
	case string(ent.CategoryPurposeLoanRepayment):
		return ent.CategoryPurposeLoanRepayment, true
	default:
		return "", false
	}
}
