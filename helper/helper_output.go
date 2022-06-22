package helper

import (
	"encoding/json"
	"net/http"

	"github.com/api-abc/internal-api-module/model/response"
)

func WriteOutput(writer http.ResponseWriter, code int, resp interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	output, _ := json.Marshal(resp)
	writer.Write([]byte(output))
}

func WriteBodyHeader(writer http.ResponseWriter, result response.BodyResponse) {
	if result.Data.Name == "" {
		result = response.BodyResponse{
			Status:  result.Status,
			Message: result.Message,
			// Data:    nil, //
		}
	}
	status := result.Status
	selectWrite(writer, status, result)
}

func WriteBodyHeaderGet(writer http.ResponseWriter, result response.BodyResponseGet) {
	if len(result.Data) == 0 {
		result = response.BodyResponseGet{
			Status:  result.Status,
			Message: result.Message,
			Data:    nil,
		}
	}
	status := result.Status
	selectWrite(writer, status, result)
}

func selectWrite(writer http.ResponseWriter, status int, result interface{}) {
	switch status {
	case 1:
		WriteOutput(writer, http.StatusOK, result)
	case 2:
		WriteOutput(writer, http.StatusCreated, result)
	case 3:
		WriteOutput(writer, http.StatusBadRequest, result)
	case 4:
		WriteOutput(writer, http.StatusNotFound, result)
	case 5:
		WriteOutput(writer, http.StatusInternalServerError, result)
	}
}
