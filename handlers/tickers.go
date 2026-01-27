// handlers/tickers.go
package handlers

import (
	"sync"
	"time"
)

var (
	centralTicker     *time.Ticker
	centralTickerChan = make(chan time.Time)
	tickerMutex       sync.RWMutex
	isTickerRunning   bool
)

// StartCentralTicker starts the central timing mechanism
// This is the function called from main.go
func StartCentralTicker() {
	tickerMutex.Lock()
	if isTickerRunning {
		tickerMutex.Unlock()
		return
	}
	isTickerRunning = true
	tickerMutex.Unlock()

	centralTicker = time.NewTicker(time.Millisecond) // 1ms base resolution
	go broadcastTicks()
}

// broadcastTicks handles the broadcasting of ticker events
func broadcastTicks() {
	for t := range centralTicker.C {
		select {
		case centralTickerChan <- t:
			// Successfully sent tick
		default:
			// Channel full, skip this tick
		}
	}
}

// StopCentralTicker stops the central ticker
func StopCentralTicker() {
	tickerMutex.Lock()
	defer tickerMutex.Unlock()

	if centralTicker != nil {
		centralTicker.Stop()
		centralTicker = nil
	}
	isTickerRunning = false
}

// GetTickerChan returns the channel for subscribing to ticks
func GetTickerChan() <-chan time.Time {
	return centralTickerChan
}

// IsTickerRunning returns the current state of the ticker
func IsTickerRunning() bool {
	tickerMutex.RLock()
	defer tickerMutex.RUnlock()
	return isTickerRunning
}