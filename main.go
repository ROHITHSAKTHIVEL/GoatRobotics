package main

import (
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/config"
	routes "github.com/ROHITHSAKTHIVEL/GoatRobotics/routes"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

func init() {
	// Read configuration before starting
	config.ReadConfig()
}

func main() {
	// Open profiling output file
	outputFile, err := os.Create("profiling_output.log")
	if err != nil {
		log.Fatal("Could not create output file:", err)
	}
	defer outputFile.Close()

	// Initialize profiling and write to the output file
	initPprof(outputFile)

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

	// Capture memory profile after server shuts down
	captureMemoryProfile(outputFile)

	log.Println("Application is exiting...")
}

func initPprof(outputFile *os.File) {
	// Start the pprof HTTP server in a separate goroutine
	go func() {
		log.Println("Starting pprof server on localhost:6060")
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Fatal("Error starting pprof server:", err)
		}
	}()

	// Start CPU Profiling and write it to the output file
	pprof.StartCPUProfile(outputFile)
	defer pprof.StopCPUProfile()

	// Optionally, start memory profiling later in the program
	// (This will be triggered before the application exits)
}

func captureMemoryProfile(outputFile *os.File) {
	// Write memory heap profile into the output file
	if err := pprof.WriteHeapProfile(outputFile); err != nil {
		log.Fatal("Could not write memory profile:", err)
	}

	log.Println("Memory profile has been written to the output file")
}
