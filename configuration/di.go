package configuration

import (
	"github.com/api-abc/internal-api-module/data/delete"
	"github.com/api-abc/internal-api-module/data/insert"
	"github.com/api-abc/internal-api-module/data/update"
	"github.com/api-abc/internal-api-module/rest"
)

type DI struct {
	config        Config
	client_insert *insert.Client
	client_delete *delete.Client
	client_update *update.Client
}

func NewDI(cfg Config) *DI {
	di := &DI{config: cfg}
	return di
}

func (di *DI) GetConfig() Config {
	return di.config
}

func (di *DI) GetClientInsert() *insert.Client {
	if di.client_insert == nil {
		client := rest.New(di.config.ApiInternal.UrlInserted)
		client_data := insert.New(client)
		di.client_insert = client_data
	}
	return di.client_insert
}

func (di *DI) GetClientDelete() *delete.Client {
	if di.client_delete == nil {
		client := rest.New(di.config.ApiInternal.UrlDeleted)
		client_data := delete.New(client)
		di.client_delete = client_data
	}
	return di.client_delete
}

func (di *DI) GetClientUpdate() *update.Client {
	if di.client_update == nil {
		client := rest.New(di.config.ApiInternal.UrlUpdated)
		client_data := update.New(client)
		di.client_update = client_data
	}
	return di.client_update
}
