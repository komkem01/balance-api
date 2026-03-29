package transactions

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type TransactionStore interface {
	entitiesinf.TransactionEntity
	entitiesinf.WalletEntity
	entitiesinf.CategoryEntity
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db TransactionStore) *Module {
	tracer := otel.Tracer("balance.modules.transaction")
	svc := newService(&Options{Config: conf, tracer: tracer, db: db})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
