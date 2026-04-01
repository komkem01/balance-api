package budgets

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type BudgetStore interface {
	entitiesinf.BudgetEntity
	entitiesinf.MemberEntity
	entitiesinf.CategoryEntity
	entitiesinf.WalletEntity
	entitiesinf.TransactionEntity
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db BudgetStore) *Module {
	tracer := otel.Tracer("balance.modules.budget")
	svc := newService(&Options{Config: conf, tracer: tracer, db: db})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
