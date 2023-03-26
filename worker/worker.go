package worker

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	stopConsumeCh chan struct{}
	stopProcessCh chan struct{}
	processCh     chan struct{}
}

func New() *Worker {
	w := new(Worker)

	w.stopConsumeCh = make(chan struct{})
	w.stopProcessCh = make(chan struct{})
	w.processCh = make(chan struct{})

	return w
}
func (w *Worker) Start() {
	go w.consume()
	go w.process()
}

func (w *Worker) Stop() {
	w.stopConsumeCh <- struct{}{}
	<-w.stopConsumeCh
	fmt.Println("consumer stoped\n.\n.")
	w.stopProcessCh <- struct{}{}
	<-w.stopProcessCh
	fmt.Println("process stoped\n.\n.")
}

func (w *Worker) consume() {
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

func (w *Worker) process() {
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
