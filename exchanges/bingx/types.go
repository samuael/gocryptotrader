package bingx

import (
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/encoding/json"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/types"
)

// ResponseWrapper represents a common server response wrapper
type ResponseWrapper struct {
	Code      uint64 `json:"code"`
	Msg       string `json:"msg"`
	Retryable int64  `json:"retryable"`
	Data      any    `json:"data"`
}

// ConfigSymbols represents a common symbols used by the market
type ConfigSymbols struct {
	Symbols []*SymbolDetail `json:"symbols"`
}

// SymbolDetail represents a spot symbol detail.
type SymbolDetail struct {
	Symbol            currency.Pair `json:"symbol"`
	MinimumQuantity   float64       `json:"minQty"`
	MaximumQuantity   float64       `json:"maxQty"`
	MinNotional       float64       `json:"minNotional"`
	MaxNotional       float64       `json:"maxNotional"`
	MaxMarketNotional float64       `json:"maxMarketNotional"`
	Status            float64       `json:"status"`
	TickSize          float64       `json:"tickSize"`
	StepSize          float64       `json:"stepSize"`
	APIStateSell      bool          `json:"apiStateSell"`
	APIStateBuy       bool          `json:"apiStateBuy"`
	TimeOnline        int           `json:"timeOnline"`
	OffTime           uint64        `json:"offTime"`
	MaintainTime      uint64        `json:"maintainTime"`
	DisplayName       string        `json:"displayName"`
}

// SpotTradeInfo represents a spot trade infomration detail
type SpotTradeInfo struct {
	ID         uint64     `json:"id"`
	Price      float64    `json:"price"`
	Quantity   float64    `json:"qty"`
	Time       types.Time `json:"time"`
	BuyerMaker bool       `json:"buyerMaker"`
}

// SpotOrderbookDetail represent an orderbook detail of a symbol
type SpotOrderbookDetail struct {
	Bids         orderbook.LevelsArrayPriceAmount `json:"bids"`
	Asks         orderbook.LevelsArrayPriceAmount `json:"asks"`
	Timestamp    types.Time                       `json:"ts"`
	LastUpdateID uint64                           `json:"lastUpdateId"`
}

// SpotCandle represents a spot kine data information
type SpotCandle struct {
	Timestamp types.Time
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Quantity  float64
	CloseTime types.Time
	Amount    float64
}

// UnmarshalJSON deserializes JSON payload into SpotCandle.
// BingX returns each candle as a positional array in the order: open time, open, high, low, close, volume, close time, quote volume.
func (s *SpotCandle) UnmarshalJSON(data []byte) error {
	target := []any{&s.Timestamp, &s.Open, &s.High, &s.Low, &s.Close, &s.Quantity, &s.CloseTime, &s.Amount}
	return json.Unmarshal(data, &target)
}

// SpotTicker24Hr represents the 24-hour rolling window price change statistics for a symbol.
type SpotTicker24Hr struct {
	Symbol             currency.Pair `json:"symbol"`
	OpenPrice          types.Number  `json:"openPrice"`
	HighPrice          types.Number  `json:"highPrice"`
	LowPrice           types.Number  `json:"lowPrice"`
	LastPrice          types.Number  `json:"lastPrice"`
	PriceChange        types.Number  `json:"priceChange"`
	PriceChangePercent string        `json:"priceChangePercent"`
	Volume             types.Number  `json:"volume"`
	QuoteVolume        types.Number  `json:"quoteVolume"`
	OpenTime           types.Time    `json:"openTime"`
	CloseTime          types.Time    `json:"closeTime"`
	Count              uint64        `json:"count"`
	BidPrice           types.Number  `json:"bidPrice"`
	BidQuantity        types.Number  `json:"bidQty"`
	AskPrice           types.Number  `json:"askPrice"`
	AskQuantity        types.Number  `json:"askQty"`
}

