package server

import (
	"AuthServer/config"
	"AuthServer/controllers"
	"AuthServer/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type serverContext struct {
	logger *log.Logger
	db     *models.UsersDatabase
	config *config.Config
}

type AuthServer struct {
	ctx    serverContext
	server *http.Server
}

func NewServer(p string, cfg *config.Config, l *log.Logger) *AuthServer {
	server := &http.Server{
		Addr:         p,
		Handler:      nil,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	return &AuthServer{
		ctx: serverContext{
			logger: l,
			db:     models.InitUsersDatabase(l),
			config: cfg,
		},
		server: server,
	}
}

func createRouter(ctx *serverContext) *mux.Router {
	router := mux.NewRouter()
	api_router := router.PathPrefix("/api/v1").Subrouter()

	getSubrouter := api_router.Methods("GET").Subrouter()
	getSubrouter.
		Path("/ping").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewPingHandler(ctx.logger)))

	getSubrouter.
		Path("/users/login").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewLoginHandler(ctx.logger, ctx.db)))

	// putSubrouter := router.Methods("PUT").Subrouter()

	postSubrouter := api_router.Methods("POST").Subrouter()
	postSubrouter.
		Path("/users/create").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewAddUserHandler(ctx.logger, ctx.db)))

	return router
}

func (s *AuthServer) Run() {
	s.ctx.logger.Println("Configuring server")

	router := createRouter(&s.ctx)
	s.server.Handler = router
	s.ctx.logger.Println("Server working on address: ", s.server.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	s.ctx.logger.Println("Received signal, graceful shutdown. Signal:", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.server.Shutdown(tc)
	if err != nil {
		return
	}

	err = s.ctx.db.Shutdown()
	if err != nil {
		return
	}

}
