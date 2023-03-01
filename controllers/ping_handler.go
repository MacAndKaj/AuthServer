package controllers

import (
	"log"
	"net/http"
)

type PingHandler struct {
	logger *log.Logger
}

func NewPingHandler(l *log.Logger) http.Handler {
	return &PingHandler{l}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Server ping from ", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
