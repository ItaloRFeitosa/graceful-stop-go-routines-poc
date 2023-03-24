package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Worker struct {
	stopConsumeCh chan struct{}
	stopProcessCh chan struct{}
	processCh     chan struct{}
}

func (w *Worker) Consume() {
	for {
		select {
		case <-w.stopConsumeCh:
			fmt.Println("stoping consume\n.\n.")
			w.stopConsumeCh <- struct{}{}
			return
		case <-time.Tick(1 * time.Second):
			fmt.Println("consuming\n.\n.")
			w.processCh <- struct{}{}
		}
	}
}

func (w *Worker) Process() {
	var wg sync.WaitGroup
	for {
		select {
		case <-w.stopProcessCh:
			fmt.Println("stoping process\n.\n.")
			wg.Wait()
			w.stopProcessCh <- struct{}{}
			return
		case <-w.processCh:
			wg.Add(1)
			go func() {
				time.Sleep(5 * time.Second)
				fmt.Println("processed\n.\n.")
				wg.Done()
			}()
		}
	}

}

func (w *Worker) Stop() {
	w.stopConsumeCh <- struct{}{}
	<-w.stopConsumeCh
	fmt.Println("consumer stoped\n.\n.")
	w.stopProcessCh <- struct{}{}
	<-w.stopProcessCh
	fmt.Println("process stoped\n.\n.")
}

func main() {
	w := new(Worker)

	w.stopConsumeCh = make(chan struct{})
	w.stopProcessCh = make(chan struct{})
	w.processCh = make(chan struct{})

	go w.Consume()
	go w.Process()

	GracefulShutdown(func(ctx context.Context) {
		w.Stop()
	}, 10)
}

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
