package util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func AwaitSignal(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalChan:
		fmt.Println("Shutting down..................")
	case <-ctx.Done():
		return
	}
}
