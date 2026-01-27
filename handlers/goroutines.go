package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GoroutineMeta struct {
	ID       string
	Waveform string
	Ticker   *time.Ticker
	StopChan chan bool
	Active   bool
}

var (
	goroutineMap = make(map[string]*GoroutineMeta)
	mu           sync.Mutex
)

// StartPatternHandler: Endpoint to start a new Goroutine
func StartPatternHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data struct {
		Waveform   string `json:"waveform"`
		TimeScript []int  `json:"time_script"`
	}
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a unique ID
	id := uuid.NewString()
	stopChan := make(chan bool)

	// Start Goroutine
	go startPatternGoroutine(id, data.Waveform, data.TimeScript, stopChan)

	// Save metadata
	mu.Lock()
	goroutineMap[id] = &GoroutineMeta{
		ID:       id,
		Waveform: data.Waveform,
		StopChan: stopChan,
		Active:   true,
	}
	mu.Unlock()

	// Respond with the Goroutine ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

// startPatternGoroutine: Function to manage a single Goroutine
// handlers/goroutines.go
// Update the startPatternGoroutine function

func startPatternGoroutine(id, waveform string, timeScript []int, stopChan chan bool) {
	log.Printf("Starting Goroutine %s with waveform: %s", id, waveform)
	
	currentIndex := 0
	nextTickTime := time.Now().Add(time.Duration(timeScript[0]) * time.Millisecond)
	
	// Get the central ticker channel
	tickerChan := GetTickerChan()

	for {
		select {
		case t := <-tickerChan:
			if t.After(nextTickTime) || t.Equal(nextTickTime) {
				// Publish waveform to RabbitMQ
				if err := PublishToRabbitMQ(waveform); err != nil {
					log.Printf("Failed to publish waveform: %v", err)
					continue
				}

				// Update index and next tick time
				currentIndex = (currentIndex + 1) % len(timeScript)
				nextTickTime = t.Add(time.Duration(timeScript[currentIndex]) * time.Millisecond)
			}
		case <-stopChan:
			log.Printf("Stopping Goroutine %s", id)
			return
		}
	}
}

// StopPatternHandler: Endpoint to stop a Goroutine
func StopPatternHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if meta, exists := goroutineMap[id]; exists {
		close(meta.StopChan)
		delete(goroutineMap, id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Stopped Goroutine " + id))
	} else {
		http.Error(w, "Goroutine not found", http.StatusNotFound)
	}
}

// StatusHandler: Endpoint to get the status of all active Goroutines
// Update this function in handlers/goroutines.go
// In handlers/goroutines.go - update this function
func StatusHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Status endpoint hit")  // Debug log

    if r.Method != http.MethodGet {
        log.Println("Method not allowed:", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    // Create a simplified response
    activePatterns := make(map[string]string)
    for id, meta := range goroutineMap {
        activePatterns[id] = meta.Waveform
    }

    response := struct {
        Status          string            `json:"status"`
        ServerTime      string            `json:"server_time"`
        ActivePatterns  map[string]string `json:"active_patterns"`
        PatternCount    int               `json:"pattern_count"`
    }{
        Status:         "running",
        ServerTime:     time.Now().Format(time.RFC3339),
        ActivePatterns: activePatterns,
        PatternCount:   len(goroutineMap),
    }

    w.Header().Set("Content-Type", "application/json")
    
    err := json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err) // Debug log
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Println("Status response sent successfully") // Debug log
}
