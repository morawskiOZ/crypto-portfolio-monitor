package main

import (
	"fmt"
	"os"

	"github.com/morawskioz/binance-monitor/configs"
	"github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/interal/signal"
	"github.com/morawskioz/binance-monitor/interal/task"
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
		mail.WithDialer(mail.MailAuthConfig{
			Port:         config.SmtpPort,
			Host:         config.SmtpHost,
			Password:     config.EmailPass,
			EmailAddress: config.EmailLogin,
		}),
	)

	so := signal.NewOsSignalObserver()
	go so.Observe()

	credentials := binance.Credentials{
		Key:    config.Key,
		Secret: config.Secret,
	}

	bc := binance.NewBinanceClient(credentials, binance.WithTestFlag())
	t := tasker.NewTasker(tasker.WithSignalChannel(so.SignalChanel))

	tasks := []tasker.Task{
		&task.BinanceTask{
			BinanceClient: bc,
			MailClient:    mc,
			Recipient:     config.EmailRecipient,
			Counter:       0,
		},
	}
	t.Run(tasks)

	exitCode := <-so.ExitChanel
	os.Exit(exitCode)
}
