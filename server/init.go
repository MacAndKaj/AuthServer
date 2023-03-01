package server

import (
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
}

type AuthServer struct {
	ctx    serverContext
	server *http.Server
}

func NewServer(p string, l *log.Logger) *AuthServer {
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
		},
		server: server,
	}
}

func createRouter(ctx *serverContext) *mux.Router {
	router := mux.NewRouter()
	getSubrouter := router.Methods("GET").Subrouter()
	getSubrouter.
		Path("/ping").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewPingHandler(ctx.logger)))

	getSubrouter.
		Path("/users/login").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewLoginHandler(ctx.logger, ctx.db)))

	putSubrouter := router.Methods("PUT").Subrouter()
	putSubrouter.
		Path("/users/create").
		Handler(handlers.LoggingHandler(os.Stdout, controllers.NewAddUserHandler(ctx.logger, ctx.db)))

	// postRouter := router.Methods("POST").Subrouter()
	// postRouter.
	// 	Path("/users/create").
	// 	Handler(controllers.NewAddUserHandler())

	return router
}

func (s *AuthServer) Run() {
	s.ctx.logger.Println("Configuring server")
	// s.ctx.logger.Println("Addr:", s.server.Addr)

	router := createRouter(&s.ctx)
	s.server.Handler = router

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
