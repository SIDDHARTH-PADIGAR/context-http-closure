package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Create root context with cancel
	_, stop := context.WithCancel(context.Background())
	defer stop()

	// Channel to catch system interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Server setup
	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Handling request...")
			select {
			case <-time.After(5 * time.Second): // simulate long work
				fmt.Fprintln(w, "Request processed.")
			case <-r.Context().Done():
				http.Error(w, "Request cancelled", http.StatusRequestTimeout)
			}
		}),
	}

	// Shutdown goroutine
	go func() {
		<-sig
		fmt.Println("Interrupt received, shutting down...")
		stop()

		ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctxShutDown); err != nil {
			fmt.Println("Shutdown error:", err)
		}
	}()

	// Start server
	fmt.Println("Server listening on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}

	fmt.Println("Server stopped cleanly.")
}
