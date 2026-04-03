package storage

import (
	"balance/internal/config"
	"balance/internal/otel"

	"go.opentelemetry.io/otel/trace"
)

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	client Client
}

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config]) *Module {
	tracer := otel.Tracer("balance.modules.storage")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		client: NewFromEnv(),
	})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
