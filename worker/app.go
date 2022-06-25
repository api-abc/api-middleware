package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/api-abc/api-middleware/configuration"
	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/internal-api-module/model/domain"
	"github.com/api-abc/internal-api-module/model/response"
)

type WorkerApp struct {
	di *configuration.DI
}

func NewWorker(conf *configuration.DI) *WorkerApp {
	return &WorkerApp{di: conf}
}

func (wa *WorkerApp) RunWorkerApp() {
	execFn := func(s string) (v interface{}, e error) {
		switch s {
		case "insert":
			v, e = insert()
			if e != nil {
				return nil, e
			}
			return v, nil
		case "update":
			v, e = update()
			if e != nil {
				return nil, e
			}
			return v, nil
		}
		return nil, nil
	}
	second := 0
	for {
		w := New(wa.di.GetConfig().Worker.NumWorker)
		t := "insert"
		var jobs []Job
		switch second {
		case 5:
			for i := 0; i < w.workersCount; i++ {
				jobs = append(jobs, Job{
					Descriptor: JobDescriptor{ID: 1, Type: "Insert"},
					ExecFn:     execFn,
					Args:       "insert",
				})
			}
		case 10:
			for i := 0; i < w.workersCount; i++ {
				jobs = append(jobs, Job{
					Descriptor: JobDescriptor{ID: 1, Type: "Update"},
					ExecFn:     execFn,
					Args:       "update",
				})
			}
		}
		w.GenerateFrom(jobs)
		if len(jobs) != 0 {
			t = jobs[0].Args
		}
		w.Run(context.Background(), t)

		time.Sleep(1 * time.Second)
		if second == 10 {
			second = 0
		} else if len(jobs) != 0 {
			if jobs[0].Args == "update" {
				second++
			}

		}
		second++
	}

}

func insert() (response.BodyResponseGet, error) {
	err := assignInsert(GenerateData())
	if err != nil {
		return response.BodyResponseGet{}, err
	}
	return response.BodyResponseGet{Status: response.StatusOK, Message: "Success"}, nil
}

func update() (response.BodyResponseGet, error) {
	check := getAllUpdate()
	if len(check) == 0 {
		return response.BodyResponseGet{}, errors.New("no update")
	} else if len(check) <= 5 {
		for a := 0; a < len(check); a++ {
			go func(data domain.Data) {
				err := assignUpdate(data)
				if err != nil {
					fmt.Println(err)
				}
			}(check[a])
		}
	} else {
		for a := 0; a < 5; a++ {
			go func(data domain.Data) {
				err := assignUpdate(data)
				if err != nil {
					fmt.Println(err)
				}
			}(check[a])
		}
	}
	time.Sleep(1 * time.Second)
	v, _ := insert()
	return v, nil
}

func assignInsert(data domain.Data) error {
	var client http.Client
	request := CreateInsertRequest(data)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err := client.Do(request)
	helper.HandlePanic(err)
	return nil
}

func assignUpdate(data domain.Data) error {
	var client http.Client
	request := CreateUpdateRequest(data)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err := client.Do(request)
	helper.HandlePanic(err)
	return nil
}

func getAllUpdate() []domain.Data {
	var respBody response.BodyResponseGet
	resp, err := http.Get("http://localhost:8090/data_updated")
	helper.HandlePanic(err)
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	helper.HandlePanic(err)
	return respBody.Data
}
