package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/api-abc/api-middleware/configuration"
	"github.com/api-abc/api-middleware/helper"
	"github.com/api-abc/internal-api-module/model/domain"
	"github.com/api-abc/internal-api-module/model/response"
)

type Worker struct {
	di *configuration.DI
}

func New(conf *configuration.DI) *Worker {
	return &Worker{di: conf}
}

func (w *Worker) RunWorker() {
	var count int
	for {
		count++
		fmt.Println("Running Worker", count)
		var wg sync.WaitGroup

		//Check Data for Update
		check := getAllUpdate()
		if len(check) != 0 && len(check) <= w.di.GetConfig().Worker.NumWorker {
			wg.Add(len(check))
			for a := 0; a < len(check); a++ {
				go func(data domain.Data) {
					err := assignUpdate(data)
					if err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}(check[a])
			}
			wg.Wait()
		}

		//Create Data
		wg.Add(w.di.GetConfig().Worker.NumWorker)
		for i := 0; i < w.di.GetConfig().Worker.NumWorker; i++ {
			go func() {
				err := assignInsert(GenerateData())
				if err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		time.Sleep(time.Duration(w.di.GetConfig().Worker.QueryDelay) * time.Second)
	}
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
