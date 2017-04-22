package main

import (
	"os"

	"github.com/ebusiness/go-disney/config"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1"
)

func main() {
	// just touch Regist(), it will be auto load all `init` function of controllers's files [v1]
	v1.Regist()

	server := os.Getenv("HOST")
	if len(server) < 1 {
		server = config.HostName
	}
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = config.HTTPPort
	}
	utils.Route.RunTLS(server+":"+port, config.CertFilePath, config.KeyFilePath)
}
