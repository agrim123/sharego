package main

import (
	"./Config"
	"./Router"
	"fmt"
	"net/http"
)

func main() {
	// Load Config file
	config := Config.ParseConfigFile()

	// Specify port
	port := config.Port
	if port == "" {
		panic("Please Specify port")
	}

	// Start server
	fmt.Println("Starting server on :" + port)
	http.ListenAndServe(":"+port, Router.Routes())
}
