package store

import (
	gen_store "github.com/oapi-codegen-multiple-packages-example/internal/gen/store"
	"github.com/oapi-codegen-multiple-packages-example/internal/service/store"
)

type Handler struct {
	service store.Service
}

func NewHandler(svc store.Service) gen_store.ServerInterface {
	return &Handler{
		service: svc,
	}
}
