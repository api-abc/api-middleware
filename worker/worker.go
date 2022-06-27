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
	var nameSlice int
	var check = make(chan domain.Data, 10)
	defer close(check)

	num_worker := w.di.GetConfig().Worker.NumWorker
	time_delay := w.di.GetConfig().Worker.QueryDelay

	count := 0
	updCount := 0

	second := time.NewTicker(997 * time.Millisecond)
	chSecond := make(chan int)
	defer close(chSecond)
	func() {
		for {
			select {
			case <-second.C:
				count++
				switch count {
				case time_delay:
					data := getAllUpdate()
					for _, v := range data {
						check <- v
					}
					fmt.Println("Check Update: ", len(check), "\t", getTime())
					// Insert Data
					if nameSlice >= 34 {
						fmt.Println("NO INSERT...\t\t", getTime())
					} else {
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
				case 10:
					// Update Data
					if len(check) != 0 && len(check) <= num_worker {
						updCount = 0
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
						updCount = 0
						length := len(check)
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

						length = length - 5
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
					} else {
						updCount++
						fmt.Println("NO UPDATE...\t\t", getTime())
						if updCount > 3 {
							chSecond <- 0
							break
						}
					}
					fmt.Println("---------------------------------")
				case 11:
					count = 0
					data := getAllUpdate()
					for _, v := range data {
						check <- v
					}
					fmt.Println("Check Update: ", len(check), "\t", getTime())
					// Insert Data
					if nameSlice >= 34 {
						fmt.Println("NO INSERT...\t\t", getTime())
					} else {
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
			case <-chSecond:
				second.Stop()
			}
		}
	}()
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
