package pkg

import (
  "time"
  "fmt"

  "github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/shopspring/decimal"
  )

type bucket struct {
	list        []string
	qty         int
	adjustedQty int
	equityAmt   float64
}

type longShortAlgo struct {
	tradeClient alpaca.Client
	dataClient  marketdata.Client
	long        bucket
	short       bucket
	allStocks   []stockField
	blacklist   []string
}

type stockField struct {
	name string
	pc   float64
}

var algo longShortAlgo

// Rebalance the portfolio every minute, making necessary trades.
func (alp longShortAlgo) run() error {

	// Figure out when the market will close so we can prepare to sell beforehand.
	clock, err := algo.tradeClient.GetClock()
	if err != nil {
		return fmt.Errorf("get clock: %w", err)
	}
	if clock.NextClose.Sub(clock.Timestamp) < 15*time.Minute {
		// Close all positions when 15 minutes til market close.
		fmt.Println("Market closing soon. Closing positions")

		positions, err := algo.tradeClient.ListPositions()
		if err != nil {
      return fmt.Errorf("get positions: %w", err)
		}
		for _, position := range positions {
			var orderSide string
			if position.Side == "long" {
				orderSide = "sell"
			} else {
				orderSide = "buy"
			}
			qty, _ := position.Qty.Float64()
			qty = math.Abs(qty)
			if err := algo.submitOrder(int(qty), position.Symbol, orderSide); err != nil {
				return fmt.Errorf("submit order: %w", err)
			}
		}
		// Run script again after market close for next trading day.
		fmt.Println("Sleeping until market close (15 minutes)")
		time.Sleep(15 * time.Minute)
	} else {
		// Rebalance the portfolio.
		if err := algo.rebalance(); err != nil {
			fmt.Println("Failed to rebalance, will try again in a minute:", err)
		}
		fmt.Println("Sleeping for 1 minute")
		time.Sleep(1 * time.Minute)
	}
	return nil
}


