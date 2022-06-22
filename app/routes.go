package app

import (
	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/api-middleware/usecase"

	"github.com/go-chi/chi"
)

func Routes(uc usecase.IDataUsecase) *chi.Mux {
	r := chi.NewRouter()

	r.Use(helper.PanicRecovery)
	r.Route(`/data_insert`, func(r chi.Router) {
		r.Get(`/`, uc.HandleGetInserted)
		r.Post(`/`, uc.HandleInsert)
	})
	r.Route(`/data_deleted`, func(r chi.Router) {
		r.Get(`/`, uc.HandleGetDeleted)
		r.Delete(`/{name}`, uc.HandleDelete)
	})
	r.Route(`/data_updated`, func(r chi.Router) {
		r.Get(`/`, uc.HandleGetUpdated)
		r.Post(`/{name}`, uc.HandleUpdate)
	})
	r.Route(`/data_all`, func(r chi.Router) {
		r.Get(`/`, uc.HandleGetAll)
	})
	return r
}
