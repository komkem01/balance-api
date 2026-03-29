package budgets

import (
	"balance/app/modules/entities/ent"
	"balance/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     BudgetStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     BudgetStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func parseBudgetPeriod(value string) (ent.BudgetPeriod, bool) {
	switch value {
	case string(ent.BudgetPeriodDaily):
		return ent.BudgetPeriodDaily, true
	case string(ent.BudgetPeriodWeekly):
		return ent.BudgetPeriodWeekly, true
	case string(ent.BudgetPeriodMonthly):
		return ent.BudgetPeriodMonthly, true
	default:
		return "", false
	}
}
