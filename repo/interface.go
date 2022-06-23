package repo

import (
	"context"

	"github.com/api-abc/internal-api-module/model/request"
	"github.com/api-abc/internal-api-module/model/response"
)

type IDataRepo interface {
	Insert(ctx context.Context, req request.InsertRequest) (response.BodyResponse, error)
	Delete(ctx context.Context, name string) (response.BodyResponse, error)
	Update(ctx context.Context, req request.UpdateRequest, name string) (response.BodyResponse, error)
	GetInserted(ctx context.Context) (response.BodyResponseGet, error)
	GetDeleted(ctx context.Context) (response.BodyResponseGet, error)
	GetUpdated(ctx context.Context) (response.BodyResponseGet, error)
}
