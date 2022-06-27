package worker

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/internal-api-module/model/domain"
	"github.com/api-abc/internal-api-module/model/request"
)

func GenerateData(nameSlice int) domain.Data {
	return domain.Data{
		Name: Names[nameSlice],
		Age:  (rand.Intn(60) + 1),
		JobDetails: domain.Job{
			Position:            Positions[rand.Intn(len(Positions)-1)],
			YearsWorkExperience: (rand.Intn(9) + 1),
			WorkStatus:          "Work",
		},
	}
}

func CreateInsertRequest(data domain.Data) *http.Request {
	newRequest := request.InsertRequest{
		Name:       data.Name,
		Age:        data.Age,
		JobDetails: data.JobDetails,
	}
	marshal, err := json.Marshal(newRequest)
	helper.HandlePanic(err)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8090/data_insert/", bytes.NewBuffer(marshal))
	helper.HandlePanic(err)
	return request
}

func CreateUpdateRequest(data domain.Data) *http.Request {
	newRequest := request.UpdateRequest{
		Age: 62,
		JobDetails: domain.Job{
			Position:            Positions[rand.Intn(len(Positions)-1)],
			YearsWorkExperience: rand.Intn(9) + 1,
			WorkStatus:          "Retired",
		},
	}
	marshal, err := json.Marshal(newRequest)
	helper.HandlePanic(err)
	request, err := http.NewRequest(http.MethodPut, "http://localhost:8090/data_updated/"+data.Name, bytes.NewBuffer(marshal))
	helper.HandlePanic(err)
	return request
}
