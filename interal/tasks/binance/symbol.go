package binance

import (
	"fmt"
)

// GenerateMonitorSymbolTask return function tak can monitor one particular symbol, one task can monitor one symbol. You are free to set up multiple tasks
// to observe multiple symbols
// Setting runningNotificationInterval to 0 will disable running periodic notification, only start and trigger notification will be sent
func GenerateMonitorSymbolTask(symbol string, threshold float64, isLongPosition bool, runningNotificationInterval int) func(*Task) error {
	return func(bt *Task) error {
		sv, err := bt.BinanceClient.GetSymbolValue(symbol)
		bt.Counter++
		if err != nil {
			return err
		}
		action := "buy"
		if !isLongPosition {
			action = "sell"
		}
		if (isLongPosition && sv > threshold) || (!isLongPosition && sv < threshold) {
			msg := fmt.Sprintf("Time to %s, %s value is: %v", action, symbol, sv)
			err := bt.sendNotification(msg, bt.Recipient, "Eth alert")
			if err != nil {
				return err
			}
		} else if runningNotificationInterval != 0 && bt.Counter%(runningNotificationInterval) == 0 {
			err := bt.sendNotification(fmt.Sprintf("%s monitor is working, be patient", symbol), bt.Recipient, fmt.Sprintf("%s monitor", symbol))
			if err != nil {
				return err
			}
		} else if bt.Counter == 1 {
			err := bt.sendNotification(fmt.Sprintf("%s monitor is working, be patient", symbol), bt.Recipient, fmt.Sprintf("%s monitor", symbol))
			if err != nil {
				return err
			}
		}

		return nil
	}

}
