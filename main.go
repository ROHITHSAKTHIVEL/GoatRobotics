package main

import (
	"log"
	"net/http"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/config"
	routes "github.com/ROHITHSAKTHIVEL/GoatRobotics/routes"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

func init() {
	// Read configuration before starting
	config.ReadConfig()
}

func main() {
	// Initialize Clients
	clients := service.NewClients()

	// Start listener in a goroutine
	go clients.Listener()

	// Setup routes using the new routes package
	handler := routes.SetupRoutes(clients)

	// Start the server with the handler
	log.Printf("Starting server on %s", config.Port)
	if err := http.ListenAndServe(config.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Application is exiting...")
}
