package task

import (
	"fmt"
)

// MonitorSymbol checks price in time intervals and sends email if needed
func MonitorSymbol(bt *BinanceTask) error {
	sv, err := bt.BinanceClient.GetSymbolValue("ETHUSDT")
	bt.Counter++
	if err != nil {
		return err
	}

	if sv > 2150 {
		msg := fmt.Sprintf("Time to sell, Eth value is: %v", sv)
		err := bt.sendNotification(msg, bt.Recipient, "Eth alert")
		if err != nil {
			return err
		}
		// msg every day - 6 times per hour * 24 = 144
	} else if bt.Counter%(24*6) == 0 {
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
