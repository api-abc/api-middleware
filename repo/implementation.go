package repo

import (
	"context"

	"github.com/api-abc/api-middleware/configuration"
	"github.com/api-abc/internal-api-module/model/request"
	"github.com/api-abc/internal-api-module/model/response"
)

type DataRepo struct {
	config *configuration.DI
}

func NewDataRepo(cfg *configuration.DI) IDataRepo {
	return &DataRepo{
		config: cfg,
	}
}

func (repo *DataRepo) Insert(ctx context.Context, req request.InsertRequest) (response.BodyResponse, error) {
	client := repo.config.GetClientInsert()

	result, err := client.Insert(context.Background(), req)
	if err != nil {
		return response.BodyResponse{}, err
	}
	return result, nil
}

func (repo *DataRepo) Delete(ctx context.Context) (response.BodyResponse, error) {
	return response.BodyResponse{}, nil
}

func (repo *DataRepo) Update(ctx context.Context, req request.InsertRequest) (response.BodyResponse, error) {
	return response.BodyResponse{}, nil
}

func (repo *DataRepo) GetInserted(ctx context.Context) (response.BodyResponseGet, error) {
	client := repo.config.GetClientInsert()

	result, err := client.GetInserted(context.Background())
	if err != nil {
		return response.BodyResponseGet{}, err
	}
	return result, nil
}

func (repo *DataRepo) GetDeleted(ctx context.Context) (response.BodyResponseGet, error) {
	return response.BodyResponseGet{}, nil
}

func (repo *DataRepo) GetUpdated(ctx context.Context) (response.BodyResponseGet, error) {
	return response.BodyResponseGet{}, nil
}
