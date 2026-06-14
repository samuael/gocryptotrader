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

// ---------------------------- Swap Trade Endpoints ----------------------

// PlaceSwapOrderRequest holds the parameters for placing or testing a perpetual futures order.
type PlaceSwapOrderRequest struct {
	Symbol          currency.Pair
	Type            string
	Side            string
	PositionSide    string
	ReduceOnly      string
	Price           float64
	Quantity        float64
	QuoteOrderQty   float64
	StopPrice       float64
	PriceRate       float64
	WorkingType     string
	StopLoss        float64
	TakeProfit      float64
	ClientOrderID   string
	TimeInForce     string
	ClosePosition   string
	ActivationPrice float64
	StopGuaranteed  string
	PositionID      int64
	RecvWindow      int64
}

// SwapOrderAck represents the acknowledgement returned when placing or testing a perpetual futures order.
type SwapOrderAck struct {
	Symbol         currency.Pair `json:"symbol"`
	OrderID        int64         `json:"orderId"`
	OrderIDString  string        `json:"orderID"`
	Side           string        `json:"side"`
	PositionSide   string        `json:"positionSide"`
	Type           string        `json:"type"`
	ReduceOnly     string        `json:"reduceOnly"`
	ClientOrderID  string        `json:"clientOrderId"`
	WorkingType    string        `json:"workingType"`
	Status         string        `json:"status"`
	StopGuaranteed string        `json:"stopGuaranteed"`
	AveragePrice   types.Number  `json:"avgPrice"`
	ExecutedQty    types.Number  `json:"executedQty"`
}

// SwapOrderAckResponse wraps a single order acknowledgement.
type SwapOrderAckResponse struct {
	Order SwapOrderAck `json:"order"`
}

// SwapOrderAckListResponse wraps a list of order acknowledgements returned by the batch order endpoint.
type SwapOrderAckListResponse struct {
	Orders []SwapOrderAck `json:"orders"`
}

// SwapOrder represents a perpetual futures order's full details.
type SwapOrder struct {
	Time           types.Time    `json:"time"`
	Symbol         currency.Pair `json:"symbol"`
	Side           string        `json:"side"`
	Type           string        `json:"type"`
	PositionSide   string        `json:"positionSide"`
	ReduceOnly     string        `json:"reduceOnly"`
	CumQuote       types.Number  `json:"cumQuote"`
	Status         string        `json:"status"`
	StopPrice      types.Number  `json:"stopPrice"`
	Price          types.Number  `json:"price"`
	OrigQty        types.Number  `json:"origQty"`
	AveragePrice   types.Number  `json:"avgPrice"`
	ExecutedQty    types.Number  `json:"executedQty"`
	OrderID        int64         `json:"orderId"`
	Profit         types.Number  `json:"profit"`
	Commission     types.Number  `json:"commission"`
	UpdateTime     types.Time    `json:"updateTime"`
	ClientOrderID  string        `json:"clientOrderId"`
	WorkingType    string        `json:"workingType"`
	StopGuaranteed string        `json:"stopGuaranteed"`
	TriggerOrderID int64         `json:"triggerOrderId"`
	ClosePosition  string        `json:"closePosition"`
	Leverage       string        `json:"leverage"`
	PostOnly       bool          `json:"postOnly"`
	IsTwap         bool          `json:"isTwap"`
	MainOrderID    string        `json:"mainOrderId"`
}

// SwapOrderDetailResponse wraps a single order's full details.
type SwapOrderDetailResponse struct {
	Order SwapOrder `json:"order"`
}

// SwapOrdersResponse wraps a list of perpetual futures orders.
type SwapOrdersResponse struct {
	Orders []SwapOrder `json:"orders"`
}

