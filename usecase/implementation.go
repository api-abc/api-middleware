package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/api-middleware/services"
	"github.com/api-abc/internal-api-module/model/request"
	"github.com/api-abc/internal-api-module/model/response"
	"github.com/go-chi/chi"
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

func (uc *DataUsecase) HandleDelete(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	name := chi.URLParam(req, "name")

	result, err := uc.service.Delete(ctx, name)
	helper.HandlePanic(err)
	helper.WriteBodyHeader(writer, result)
}
func (uc *DataUsecase) HandleUpdate(writer http.ResponseWriter, req *http.Request) {
	var bodyReq request.UpdateRequest
	ctx := context.Background()
	name := chi.URLParam(req, "name")
	fmt.Println("UpdateCase - Get Param, Param:", name)

	fmt.Println("UpdateCase - Decode")
	err := json.NewDecoder(req.Body).Decode(&bodyReq)
	helper.HandlePanic(err)
	fmt.Println("UpdateCase - Decode Done, Request:", bodyReq)

	fmt.Println("UpdateCase - Process to Service")
	result, err := uc.service.Update(ctx, bodyReq, name)
	helper.HandlePanic(err)
	fmt.Println("UpdateCase - Process to Service Done")
	helper.WriteBodyHeader(writer, result)
}

func (uc *DataUsecase) HandleGetInserted(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	result := uc.service.GetInserted(ctx)

	helper.WriteBodyHeaderGet(writer, result)
}

func (uc *DataUsecase) HandleGetDeleted(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	result := uc.service.GetDeleted(ctx)

	helper.WriteBodyHeaderGet(writer, result)
}
func (uc *DataUsecase) HandleGetUpdated(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	result := uc.service.GetUpdated(ctx)

	helper.WriteBodyHeaderGet(writer, result)
}

func (uc *DataUsecase) HandleGetAll(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	var result response.BodyResponseGet
	var wg sync.WaitGroup

	wg.Add(3)
	go func(res *response.BodyResponseGet) {
		insert := uc.service.GetInserted(ctx)
		res.Data = append(res.Data, insert.Data...)
		wg.Done()
	}(&result)

	go func(res *response.BodyResponseGet) {
		delete := uc.service.GetDeleted(ctx)
		res.Data = append(res.Data, delete.Data...)
		wg.Done()
	}(&result)

	wg.Wait()
	result.Status = response.StatusOK
	result.Message = "Test Success"
	helper.WriteBodyHeaderGet(writer, result)
}
