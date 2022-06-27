package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/api-abc/api-middleware/configuration"
	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/internal-api-module/model/domain"
	"github.com/api-abc/internal-api-module/model/response"
)

var nameSlice int = 0
var check = make(chan domain.Data, 10)

type WorkerApp struct {
	di *configuration.DI
}

func NewWorker(conf *configuration.DI) *WorkerApp {
	return &WorkerApp{di: conf}
}

func (wa *WorkerApp) RunWorkerApp() {
	defer close(check)
	execFn := func(s string) (v interface{}, e error) {
		switch s {
		case "insert":
			v, e = wa.insert()
			if e != nil {
				return nil, e
			}
			return v, nil
		case "update":
			v, e = wa.update()
			if e != nil {
				return nil, e
			}
			return v, nil
		}
		return nil, nil
	}
	second := 0

	for {
		if nameSlice >= 35 {
			break
		}
		fmt.Println(second, time.Now().Minute(), time.Now().Second())
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
			data := getAllUpdate()
			for _, v := range data {
				check <- v
			}
			for i := 0; i < w.workersCount; i++ {
				jobs = append(jobs, Job{
					Descriptor: JobDescriptor{ID: 1, Type: "Insert"},
					ExecFn:     execFn,
					Args:       "insert",
				})
			}
		case 15:
			data := getAllUpdate()
			for _, v := range data {
				check <- v
			}
			fmt.Println("I'm on app.go case 15")
			for i := 0; i < 1; i++ {
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

		time.Sleep(850 * time.Millisecond)
		if second == 15 {
			second = 5
		} else if len(jobs) != 0 {
			if jobs[0].Args == "update" {
				second++
			}

		}
		second++
	}

}

func (wa *WorkerApp) insert() (response.BodyResponseGet, error) {
	err := assignInsert(GenerateData(nameSlice))
	if err != nil {
		return response.BodyResponseGet{}, err
	}
	return response.BodyResponseGet{Status: response.StatusOK, Message: "Success"}, nil
}

func (wa *WorkerApp) update() (response.BodyResponseGet, error) {
	var wg sync.WaitGroup
	fmt.Println("I'm on app.go update")
	fmt.Println("I'm on app.go len checking ", len(check))
	if len(check) == 0 {
		for b := 0; b < wa.di.GetConfig().Worker.NumWorker; b++ {
			wg.Add(1)
			go func() {
				_, err := wa.insert()
				if err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		return response.BodyResponseGet{}, errors.New("no update")
	} else if len(check) <= 5 {
		for a := 0; a < len(check); a++ {
			go func(data domain.Data) {
				err := assignUpdate(data)
				if err != nil {
					fmt.Println(err)
				}
			}(<-check)
		}
	} else {
		length := len(check)
		for i := 0; i < 2; i++ {
			var counData int
			if length > 5 {
				counData = 5
			} else {
				counData = length
			}
			wg.Add(counData)
			for a := 0; a < counData; a++ {
				go func(data domain.Data) {
					err := assignUpdate(data)
					if err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}(<-check)
				length--
			}
			wg.Wait()
			if i == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}
	fmt.Println("I'm on app.go len checking ", len(check))
	time.Sleep(1 * time.Second)
	for b := 0; b < 5; b++ {
		go func() {
			_, err := wa.insert()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	return response.BodyResponseGet{Status: response.StatusOK, Message: "Success"}, nil
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
