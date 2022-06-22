package usecase

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/api-middleware/services"
	"github.com/api-abc/internal-api-module/model/request"
)

type DataUsecase struct {
	service services.IDataService
}

func NewDataUsecase(serv services.IDataService) IDataUsecase {
	return &DataUsecase{
		service: serv,
	}
}

func (uc *DataUsecase) HandleInsert(writer http.ResponseWriter, req *http.Request) {
	var bodyReq request.InsertRequest
	ctx := context.Background()

	err := json.NewDecoder(req.Body).Decode(&bodyReq)
	helper.HandlePanic(err)

	result, err := uc.service.Insert(ctx, bodyReq)
	helper.HandlePanic(err)

	helper.WriteBodyHeader(writer, result)
}

func (uc *DataUsecase) HandleDelete(writer http.ResponseWriter, req *http.Request) {}
func (uc *DataUsecase) HandleUpdate(writer http.ResponseWriter, req *http.Request) {}

func (uc *DataUsecase) HandleGetInserted(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	result := uc.service.GetInserted(ctx)

	helper.WriteBodyHeaderGet(writer, result)
}

func (uc *DataUsecase) HandleGetDeleted(writer http.ResponseWriter, req *http.Request) {}
func (uc *DataUsecase) HandleGetUpdated(writer http.ResponseWriter, req *http.Request) {}
