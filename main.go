package trade

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/shopspring/decimal"
)

func main() {
	// err := godotenv.Load("~/projects/paca-go/.env")
	err := godotenv.Load()
	if err != nil {
		t.Log("error: file not found")
	}
	// stream.New
	// Creating a client that connexts to iex
	c := stream.NewStocksClient("iex")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// setting up cancelling upon interrupt
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	go func() {
		<-s
		cancel()
	}()

	if err := c.Connect(ctx); err != nil {
		log.Fatalf("could not establish connection, error: %s", err)
	}

	fmt.Print("Cancelling all open orders so they don't impact our buying power... ")
	status, until, limit := "open", time.Now(), 100
	orders, err := algo.tradeClient.ListOrders(&status, &until, &limit, nil)
	if err != nil {
		log.Fatalf("Failed to list orders: %v", err)
	}

	for _, order := range orders {
		if err := algo.tradeClient.CancelOrder(order.ID); err != nil {
			log.Fatalf("Failed to cancel order %s: %v", order.ID, err)
		}
	}
	fmt.Printf("%d order(s) cancelled\n", len(orders))

	for {
		isOpen, err := algo.awaitMarketOpen()
		if err != nil {
			log.Fatalf("Failed to wait for market open: %v", err)
		}
		if !isOpen {
			time.Sleep(1 * time.Minute)
			continue
		}
		if err := algo.run(); err != nil {
			log.Fatalf("Run error: %v", err)
		}
	}
}

}
