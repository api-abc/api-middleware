package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func getTime() string {
	timeNow := strings.Split(time.Now().String(), " ")[1]
	time := strings.Split(timeNow, ".")[0]
	return time
}

func (w *Worker) RunWorker() {
	/*
		INSERT 5 -> sleep 5 detik -> INSERT 5 -> delay 1 detik -> UPDATE 5 -> delay 1 detik -> UPDATE 5 ->
		sleep 5 detik -> INSERT 5 -> ...
		10.55 WORKER RUN
		10.00 INSERT 5
		10.05 INSERT 5 -> delay 1 detik
		10.06 UPDATE 5 (if not nil) -> delay 1 detik
		10.07 UPDATE 5 (if not nil)
	*/

	var wg sync.WaitGroup
	var count int
	var nameSlice int
	var check = make(chan domain.Data, 10)
	defer close(check)

	num_worker := w.di.GetConfig().Worker.NumWorker
	time_delay := w.di.GetConfig().Worker.QueryDelay

	fmt.Println("WORKER RUN IN\t\t", getTime())

	for {
		count++
		fmt.Println("Running Worker for", count)

		// Create Data
		for v := 0; v < 2; v++ {
			if nameSlice >= 35 {
				time.Sleep(time.Duration(time_delay) * time.Second)
				fmt.Println("NO INSERT...\t\t", getTime())
			} else {
				time.Sleep(time.Duration(time_delay) * time.Second)
				wg.Add(num_worker)
				insert := 0
				for i := 0; i < num_worker; i++ {
					go func() {
						err := assignInsert(GenerateData(nameSlice))
						if err != nil {
							fmt.Println(err)
						}
						wg.Done()
					}()
					insert++
					fmt.Println("Insert data", Names[nameSlice])
					nameSlice++
				}
				wg.Wait()
				fmt.Println("Insert", insert, "data", "\t\t", getTime())
			}
		}

		// Check Data for Update
		ticker := time.NewTicker(time.Duration(time_delay) * time.Second)
		channel := make(chan int)
		defer close(channel)
		go func() {
			for {
				select {
				case <-ticker.C:
					data := getAllUpdate()
					for _, v := range data {
						check <- v
					}
				case <-channel:
					ticker.Stop()
				}
			}
		}()

		// Update Data
		time.Sleep(1 * time.Second)
		fmt.Println("Check Update: ", len(check))
		if len(check) != 0 && len(check) <= num_worker {
			time.Sleep(1 * time.Second)
			length := len(check)
			wg.Add(length)
			for a := 0; a < length; a++ {
				go func(data domain.Data) {
					err := assignUpdate(data)
					if err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}(<-check)
			}
			wg.Wait()
			fmt.Println("Done Update: ", length, "\t", getTime())
		} else if len(check) != 0 && len(check) > num_worker {
			for i := 0; i < 2; i++ {
				wg.Add(5)
				for a := 0; a < 5; a++ {
					go func(data domain.Data) {
						err := assignUpdate(data)
						if err != nil {
							fmt.Println(err)
						}
						wg.Done()
					}(<-check)
				}
				wg.Wait()
				fmt.Println("Done Update: ", 5, "\t", getTime())
				if i == 0 {
					time.Sleep(1 * time.Second)
				}
			}
		} else {
			fmt.Println("NO UPDATE...\t\t", getTime())
		}

		fmt.Println("---------------------------------")

		if count >= 7 {
			if len(check) == 0 {
				break
			}
		}
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
