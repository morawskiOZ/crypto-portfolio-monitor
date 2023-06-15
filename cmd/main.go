package main

import (
	"fmt"
	"github.com/morawskioz/binance-monitor/interal/health"
	"github.com/morawskioz/binance-monitor/interal/tasks/binance"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"

	binanceAPI "github.com/morawskioz/binance-monitor/interal/binance"
	"github.com/morawskioz/binance-monitor/interal/signal"
	"github.com/morawskioz/binance-monitor/interal/tasker"
	"github.com/morawskioz/binance-monitor/pkg/mail"
)

func main() {
	if appEnv := os.Getenv("APP_ENV"); appEnv != "production" {
		viper.SetConfigName("dev")
		viper.AddConfigPath("../")
	} else {
		viper.SetConfigName("prod")
		viper.AddConfigPath("./")

	}

	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Info("no .env file found")
	}

	viper.AutomaticEnv()
	fmt.Println(viper.GetString("EMAIL_RECIPIENT"))
	go health.StartHealthCheck()
	fmt.Println("Starting binance monitor", "env:", viper.GetString("APP_ENV"))
	mc := mail.NewMailClient(
		mail.WithDialer(mail.AuthConfig{
			Port:         viper.GetInt("SMTP_PORT"),
			Host:         viper.GetString("SMTP_HOST"),
			Password:     viper.GetString("EMAIL_PASS"),
			EmailAddress: viper.GetString("EMAIL_LOGIN"),
		}),
	)

	so := signal.NewOsSignalObserver()
	go so.Observe()

	binanceAPI.WithTestFlag()
	credentials := binanceAPI.Credentials{
		Key:    viper.GetString("BINANCE_KEY"),
		Secret: viper.GetString("BINANCE_SECRET"),
	}

	bc := binanceAPI.NewBinanceClient(credentials)

	t := tasker.NewTasker(tasker.WithSignalChannel(so.SignalChanel))
	recipient := viper.GetString("EMAIL_RECIPIENT")
	//Tasker can run any type of task as long as it satisfies Task interface
	tasks := []tasker.Task{
		&binance.Task{
			BinanceClient:  bc,
			MailClient:     mc,
			Recipient:      recipient,
			TickerDuration: time.Minute * 30,
			// every week so ticker is every 30 minutes * 2 (1 hour) * 24 (1 day)
			Task: binance.GenerateMonitorSymbolTask("BTCUSDT", 17001, false, 2*24),
		},
		&binance.Task{
			BinanceClient:  bc,
			MailClient:     mc,
			Recipient:      recipient,
			TickerDuration: time.Minute * 30,
			Task:           binance.GenerateMonitorSymbolTask("BTCUSDT", 15001, false, 0),
		},
	}
	t.Run(tasks)

	exitCode := <-so.ExitChanel
	os.Exit(exitCode)
}
