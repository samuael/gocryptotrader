package bingx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/margin"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

// GetCoinMContracts retrieves coin-margined (Coin-M) futures contracts and their trading rules.
func (e *Exchange) GetCoinMContracts(ctx context.Context, symbol currency.Pair) ([]*CoinMContract, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMContract
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/contracts", params, &resp, coinMV1MarketContractsEPL)
}

// GetCoinMMarkPriceAndFundingRate retrieves the latest mark price, index price and funding rate for coin-margined futures symbols.
func (e *Exchange) GetCoinMMarkPriceAndFundingRate(ctx context.Context, symbol currency.Pair) ([]*CoinMMarkPriceFundingRate, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMMarkPriceFundingRate
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/premiumIndex", params, &resp, coinMV1MarketPremiumIndexEPL)
}

// GetCoinMOpenInterest retrieves the open interest of coin-margined futures symbols.
func (e *Exchange) GetCoinMOpenInterest(ctx context.Context, symbol currency.Pair) ([]*CoinMOpenInterest, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMOpenInterest
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/openInterest", params, &resp, coinMV1MarketOpenInterestEPL)
}

// GetCoinMKlineData retrieves candlestick data for a coin-margined futures symbol. Only the last 30 days of data is available.
func (e *Exchange) GetCoinMKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, timeZone, limit uint64) ([]*CoinMCandle, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if interval == kline.Interval(0) {
		return nil, kline.ErrInvalidInterval
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("interval", intervalToString(interval))
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if timeZone != 0 {
		params.Set("timeZone", strconv.FormatUint(timeZone, 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*CoinMCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/klines", params, &resp, coinMV1MarketKlinesEPL)
}

// GetCoinMOrderbookDepth retrieves the order book depth for a coin-margined futures symbol.
func (e *Exchange) GetCoinMOrderbookDepth(ctx context.Context, symbol currency.Pair, limit uint64) (*CoinMOrderbookDepth, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp *CoinMOrderbookDepth
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/depth", params, &resp, coinMV1MarketDepthEPL)
}

// GetCoinM24HrTickerPriceChange retrieves the 24-hour rolling window price change statistics for coin-margined futures symbols.
func (e *Exchange) GetCoinM24HrTickerPriceChange(ctx context.Context, symbol currency.Pair) ([]*CoinMTicker24Hr, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMTicker24Hr
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/ticker", params, &resp, coinMV1MarketTickerEPL)
}

func setFloat(params url.Values, key string, value float64) {
	if value != 0 {
		params.Set(key, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

// PlaceCoinMOrder submits a coin-margined futures order.
func (e *Exchange) PlaceCoinMOrder(ctx context.Context, arg *PlaceCoinMOrderRequest) (*CoinMOrderAck, error) {
	if err := common.NilGuard(arg); err != nil {
		return nil, err
	}
	if arg.Symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if arg.Type == "" {
		return nil, order.ErrTypeIsInvalid
	}
	if arg.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	params := url.Values{}
	params.Set("symbol", arg.Symbol.String())
	params.Set("type", arg.Type)
	params.Set("side", arg.Side)
	if arg.PositionSide != "" {
		params.Set("positionSide", arg.PositionSide)
	}
	setFloat(params, "price", arg.Price)
	setFloat(params, "quantity", arg.Quantity)
	setFloat(params, "stopPrice", arg.StopPrice)
	setFloat(params, "stopLoss", arg.StopLoss)
	setFloat(params, "takeProfit", arg.TakeProfit)
	if arg.WorkingType != "" {
		params.Set("workingType", arg.WorkingType)
	}
	if arg.ClientOrderID != "" {
		params.Set("clientOrderId", arg.ClientOrderID)
	}
	if arg.TimeInForce != "" {
		params.Set("timeInForce", arg.TimeInForce)
	}
	if arg.RecvWindow != 0 {
		params.Set("recvWindow", strconv.FormatInt(arg.RecvWindow, 10))
	}
	var resp *CoinMOrderAck
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "cswap/v1/trade/order", params, &resp, coinMV1TradeOrderEPL)
}

// GetCoinMCommissionRate retrieves the maker and taker commission rates for coin-margined futures.
func (e *Exchange) GetCoinMCommissionRate(ctx context.Context) (*CoinMCommissionRate, error) {
	var resp *CoinMCommissionRate
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/user/commissionRate", nil, &resp, coinMV1UserCommissionRateEPL)
}

// GetCoinMLeverage retrieves the leverage and available positions for a coin-margined futures symbol.
func (e *Exchange) GetCoinMLeverage(ctx context.Context, symbol currency.Pair) (*CoinMLeverage, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *CoinMLeverage
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/leverage", params, &resp, coinMV1TradeLeverageEPL)
}

// SetCoinMLeverage sets the leverage for the long or short side of a coin-margined futures symbol.
func (e *Exchange) SetCoinMLeverage(ctx context.Context, symbol currency.Pair, side string, leverage int64) (*CoinMLeverage, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if leverage <= 0 {
		return nil, errLeverageRequired
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("side", side)
	params.Set("leverage", strconv.FormatInt(leverage, 10))
	var resp *CoinMLeverage
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "cswap/v1/trade/leverage", params, &resp, coinMV1TradeLeverageEPL)
}

// CancelAllCoinMOrders cancels all open coin-margined futures orders, optionally filtered by symbol.
func (e *Exchange) CancelAllCoinMOrders(ctx context.Context, symbol currency.Pair) (*CoinMCancelOrdersResponse, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *CoinMCancelOrdersResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, "cswap/v1/trade/allOpenOrders", params, &resp, coinMV1TradeAllOpenOrdersEPL)
}

// CloseAllCoinMPositions closes all open coin-margined futures positions in bulk, optionally filtered by symbol.
func (e *Exchange) CloseAllCoinMPositions(ctx context.Context, symbol currency.Pair) (*CoinMCloseAllPositionsResponse, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *CoinMCloseAllPositionsResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "cswap/v1/trade/closeAllPositions", params, &resp, coinMV1TradeCloseAllPositionsEPL)
}

// GetCoinMPositions retrieves the open coin-margined futures positions, optionally filtered by symbol.
func (e *Exchange) GetCoinMPositions(ctx context.Context, symbol currency.Pair) ([]*CoinMPosition, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMPosition
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/user/positions", params, &resp, coinMV1UserPositionsEPL)
}

// GetCoinMAccountAssets retrieves the coin-margined futures account asset balances.
func (e *Exchange) GetCoinMAccountAssets(ctx context.Context, symbol currency.Pair) ([]*CoinMAccountAsset, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMAccountAsset
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/user/balance", params, &resp, coinMV1UserBalanceEPL)
}

// GetCoinMForceOrders retrieves the coin-margined futures liquidation (force) order history.
func (e *Exchange) GetCoinMForceOrders(ctx context.Context, symbol currency.Pair, autoCloseType string, startTime, endTime time.Time, limit uint64) ([]*CoinMOrder, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if autoCloseType != "" {
		params.Set("autoCloseType", autoCloseType)
	}
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*CoinMOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/forceOrders", params, &resp, coinMV1TradeForceOrdersEPL)
}

// GetCoinMOrderTradeDetail retrieves the trade fills for a coin-margined futures order.
func (e *Exchange) GetCoinMOrderTradeDetail(ctx context.Context, orderID string, pageIndex, pageSize uint64) ([]*CoinMTradeDetail, error) {
	if orderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("orderId", orderID)
	if pageIndex != 0 {
		params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	}
	var resp []*CoinMTradeDetail
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/allFillOrders", params, &resp, coinMV1TradeAllFillOrdersEPL)
}

// CancelCoinMOrder cancels an active coin-margined futures order.
func (e *Exchange) CancelCoinMOrder(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*CoinMOrder, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if orderID == 0 && clientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp *CoinMOrderDetailResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, "cswap/v1/trade/cancelOrder", params, &resp, coinMV1TradeCancelOrderEPL); err != nil {
		return nil, err
	}
	return &resp.Order, nil
}

// GetCoinMOpenOrders retrieves all current open coin-margined futures orders, optionally filtered by symbol.
func (e *Exchange) GetCoinMOpenOrders(ctx context.Context, symbol currency.Pair) ([]CoinMOrder, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *CoinMOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/openOrders", params, &resp, coinMV1TradeOpenOrdersEPL); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetCoinMOrderDetail retrieves the details of a coin-margined futures order.
func (e *Exchange) GetCoinMOrderDetail(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*CoinMOrder, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if orderID == 0 && clientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp *CoinMOrderDetailResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/orderDetail", params, &resp, coinMV1TradeOrderDetailEPL); err != nil {
		return nil, err
	}
	return &resp.Order, nil
}

// GetCoinMOrderHistory retrieves the coin-margined futures order history.
func (e *Exchange) GetCoinMOrderHistory(ctx context.Context, symbol currency.Pair, orderID int64, startTime, endTime time.Time, limit uint64) ([]CoinMOrder, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp *CoinMOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/orderHistory", params, &resp, coinMV1TradeOrderHistoryEPL); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetCoinMMarginType retrieves the margin mode of coin-margined futures symbols.
func (e *Exchange) GetCoinMMarginType(ctx context.Context, symbol currency.Pair) ([]*CoinMMarginType, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp []*CoinMMarginType
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "cswap/v1/trade/marginType", params, &resp, coinMV1TradeMarginTypeEPL)
}

// SetCoinMMarginType sets the margin mode of a coin-margined futures symbol.
func (e *Exchange) SetCoinMMarginType(ctx context.Context, symbol currency.Pair, marginType string) error {
	if symbol.IsEmpty() {
		return currency.ErrCurrencyPairEmpty
	}
	if marginType == "" {
		return margin.ErrInvalidMarginType
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("marginType", marginType)
	return e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "cswap/v1/trade/marginType", params, &struct{}{}, coinMV1TradeMarginTypeEPL)
}

// AdjustCoinMIsolatedMargin increases or decreases the isolated margin of a coin-margined futures position.
func (e *Exchange) AdjustCoinMIsolatedMargin(ctx context.Context, symbol currency.Pair, amount float64, adjustmentType int64, positionSide string) error {
	if symbol.IsEmpty() {
		return currency.ErrCurrencyPairEmpty
	}
	if amount <= 0 {
		return errMarginAmountRequired
	}
	if adjustmentType == 0 {
		return errAdjustmentTypeRequired
	}
	if positionSide == "" {
		return order.ErrSideIsInvalid
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	setFloat(params, "amount", amount)
	params.Set("type", strconv.FormatInt(adjustmentType, 10))
	params.Set("positionSide", positionSide)
	return e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "cswap/v1/trade/positionMargin", params, &struct{}{}, coinMV1TradePositionMarginEPL)
}
