package main

import (
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

type AlgoPipeLine struct{
  clients AlgoClient
  job1 Job1DataExtract 
  job2 Job2Train 
}

type AlgoClient struct {
	tradeClient alpaca.Client
	dataClient  marketdata.Client
}

type bucket struct {
	list        []string
	qty         int
	adjustedQty int
	equityAmt   float64
}

type stockField struct {
	name string
	pc   float64
}


// Set this to true if you have unlimited subscription!
var hasSipAccess bool = false

func NewAlgoClient() AlgoClient{
	// You can set your API key/secret here or you can use environment variables!
	apiKey := ""
	apiSecret := ""
	// Change baseURL to https://paper-api.alpaca.markets if you want use paper!
	baseURL := ""

  var ac AlgoClient
	// Format the allStocks variable for use in the class.
	allStocks := []stockField{}
	stockList := []string{"DOMO", "SQ", "MRO", "AAPL", "GM", "SNAP", "SHOP", "SPLK", "BA", "AMZN", "SUI", "SUN", "TSLA", "CGC", "SPWR", "NIO", "CAT", "MSFT", "PANW", "OKTA", "TWTR", "TM", "GE", "ATVI", "GS", "BAC", "MS", "TWLO", "QCOM", "IBM"}
	for _, stock := range stockList {
		allStocks = append(allStocks, stockField{stock, 0})
	}

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

func NewAlgoPipeLine (ac AlgoClient) AlgoPipeLine {
    ap := AlgoPipeLine{
      clients: ac,
    }
    return ap
}
