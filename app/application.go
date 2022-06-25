package app

import (
	"net/http"

	"github.com/api-abc/api-middleware/configuration"
	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/api-middleware/repo"
	"github.com/api-abc/api-middleware/services"
	"github.com/api-abc/api-middleware/usecase"
	"github.com/api-abc/api-middleware/worker"
)

func Run(di *configuration.DI) {
	repo := repo.NewDataRepo(di)
	serv := services.NewDataService(repo)
	usecase := usecase.NewDataUsecase(serv)

	port := di.GetConfig().Host.Port
	go func() {
		err := http.ListenAndServe(port, Routes(usecase))
		helper.HandlePanic(err)
	}()
	w := worker.NewWorker(di)
	w.RunWorkerApp()

}
