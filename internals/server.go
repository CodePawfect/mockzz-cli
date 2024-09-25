package internals

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Server *http.Server

func StartHttpServer(mux *http.ServeMux, port string) {
	Server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Channel to listen for system interrupt signals (e.g., Ctrl+C)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine to start the server
	go func() {
		fmt.Printf("\nHTTP Server started on port %s\n", port)
		fmt.Printf("Your mocks are ready to be used!")

		if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Block until we receive an interrupt signal
	<-stopChan

	// Stop listening for further signals
	signal.Stop(stopChan)

	fmt.Println("Shutting down the server...")

	// Create a context with a timeout to allow graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := Server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}

func StopHttpServer() {
	// Create a context with a timeout to allow graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := Server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}