// SpotPriceTicker represents the latest market price information for a symbol.
type SpotPriceTicker struct {
	Symbol currency.Pair     `json:"symbol"`
	Trades []*SpotPriceTrade `json:"trades"`
}

// SpotPriceTrade represents a single latest-price trade entry within a SpotPriceTicker.
type SpotPriceTrade struct {
	Timestamp types.Time   `json:"timestamp"`
	TradeID   string       `json:"tradeId"`
	Price     types.Number `json:"price"`
	Amount    types.Number `json:"amount"`
	Type      int64        `json:"type"`
	Volume    types.Number `json:"volume"`
}

// SpotBookTicker represents the best bid and ask price and quantity for a symbol.
type SpotBookTicker struct {
	EventType string        `json:"eventType"`
	Time      types.Time    `json:"time"`
	Symbol    currency.Pair `json:"symbol"`
	BidPrice  types.Number  `json:"bidPrice"`
	BidVolume types.Number  `json:"bidVolume"`
	AskPrice  types.Number  `json:"askPrice"`
	AskVolume types.Number  `json:"askVolume"`
}

// HistoricalSpotTrade represents a historical trade record returned by the old trade lookup endpoint.
type HistoricalSpotTrade struct {
	TradeID   string        `json:"tid"`
	Time      types.Time    `json:"t"`
	MakerSide uint8         `json:"ms"`
	Symbol    currency.Pair `json:"s"`
	Price     types.Number  `json:"p"`
	Volume    types.Number  `json:"v"`
}

// SwapContractDetail represents a USDT-M perpetual futures contract and its trading rules.
type SwapContractDetail struct {
	ContractID        string        `json:"contractId"`
	Symbol            currency.Pair `json:"symbol"`
	Size              types.Number  `json:"size"`
	QuantityPrecision uint64        `json:"quantityPrecision"`
	PricePrecision    uint64        `json:"pricePrecision"`
	FeeRate           float64       `json:"feeRate"`
	MakerFeeRate      float64       `json:"makerFeeRate"`
	TakerFeeRate      float64       `json:"takerFeeRate"`
	TradeMinLimit     float64       `json:"tradeMinLimit"`
	TradeMinQuantity  float64       `json:"tradeMinQuantity"`
	TradeMinUSDT      float64       `json:"tradeMinUSDT"`
	Currency          currency.Code `json:"currency"`
	Asset             currency.Code `json:"asset"`
	Status            uint64        `json:"status"`
	APIStateOpen      string        `json:"apiStateOpen"`
	APIStateClose     string        `json:"apiStateClose"`
	EnsureTrigger     bool          `json:"ensureTrigger"`
	TriggerFeeRate    types.Number  `json:"triggerFeeRate"`
	BrokerState       bool          `json:"brokerState"`
	LaunchTime        types.Time    `json:"launchTime"`
	MaintainTime      types.Time    `json:"maintainTime"`
	OffTime           types.Time    `json:"offTime"`
	DisplayName       string        `json:"displayName"`
}

// SwapOrderbookDepth represents the order book depth of a perpetual futures symbol
type SwapOrderbookDepth struct {
	Timestamp types.Time                       `json:"T"`
	Bids      orderbook.LevelsArrayPriceAmount `json:"bids"`
	Asks      orderbook.LevelsArrayPriceAmount `json:"asks"`
	BidsCoin  orderbook.LevelsArrayPriceAmount `json:"bidsCoin"`
	AsksCoin  orderbook.LevelsArrayPriceAmount `json:"asksCoin"`
}

// SwapTrade represents a perpetual futures trade
type SwapTrade struct {
	TradeID       string       `json:"id"`
	Time          types.Time   `json:"time"`
	IsBuyerMaker  bool         `json:"isBuyerMaker"`
	Price         types.Number `json:"price"`
	Quantity      types.Number `json:"qty"`
	QuoteQuantity types.Number `json:"quoteQty"`
}

