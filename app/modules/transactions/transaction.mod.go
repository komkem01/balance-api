package transactions

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/app/modules/storage"
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type TransactionStore interface {
	entitiesinf.TransactionEntity
	entitiesinf.WalletEntity
	entitiesinf.CategoryEntity
	entitiesinf.MemberEntity
	entitiesinf.BudgetEntity
	entitiesinf.GoalEntity
	entitiesinf.GoalEntryEntity
	entitiesinf.MemberNotificationEntity
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db TransactionStore, sto storage.Client) *Module {
	tracer := otel.Tracer("balance.modules.transaction")
	svc := newService(&Options{Config: conf, tracer: tracer, db: db, sto: sto})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
