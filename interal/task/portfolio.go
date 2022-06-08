package task

import (
	"fmt"
)

// MonitorPortfolio checks price in time intervals and sends email if needed
func MonitorPortfolio(bt *BinanceTask) error {
	portfolioValue, err := bt.BinanceClient.GetPortfolioTotalValue()
	bt.Counter++
	if err != nil {
		return err
	}

	if portfolioValue > 26987 {
		msg := fmt.Sprintf("Time to sell, portfolio value is: %v", portfolioValue)
		err := bt.sendNotification(msg, bt.Recipient, "Portfolio alert")
		if err != nil {
			return err
		}
		// msg every week 24h*7=168
	} else if bt.Counter%(24*7) == 0 {
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
