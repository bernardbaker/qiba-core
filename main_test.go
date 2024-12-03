package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// For testing multiple instances
func StartServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if isPortInUse("0.0.0.0", port) {
		return fmt.Errorf("port %s is already in use", port)
	}

	go main()
	return nil
}

// Example test
func TestMultipleServerStarts(t *testing.T) {
	// First start should succeed
	err := StartServer()
	if err != nil {
		t.Fatalf("First server should start: %v", err)
	}

	// Give the server time to start
	time.Sleep(time.Second)

	// Second start should fail
	err = StartServer()
	if err == nil {
		t.Fatal("Second server should not start")
	}

	// Try with different port
	os.Setenv("PORT", "8081")
	err = StartServer()
	if err != nil {
		t.Fatalf("Server on different port should start: %v", err)
	}
}

func TestMainMultipleTimes(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
	}{
		{"First run", 1},
		{"Second run", 2},
		{"Third run", 3},
		{"Multiple runs", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.iterations; i++ {
				t.Logf("Iteration %d of %d", i+1, tt.iterations)
				StartServer()
			}
		})
	}
}

// If you need to test with delays between runs
func TestMainWithDelay(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
		delay      time.Duration
	}{
		{"Quick runs", 3, time.Millisecond * 100},
		{"Slower runs", 2, time.Second * 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.iterations; i++ {
				StartServer()
				time.Sleep(tt.delay)
			}
		})
	}
}
