package helper

import (
	"fmt"
	"net/http"

	"github.com/api-abc/internal-api-module/model/response"
)

func HandlePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				resp := response.BodyResponse{
					Status:  500,
					Message: "Internal Server Error",
				}
				WriteOutput(w, http.StatusInternalServerError, resp)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
