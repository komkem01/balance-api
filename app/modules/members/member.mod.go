package members

import (
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], db MemberStore) *Module {
	tracer := otel.Tracer("balance.modules.member")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		db:     db,
	})
	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
