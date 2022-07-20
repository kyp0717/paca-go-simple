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
	//  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	// "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	// "github.com/shopspring/decimal"
	// "/home/phage/work/projects/paca-go-simple/pkg"
)

func main() {
	// err := godotenv.Load("~/projects/paca-go/.env")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error: file not found")
	}

  algoClient := NewAlgoClient()
  algoPL:= NewAlgoPipeLine(algoClient)
  

  for {
    algoPL.Extract()
  }


}
