package controllers

import (
	"AuthServer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AddUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type AddUserHandler struct {
	logger *log.Logger
	db     *models.UsersDatabase
}

func NewAddUserHandler(l *log.Logger, d *models.UsersDatabase) http.Handler {
	return &AddUserHandler{l, d}
}

func (h *AddUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req := &AddUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(createErrorResponse("JSON decoding failed"))
		h.logger.Println("JSON decoding error: ", err)
		return
	}

	statusCode, err := h.addUserToDb(req)
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write(createErrorResponse(err.Error()))
	}
}

func (h *AddUserHandler) addUserToDb(req *AddUserRequest) (int, error) {
	loginExists, err1 := h.db.LoginExists(req.Username)
	emailExists, err2 := h.db.EmailExists(req.Email)

	if err1 != nil || err2 != nil {
		return http.StatusInternalServerError, fmt.Errorf("Error during checking request validity")
	}

	if loginExists {
		return http.StatusBadRequest, fmt.Errorf("Username already in use")
	}

	if emailExists {
		return http.StatusBadRequest, fmt.Errorf("Email already in use")
	}

	h.db.AddNewUser(models.User{
		Nickname:     req.Username,
		FirstName:    req.Firstname,
		Email:        req.Email,
		Password:     req.Password,
		CreationDate: time.Now().Format(time.DateTime),
	})

	return http.StatusOK, nil
}
