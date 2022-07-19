package pkg

import (
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

type AlgoClient struct {
	tradeClient alpaca.Client
	dataClient  marketdata.Client
	long        bucket
	short       bucket
	allStocks   []stockField
	blacklist   []string
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

func InitAlgoClient() AlgoClient {
	// You can set your API key/secret here or you can use environment variables!
	apiKey := ""
	apiSecret := ""
	// Change baseURL to https://paper-api.alpaca.markets if you want use paper!
	baseURL := ""

  var algo AlgoClient
	// Format the allStocks variable for use in the class.
	allStocks := []stockField{}
	stockList := []string{"DOMO", "SQ", "MRO", "AAPL", "GM", "SNAP", "SHOP", "SPLK", "BA", "AMZN", "SUI", "SUN", "TSLA", "CGC", "SPWR", "NIO", "CAT", "MSFT", "PANW", "OKTA", "TWTR", "TM", "GE", "ATVI", "GS", "BAC", "MS", "TWLO", "QCOM", "IBM"}
	for _, stock := range stockList {
		allStocks = append(allStocks, stockField{stock, 0})
	}

	algo = AlgoClient{
		tradeClient: alpaca.NewClient(alpaca.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
			BaseURL:   baseURL,
		}),
		dataClient: marketdata.NewClient(marketdata.ClientOpts{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		}),
		long: bucket{
			qty:         -1,
			adjustedQty: -1,
		},
		short: bucket{
			qty:         -1,
			adjustedQty: -1,
		},
		allStocks: allStocks,
		blacklist: []string{},
	}
  return algo
}
