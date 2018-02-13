package exchange

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func testMarketRecord() MarketRecord {
	return MarketRecord{
		Timestamp: time.Now(),
		Symbol:    "SKY/BTC",
		Asks: []MarketOrder{
			// this person is trying to sell 9 skycoins at a price of 5
			// bitcoins per skycoin
			// if their order is fulfilled, they lose 9 SKY and gain 45 BTC
			{
				Price:  decimal.NewFromFloat(5.0),
				Volume: decimal.NewFromFloat(9.0),
			},
			// this person is trying to sell 8 skycoins at a price of 4
			// bitcoins per skycoin
			// if their order is fulfilled, they lose 8 SKY and gain 32 BTC
			{
				Price:  decimal.NewFromFloat(4.0),
				Volume: decimal.NewFromFloat(8.0),
			},
		},
		Bids: []MarketOrder{
			// this person is trying to buy 7 skycoins at a price of 3
			// bitcoins per skycoin
			// if their order is fulfilled, they gain 7 SKY and lose 21 BTC
			{
				Price:  decimal.NewFromFloat(3.0),
				Volume: decimal.NewFromFloat(7.0),
			},
			// this person is trying to buy 6 skycoins at a price of 2
			// bitcoins per skycoin
			// if their order is fulfilled, they gain 6 SKY and lose 12 BTC
			{
				Price:  decimal.NewFromFloat(2.0),
				Volume: decimal.NewFromFloat(6.0),
			},
		},
	}
}

func TestSpendItAll_Success(t *testing.T) {
	marketRecord := testMarketRecord()
	bankroll := decimal.NewFromFloat(40.0)

	orders, remaining, err := marketRecord.SpendItAll(bankroll)

	if err != nil {
		t.Fatal(err)
	}

	if !remaining.Equal(decimal.NewFromFloat(0.0)) {
		t.Fatalf("Failed to spend all of the bankroll")
	}

	totalSpent := decimal.NewFromFloat(0.0)
	totalPurchased := decimal.NewFromFloat(0.0)

	for _, order := range orders {
		totalSpent = totalSpent.Add(order.Price.Mul(order.Volume))
		totalPurchased = totalPurchased.Add(order.Volume)
	}

	if !totalSpent.Equal(bankroll) {
		t.Fatalf("Returned orders don't equal bankroll")
	}

	if !totalPurchased.Equal(decimal.NewFromFloat(9.6)) {
		// that's 8.0 skycoins at 4 btc per sky = 32 btc
		// plus   1.6 skycoins at 5 btc per sky =  8 btc
		// ===========================================
		//        9.6                             40 btc (our original bankroll)
		t.Fatalf("Failed to get cheapest price")
	}
}

func TestSpendItAll_ErrNegativeAmount(t *testing.T) {
	marketRecord := testMarketRecord()
	bankroll := decimal.NewFromFloat(-5.0)

	_, _, err := marketRecord.SpendItAll(bankroll)

	if err != ErrNegativeAmount {
		t.Fatalf("Shouldn't be able to spend negative currency")
	}
}

func TestSpendItAll_ErrOrdersRanOut(t *testing.T) {
	marketRecord := testMarketRecord()
	bankroll := decimal.NewFromFloat(78.0)

	_, remaining, err := marketRecord.SpendItAll(bankroll)

	if err != ErrOrdersRanOut {
		t.Fatalf("Orders were supposed to run out, but didn't")
	}

	if remaining.LessThanOrEqual(decimal.NewFromFloat(0.0)) {
		t.Fatal("Remaining should be positive if orders ran out")
	}
}
