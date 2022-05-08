package binance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/morawskioz/binance-monitor/interal/price"
	"github.com/pkg/errors"
)

// Client is the wrapper for binance API client
type Client struct {
	binanceClient *binance.Client
	priceService  *binance.ListPricesService
}

// Credentials is the credentials for binance API
type Credentials struct {
	Key    string
	Secret string
}

// Option allows to set options for the client during creation
type Option func(*Client)

// Coin represents a coin with its balance, price and value
type Coin struct {
	asset        string
	price        float64
	totalBalance float64
	value        float64
}

type portfolio []Coin

// WithTestFlag sets the test flag to true
func WithTestFlag() Option {
	return func(c *Client) {
		binance.UseTestnet = true
	}
}

// NewBinanceClient creates a new binance client
func NewBinanceClient(cr Credentials, options ...Option) *Client {
	c := &Client{}
	for _, opt := range options {
		opt(c)
	}
	c.binanceClient = binance.NewClient(cr.Key, cr.Secret)
	c.priceService = c.binanceClient.NewListPricesService()
	return c
}

func (c *Client) getAccountBalances() []binance.Balance {
	res, err := c.binanceClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res.Balances
}

func (c *Client) getPortfolio() (portfolio, error) {
	balances := c.getAccountBalances()
	portfolio := []Coin{}
	for _, balance := range balances {
		s, ok := price.ConvertToUSDTPair(balance.Asset)
		if !ok {
			continue
		}

		p, err := c.priceService.Symbol(s).Do(context.Background())
		price, _ := strconv.ParseFloat(p[0].Price, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get price")
		}

		freeBalance, _ := strconv.ParseFloat(balance.Free, 64)
		lockedBalance, _ := strconv.ParseFloat(balance.Locked, 64)
		totalBalance := freeBalance + lockedBalance
		v := totalBalance * price
		portfolio = append(portfolio, Coin{
			asset:        p[0].Symbol,
			price:        price,
			totalBalance: totalBalance,
			value:        v,
		})

	}

	return portfolio, nil
}

func (c *Client) GetPortfolioTotalValue() (float64, error) {
	p, err := c.getPortfolio()
	if err != nil {
		return 0, errors.Wrap(err, "Can't calculate portfolio value, porfolios are not available")
	}
	var portfolioValue float64
	for _, coin := range p {
		portfolioValue += coin.value
	}
	return portfolioValue, nil
}
