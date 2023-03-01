package controllers

import (
	"AuthServer/models"
	"encoding/json"
	"log"
	"net/http"
)

type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	SessionToken   string `json:"sessionToken"`
	ExpirationDate string `json:"expirationDate"`
}

type LoginHandler struct {
	logger *log.Logger
	db     *models.UsersDatabase
}

func NewLoginHandler(l *log.Logger, d *models.UsersDatabase) http.Handler {
	return &LoginHandler{l, d}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(createErrorResponse("JSON decoding failed"))
		h.logger.Println("JSON decoding error: ", err)
		return
	}

	loginExists, err1 := h.db.LoginExists(req.UsernameOrEmail)
	emailExists, err2 := h.db.EmailExists(req.UsernameOrEmail)

	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(createErrorResponse("Error during checking request validity"))
	}

	if emailExists {
		if !h.db.VerifyEmail(req.UsernameOrEmail, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(createErrorResponse("Incorrect password"))
		}
	}

	if loginExists {
		if !h.db.VerifyLogin(req.UsernameOrEmail, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(createErrorResponse("Incorrect password"))
		}
	}

	// h.db.CreateTokenFor(req.UsernameOrEmail)

	w.WriteHeader(http.StatusOK)
}
