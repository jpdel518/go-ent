package handler

import (
	"encoding/json"
	"github.com/jpdel518/go-ent/usecase"
	"log"
	"net/http"
)

func NewHandler(userUsecase usecase.UserUsecase) {
	userHandler := NewUserHandler(userUsecase)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World"))
	})
	http.HandleFunc("/user/fetch", userHandler.Fetch)
	http.HandleFunc("/user/get-by-id/", userHandler.GetById)
	http.HandleFunc("/user/update", userHandler.Update)
	http.HandleFunc("/user/create", userHandler.Create)
	http.HandleFunc("/user/delete", userHandler.Delete)

	http.ListenAndServe(":8080", nil)
}

// ApiRequestResponse response json
type ApiRequestResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// CreateResponseJson create response data as json
func CreateResponseJson(a *ApiRequestResponse) []byte {
	js, err := json.Marshal(a)
	if err != nil {
		log.Fatalf("create response json error: %v", err)
	}
	return js
}
