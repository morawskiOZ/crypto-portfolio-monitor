package binance

import (
	"fmt"
)

// GenerateMonitorPortfolioTask return function which is executable tasks for portfolio monitoring
func GenerateMonitorPortfolioTask(threshold float64, runningNotificationInterval int) func(*Task) error {
	return func(bt *Task) error {
		portfolioValue, err := bt.BinanceClient.GetPortfolioTotalValue()
		bt.Counter++
		if err != nil {
			return err
		}
		if portfolioValue > threshold {
			msg := fmt.Sprintf("Time to sell, portfolio value is: %v", portfolioValue)
			err := bt.sendNotification(msg, bt.Recipient, "Portfolio alert")
			if err != nil {
				return err
			}
			// msg every week 24h*7=168
		} else if runningNotificationInterval != 0 && bt.Counter%(runningNotificationInterval) == 0 {
			err := bt.sendNotification("Crypto monitor is working, be patient", bt.Recipient, "Portfolio alert")
			if err != nil {
				return err
			}
		} else if bt.Counter == 1 {
			err := bt.sendNotification("Crypto monitor is started", bt.Recipient, "Portfolio alert")
			if err != nil {
				return err
			}
		}

		return nil
	}
}
