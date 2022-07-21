package main

import (
	// "fmt"
	"time"

	// "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type Job1 struct {
  paca PacaClient
  // features Features
  stock string
  enterTrade bool
  exitTrade float64
}

type Outcome int
const (
  Buy Outcome = iota + 1
  Sell
  Hold
  )

type Job2Chan chan marketdata.Snapshot // predict, train 
type Job3Chan chan float64 // Result
type Job4Chan chan Outcome

func NewJob(ac PacaClient, s string) Job1 {
  j := Job1{
    paca: ac,
    stock: s,
  }
  return j
}

// return channel
func (j Job1) GetData() Job2Chan {
  j2 := make(chan marketdata.Snapshot)
  go func() {
    for {
      snapshot, err := j.paca.data.GetSnapshot(j.stock)
      // data 
      if err != nil {
        j2<- *snapshot
      }
      // throughput is control in this goroutine
      time.Sleep(3*time.Second)
    }
  }()
  return j2
} 

// predict, train, calculate
func (j2 Job2Chan) Train() Job3Chan {
  j3 := make(chan float64)
  go func() {
    for {
    snapshot := <-j2
    priceChange := snapshot.MinuteBar.Open - snapshot.MinuteBar.Close 
    pctChange := priceChange/snapshot.MinuteBar.Open
    j3 <- pctChange
    }
  }()
  return j3
}
// infer decision based on result from Job3 Channel
func (j3 Job3Chan) Infer() Job4Chan {
  j4 := make(chan Outcome)
  go func() {
    for {
      result := <-j3
      if result <0.05 {
        j4 <- Hold
      } else {
        j4 <- Sell
      }
    }
  }()
  return j4
}

func (j4 Job4Chan) Trade(paca PacaClient, done chan bool) {
  go func() {
    for {
    decision := <- j4
    switch decision {
        case Hold: // resume pipeline
        case Buy: paca.SubmitOrder(100, "AMD", "buy")
        case Sell: {
          paca.SubmitOrder(100, "AMD", "sell")
          done<-true
          }
      }
    }
  }()
}

