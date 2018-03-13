package exchange

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
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
			// this person is trying to sell 5 skycoins at a price of 4
			// bitcoins per skycoin
			// if their order is fulfilled, they lose 5 SKY and gain 20 BTC
			{
				Price:  decimal.NewFromFloat(4.0),
				Volume: decimal.NewFromFloat(5.0),
			},
			// this person is trying to sell 3 skycoins at a price of 4
			// bitcoins per skycoin
			// if their order is fulfilled, they lose 3 SKY and gain 12 BTC
			{
				Price:  decimal.NewFromFloat(4.0),
				Volume: decimal.NewFromFloat(3.0),
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

	orders, err := marketRecord.SpendItAll(bankroll)

	require.NoError(t, err)

	totalSpent := decimal.Zero

	for _, order := range orders {
		totalSpent = totalSpent.Add(order.Price.Mul(order.Volume))
	}

	require.True(t, totalSpent.Equal(bankroll), "Returned orders don't equal bankroll")

	totalPurchased := orders.Volume()

	require.True(t, totalPurchased.Equal(decimal.NewFromFloat(9.6)),
		// that's 8.0 skycoins at 4 btc per sky = 32 btc
		// plus   1.6 skycoins at 5 btc per sky =  8 btc
		// =============================================
		//        9.6                             40 btc (our original bankroll)
		"Failed to get cheapest price")
}

func TestSpendItAll_ErrNegativeAmount(t *testing.T) {
	marketRecord := testMarketRecord()
	bankroll := decimal.NewFromFloat(-5.0)

	_, err := marketRecord.SpendItAll(bankroll)

	require.Equal(t, ErrNegativeAmount, err)
}

func TestSpendItAll_ErrOrdersRanOut(t *testing.T) {
	marketRecord := testMarketRecord()
	bankroll := decimal.NewFromFloat(78.0)

	_, err := marketRecord.SpendItAll(bankroll)

	require.Equal(t, ErrOrdersRanOut, err)
}

func TestCheapestAsk_Success(t *testing.T) {
	marketRecord := testMarketRecord()

	order := marketRecord.CheapestAsk()

	require.NotNil(t, order)

	// CheapestAsk() should give us 20btc in this case because of the two lowest prices, one was selling 5 sky and the other 3 sky
	require.True(t, order.TotalCost().Equal(decimal.NewFromFloat(20.0)))
}

func TestCheapestAsk_NoAsks(t *testing.T) {
	marketRecord := testMarketRecord()

	marketRecord.Asks = []MarketOrder{}

	order := marketRecord.CheapestAsk()

	require.Nil(t, order)
}
