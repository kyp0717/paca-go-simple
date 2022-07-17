
type longShortAlgo struct {
	tradeClient alpaca.Client
	dataClient  marketdata.Client
	long        bucket
	short       bucket
	allStocks   []stockField
	blacklist   []string
}

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
