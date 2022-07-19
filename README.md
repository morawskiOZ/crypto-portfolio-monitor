# Crypto monitor

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project
Go program that can monitor portfolio and coins' value. Currently, it has build in Binance API. Feel free to extend it

> :warning: **By default tasker is using Binance TESTING API, add prod envs and remove `binanceAPI.WithTestFlag()` from `main.go` file`**
<p align="right">(<a href="#top">back to top</a>)</p>



### Built With

* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)
<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

Go version > 1.16 installed.

### Starting

1. Clone repo and start with: 
```sh
   go run main.go
   ```
In the `cmd/main.go` file you can see an example of how to simply schedule tasks`
```go
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
```
You can extend this monitor to any other exchange or even schedule non crypto related tasks. Feel free to fork it and play with it.


<p align="right">(<a href="#top">back to top</a>)</p>

### Running docker on server with custom envs and in detached mode
```
docker run -e EMAIL_PASS="" -e EMAIL_LOGIN="" -e EMAIL_RECIPIENT="" -e KEY="" -e SECRET="" -d dockerTagOrID
```
