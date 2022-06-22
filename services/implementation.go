package services

import (
	"context"

	"github.com/api-abc/api-middleware/repo"
	"github.com/api-abc/internal-api-module/model/request"
	"github.com/api-abc/internal-api-module/model/response"
)

type DataService struct {
	repo repo.IDataRepo
}

func NewDataService(datarepo repo.IDataRepo) IDataService {
	return &DataService{
		repo: datarepo,
	}
}

func (service *DataService) Insert(ctx context.Context, req request.InsertRequest) (response.BodyResponse, error) {
	result, err := service.repo.Insert(ctx, req)
	if err != nil {
		return response.BodyResponse{}, err
	}
	return result, nil
}

func (service *DataService) Delete(ctx context.Context, name string) (response.BodyResponse, error) {
	return response.BodyResponse{}, nil
}

func (service *DataService) Update(ctx context.Context, req request.UpdateRequest, name string) (response.BodyResponse, error) {
	return response.BodyResponse{}, nil
}

func (service *DataService) GetInserted(ctx context.Context) response.BodyResponseGet {
	result, err := service.repo.GetInserted(ctx)
	if err != nil {
		return response.BodyResponseGet{}
	}
	return result
}

func (service *DataService) GetDeleted(ctx context.Context) response.BodyResponseGet {
	return response.BodyResponseGet{}
}

func (service *DataService) Getupdated(ctx context.Context) response.BodyResponseGet {
	return response.BodyResponseGet{}
}