// SwapMarkPriceFundingRate represents the current mark price, index price and funding rate information for a perpetual futures symbol.
type SwapMarkPriceFundingRate struct {
	Symbol               currency.Pair `json:"symbol"`
	MarkPrice            types.Number  `json:"markPrice"`
	IndexPrice           types.Number  `json:"indexPrice"`
	LastFundingRate      types.Number  `json:"lastFundingRate"`
	NextFundingTime      types.Time    `json:"nextFundingTime"`
	FundingIntervalHours uint64        `json:"fundingIntervalHours"`
	MinFundingRate       types.Number  `json:"minFundingRate"`
	MaxFundingRate       types.Number  `json:"maxFundingRate"`
	UpdateTime           types.Time    `json:"updateTime"`
}

// SwapFundingRate represents a historical funding rate record for a perpetual futures symbol.
type SwapFundingRate struct {
	Symbol      currency.Pair `json:"symbol"`
	FundingRate types.Number  `json:"fundingRate"`
	FundingTime types.Time    `json:"fundingTime"`
}

// SwapCandle represents a perpetual futures kline entry.
type SwapCandle struct {
	Open   types.Number `json:"open"`
	Close  types.Number `json:"close"`
	High   types.Number `json:"high"`
	Low    types.Number `json:"low"`
	Volume types.Number `json:"volume"`
	Time   types.Time   `json:"time"`
}

// SwapOpenInterest represents the open interest of a perpetual futures symbol.
type SwapOpenInterest struct {
	OpenInterest types.Number  `json:"openInterest"`
	Symbol       currency.Pair `json:"symbol"`
	Time         types.Time    `json:"time"`
}

// SwapTicker24Hr represents the 24-hour rolling window price change statistics for a perpetual futures symbol.
type SwapTicker24Hr struct {
	Symbol             currency.Pair `json:"symbol"`
	PriceChange        types.Number  `json:"priceChange"`
	PriceChangePercent types.Number  `json:"priceChangePercent"`
	LastPrice          types.Number  `json:"lastPrice"`
	LastQuantity       types.Number  `json:"lastQty"`
	HighPrice          types.Number  `json:"highPrice"`
	LowPrice           types.Number  `json:"lowPrice"`
	Volume             types.Number  `json:"volume"`
	QuoteVolume        types.Number  `json:"quoteVolume"`
	OpenPrice          types.Number  `json:"openPrice"`
	OpenTime           types.Time    `json:"openTime"`
	CloseTime          types.Time    `json:"closeTime"`
	AskPrice           types.Number  `json:"askPrice"`
	AskQuantity        types.Number  `json:"askQty"`
	BidPrice           types.Number  `json:"bidPrice"`
	BidQuantity        types.Number  `json:"bidQty"`
}

// SwapBookTickerResponse wraps the best order book entry returned by the symbol order book ticker endpoint.
type SwapBookTickerResponse struct {
	BookTicker SwapBookTicker `json:"book_ticker"`
}

// SwapBookTicker represents the best bid and ask price and quantity for a perpetual futures symbol.
type SwapBookTicker struct {
	Symbol       currency.Pair `json:"symbol"`
	BidPrice     float64       `json:"bid_price"`
	BidQuantity  float64       `json:"bid_qty"`
	AskPrice     float64       `json:"ask_price"`
	AskQuantity  float64       `json:"ask_qty"`
	LastUpdateID uint64        `json:"lastUpdateId"`
	Time         types.Time    `json:"time"`
}

// SwapMarkPriceCandle represents a mark price kline entry for a perpetual futures symbol.
type SwapMarkPriceCandle struct {
	Open      types.Number `json:"open"`
	Close     types.Number `json:"close"`
	High      types.Number `json:"high"`
	Low       types.Number `json:"low"`
	Volume    types.Number `json:"volume"`
	OpenTime  types.Time   `json:"openTime"`
	CloseTime types.Time   `json:"closeTime"`
}

// SwapPriceTicker represents the latest transaction price for a perpetual futures symbol.
type SwapPriceTicker struct {
	Symbol currency.Pair `json:"symbol"`
	Price  types.Number  `json:"price"`
	Time   types.Time    `json:"time"`
}