// Rebalance our position after an update.
func (alp longShortAlgo) rebalance() error {
	if err := algo.rerank(); err != nil {
		return fmt.Errorf("rerank: %w", err)
	}

	fmt.Printf("We are taking a long position in: %v\n", algo.long.list)
	fmt.Printf("We are taking a short position in: %v\n", algo.short.list)

	// Clear existing orders again.
	status, until, limit := "open", time.Now(), 100
	orders, err := algo.tradeClient.ListOrders(&status, &until, &limit, nil)
	if err != nil {
		return fmt.Errorf("list orders: %w", err)
	}
	for _, order := range orders {
		if err := algo.tradeClient.CancelOrder(order.ID); err != nil {
			return fmt.Errorf("cancel order %s: %w", order.ID, err)
		}
	}

	// Remove positions that are no longer in the short or long list, and make a list of positions that do not need to change.  Adjust position quantities if needed.
	algo.blacklist = nil
	var executed [2][]string
	positions, err := algo.tradeClient.ListPositions()
	if err != nil {
		return fmt.Errorf("list positions: %w", err)
	}
	for _, position := range positions {
		indLong := indexOf(algo.long.list, position.Symbol)
		indShort := indexOf(algo.short.list, position.Symbol)

		rawQty, _ := position.Qty.Float64()
		qty := int(math.Abs(rawQty))
		side := "buy"
		if indLong < 0 {
			// Position is not in long list.
			if indShort < 0 {
				// Position not in short list either.  Clear position.
				if position.Side == "long" {
					side = "sell"
				} else {
					side = "buy"
				}
				if err := algo.submitOrder(int(math.Abs(float64(qty))), position.Symbol, side); err != nil {
					return fmt.Errorf("submit order for %d %s: %w", qty, position.Symbol, err)
				}
			} else {
				if position.Side == "long" {
					// Position changed from long to short.  Clear long position to prep for short sell.
					side = "sell"
					if err := algo.submitOrder(qty, position.Symbol, side); err != nil {
						return fmt.Errorf("submit order for %d %s: %w", qty, position.Symbol, err)
					}
				} else {
					// Position in short list
					if qty == algo.short.qty {
						// Position is where we want it.  Pass for now
					} else {
						// Need to adjust position amount.
						diff := qty - algo.short.qty
						if diff > 0 {
							// Too many short positions.  Buy some back to rebalance.
							side = "buy"
						} else {
							// Too little short positions.  Sell some more.
							diff = int(math.Abs(float64(diff)))
							side = "sell"
						}
						qty = diff
						if err := algo.submitOrder(qty, position.Symbol, side); err != nil {
							return fmt.Errorf("submit order for %d %s: %w", qty, position.Symbol, err)
						}
					}
					executed[1] = append(executed[1], position.Symbol)
					algo.blacklist = append(algo.blacklist, position.Symbol)
				}
			}
		} else {
			// Position in long list.
			if position.Side == "short" {
				// Position changed from short to long.  Clear short position to prep for long purchase.
				side = "buy"
				if err := algo.submitOrder(qty, position.Symbol, side); err != nil {
					return fmt.Errorf("submit order for %d %s: %w", qty, position.Symbol, err)
				}
			} else {
				if qty == algo.long.qty {
					// Position is where we want it.  Pass for now.
				} else {
					// Need to adjust position amount
					diff := qty - algo.long.qty
					if diff > 0 {
						// Too many long positions.  Sell some to rebalance.
						side = "sell"
					} else {
						diff = int(math.Abs(float64(diff)))
						side = "buy"
					}
					qty = diff
					if err := algo.submitOrder(qty, position.Symbol, side); err != nil {
						return fmt.Errorf("submit order for %d %s: %w", qty, position.Symbol, err)
					}
				}
				executed[0] = append(executed[0], position.Symbol)
				algo.blacklist = append(algo.blacklist, position.Symbol)
			}
		}
	}

	// Send orders to all remaining stocks in the long and short list.
	longBOResp := algo.sendBatchOrder(algo.long.qty, algo.long.list, "buy")
	executed[0] = append(executed[0], longBOResp[0][:]...)
	if len(longBOResp[1][:]) > 0 {
		// Handle rejected/incomplete orders and determine new quantities to purchase.

		longTPResp, err := algo.getTotalPrice(executed[0])
		if err != nil {
			return fmt.Errorf("get total long price: %w", err)
		}
		if longTPResp > 0 {
			algo.long.adjustedQty = int(algo.long.equityAmt / longTPResp)
		} else {
			algo.long.adjustedQty = -1
		}
	} else {
		algo.long.adjustedQty = -1
	}

	shortBOResp := algo.sendBatchOrder(algo.short.qty, algo.short.list, "sell")
	executed[1] = append(executed[1], shortBOResp[0][:]...)
	if len(shortBOResp[1][:]) > 0 {
		// Handle rejected/incomplete orders and determine new quantities to purchase.
		shortTPResp, err := algo.getTotalPrice(executed[1])
		if err != nil {
			return fmt.Errorf("get total short price: %w", err)
		}
		if shortTPResp > 0 {
			algo.short.adjustedQty = int(algo.short.equityAmt / shortTPResp)
		} else {
			algo.short.adjustedQty = -1
		}
	} else {
		algo.short.adjustedQty = -1
	}

	// Reorder stocks that didn't throw an error so that the equity quota is reached.
	if algo.long.adjustedQty > -1 {
		algo.long.qty = algo.long.adjustedQty - algo.long.qty
		for _, stock := range executed[0] {
			if err := algo.submitOrder(algo.long.qty, stock, "buy"); err != nil {
				return fmt.Errorf("submit order for %d %s: %w", algo.long.qty, stock, err)
			}
		}
	}

	if algo.short.adjustedQty > -1 {
		algo.short.qty = algo.short.adjustedQty - algo.short.qty
		for _, stock := range executed[1] {
			if err := algo.submitOrder(algo.short.qty, stock, "sell"); err != nil {
				return fmt.Errorf("submit order for %d %s: %w", algo.long.qty, stock, err)
			}
		}
	}

	return nil
}

