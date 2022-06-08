package task

import (
	"sync"
	"time"

	"github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/pkg/mail"
)


type exetuableTask func(bt *BinanceTask) error

// BinanceTask is a task that implements the tasker interface so it can be scheduled
type BinanceTask struct {
	BinanceClient *binance.Client
	MailClient    *mail.Client
	Recipient     string
	Counter       int
	TickerDuration time.Duration
	Task          exetuableTask
}

// RunTask runs the task
func (bt *BinanceTask) RunTask(wg *sync.WaitGroup) error {
	// TODO move the wg.Done inside tasker
	defer wg.Done()

	return bt.Task(bt)
}

// SetupTicker returns a ticker that will run the task every 60 minutes
func (bt *BinanceTask) SetupTicker() *time.Ticker {
	return time.NewTicker(bt.TickerDuration)
}

func (bt *BinanceTask) sendNotification(msg string, recipient string, subject string) error {
	if err := bt.MailClient.Send(recipient, subject, msg); err != nil {
		return err
	}

	return nil
}
