package main

import (
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
  )

type AlgoPipeLine struct{
  clients AlgoClient
  job1 Job1 
}


func NewAlgoPipeLine (ac AlgoClient) AlgoPipeLine {
    ap := AlgoPipeLine{
      clients: ac,
    }
    return ap
}
