package config

import "AuthServer/utils"

type ServerArguments struct {
	Port string
}

func ParseArguments(args []string) *ServerArguments {
	temp := &ServerArguments{}
	index := utils.Find(args, "port")
	if index < 0 {
		temp.Port = ":8080"
	} else {
		temp.Port = ":" + args[index+1]
	}

	return temp
}
