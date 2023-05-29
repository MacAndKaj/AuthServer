package main

import (
	"AuthServer/config"
	"AuthServer/server"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "[AuthServer]", log.LstdFlags)

	serverArguments := config.ParseArguments(os.Args)
	cfg := config.NewConfig()

	s := server.NewServer(serverArguments.Port, cfg, l)
	s.Run()
}
