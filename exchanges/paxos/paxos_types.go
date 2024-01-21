package paxos

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/types"
)

// Profiles holds profile information attached to the account
type Profiles struct {
	Items []struct {
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Type     string `json:"type"`
	} `json:"items"`
}

// MarketInstruments represents a list of market instrument.
type MarketInstruments struct {
	Markets []struct {
		Market     string       `json:"market"`
		BaseAsset  string       `json:"base_asset"`
		QuoteAsset string       `json:"quote_asset"`
		TickRate   types.Number `json:"tick_rate"`
	} `json:"markets"`
}

// Orderbook represents an orderbook data for a market.
type Orderbook struct {
	Market string `json:"market"`
	Asks   []struct {
		Price  types.Number `json:"price"`
		Amount types.Number `json:"amount"`
	} `json:"asks"`
	Bids []struct {
		Price  types.Number `json:"price"`
		Amount types.Number `json:"amount"`
	} `json:"bids"`
}

// RecentExecutions represents most recent executions by all users
type RecentExecutions struct {
	Items []struct {
		MatchNumber string       `json:"match_number"`
		Price       types.Number `json:"price"`
		Amount      types.Number `json:"amount"`
		ExecutedAt  time.Time    `json:"executed_at"`
	} `json:"items"`
}

// TckerDetail represents an orderbook statistics
type TckerDetail struct {
	Market        string          `json:"market"`
	BestBid       AmountPriceInfo `json:"best_bid"`
	BestAsk       AmountPriceInfo `json:"best_ask"`
	LastExecution AmountPriceInfo `json:"last_execution"`
	LastDay       TickerInfo      `json:"last_day"`
	Today         TickerInfo      `json:"today"`
	SnapshotAt    time.Time       `json:"snapshot_at"`
}

// AmountPriceInfo represents a price and amount infor
type AmountPriceInfo struct {
	Price  types.Number `json:"price"`
	Amount types.Number `json:"amount"`
}

// TickerInfo represents details of ticker information.
type TickerInfo struct {
	High      types.Number `json:"high"`
	Low       types.Number `json:"low"`
	Open      types.Number `json:"open"`
	Volume    types.Number `json:"volume"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
}

// MarketPricing hodls current prices, as well as 24 hour prior (yesterday) prices
type MarketPricing struct {
	Prices []struct {
		Market         string       `json:"market"`
		CurrentPrice   types.Number `json:"current_price"`
		YesterdayPrice types.Number `json:"yesterday_price"`
		SnapshotAt     time.Time    `json:"snapshot_at"`
	} `json:"prices"`
}
