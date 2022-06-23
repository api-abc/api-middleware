package usecase

import "net/http"

type IDataUsecase interface {
	HandleInsert(writer http.ResponseWriter, req *http.Request)
	HandleDelete(writer http.ResponseWriter, req *http.Request)
	HandleUpdate(writer http.ResponseWriter, req *http.Request)
	HandleGetInserted(writer http.ResponseWriter, req *http.Request)
	HandleGetDeleted(writer http.ResponseWriter, req *http.Request)
	HandleGetUpdated(writer http.ResponseWriter, req *http.Request)
	HandleGetAll(writer http.ResponseWriter, req *http.Request)
}
