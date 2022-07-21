package main

import (
  "fmt"

	"github.com/shopspring/decimal"
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )


type PacaClient struct {
	trade alpaca.Client
	data  marketdata.Client
}

func NewPacaClient() PacaClient{
	// You can set your API key/secret here or you can use environment variables!
	apiKey := ""
	apiSecret := ""
	// Change baseURL to https://paper-api.alpaca.markets if you want use paper!
	baseURL := ""

  ac := PacaClient{
		trade: alpaca.NewClient(alpaca.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
			BaseURL:   baseURL,
		}),
		data: marketdata.NewClient(marketdata.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		}),
    }
  return ac
}

func (paca PacaClient) SubmitOrder(qty int, symbol string, side string) error {
	account, err := paca.trade.GetAccount()
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}
	if qty > 0 {
		adjSide := alpaca.Side(side)
		decimalQty := decimal.NewFromInt(int64(qty))
		_, err := paca.trade.PlaceOrder(alpaca.PlaceOrderRequest{
			AccountID:   account.ID,
			AssetKey:    &symbol,
			Qty:         &decimalQty,
			Side:        adjSide,
			Type:        "market",
			TimeInForce: "day",
		})
		if err == nil {
			fmt.Printf("Market order of | %d %s %s | completed\n", qty, symbol, side)
		} else {
			fmt.Printf("Order of | %d %s %s | did not go through: %s\n", qty, symbol, side, err)
		}
		return err
	}
	fmt.Printf("Quantity is <= 0, order of | %d %s %s | not sent\n", qty, symbol, side)
	return nil
}
