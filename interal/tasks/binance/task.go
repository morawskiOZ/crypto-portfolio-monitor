package binance

import (
	"time"

	"github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/pkg/mail"
)

type executableTask func(bt *Task) error

// Task is a tasks that implements the tasker interface, so it can be scheduled
type Task struct {
	BinanceClient  *binance.Client
	MailClient     *mail.Client
	Recipient      string
	Counter        int
	TickerDuration time.Duration
	Task           executableTask
}

// RunTask runs the tasks
func (bt *Task) RunTask() error {
	return bt.Task(bt)
}

// SetupTicker returns a ticker that will run the tasks every 60 minutes
func (bt *Task) SetupTicker() *time.Ticker {
	return time.NewTicker(bt.TickerDuration)
}

func (bt *Task) sendNotification(msg string, recipient string, subject string) error {
	if err := bt.MailClient.Send(recipient, subject, msg); err != nil {
		return err
	}

	return nil
}
