package binance

import (
	"fmt"
)

// GenerateMonitorSymbolTask return function tak can monitor one particular symbol, one task can monitor one symbol. You are free to set up multiple tasks
// to observe multiple symbols
func GenerateMonitorSymbolTask(symbol string, threshold float64, isLongPosition bool, runningNotificationInterval int) func(*Task) error {
	return func(bt *Task) error {
		sv, err := bt.BinanceClient.GetSymbolValue(symbol)
		bt.Counter++
		if err != nil {
			return err
		}
		if (isLongPosition && sv > threshold) || (!isLongPosition && sv < threshold) {
			msg := fmt.Sprintf("Time to sell, Eth value is: %v", sv)
			err := bt.sendNotification(msg, bt.Recipient, "Eth alert")
			if err != nil {
				return err
			}
			// msg every day - 6 times per hour * 24 = 144
		} else if runningNotificationInterval != 0 && bt.Counter%(runningNotificationInterval) == 0 {
			err := bt.sendNotification("Eth monitor is working, be patient", bt.Recipient, "Eth alert")
			if err != nil {
				return err
			}
		} else if bt.Counter == 1 {
			err := bt.sendNotification("Eth monitor is started", bt.Recipient, "Eth alert")
			if err != nil {
				return err
			}
		}

		return nil
	}

}
