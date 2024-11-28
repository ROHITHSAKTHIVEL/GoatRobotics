package main

import (
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	apihandlers "github.com/ROHITHSAKTHIVEL/GoatRobotics/apiHandlers"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/config"
	"github.com/ROHITHSAKTHIVEL/GoatRobotics/service"
)

func init() {
	config.ReadConfig()
}

func main() {

	outputFile, err := os.Create("profiling_output.log")
	if err != nil {
		log.Fatal("Could not create output file:", err)
	}
	defer outputFile.Close()

	// Initialize profiling and write to the output file
	initPprof(outputFile)

	clients := service.NewClients()

	go clients.Listener()

	handler := apihandlers.NewChatHandler(clients)

	commonMiddleware := apihandlers.NewMiddleware(10, 20) // TODO set new value in config for rate limit value

	http.HandleFunc("/join", commonMiddleware.MiddleWare(handler.JoinHandler))
	http.HandleFunc("/leave", commonMiddleware.MiddleWare(handler.LeaveHandler))
	http.HandleFunc("/send-message", commonMiddleware.MiddleWare(handler.SendMessageHandler))
	http.HandleFunc("/ping", apihandlers.Ping)

	log.Printf("Starting server on %s", config.Port)
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

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
