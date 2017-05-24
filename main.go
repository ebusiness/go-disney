package main

import (
	"os"

	"github.com/ebusiness/go-disney/config"
	"github.com/ebusiness/go-disney/utils"
	_ "github.com/ebusiness/go-disney/v1"
)

func main() {
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
