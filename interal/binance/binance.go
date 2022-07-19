package binance

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/morawskioz/binance-monitor/interal/price"
	"github.com/pkg/errors"
)

// Client is the wrapper for binance API client
type Client struct {
	binanceClient  *binance.Client
	priceService   *binance.ListPricesService
	accountService *binance.GetAccountService
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
	c.accountService = c.binanceClient.NewGetAccountService()

	return c
}

func (c *Client) getAccountBalances() ([]binance.Balance, error) {
	res, err := c.accountService.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res.Balances, nil
}

func (c *Client) getPortfolio() (portfolio, error) {
	balances, err := c.getAccountBalances()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account balance")
	}
	var portfolio []Coin
	for _, balance := range balances {
		freeBalance, _ := strconv.ParseFloat(balance.Free, 64)
		lockedBalance, _ := strconv.ParseFloat(balance.Locked, 64)
		totalBalance := freeBalance + lockedBalance
		if totalBalance == 0 {
			continue
		}

		s, ok := price.ConvertToUSDTPair(balance.Asset)
		if !ok {
			continue
		}

		p, err := c.priceService.Symbol(s).Do(context.Background())
		if err != nil || len(p) == 0 {
			return nil, errors.Wrap(err, "Failed to get prf")
		}

		prf, _ := strconv.ParseFloat(p[0].Price, 64)
		v := totalBalance * prf
		portfolio = append(portfolio, Coin{
			asset:        p[0].Symbol,
			price:        prf,
			totalBalance: totalBalance,
			value:        v,
		})
	}

	return portfolio, nil
}

func (c *Client) GetPortfolioTotalValue() (float64, error) {
	p, err := c.getPortfolio()
	if err != nil {
		return 0, errors.Wrap(err, "Can't calculate portfolio value, portfolios are not available")
	}
	if len(p) == 0 {
		return 0, nil
	}
	var portfolioValue float64
	for _, coin := range p {
		portfolioValue += coin.value
	}

	return portfolioValue, nil
}

// GetSymbolValue allows to get price for symbol
func (c *Client) GetSymbolValue(s string) (float64, error) {
	p, err := c.priceService.Symbol(s).Do(context.Background())
	if err != nil || len(p) == 0 {
		return 0, errors.Wrap(err, "Failed to get price")
	}

	return strconv.ParseFloat(p[0].Price, 64)
}
