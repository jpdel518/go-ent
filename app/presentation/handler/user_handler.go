package handler

import (
	"encoding/json"
	"github.com/jpdel518/go-ent/domain/model"
	"github.com/jpdel518/go-ent/usecase"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Handler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *Handler {
	return &Handler{usecase}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	// get query parameters
	num, err := strconv.Atoi(r.URL.Query().Get("num"))
	if err != nil {
		log.Println(err)
		num = 10
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
	}

	// fetch user data
	users, err := h.usecase.Fetch(r.Context(), num)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3000, Data: err.Error()}))
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 2000, Data: users}))
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// get path parameters
	sub := strings.TrimPrefix(r.URL.Path, "/user")
	id, err := strconv.Atoi(filepath.Base(sub))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3101, Data: err.Error()}))
	}

	// fetch user data
	user, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3100, Data: err.Error()}))
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 2000, Data: user}))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if r.Header.Get("Content-Type") != "application/json" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")

	// get parameters
	var user *model.User
	var file multipart.File
	var fileHeader *multipart.FileHeader
	if r.Header.Get("Content-Type") != "application/json" {
		// multipart/form-data or application/x-www-form-urlencoded
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		age, _ := strconv.Atoi(r.FormValue("age")) // optionalなのでエラーは無視
		carFormValue := r.FormValue("car_ids")
		var carIDs []int
		if carFormValue != "" {
			carStrings := strings.Split(strings.Trim(strings.ReplaceAll(carFormValue, " ", ""), "[]"), ",")
			for _, carString := range carStrings {
				car, err := strconv.Atoi(carString)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3201, Data: err.Error()}))
					return
				}
				carIDs = append(carIDs, car)
			}
		}
		var err error
		file, fileHeader, err = r.FormFile("avatar")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3205, Data: err.Error()}))
			return
		}

		user = &model.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Age:       age,
			CarIDs:    carIDs,
		}
	} else {
		// application/json
		user = &model.User{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3202, Data: err.Error()}))
			return
		}
		err = json.Unmarshal(body, user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3203, Data: err.Error()}))
			return
		}
	}

	// validation
	err := user.Validate()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3204, Data: err.Error()}))
		return
	}

	// create user
	err = h.usecase.Create(r.Context(), user, file, fileHeader)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3200, Data: err.Error()}))
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 2000, Data: user}))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// get parameters
	var user *model.User
	var file multipart.File
	var fileHeader *multipart.FileHeader
	if r.Header.Get("Content-Type") != "application/json" {
		// multipart/form-data or application/x-www-form-urlencoded
		ID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3301, Data: err.Error()}))
			return
		}
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		age, _ := strconv.Atoi(r.FormValue("age")) // optionalなのでエラーは無視
		carStrings := strings.Split(strings.Trim(strings.ReplaceAll(r.FormValue("cars"), " ", ""), "[]"), ",")
		var carIDs []int
		for _, carString := range carStrings {
			car, err := strconv.Atoi(carString)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3302, Data: err.Error()}))
				return
			}
			carIDs = append(carIDs, car)
		}
		file, fileHeader, err = r.FormFile("avatar")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3303, Data: err.Error()}))
			return
		}

		user = &model.User{
			ID:        ID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Age:       age,
			CarIDs:    carIDs,
		}
	} else {
		// application/json
		user = &model.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3204, Data: err.Error()}))
			return
		}
	}

	// validation
	err := user.Validate()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3305, Data: err.Error()}))
		return
	}

	// update user
	err = h.usecase.Update(r.Context(), user, file, fileHeader)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3300, Data: err.Error()}))
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 2000, Data: user}))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// get parameters
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3401, Data: err.Error()}))
		return
	}

	// delete user
	err = h.usecase.Delete(r.Context(), id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 3400, Data: err.Error()}))
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(CreateResponseJson(&ApiRequestResponse{Code: 2000, Data: "success"}))
}
