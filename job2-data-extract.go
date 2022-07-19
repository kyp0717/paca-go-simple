package main

import (
  // "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	// "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

func (ar AlgoRunner) Job2DataExtract() {
  stocklist := []string{"aaa", "bbb"}
  prices, err := ar.dataClient.GetSnapshots(stocklist)

  
  


}
