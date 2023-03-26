package process

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(shutdownCallback func(context.Context), dur time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	fmt.Println("gracefully shutdown process...")

	ctx, cancel := context.WithTimeout(context.Background(), dur*time.Second)
	defer cancel()
	defer signal.Stop(quit)

	go shutdownCallback(ctx)

	<-ctx.Done()

	fmt.Println("exiting process...")
}
