package goals

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type GoalStore interface {
	entitiesinf.GoalEntity
	entitiesinf.GoalEntryEntity
	entitiesinf.MemberEntity
	entitiesinf.WalletEntity
	entitiesinf.LoanEntity
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db GoalStore) *Module {
	tracer := otel.Tracer("balance.modules.goal")
	svc := newService(&Options{Config: conf, tracer: tracer, db: db})
	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}
