package controllers

import (
	"AuthServer/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		return
	}

	var userId uint64

	if emailExists {
		if !h.db.VerifyEmail(req.UsernameOrEmail, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(createErrorResponse("Incorrect password"))
			return
		}

		userId, err = h.db.GetUserIdForLogin(req.UsernameOrEmail)

	} else if loginExists {
		if !h.db.VerifyLogin(req.UsernameOrEmail, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(createErrorResponse("Incorrect password"))
			return
		}

		userId, err = h.db.GetUserIdForLogin(req.UsernameOrEmail)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(createErrorResponse("Internal error (user)"))
		return
	}

	t := models.Token{
		UserId: int(userId),
		Hash:   generateNewToken(int(userId)),
		// ExpirationDate: time.,
		Permissions: "wr",
	}

	h.logger.Println("New token: ", t)

	err = h.db.AddNewTokenForUser(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(createErrorResponse("Internal error (token)"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateNewToken(id int) string {
	claims := jwt.RegisteredClaims{
		ID:        string(id),
		Issuer:    "SAS", // SimpleAuthServer
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate((time.Now().Add(500 * 24 * time.Hour))),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.Raw
}
