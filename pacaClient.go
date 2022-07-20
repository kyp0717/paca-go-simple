package main

import (
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

  var ac PacaClient

	ac = PacaClient{
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

