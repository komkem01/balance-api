package exampletwo

import (
	entitiesinf "balance/app/modules/entities/inf"
	"balance/internal/config"
)

type (
	Module struct {
		Svc *Service
		Ctl *Controller
	}
	Service    struct{}
	Controller struct{}

	Config struct{}
)

func New(conf *config.Config[Config], db entitiesinf.ExampleEntity) *Module {
	return &Module{
		Svc: &Service{},
		Ctl: &Controller{},
	}
}
