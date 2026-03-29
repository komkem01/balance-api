package wallets

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type WalletStore interface {
	entitiesinf.WalletEntity
	entitiesinf.MemberEntity
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db WalletStore) *Module {
	tracer := otel.Tracer("balance.modules.wallet")
	svc := newService(&Options{Config: conf, tracer: tracer, db: db})
	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}
