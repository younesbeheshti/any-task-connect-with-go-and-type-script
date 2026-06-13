package main

import (
	"log"

	"github.com/younesbeheshti/any-task-connect/backend/internal/bootstrap"
)

// @title           TaskBridge API
// @version         1.0
// @description     Escrow-based marketplace backend API.
// @host            localhost:8080
// @BasePath        /
func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
