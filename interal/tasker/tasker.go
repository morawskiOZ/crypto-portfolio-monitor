package tasker

import (
	"fmt"
	"sync"
	"time"
)

type client struct {
	wg              *sync.WaitGroup
	close           chan bool
	signalChannel   chan bool
	observableTasks []*observableTask
}

type observableTask struct {
	Task
	ticker *time.Ticker
}

type Task interface {
	RunTask(wg *sync.WaitGroup) error
	SetupTicker() *time.Ticker
}

type option func(*client)

func WithSignalChannel(sch chan bool) option {
	return func(pc *client) {
		pc.signalChannel = sch
	}
}

func NewTasker(options ...option) *client {
	c := &client{
		wg: &sync.WaitGroup{},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *client) scheduleTask(ot *observableTask) {
	// Run the task for the first time before the ticker signals
	c.wg.Add(1)
	if err := ot.RunTask(c.wg); err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	// Run the task every time the ticker signals or the close channel is closed
	for {
		select {
		case <-ot.ticker.C:
			c.wg.Add(1)
			if err := ot.RunTask(c.wg); err != nil {
				fmt.Printf("Error: %+v\n", err)
			}
		case <-c.close:
			fmt.Println("Error: close signal, all tasks cancelled")
			return
		}
	}

}

func (pc *client) newObservableTask(t Task) *observableTask {
	ot := observableTask{
		ticker: t.SetupTicker(),
		Task:   t,
	}

	return &ot
}
func (c *client) Run(tasks []Task) {
	close := make(chan bool, len(tasks))
	c.close = close

	for _, t := range tasks {
		ot := c.newObservableTask(t)
		go c.scheduleTask(ot)
		c.observableTasks = append(c.observableTasks, ot)
	}

	exit := <-c.signalChannel
	if exit {
		for _, ot := range c.observableTasks {
			ot.ticker.Stop()
			c.close <- true
		}
	}
	c.wg.Wait()
}