// SwapFailedOrder represents an order that could not be cancelled.
type SwapFailedOrder struct {
	OrderID      int64  `json:"orderId"`
	ClientID     string `json:"clientOrderId"`
	ErrorCode    int64  `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// SwapCancelOrdersResponse represents the result of cancelling multiple or all open orders.
type SwapCancelOrdersResponse struct {
	Success []SwapOrder       `json:"success"`
	Failed  []SwapFailedOrder `json:"failed"`
}

// SwapModifyOrderResponse represents the result of amending an order's quantity.
type SwapModifyOrderResponse struct {
	OrderID       int64         `json:"orderId"`
	ClientOrderID string        `json:"clientOrderId"`
	Symbol        currency.Pair `json:"symbol"`
	Quantity      types.Number  `json:"quantity"`
}

// SwapCloseAllPositionsResponse represents the order numbers generated by a one-click liquidation.
type SwapCloseAllPositionsResponse struct {
	Success []int64 `json:"success"`
	Failed  []int64 `json:"failed"`
}

// SwapMarginTypeResponse represents the margin mode of a perpetual futures symbol.
type SwapMarginTypeResponse struct {
	MarginType string `json:"marginType"`
}

// SwapLeverage represents the leverage and available position information for a perpetual futures symbol.
type SwapLeverage struct {
	Leverage            int64        `json:"leverage"`
	Symbol              string       `json:"symbol"`
	LongLeverage        int64        `json:"longLeverage"`
	ShortLeverage       int64        `json:"shortLeverage"`
	MaxLongLeverage     int64        `json:"maxLongLeverage"`
	MaxShortLeverage    int64        `json:"maxShortLeverage"`
	AvailableLongVol    types.Number `json:"availableLongVol"`
	AvailableShortVol   types.Number `json:"availableShortVol"`
	AvailableLongVal    types.Number `json:"availableLongVal"`
	AvailableShortVal   types.Number `json:"availableShortVal"`
	MaxPositionLongVal  types.Number `json:"maxPositionLongVal"`
	MaxPositionShortVal types.Number `json:"maxPositionShortVal"`
}

// SwapPositionMarginResponse represents the result of modifying isolated position margin.
type SwapPositionMarginResponse struct {
	Amount     float64 `json:"amount"`
	Type       int64   `json:"type"`
	PositionID int64   `json:"positionId"`
}

// SwapFillOrder represents a historical transaction order for a perpetual futures symbol.
type SwapFillOrder struct {
	FilledTime            types.Time    `json:"filledTime"`
	FilledTimestamp       string        `json:"filledTm"`
	Symbol                currency.Pair `json:"symbol"`
	Volume                types.Number  `json:"volume"`
	Price                 types.Number  `json:"price"`
	Amount                types.Number  `json:"amount"`
	Commission            types.Number  `json:"commission"`
	Currency              currency.Code `json:"currency"`
	OrderID               string        `json:"orderId"`
	LiquidatedPrice       types.Number  `json:"liquidatedPrice"`
	LiquidatedMarginRatio types.Number  `json:"liquidatedMarginRatio"`
	WorkingType           string        `json:"workingType"`
	Side                  string        `json:"side"`
	Type                  string        `json:"type"`
	PositionSide          string        `json:"positionSide"`
	ClientOrderID         string        `json:"clientOrderId"`
	OnlyOnePosition       bool          `json:"onlyOnePosition"`
}

// SwapFillOrdersResponse wraps a list of historical transaction orders.
type SwapFillOrdersResponse struct {
	FillOrders []SwapFillOrder `json:"fill_orders"`
}

// SwapPositionModeResponse represents the dual side position mode setting.
type SwapPositionModeResponse struct {
	DualSidePosition string `json:"dualSidePosition"`
}

// SwapCancelReplaceOrder represents a cancelled or newly created order within a cancel-replace operation.
type SwapCancelReplaceOrder struct {
	CancelClientOrderID string        `json:"cancelClientOrderId"`
	CancelOrderID       int64         `json:"cancelOrderId"`
	Symbol              currency.Pair `json:"symbol"`
	OrderID             int64         `json:"orderId"`
	Side                string        `json:"side"`
	PositionSide        string        `json:"positionSide"`
	Type                string        `json:"type"`
	OrigQty             types.Number  `json:"origQty"`
	Price               types.Number  `json:"price"`
	ExecutedQty         types.Number  `json:"executedQty"`
	AveragePrice        types.Number  `json:"avgPrice"`
	CumQuote            types.Number  `json:"cumQuote"`
	StopPrice           types.Number  `json:"stopPrice"`
	Profit              types.Number  `json:"profit"`
}

// SwapCancelReplaceResponse represents the result of cancelling an existing order and sending a new one.
type SwapCancelReplaceResponse struct {
	CancelResult     string                  `json:"cancelResult"`
	CancelMsg        string                  `json:"cancelMsg"`
	CancelResponse   *SwapCancelReplaceOrder `json:"cancelResponse"`
	ReplaceResult    string                  `json:"replaceResult"`
	ReplaceMsg       string                  `json:"replaceMsg"`
	NewOrderResponse *SwapCancelReplaceOrder `json:"newOrderResponse"`
}

// SwapCancelAllAfterResponse represents the result of setting a countdown to cancel all orders.
type SwapCancelAllAfterResponse struct {
	TriggerTime types.Time `json:"triggerTime"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
}

