package main

import (
	"fmt"
	"github.com/morawskioz/binance-monitor/interal/tasks/binance"
	"os"
	"time"

	"github.com/morawskioz/binance-monitor/configs"
	binanceAPI "github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/interal/signal"
	"github.com/morawskioz/binance-monitor/interal/tasker"
	"github.com/morawskioz/binance-monitor/pkg/mail"
)

func main() {
	config, err := configs.LoadEnvConfig()
	if err != nil {
		fmt.Printf("Fatal error: %+v\n", err)
		os.Exit(1)
	}

	mc := mail.NewMailClient(
		mail.WithDialer(mail.AuthConfig{
			Port:         config.SMTPPort,
			Host:         config.SMTPHost,
			Password:     config.EmailPass,
			EmailAddress: config.EmailLogin,
		}),
	)

	so := signal.NewOsSignalObserver()
	go so.Observe()

	binanceAPI.WithTestFlag()
	credentials := binanceAPI.Credentials{
		Key:    config.Key,
		Secret: config.Secret,
	}

	bc := binanceAPI.NewBinanceClient(credentials)
	// Delete next line to use prod API (you have to provide envs)

	t := tasker.NewTasker(tasker.WithSignalChannel(so.SignalChanel))

	//Tasker can run any type of task as long as it satisfies Task interface
	tasks := []tasker.Task{
		&binance.Task{
			BinanceClient:  bc,
			MailClient:     mc,
			Recipient:      config.EmailRecipient,
			Counter:        0,
			TickerDuration: time.Second * 10,
			Task:           binance.GenerateMonitorPortfolioTask(24000, 24),
		},
		&binance.Task{
			BinanceClient:  bc,
			MailClient:     mc,
			Recipient:      config.EmailRecipient,
			Counter:        0,
			TickerDuration: time.Second * 5,
			Task:           binance.GenerateMonitorSymbolTask("ETHUSDT", 1981, true, 6*24),
		},
		&binance.Task{
			BinanceClient:  bc,
			MailClient:     mc,
			Recipient:      config.EmailRecipient,
			Counter:        0,
			TickerDuration: time.Second * 15,
			Task:           binance.GenerateMonitorSymbolTask("ETHUSDT", 780, false, 6*24),
		},
	}
	t.Run(tasks)

	exitCode := <-so.ExitChanel
	os.Exit(exitCode)
}
