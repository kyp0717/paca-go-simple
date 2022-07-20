package main


import (
  // "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

// define the features for your algo
type Features struct {
  // only 1 feature for this algo 
  // define more if needed
  stock string
  GetQuotes marketdata.GetQuotesParams
}