// SwapTradingRules represents the trading rules of a perpetual futures symbol
type SwapTradingRules struct {
	MinSizeCoin         types.Number `json:"minSizeCoin"`
	MinSizeUSD          types.Number `json:"minSizeUsd"`
	MaxNumOrder         types.Number `json:"maxNumOrder"`
	ProtectionThreshold types.Number `json:"protectionThreshold"`
	BuyMaxPrice         types.Number `json:"buyMaxPrice"`
	BuyMinPrice         types.Number `json:"buyMinPrice"`
	SellMaxPrice        types.Number `json:"sellMaxPrice"`
	SellMinPrice        types.Number `json:"sellMinPrice"`
	MarketRatio         types.Number `json:"marketRatio"`
}

// CoinMContract represents a coin-margined (Coin-M) futures contract and its trading rules.
type CoinMContract struct {
	Symbol         currency.Pair `json:"symbol"`
	PricePrecision uint64        `json:"pricePrecision"`
	MinTickSize    types.Number  `json:"minTickSize"`
	MinTradeValue  types.Number  `json:"minTradeValue"`
	MinQuantity    types.Number  `json:"minQty"`
	Status         uint64        `json:"status"`
	TimeOnline     types.Time    `json:"timeOnline"`
	DisplayName    string        `json:"displayName"`
}

// CoinMMarkPriceFundingRate represents the latest mark price, index price and funding rate for a coin-margined futures symbol.
type CoinMMarkPriceFundingRate struct {
	Symbol          currency.Pair `json:"symbol"`
	MarkPrice       types.Number  `json:"markPrice"`
	IndexPrice      types.Number  `json:"indexPrice"`
	LastFundingRate types.Number  `json:"lastFundingRate"`
	NextFundingTime types.Time    `json:"nextFundingTime"`
}

// CoinMOpenInterest represents the open interest of a coin-margined futures symbol.
type CoinMOpenInterest struct {
	OpenInterest types.Number  `json:"openInterest"`
	Symbol       currency.Pair `json:"symbol"`
	Timestamp    types.Time    `json:"timestamp"`
}

// CoinMCandle represents a coin-margined futures kline entry.
type CoinMCandle struct {
	Open   types.Number `json:"open"`
	Close  types.Number `json:"close"`
	High   types.Number `json:"high"`
	Low    types.Number `json:"low"`
	Volume types.Number `json:"volume"`
	Time   types.Time   `json:"time"`
}

// CoinMOrderbookDepth represents the order book depth of a coin-margined futures symbol.
type CoinMOrderbookDepth struct {
	Timestamp types.Time                       `json:"T"`
	Bids      orderbook.LevelsArrayPriceAmount `json:"bids"`
	Asks      orderbook.LevelsArrayPriceAmount `json:"asks"`
}

// CoinMTicker24Hr represents the 24-hour rolling window price change statistics for a coin-margined futures symbol.
type CoinMTicker24Hr struct {
	Symbol             currency.Pair `json:"symbol"`
	PriceChange        types.Number  `json:"priceChange"`
	PriceChangePercent string        `json:"priceChangePercent"`
	LastPrice          types.Number  `json:"lastPrice"`
	LastQuantity       types.Number  `json:"lastQty"`
	HighPrice          types.Number  `json:"highPrice"`
	LowPrice           types.Number  `json:"lowPrice"`
	Volume             types.Number  `json:"volume"`
	QuoteVolume        types.Number  `json:"quoteVolume"`
	OpenPrice          types.Number  `json:"openPrice"`
	CloseTime          types.Time    `json:"closeTime"`
	BidPrice           types.Number  `json:"bidPrice"`
	BidQuantity        types.Number  `json:"bidQty"`
	AskPrice           types.Number  `json:"askPrice"`
	AskQuantity        types.Number  `json:"askQty"`
}
