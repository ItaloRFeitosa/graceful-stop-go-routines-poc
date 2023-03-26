package main

import (
	"context"

	"github.com/italorfeitosa/graceful-stop-go-routines-poc/process"
	"github.com/italorfeitosa/graceful-stop-go-routines-poc/worker"
)

func main() {
	w := worker.New()

	w.Start()

	process.GracefulShutdown(func(ctx context.Context) {
		w.Stop()
	}, 10)
}