// SwapClosePositionResponse represents the order created when closing a position by ID.
type SwapClosePositionResponse struct {
	OrderID      int64         `json:"orderId"`
	PositionID   string        `json:"positionId"`
	Symbol       currency.Pair `json:"symbol"`
	Side         string        `json:"side"`
	Type         string        `json:"type"`
	PositionSide string        `json:"positionSide"`
	OrigQty      types.Number  `json:"origQty"`
}

// SwapMaintMarginRatio represents a position and maintenance margin ratio tier for a perpetual futures symbol.
type SwapMaintMarginRatio struct {
	Tier             string        `json:"tier"`
	Symbol           currency.Pair `json:"symbol"`
	MinPositionVal   types.Number  `json:"minPositionVal"`
	MaxPositionVal   types.Number  `json:"maxPositionVal"`
	MaintMarginRatio types.Number  `json:"maintMarginRatio"`
	MaintAmount      types.Number  `json:"maintAmount"`
}

// SwapFillDetail represents a historical transaction detail for a perpetual futures symbol.
type SwapFillDetail struct {
	Symbol          currency.Pair `json:"symbol"`
	Quantity        types.Number  `json:"qty"`
	QuoteQty        types.Number  `json:"quoteQty"`
	Price           types.Number  `json:"price"`
	Commission      types.Number  `json:"commission"`
	CommissionAsset currency.Code `json:"commissionAsset"`
	TradeID         string        `json:"tradeId"`
	OrderID         string        `json:"orderId"`
	FilledTime      string        `json:"filledTime"`
	Side            string        `json:"side"`
	PositionSide    string        `json:"positionSide"`
	Role            string        `json:"role"`
}

// SwapFillHistoryResponse wraps a list of historical transaction details.
type SwapFillHistoryResponse struct {
	FillHistoryOrders []SwapFillDetail `json:"fill_history_orders"`
	Total             int64            `json:"total"`
}

// SwapPositionHistory represents a closed position record for a perpetual futures symbol.
type SwapPositionHistory struct {
	Symbol             currency.Pair `json:"symbol"`
	PositionID         string        `json:"positionId"`
	PositionSide       string        `json:"positionSide"`
	Isolated           bool          `json:"isolated"`
	CloseAllPositions  bool          `json:"closeAllPositions"`
	PositionAmt        types.Number  `json:"positionAmt"`
	ClosePositionAmt   types.Number  `json:"closePositionAmt"`
	RealisedProfit     types.Number  `json:"realisedProfit"`
	NetProfit          types.Number  `json:"netProfit"`
	AverageClosePrice  types.Number  `json:"avgClosePrice"`
	AveragePrice       types.Number  `json:"avgPrice"`
	Leverage           int64         `json:"leverage"`
	PositionCommission types.Number  `json:"positionCommission"`
	TotalFunding       types.Number  `json:"totalFunding"`
	OpenTime           types.Time    `json:"openTime"`
	UpdateTime         types.Time    `json:"updateTime"`
}

// SwapMarginChange represents a single isolated margin change record.
type SwapMarginChange struct {
	Symbol            currency.Pair `json:"symbol"`
	PositionID        string        `json:"positionId"`
	ChangeReason      string        `json:"changeReason"`
	MarginChange      types.Number  `json:"marginChange"`
	MarginAfterChange types.Number  `json:"marginAfterChange"`
	Time              types.Time    `json:"time"`
}

// SwapMarginChangeHistoryResponse wraps a list of isolated margin change records.
type SwapMarginChangeHistoryResponse struct {
	Records []SwapMarginChange `json:"records"`
	Total   int64              `json:"total"`
}

// SwapVSTResponse represents the result of applying for VST testnet funds.
type SwapVSTResponse struct {
	AdjustType string       `json:"adjustType"`
	Amount     types.Number `json:"amount"`
}

