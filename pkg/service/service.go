package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func RunWithGracefulShutdown(ctx context.Context, runFunc func(ctx context.Context)) {
	// Create a channel to receive OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure cancellation on function exit

	// Run the provided function in a goroutine
	go func() {
		runFunc(ctx)
	}()

	// Wait for an OS signal
	<-stop

	// Cancel the context to signal shutdown
	cancel()

	// Optionally, wait for the function to complete if needed
	// You may want to add a timeout or additional waiting logic here
	// to ensure the runFunc completes gracefully.

	fmt.Println("Shutdown signal received, exiting...")
}
