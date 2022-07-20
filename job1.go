package main

import (
  "fmt"
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

type Job1 struct {
  paca AlgoClient
  stockprice float64
}

type AlgoClient struct {
	tradeClient alpaca.Client
	dataClient  marketdata.Client
}

func NewAlgoClient() AlgoClient{
	// You can set your API key/secret here or you can use environment variables!
	apiKey := ""
	apiSecret := ""
	// Change baseURL to https://paper-api.alpaca.markets if you want use paper!
	baseURL := ""

  var ac AlgoClient

	ac = AlgoClient{
		tradeClient: alpaca.NewClient(alpaca.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
			BaseURL:   baseURL,
		}),
		dataClient: marketdata.NewClient(marketdata.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		}),
    }
  return ac
}

func NewJob(ac AlgoClient) Job1 {
  j := Job1{
    paca: ac,
  }
  return j
}

func (j Job1) GetPrice(symbol string) {
  snapshot, err := j.paca.dataClient.GetSnapshot(symbol)
  if err != nil {
    j.stockprice = snapshot.MinuteBar.Close
  } else {
    fmt.Println("error")
  }
} 



