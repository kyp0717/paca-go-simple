package main

import (
	// "context"
	// "fmt"
	// "log"
	// "os"
	// "os/signal"
	// "time"
	//
	"fmt"

	"github.com/joho/godotenv"
	// "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	// "/home/phage/work/projects/paca-go-simple/pkg"
)

func main() {
  exitTrade := make(chan bool)
	// err := godotenv.Load("~/projects/paca-go/.env")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error: file not found")
	}

  paca:= NewPacaClient()
  job1 := NewJob(paca, "AMD") // get prices from alpaca
  job2 := job1.StreamData()
  job2.EnterTrade(paca) // hang here until order is submitted
  job3 := job2.Train()
  job4 := job3.Infer()
  job4.Trade(paca, exitTrade)
  <-exitTrade
}
