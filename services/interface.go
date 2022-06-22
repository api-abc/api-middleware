package services

import (
	"context"

	"github.com/api-abc/internal-api-module/model/request"
	"github.com/api-abc/internal-api-module/model/response"
)

type IDataService interface {
	Insert(ctx context.Context, req request.InsertRequest) (response.BodyResponse, error)
	Delete(ctx context.Context, name string) (response.BodyResponse, error)
	Update(ctx context.Context, req request.UpdateRequest, name string) (response.BodyResponse, error)
	GetInserted(ctx context.Context) response.BodyResponseGet
	GetDeleted(ctx context.Context) response.BodyResponseGet
	Getupdated(ctx context.Context) response.BodyResponseGet
}
