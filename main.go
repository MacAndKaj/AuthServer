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

	s := server.NewServer(serverArguments.Port, l)
	s.Run()
}