// PlaceTWAPOrderRequest holds the parameters for placing a TWAP order.
type PlaceTWAPOrderRequest struct {
	Symbol         currency.Pair
	Side           string
	PositionSide   string
	PriceType      string
	PriceVariance  string
	TriggerPrice   string
	Interval       int64
	AmountPerOrder string
	TotalAmount    string
	RecvWindow     int64
}

// SwapTWAPOrderResponse represents the acknowledgement returned when placing a TWAP order.
type SwapTWAPOrderResponse struct {
	MainOrderID string `json:"mainOrderId"`
}

// SwapTWAPOrder represents a TWAP order's details.
type SwapTWAPOrder struct {
	Symbol         currency.Pair `json:"symbol"`
	MainOrderID    string        `json:"mainOrderId"`
	Side           string        `json:"side"`
	PositionSide   string        `json:"positionSide"`
	PriceType      string        `json:"priceType"`
	PriceVariance  types.Number  `json:"priceVariance"`
	TriggerPrice   types.Number  `json:"triggerPrice"`
	Interval       int64         `json:"interval"`
	AmountPerOrder types.Number  `json:"amountPerOrder"`
	TotalAmount    types.Number  `json:"totalAmount"`
	OrderStatus    string        `json:"orderStatus"`
	ExecutedQty    types.Number  `json:"executedQty"`
	Duration       int64         `json:"duration"`
	MaxDuration    int64         `json:"maxDuration"`
	CreatedTime    types.Time    `json:"createdTime"`
	UpdateTime     types.Time    `json:"updateTime"`
}

// SwapTWAPOrdersResponse wraps a list of TWAP orders.
type SwapTWAPOrdersResponse struct {
	List  []SwapTWAPOrder `json:"list"`
	Total int64           `json:"total"`
}

// SwapAssetModeResponse represents the multi-assets margin mode setting.
type SwapAssetModeResponse struct {
	AssetMode string `json:"assetMode"`
}

// SwapMultiAssetsRule represents a margin asset rule under multi-assets mode.
type SwapMultiAssetsRule struct {
	MarginAssets         currency.Code `json:"marginAssets"`
	LoanToValue          string        `json:"ltv"`
	CollateralValueRatio string        `json:"collateralValueRatio"`
	MaxTransfer          types.Number  `json:"maxTransfer"`
	IndexPrice           types.Number  `json:"indexPrice"`
}

// SwapMultiAssetsMargin represents a margin asset balance under multi-assets mode.
type SwapMultiAssetsMargin struct {
	Currency             currency.Code `json:"currency"`
	TotalAmount          types.Number  `json:"totalAmount"`
	AvailableTransfer    types.Number  `json:"availableTransfer"`
	LatestMortgageAmount types.Number  `json:"latestMortgageAmount"`
}

// SwapReversePositionResponse represents the result of a one-click reverse position operation.
type SwapReversePositionResponse struct {
	Type               string        `json:"type"`
	PositionID         string        `json:"positionId"`
	NewPositionID      string        `json:"newPositionId"`
	Symbol             currency.Pair `json:"symbol"`
	PositionSide       string        `json:"positionSide"`
	Isolated           bool          `json:"isolated"`
	PositionAmt        types.Number  `json:"positionAmt"`
	AvailableAmt       types.Number  `json:"availableAmt"`
	UnrealizedProfit   types.Number  `json:"unrealizedProfit"`
	RealisedProfit     types.Number  `json:"realisedProfit"`
	InitialMargin      types.Number  `json:"initialMargin"`
	Margin             types.Number  `json:"margin"`
	LiquidationPrice   float64       `json:"liquidationPrice"`
	AveragePrice       types.Number  `json:"avgPrice"`
	Leverage           int64         `json:"leverage"`
	PositionValue      types.Number  `json:"positionValue"`
	MarkPrice          types.Number  `json:"markPrice"`
	RiskRate           types.Number  `json:"riskRate"`
	MaxMarginReduction types.Number  `json:"maxMarginReduction"`
	PnlRatio           types.Number  `json:"pnlRatio"`
	UpdateTime         types.Time    `json:"updateTime"`
}

// SwapAutoAddMarginResponse represents the result of toggling automatic margin addition for a hedge-mode position.
type SwapAutoAddMarginResponse struct {
	Symbol         currency.Pair `json:"symbol"`
	PositionID     int64         `json:"positionId"`
	FunctionSwitch string        `json:"functionSwitch"`
	Amount         types.Number  `json:"amount"`
}
