package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"GPSG/handlers"
)

func main() {
	// Initialize Valkey DB and RabbitMQ connections
	err := handlers.InitValkey()
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	err = handlers.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Start the central ticker Goroutine
	go handlers.StartCentralTicker()

	// Define API routes
	http.HandleFunc("/start-pattern", handlers.StartPatternHandler)
	http.HandleFunc("/stop-pattern", handlers.StopPatternHandler)
	http.HandleFunc("/status", handlers.StatusHandler)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting GPSG server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")
	
	// Cleanup
	handlers.StopCentralTicker()
}
//curl -X POST http://localhost:8080/start-pattern -H "Content-Type: application/json" -d "{\"waveform\": \"sine\", \"time_script\": [100, 200, 300]}"