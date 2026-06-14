package bingx

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/encoding/json"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/margin"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

var (
	errLeverageRequired        = errors.New("leverage is required")
	errPositionIDRequired      = errors.New("position id is required")
	errAdjustmentTypeRequired  = errors.New("adjustment direction type is required")
	errMarginAmountRequired    = errors.New("margin amount is required")
	errTradingUnitRequired     = errors.New("trading unit is required")
	errPositionModeRequired    = errors.New("dual side position mode is required")
	errCancelReplaceModeEmpty  = errors.New("cancel replace mode is required")
	errBatchOrdersRequired     = errors.New("batch orders payload is required")
	errCountdownTypeRequired   = errors.New("countdown request type is required")
	errCountdownTimeoutInvalid = errors.New("countdown timeout must be between 10 and 120 seconds")
	errAssetModeRequired       = errors.New("multi-assets mode is required")
	errReverseTypeRequired     = errors.New("reverse type is required")
	errFunctionSwitchRequired  = errors.New("function switch is required")
	errTriggerPriceRequired    = errors.New("trigger price is required")
)

// Exchange implements exchange.IBotExchange for BingX exchange
type Exchange struct {
	exchange.Base
}

const (
	apiURL     = "https://open-api.bingx.com/openApi/"
	apiVersion = ""
)

// GetSpotSymbols query spot trading symbol information, including trading rules such as minimum/maximum notional amounts, price tick size, and quantity step size.
func (e *Exchange) GetSpotSymbols(ctx context.Context, symbol currency.Pair) (*ConfigSymbols, error) {
	params := url.Values{}
	if symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *ConfigSymbols
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v1/common/symbols", params, &resp)
}

// GetMarketTrades retrieve the most recent market trades for a specified symbol.
func (e *Exchange) GetMarketTrades(ctx context.Context, symbol currency.Pair, limit uint64) ([]*SpotTradeInfo, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SpotTradeInfo
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v1/market/trades", params, &resp)
}

// GetSpotOrderbookDepth get the order book depth for a specified symbol
func (e *Exchange) GetSpotOrderbookDepth(ctx context.Context, symbol currency.Pair, limit uint64) (*SpotOrderbookDetail, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp *SpotOrderbookDetail
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v1/market/depth", params, &resp)
}

// GetSpotKlineData load candlestick data for aspecific symbol.
func (e *Exchange) GetSpotKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, timezone, limit uint64) ([]*SpotCandle, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if timezone != 0 {
		params.Set("timezone", strconv.FormatUint(timezone, 10))
	}
	if interval != kline.Interval(0) {
		params.Set("interval", intervalToString(interval))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SpotCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v2/market/kline", params, &resp)
}

// intervalToString converts a kline.Interval into the string representation expected by the BingX API.
func intervalToString(interval kline.Interval) string {
	switch interval {
	case kline.OneDay:
		return "1d"
	case kline.ThreeDay:
		return "3d"
	case kline.OneWeek:
		return "1w"
	case kline.OneMonth:
		return "1M"
	default:
		return interval.Short()
	}
}

// Get24HrTickerPriceChange retrieves the 24-hour rolling window price change statistics
func (e *Exchange) Get24HrTickerPriceChange(ctx context.Context, symbol currency.Pair) ([]*SpotTicker24Hr, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*SpotTicker24Hr
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v1/ticker/24hr", params, &resp)
}

// GetSpotOrderbookAggregation retrieves the aggregated order book depth for a symbol
func (e *Exchange) GetSpotOrderbookAggregation(ctx context.Context, symbol currency.Pair, depth uint64, depthType string) (*SpotOrderbookDetail, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if depth != 0 {
		params.Set("depth", strconv.FormatUint(depth, 10))
	}
	if depthType != "" {
		params.Set("type", depthType)
	}
	var resp *SpotOrderbookDetail
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v2/market/depth", params, &resp)
}

// GetSpotSymbolPriceTicker retrieves the latest market price for a symbol.
func (e *Exchange) GetSpotSymbolPriceTicker(ctx context.Context, symbol currency.Pair) ([]*SpotPriceTicker, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp []*SpotPriceTicker
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v2/ticker/price", params, &resp)
}

// GetSpotSymbolOrderbookTicker retrieves the best bid and ask price and quantity for a symbol.
func (e *Exchange) GetSpotSymbolOrderbookTicker(ctx context.Context, symbol currency.Pair) ([]*SpotBookTicker, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp []*SpotBookTicker
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v1/ticker/bookTicker", params, &resp)
}

// GetHistoricalKlineData retrieves historical candlestick data for a symbol over different time intervals.
func (e *Exchange) GetHistoricalKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, limit uint64) ([]*SpotCandle, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if interval != kline.Interval(0) {
		params.Set("interval", intervalToString(interval))
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
	var resp []*SpotCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "market/his/v1/kline", params, &resp)
}

// GetOldSpotTrades retrieves historical trade records for a symbol.
func (e *Exchange) GetOldSpotTrades(ctx context.Context, symbol currency.Pair, limit uint64, fromID string) ([]*HistoricalSpotTrade, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	if fromID != "" {
		params.Set("fromId", fromID)
	}
	var resp []*HistoricalSpotTrade
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "market/his/v1/trade", params, &resp)
}

// params validates and converts a PlaceSpotOrderRequest into url.Values
func (r *PlaceSpotOrderRequest) params() (url.Values, error) {
	if r == nil {
		return nil, common.ErrNilPointer
	}
	if r.Symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if r.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if r.OrderType == "" {
		return nil, order.ErrTypeIsInvalid
	}
	params := url.Values{}
	params.Set("symbol", r.Symbol.String())
	params.Set("side", r.Side)
	params.Set("type", r.OrderType)
	setFloat(params, "stopPrice", r.StopPrice)
	setFloat(params, "quantity", r.Quantity)
	setFloat(params, "quoteOrderQty", r.QuoteOrderQty)
	setFloat(params, "price", r.Price)
	if r.NewClientOrderID != "" {
		params.Set("newClientOrderId", r.NewClientOrderID)
	}
	if r.TimeInForce != "" {
		params.Set("timeInForce", r.TimeInForce)
	}
	if r.RecvWindow != 0 {
		params.Set("recvWindow", strconv.FormatInt(r.RecvWindow, 10))
	}
	return params, nil
}

// PlaceSpotOrder submits a spot order.
func (e *Exchange) PlaceSpotOrder(ctx context.Context, arg *PlaceSpotOrderRequest) (*SpotOrder, error) {
	params, err := arg.params()
	if err != nil {
		return nil, err
	}
	var resp *SpotOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/order", params, &resp)
}

// PlaceMultipleSpotOrders submits up to five spot orders in a single request.
func (e *Exchange) PlaceMultipleSpotOrders(ctx context.Context, args []PlaceSpotOrderRequest, sync bool) ([]*SpotOrder, error) {
	if len(args) == 0 {
		return nil, errBatchOrdersRequired
	}
	orders := make([]map[string]string, 0, len(args))
	for i := range args {
		params, err := args[i].params()
		if err != nil {
			return nil, err
		}
		orderParams := make(map[string]string, len(params))
		for key := range params {
			orderParams[key] = params.Get(key)
		}
		orders = append(orders, orderParams)
	}
	payload, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("data", string(payload))
	if sync {
		params.Set("sync", "true")
	}
	var resp *SpotOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/batchOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// CancelSpotOrder cancels an active spot order.
func (e *Exchange) CancelSpotOrder(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID, cancelRestrictions string) (*SpotOrder, error) {
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
		params.Set("clientOrderID", clientOrderID)
	}
	if cancelRestrictions != "" {
		params.Set("cancelRestrictions", cancelRestrictions)
	}
	var resp *SpotOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/cancel", params, &resp)
}

// CancelMultipleSpotOrders cancels multiple spot orders by order ID or client order ID
func (e *Exchange) CancelMultipleSpotOrders(ctx context.Context, symbol currency.Pair, process uint64, orderIDs, clientOrderIDs string) (*SpotCancelOrdersResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if orderIDs == "" && clientOrderIDs == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if process != 0 {
		params.Set("process", strconv.FormatUint(process, 10))
	}
	if orderIDs != "" {
		params.Set("orderIds", orderIDs)
	}
	if clientOrderIDs != "" {
		params.Set("clientOrderIDs", clientOrderIDs)
	}
	var resp *SpotCancelOrdersResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/cancelOrders", params, &resp)
}

// CancelAllSpotOpenOrders cancels all open spot orders
func (e *Exchange) CancelAllSpotOpenOrders(ctx context.Context, symbol currency.Pair) ([]*SpotOrder, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *SpotOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/cancelOpenOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// CancelReplaceSpotOrder cancels an existing spot order and submits a new one.
func (e *Exchange) CancelReplaceSpotOrder(ctx context.Context, cancelReplaceMode string, cancelOrderID int64, cancelClientOrderID, cancelRestrictions string, arg *PlaceSpotOrderRequest) (*SpotCancelReplaceResponse, error) {
	if cancelReplaceMode == "" {
		return nil, errCancelReplaceModeEmpty
	}
	if cancelOrderID == 0 && cancelClientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params, err := arg.params()
	if err != nil {
		return nil, err
	}
	params.Set("cancelReplaceMode", cancelReplaceMode)
	if cancelOrderID != 0 {
		params.Set("cancelOrderId", strconv.FormatInt(cancelOrderID, 10))
	}
	if cancelClientOrderID != "" {
		params.Set("cancelClientOrderID", cancelClientOrderID)
	}
	if cancelRestrictions != "" {
		params.Set("cancelRestrictions", cancelRestrictions)
	}
	var resp *SpotCancelReplaceResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/order/cancelReplace", params, &resp)
}

// GetSpotOrderDetails retrieves the details of a spot order
func (e *Exchange) GetSpotOrderDetails(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*SpotOrderDetail, error) {
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
		params.Set("clientOrderID", clientOrderID)
	}
	var resp *SpotOrderDetail
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/trade/query", params, &resp)
}

// GetSpotOpenOrders retrieves all current open spot orders.
func (e *Exchange) GetSpotOpenOrders(ctx context.Context, symbol currency.Pair) ([]*SpotOrderDetail, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *SpotOrderDetailsResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/trade/openOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetSpotOrderHistory retrieves the spot order history
func (e *Exchange) GetSpotOrderHistory(ctx context.Context, symbol currency.Pair, orderID int64, status, orderType string, startTime, endTime time.Time, pageIndex, pageSize uint64) ([]*SpotOrderDetail, error) {
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
	if pageIndex != 0 {
		params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	}
	if status != "" {
		params.Set("status", status)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	var resp *SpotOrderDetailsResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/trade/historyOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetSpotTransactionDetails retrieves the spot trade fills
func (e *Exchange) GetSpotTransactionDetails(ctx context.Context, symbol currency.Pair, orderID int64, startTime, endTime time.Time, fromID, limit uint64) ([]*SpotTrade, error) {
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
	if fromID != 0 {
		params.Set("fromId", strconv.FormatUint(fromID, 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp *SpotTradesResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/trade/myTrades", params, &resp); err != nil {
		return nil, err
	}
	return resp.Fills, nil
}

// GetSpotCommissionRate retrieves the maker and taker commission rates for a spot symbol.
func (e *Exchange) GetSpotCommissionRate(ctx context.Context, symbol currency.Pair) (*SpotCommissionRate, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SpotCommissionRate
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/user/commissionRate", params, &resp)
}

// SetSpotCancelAllAfter sets or clears a countdown that cancels all open spot orders after a timeout.
func (e *Exchange) SetSpotCancelAllAfter(ctx context.Context, requestType string, timeout int64) (*SpotCancelAllAfterResponse, error) {
	if requestType == "" {
		return nil, errCountdownTypeRequired
	}
	if timeout < 10 || timeout > 120 {
		return nil, errCountdownTimeoutInvalid
	}
	params := url.Values{}
	params.Set("type", requestType)
	params.Set("timeOut", strconv.FormatInt(timeout, 10))
	var resp *SpotCancelAllAfterResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/trade/cancelAllAfter", params, &resp)
}

// CreateSpotOCOOrder creates a one-cancels-the-other (OCO) order.
func (e *Exchange) CreateSpotOCOOrder(ctx context.Context, arg *CreateSpotOCOOrderRequest) ([]*SpotOCOOrder, error) {
	if err := common.NilGuard(arg); err != nil {
		return nil, err
	}
	if arg.Symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if arg.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if arg.Quantity <= 0 {
		return nil, order.ErrAmountIsInvalid
	}
	if arg.LimitPrice <= 0 {
		return nil, fmt.Errorf("%w: limit price is required", order.ErrPriceMustBeSetIfLimitOrder)
	}
	if arg.OrderPrice <= 0 {
		return nil, fmt.Errorf("%w: order price is required", order.ErrPriceMustBeSetIfLimitOrder)
	}
	if arg.TriggerPrice <= 0 {
		return nil, errTriggerPriceRequired
	}
	params := url.Values{}
	params.Set("symbol", arg.Symbol.String())
	params.Set("side", arg.Side)
	setFloat(params, "quantity", arg.Quantity)
	setFloat(params, "limitPrice", arg.LimitPrice)
	setFloat(params, "orderPrice", arg.OrderPrice)
	setFloat(params, "triggerPrice", arg.TriggerPrice)
	if arg.ListClientOrderID != "" {
		params.Set("listClientOrderId", arg.ListClientOrderID)
	}
	if arg.AboveClientOrderID != "" {
		params.Set("aboveClientOrderId", arg.AboveClientOrderID)
	}
	if arg.BelowClientOrderID != "" {
		params.Set("belowClientOrderId", arg.BelowClientOrderID)
	}
	if arg.RecvWindow != 0 {
		params.Set("recvWindow", strconv.FormatInt(arg.RecvWindow, 10))
	}
	var resp []*SpotOCOOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/oco/order", params, &resp)
}

// CancelSpotOCOOrder cancels an OCO order list by order ID or client order ID.
func (e *Exchange) CancelSpotOCOOrder(ctx context.Context, orderID, clientOrderID string) (*SpotOCOCancelResponse, error) {
	if orderID == "" && clientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp *SpotOCOCancelResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "spot/v1/oco/cancel", params, &resp)
}

// GetSpotOCOOrderList retrieves a single OCO order list by group ID or client order ID.
func (e *Exchange) GetSpotOCOOrderList(ctx context.Context, orderListID, clientOrderID string) ([]*SpotOCOOrder, error) {
	if orderListID == "" && clientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	if orderListID != "" {
		params.Set("orderListId", orderListID)
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp []*SpotOCOOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/oco/orderList", params, &resp)
}

// GetSpotOpenOCOOrders retrieves all open OCO order lists.
func (e *Exchange) GetSpotOpenOCOOrders(ctx context.Context, pageIndex, pageSize uint64) ([]*SpotOCOOrder, error) {
	params := url.Values{}
	params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	var resp []*SpotOCOOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/oco/openOrderList", params, &resp)
}

// GetSpotOCOHistory retrieves the OCO order list history.
func (e *Exchange) GetSpotOCOHistory(ctx context.Context, startTime, endTime time.Time, pageIndex, pageSize uint64) ([]*SpotOCOOrder, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []*SpotOCOOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "spot/v1/oco/historyOrderList", params, &resp)
}

// --------------------- SWAP Endpoints -----------------------------

// GetSwapContracts retrieves USDT-M perpetual futures contracts and their trading rules.
func (e *Exchange) GetSwapContracts(ctx context.Context, symbol currency.Pair) ([]*SwapContractDetail, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*SwapContractDetail
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/contracts", params, &resp)
}

// GetSwapOrderbookDepth retrieves the order book depth for a perpetual futures symbol.
func (e *Exchange) GetSwapOrderbookDepth(ctx context.Context, symbol currency.Pair, limit uint64) (*SwapOrderbookDepth, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp *SwapOrderbookDepth
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/depth", params, &resp)
}

// GetSwapRecentTrades retrieves the most recent trades for a perpetual futures symbol.
func (e *Exchange) GetSwapRecentTrades(ctx context.Context, symbol currency.Pair, limit uint64) ([]*SwapTrade, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SwapTrade
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/trades", params, &resp)
}

// GetSwapMarkPriceAndFundingRate retrieves the current mark price, index price and funding rate.
func (e *Exchange) GetSwapMarkPriceAndFundingRate(ctx context.Context, symbol currency.Pair) ([]*SwapMarkPriceFundingRate, error) {
	if !symbol.IsEmpty() {
		params := url.Values{}
		params.Set("symbol", symbol.String())
		var resp *SwapMarkPriceFundingRate
		if err := e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/premiumIndex", params, &resp); err != nil {
			return nil, err
		}
		return []*SwapMarkPriceFundingRate{resp}, nil
	}
	var resp []*SwapMarkPriceFundingRate
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/premiumIndex", nil, &resp)
}

// GetSwapFundingRateHistory retrieves historical funding rates.
func (e *Exchange) GetSwapFundingRateHistory(ctx context.Context, symbol currency.Pair, startTime, endTime time.Time, limit uint64) ([]*SwapFundingRate, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
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
	var resp []*SwapFundingRate
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/fundingRate", params, &resp)
}

// GetSwapKlineData retrieves candlestick data for a perpetual futures symbol.
func (e *Exchange) GetSwapKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, timeZone, limit uint64) ([]*SwapCandle, error) {
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
	var resp []*SwapCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v3/quote/klines", params, &resp)
}

// GetSwapOpenInterest retrieves the open interest of a perpetual futures symbol.
func (e *Exchange) GetSwapOpenInterest(ctx context.Context, symbol currency.Pair) (*SwapOpenInterest, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SwapOpenInterest
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/openInterest", params, &resp)
}

// GetSwap24HrTickerPriceChange retrieves the 24-hour rolling window price change statistics.
func (e *Exchange) GetSwap24HrTickerPriceChange(ctx context.Context, symbol currency.Pair) ([]*SwapTicker24Hr, error) {
	// The endpoint returns a single object when a symbol is provided and an array otherwise
	if !symbol.IsEmpty() {
		params := url.Values{}
		params.Set("symbol", symbol.String())
		var resp *SwapTicker24Hr
		if err := e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/ticker", params, &resp); err != nil {
			return nil, err
		}
		return []*SwapTicker24Hr{resp}, nil
	}
	var resp []*SwapTicker24Hr
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/ticker", nil, &resp)
}

// GetSwapHistoricalTrades retrieves historical trade records of all users for a perpetual futures symbol.
func (e *Exchange) GetSwapHistoricalTrades(ctx context.Context, symbol currency.Pair, fromID string, limit uint64) ([]*SwapTrade, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if fromID != "" {
		params.Set("fromId", fromID)
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SwapTrade
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v1/market/historicalTrades", params, &resp)
}

// GetSwapSymbolOrderbookTicker retrieves the best bid and ask price and quantity for a perpetual futures symbol.
func (e *Exchange) GetSwapSymbolOrderbookTicker(ctx context.Context, symbol currency.Pair) (*SwapBookTicker, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SwapBookTickerResponse
	if err := e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v2/quote/bookTicker", params, &resp); err != nil {
		return nil, err
	}
	return &resp.BookTicker, nil
}

// GetSwapMarkPriceKlineData retrieves mark price candlestick data for a perpetual futures symbol.
func (e *Exchange) GetSwapMarkPriceKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, limit uint64) ([]*SwapMarkPriceCandle, error) {
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
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SwapMarkPriceCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v1/market/markPriceKlines", params, &resp)
}

// GetSwapSymbolPriceTicker retrieves the latest transaction price
func (e *Exchange) GetSwapSymbolPriceTicker(ctx context.Context, symbol currency.Pair) ([]*SwapPriceTicker, error) {
	if !symbol.IsEmpty() {
		params := url.Values{}
		params.Set("symbol", symbol.String())
		var resp *SwapPriceTicker
		if err := e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v1/ticker/price", params, &resp); err != nil {
			return nil, err
		}
		return []*SwapPriceTicker{resp}, nil
	}
	var resp []*SwapPriceTicker
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v1/ticker/price", nil, &resp)
}

// GetSwapTradingRules retrieves the trading rules for a perpetual futures symbol
func (e *Exchange) GetSwapTradingRules(ctx context.Context, symbol currency.Pair) (*SwapTradingRules, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SwapTradingRules
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "swap/v1/tradingRules", params, &resp)
}

// ---------------------------- Coin Marginid Futures ----------------------

// ------------------------- Trade Endpoint ---------------------------

// params validates and converts a PlaceSwapOrderRequest into url.Values shared by the order, test order and cancel-replace endpoints.
func (r *PlaceSwapOrderRequest) params() (url.Values, error) {
	if r == nil {
		return nil, common.ErrNilPointer
	}
	if r.Symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if r.Type == "" {
		return nil, order.ErrTypeIsInvalid
	}
	if r.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	params := url.Values{}
	params.Set("symbol", r.Symbol.String())
	params.Set("type", r.Type)
	params.Set("side", r.Side)
	if r.PositionSide != "" {
		params.Set("positionSide", r.PositionSide)
	}
	if r.ReduceOnly != "" {
		params.Set("reduceOnly", r.ReduceOnly)
	}
	setFloat(params, "price", r.Price)
	setFloat(params, "quantity", r.Quantity)
	setFloat(params, "quoteOrderQty", r.QuoteOrderQty)
	setFloat(params, "stopPrice", r.StopPrice)
	setFloat(params, "priceRate", r.PriceRate)
	setFloat(params, "activationPrice", r.ActivationPrice)
	setFloat(params, "stopLoss", r.StopLoss)
	setFloat(params, "takeProfit", r.TakeProfit)
	if r.WorkingType != "" {
		params.Set("workingType", r.WorkingType)
	}
	if r.ClientOrderID != "" {
		params.Set("clientOrderId", r.ClientOrderID)
	}
	if r.TimeInForce != "" {
		params.Set("timeInForce", r.TimeInForce)
	}
	if r.ClosePosition != "" {
		params.Set("closePosition", r.ClosePosition)
	}
	if r.StopGuaranteed != "" {
		params.Set("stopGuaranteed", r.StopGuaranteed)
	}
	if r.PositionID != 0 {
		params.Set("positionId", strconv.FormatInt(r.PositionID, 10))
	}
	if r.RecvWindow != 0 {
		params.Set("recvWindow", strconv.FormatInt(r.RecvWindow, 10))
	}
	return params, nil
}

// TestSwapOrder validates the parameters of a perpetual futures order without actually submitting it.
func (e *Exchange) TestSwapOrder(ctx context.Context, arg *PlaceSwapOrderRequest) (*SwapOrderAck, error) {
	params, err := arg.params()
	if err != nil {
		return nil, err
	}
	var resp *SwapOrderAckResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/order/test", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Order, nil
}

// PlaceSwapOrder submits a perpetual futures order.
func (e *Exchange) PlaceSwapOrder(ctx context.Context, arg *PlaceSwapOrderRequest) (*SwapOrderAck, error) {
	params, err := arg.params()
	if err != nil {
		return nil, err
	}
	var resp *SwapOrderAckResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/order", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Order, nil
}

// ModifySwapOrder amends the quantity of an existing perpetual futures order.
func (e *Exchange) ModifySwapOrder(ctx context.Context, symbol currency.Pair, orderID, clientOrderID string, quantity float64) (*SwapModifyOrderResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if orderID == "" && clientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	if quantity <= 0 {
		return nil, order.ErrAmountIsInvalid
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	setFloat(params, "quantity", quantity)
	var resp *SwapModifyOrderResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/amend", params, &resp)
}

// PlaceMultipleSwapOrders submits up to five perpetual futures orders in a single request.
func (e *Exchange) PlaceMultipleSwapOrders(ctx context.Context, args []PlaceSwapOrderRequest) ([]*SwapOrderAck, error) {
	if len(args) == 0 {
		return nil, errBatchOrdersRequired
	}
	orders := make([]map[string]string, 0, len(args))
	for i := range args {
		params, err := args[i].params()
		if err != nil {
			return nil, err
		}
		orderParams := make(map[string]string, len(params))
		for key := range params {
			orderParams[key] = params.Get(key)
		}
		orders = append(orders, orderParams)
	}
	payload, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("batchOrders", string(payload))
	var resp *SwapOrderAckListResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/batchOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// CloseAllSwapPositions closes all open positions, optionally limited to a single symbol, at market price.
func (e *Exchange) CloseAllSwapPositions(ctx context.Context, symbol currency.Pair) (*SwapCloseAllPositionsResponse, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *SwapCloseAllPositionsResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/closeAllPositions", params, &resp)
}

// CancelSwapOrder cancels an active perpetual futures order.
func (e *Exchange) CancelSwapOrder(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*SwapOrder, error) {
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
	var resp *SwapOrderDetailResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, "swap/v2/trade/order", params, &resp); err != nil {
		return nil, err
	}
	return resp.Order, nil
}

// CancelMultipleSwapOrders cancels up to ten perpetual futures orders by order ID or client order ID.
func (e *Exchange) CancelMultipleSwapOrders(ctx context.Context, symbol currency.Pair, orderIDList []int64, clientOrderIDList []string) (*SwapCancelOrdersResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if len(orderIDList) == 0 && len(clientOrderIDList) == 0 {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if len(orderIDList) > 0 {
		ids := make([]string, len(orderIDList))
		for i := range orderIDList {
			ids[i] = strconv.FormatInt(orderIDList[i], 10)
		}
		params.Set("orderIdList", "["+strings.Join(ids, ",")+"]")
	}
	if len(clientOrderIDList) > 0 {
		payload, err := json.Marshal(clientOrderIDList)
		if err != nil {
			return nil, err
		}
		params.Set("clientOrderIdList", string(payload))
	}
	var resp *SwapCancelOrdersResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, "swap/v2/trade/batchOrders", params, &resp)
}

// CancelAllSwapOpenOrders cancels all open perpetual futures orders, optionally filtered by symbol and order type.
func (e *Exchange) CancelAllSwapOpenOrders(ctx context.Context, symbol currency.Pair, orderType string) (*SwapCancelOrdersResponse, error) {
	if orderType == "" {
		return nil, order.ErrTypeIsInvalid
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	params.Set("type", orderType)
	var resp *SwapCancelOrdersResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, "swap/v2/trade/allOpenOrders", params, &resp)
}

// GetSwapOpenOrders retrieves all current open perpetual futures orders, optionally filtered by symbol and order type.
func (e *Exchange) GetSwapOpenOrders(ctx context.Context, symbol currency.Pair, orderType string) ([]*SwapOrder, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	var resp *SwapOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/openOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetSwapPendingOrderStatus retrieves the status of a pending perpetual futures order.
func (e *Exchange) GetSwapPendingOrderStatus(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*SwapOrder, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp *SwapOrderDetailResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/openOrder", params, &resp); err != nil {
		return nil, err
	}
	return resp.Order, nil
}

// GetSwapOrderDetails retrieves the full details of a perpetual futures order.
func (e *Exchange) GetSwapOrderDetails(ctx context.Context, symbol currency.Pair, orderID int64, clientOrderID string) (*SwapOrder, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	var resp *SwapOrderDetailResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/order", params, &resp); err != nil {
		return nil, err
	}
	return resp.Order, nil
}

// GetSwapMarginType retrieves the margin mode of a perpetual futures symbol.
func (e *Exchange) GetSwapMarginType(ctx context.Context, symbol currency.Pair) (*SwapMarginTypeResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SwapMarginTypeResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/marginType", params, &resp)
}

// ChangeSwapMarginType changes the margin mode of a perpetual futures symbol.
func (e *Exchange) ChangeSwapMarginType(ctx context.Context, symbol currency.Pair, marginType string) error {
	if symbol.IsEmpty() {
		return currency.ErrCurrencyPairEmpty
	}
	if marginType == "" {
		return fmt.Errorf("%w; empty margin type", margin.ErrInvalidMarginType)
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("marginType", marginType)
	return e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/marginType", params, &struct{}{})
}

// GetSwapLeverage retrieves the leverage and available positions for a perpetual futures symbol.
func (e *Exchange) GetSwapLeverage(ctx context.Context, symbol currency.Pair) (*SwapLeverage, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *SwapLeverage
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/leverage", params, &resp)
}

// SetSwapLeverage sets the leverage for the long or short side of a perpetual futures symbol.
func (e *Exchange) SetSwapLeverage(ctx context.Context, symbol currency.Pair, side string, leverage int64) (*SwapLeverage, error) {
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
	var resp *SwapLeverage
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/leverage", params, &resp)
}

// GetSwapForceOrders retrieves the user's liquidation (force) order history.
func (e *Exchange) GetSwapForceOrders(ctx context.Context, symbol currency.Pair, currencyCode, autoCloseType string, startTime, endTime time.Time, limit uint64) ([]*SwapOrder, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if currencyCode != "" {
		params.Set("currency", currencyCode)
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
	var resp *SwapOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/forceOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetSwapOrderHistory retrieves the perpetual futures order history.
func (e *Exchange) GetSwapOrderHistory(ctx context.Context, symbol currency.Pair, currencyCode string, orderID int64, startTime, endTime time.Time, limit uint64) ([]*SwapOrder, error) {
	if !startTime.IsZero() && !endTime.IsZero() {
		if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
			return nil, err
		}
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if currencyCode != "" {
		params.Set("currency", currencyCode)
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
	var resp *SwapOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/allOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// ModifyIsolatedPositionMargin increases or decreases the isolated margin of a position.
func (e *Exchange) ModifyIsolatedPositionMargin(ctx context.Context, symbol currency.Pair, amount float64, adjustmentType, positionSide string, positionID int64) (*SwapPositionMarginResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if amount <= 0 {
		return nil, errMarginAmountRequired
	}
	if adjustmentType == "" {
		return nil, errAdjustmentTypeRequired
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	setFloat(params, "amount", amount)
	params.Set("type", adjustmentType)
	if positionSide != "" {
		params.Set("positionSide", positionSide)
	}
	if positionID != 0 {
		params.Set("positionId", strconv.FormatInt(positionID, 10))
	}
	var resp *SwapPositionMarginResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/positionMargin", params, &resp)
}

// GetSwapHistoricalTransactionOrders retrieves historical transaction orders for a perpetual futures symbol.
func (e *Exchange) GetSwapHistoricalTransactionOrders(ctx context.Context, symbol currency.Pair, currencyCode, tradingUnit string, orderID int64, startTime, endTime time.Time) ([]*SwapFillOrder, error) {
	if tradingUnit == "" {
		return nil, errTradingUnitRequired
	}
	if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
		return nil, err
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if currencyCode != "" {
		params.Set("currency", currencyCode)
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	params.Set("tradingUnit", tradingUnit)
	params.Set("startTs", strconv.FormatInt(startTime.UnixMilli(), 10))
	params.Set("endTs", strconv.FormatInt(endTime.UnixMilli(), 10))
	var resp *SwapFillOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/allFillOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.FillOrders, nil
}

// SetSwapPositionMode sets the position mode to either hedge (dual side) or one-way.
func (e *Exchange) SetSwapPositionMode(ctx context.Context, dualSidePosition string) (*SwapPositionModeResponse, error) {
	if dualSidePosition == "" {
		return nil, errPositionModeRequired
	}
	params := url.Values{}
	params.Set("dualSidePosition", dualSidePosition)
	var resp *SwapPositionModeResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/positionSide/dual", params, &resp)
}

// GetSwapPositionMode retrieves the current position mode.
func (e *Exchange) GetSwapPositionMode(ctx context.Context) (*SwapPositionModeResponse, error) {
	var resp *SwapPositionModeResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/positionSide/dual", nil, &resp)
} // swap/v1/positionSide/dual

// CancelReplaceSwapOrder cancels an existing perpetual futures order and submits a new one.
func (e *Exchange) CancelReplaceSwapOrder(ctx context.Context, cancelReplaceMode, cancelClientOrderID string, cancelOrderID int64, cancelRestrictions string, arg *PlaceSwapOrderRequest) (*SwapCancelReplaceResponse, error) {
	if cancelReplaceMode == "" {
		return nil, errCancelReplaceModeEmpty
	}
	if cancelOrderID == 0 && cancelClientOrderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params, err := arg.params()
	if err != nil {
		return nil, err
	}
	params.Set("cancelReplaceMode", cancelReplaceMode)
	if cancelClientOrderID != "" {
		params.Set("cancelClientOrderId", cancelClientOrderID)
	}
	if cancelOrderID != 0 {
		params.Set("cancelOrderId", strconv.FormatInt(cancelOrderID, 10))
	}
	if cancelRestrictions != "" {
		params.Set("cancelRestrictions", cancelRestrictions)
	}
	var resp *SwapCancelReplaceResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/cancelReplace", params, &resp)
}

// BatchCancelReplaceSwapOrders cancels a batch of perpetual futures orders and places new ones.
func (e *Exchange) BatchCancelReplaceSwapOrders(ctx context.Context, batchOrders string) ([]SwapCancelReplaceResponse, error) {
	if batchOrders == "" {
		return nil, errBatchOrdersRequired
	}
	params := url.Values{}
	params.Set("batchOrders", batchOrders)
	var resp []SwapCancelReplaceResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/batchCancelReplace", params, &resp)
}

// SetSwapCancelAllAfter sets or clears a countdown that cancels all open orders after a timeout.
func (e *Exchange) SetSwapCancelAllAfter(ctx context.Context, requestType string, timeout int64) (*SwapCancelAllAfterResponse, error) {
	if requestType == "" {
		return nil, errCountdownTypeRequired
	}
	if timeout < 10 || timeout > 120 {
		return nil, errCountdownTimeoutInvalid
	}
	params := url.Values{}
	params.Set("type", requestType)
	params.Set("timeOut", strconv.FormatInt(timeout, 10))
	var resp *SwapCancelAllAfterResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/cancelAllAfter", params, &resp)
}

// CloseSwapPositionByID closes a position at market price using its position ID.
func (e *Exchange) CloseSwapPositionByID(ctx context.Context, positionID string) (*SwapClosePositionResponse, error) {
	if positionID == "" {
		return nil, errPositionIDRequired
	}
	params := url.Values{}
	params.Set("positionId", positionID)
	var resp *SwapClosePositionResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/closePosition", params, &resp)
}

// GetAllSwapOrders retrieves the full perpetual futures order history including conditional orders.
func (e *Exchange) GetAllSwapOrders(ctx context.Context, symbol currency.Pair, orderID int64, startTime, endTime time.Time, limit uint64) ([]*SwapOrder, error) {
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
	var resp *SwapOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/trade/fullOrder", params, &resp); err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// GetSwapMaintMarginRatio retrieves the position and maintenance margin ratio tiers for a perpetual futures symbol.
func (e *Exchange) GetSwapMaintMarginRatio(ctx context.Context, symbol currency.Pair) ([]*SwapMaintMarginRatio, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp []*SwapMaintMarginRatio
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/maintMarginRatio", params, &resp)
}

// GetSwapTransactionDetails retrieves historical transaction details for a perpetual futures symbol.
func (e *Exchange) GetSwapTransactionDetails(ctx context.Context, symbol currency.Pair, currencyCode string, orderID, lastFillID int64, pageIndex, pageSize uint64, startTime, endTime time.Time) (*SwapFillHistoryResponse, error) {
	if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
		return nil, err
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	if currencyCode != "" {
		params.Set("currency", currencyCode)
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if lastFillID != 0 {
		params.Set("lastFillId", strconv.FormatInt(lastFillID, 10))
	}
	if pageIndex != 0 {
		params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	}
	params.Set("startTs", strconv.FormatInt(startTime.UnixMilli(), 10))
	params.Set("endTs", strconv.FormatInt(endTime.UnixMilli(), 10))
	var resp *SwapFillHistoryResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v2/trade/fillHistory", params, &resp)
}

// GetSwapPositionHistory retrieves the closed position history for a perpetual futures symbol.
func (e *Exchange) GetSwapPositionHistory(ctx context.Context, symbol currency.Pair, currencyCode string, positionID int64, startTime, endTime time.Time, pageIndex, pageSize uint64) ([]*SwapPositionHistory, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if currencyCode != "" {
		params.Set("currency", currencyCode)
	}
	if positionID != 0 {
		params.Set("positionId", strconv.FormatInt(positionID, 10))
	}
	params.Set("startTs", strconv.FormatInt(startTime.UnixMilli(), 10))
	params.Set("endTs", strconv.FormatInt(endTime.UnixMilli(), 10))
	if pageIndex != 0 {
		params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	}
	var resp []*SwapPositionHistory
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/trade/positionHistory", params, &resp)
}

// GetSwapIsolatedMarginChangeHistory retrieves the isolated margin change history for a position.
func (e *Exchange) GetSwapIsolatedMarginChangeHistory(ctx context.Context, symbol currency.Pair, positionID string, startTime, endTime time.Time, pageIndex, pageSize uint64) (*SwapMarginChangeHistoryResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if positionID == "" {
		return nil, errPositionIDRequired
	}
	if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("positionId", positionID)
	params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	var resp *SwapMarginChangeHistoryResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/positionMargin/history", params, &resp)
}

// ApplyVST applies for VST funds on the BingX demo trading account.
func (e *Exchange) ApplyVST(ctx context.Context, adjustType string, amount int64) (*SwapVSTResponse, error) {
	params := url.Values{}
	if adjustType != "" {
		params.Set("adjustType", adjustType)
	}
	if amount != 0 {
		params.Set("amount", strconv.FormatInt(amount, 10))
	}
	var resp *SwapVSTResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v2/trade/getVst", params, &resp)
}

// PlaceTWAPOrder submits a time-weighted average price (TWAP) order.
func (e *Exchange) PlaceTWAPOrder(ctx context.Context, arg *PlaceTWAPOrderRequest) (*SwapTWAPOrderResponse, error) {
	if err := common.NilGuard(arg); err != nil {
		return nil, err
	}
	if arg.Symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if arg.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if arg.PositionSide == "" {
		return nil, order.ErrSideIsInvalid
	}
	if arg.PriceType == "" {
		return nil, fmt.Errorf("%w price type is required", order.ErrUnknownPriceType)
	}
	params := url.Values{}
	params.Set("symbol", arg.Symbol.String())
	params.Set("side", arg.Side)
	params.Set("positionSide", arg.PositionSide)
	params.Set("priceType", arg.PriceType)
	params.Set("priceVariance", arg.PriceVariance)
	params.Set("triggerPrice", arg.TriggerPrice)
	params.Set("interval", strconv.FormatInt(arg.Interval, 10))
	params.Set("amountPerOrder", arg.AmountPerOrder)
	params.Set("totalAmount", arg.TotalAmount)
	if arg.RecvWindow != 0 {
		params.Set("recvWindow", strconv.FormatInt(arg.RecvWindow, 10))
	}
	var resp *SwapTWAPOrderResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/twap/order", params, &resp)
}

// GetTWAPOpenOrders retrieves the open TWAP orders, optionally filtered by symbol.
func (e *Exchange) GetTWAPOpenOrders(ctx context.Context, symbol currency.Pair) ([]*SwapTWAPOrder, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp *SwapTWAPOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/twap/openOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}

// GetTWAPHistoricalOrders retrieves the historical TWAP orders.
func (e *Exchange) GetTWAPHistoricalOrders(ctx context.Context, symbol currency.Pair, startTime, endTime time.Time, pageIndex, pageSize uint64) ([]*SwapTWAPOrder, error) {
	if err := common.StartEndTimeCheck(startTime, endTime); err != nil {
		return nil, err
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	params.Set("pageIndex", strconv.FormatUint(pageIndex, 10))
	params.Set("pageSize", strconv.FormatUint(pageSize, 10))
	params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	var resp *SwapTWAPOrdersResponse
	if err := e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/twap/historyOrders", params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}

// GetTWAPOrderDetails retrieves the details of a TWAP order.
func (e *Exchange) GetTWAPOrderDetails(ctx context.Context, mainOrderID string) (*SwapTWAPOrder, error) {
	if mainOrderID == "" {
		return nil, fmt.Errorf("%w twap main order id is required", order.ErrOrderIDNotSet)
	}
	params := url.Values{}
	params.Set("mainOrderId", mainOrderID)
	var resp *SwapTWAPOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/twap/orderDetail", params, &resp)
}

// CancelTWAPOrder cancels a TWAP order.
func (e *Exchange) CancelTWAPOrder(ctx context.Context, mainOrderID string) (*SwapTWAPOrder, error) {
	if mainOrderID == "" {
		return nil, fmt.Errorf("%w twap main order id is required", order.ErrOrderIDNotSet)
	}
	params := url.Values{}
	params.Set("mainOrderId", mainOrderID)
	var resp *SwapTWAPOrder
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/twap/cancelOrder", params, &resp)
}

// SwitchMultiAssetsMode switches between single-asset and multi-assets margin mode.
func (e *Exchange) SwitchMultiAssetsMode(ctx context.Context, assetMode string) (*SwapAssetModeResponse, error) {
	if assetMode == "" {
		return nil, errAssetModeRequired
	}
	params := url.Values{}
	params.Set("assetMode", assetMode)
	var resp *SwapAssetModeResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/assetMode", params, &resp)
}

// GetMultiAssetsMode retrieves the current multi-assets margin mode.
func (e *Exchange) GetMultiAssetsMode(ctx context.Context) (*SwapAssetModeResponse, error) {
	var resp *SwapAssetModeResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/trade/assetMode", nil, &resp)
}

// GetMultiAssetsRules retrieves the margin asset rules under multi-assets mode.
func (e *Exchange) GetMultiAssetsRules(ctx context.Context) ([]*SwapMultiAssetsRule, error) {
	var resp []*SwapMultiAssetsRule
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/trade/multiAssetsRules", nil, &resp)
}

// GetMultiAssetsMargin retrieves the margin asset balances under multi-assets mode.
func (e *Exchange) GetMultiAssetsMargin(ctx context.Context) ([]*SwapMultiAssetsMargin, error) {
	var resp []*SwapMultiAssetsMargin
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "swap/v1/user/marginAssets", nil, &resp)
}

// OneClickReversePosition reverses an open position immediately or via a planned trigger.
func (e *Exchange) OneClickReversePosition(ctx context.Context, reverseType string, symbol currency.Pair, triggerPrice, workingType string) (*SwapReversePositionResponse, error) {
	if reverseType == "" {
		return nil, errReverseTypeRequired
	}
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	params := url.Values{}
	params.Set("type", reverseType)
	params.Set("symbol", symbol.String())
	if triggerPrice != "" {
		params.Set("triggerPrice", triggerPrice)
	}
	if workingType != "" {
		params.Set("workingType", workingType)
	}
	var resp *SwapReversePositionResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/reverse", params, &resp)
}

// SetHedgeModeAutoAddMargin toggles automatic margin addition for a hedge-mode position.
func (e *Exchange) SetHedgeModeAutoAddMargin(ctx context.Context, symbol currency.Pair, positionID int64, functionSwitch string, amount float64) (*SwapAutoAddMarginResponse, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if positionID == 0 {
		return nil, errPositionIDRequired
	}
	if functionSwitch == "" {
		return nil, errFunctionSwitchRequired
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("positionId", strconv.FormatInt(positionID, 10))
	params.Set("functionSwitch", functionSwitch)
	setFloat(params, "amount", amount)
	var resp *SwapAutoAddMarginResponse
	return resp, e.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "swap/v1/trade/autoAddMargin", params, &resp)
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (e *Exchange) SendHTTPRequest(ctx context.Context, ep exchange.URL, path string, params url.Values, result any) error {
	endpoint, err := e.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	if params == nil {
		params = url.Values{}
	}
	resp := &ResponseWrapper{
		Data: result,
	}
	if params.Get("timestamp") == "" {
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}
	item := &request.Item{
		Method:                 http.MethodGet,
		Path:                   endpoint + common.EncodeURLValues(path, params),
		Result:                 resp,
		Verbose:                e.Verbose,
		HTTPDebugging:          e.HTTPDebugging,
		HTTPRecording:          e.HTTPRecording,
		HTTPMockDataSliceLimit: e.HTTPMockDataSliceLimit,
	}
	if err := e.SendPayload(ctx, request.Unset, func() (*request.Item, error) {
		return item, nil
	}, request.UnauthenticatedRequest); err != nil {
		return err
	}
	if resp.Code != 0 || resp.Msg != "" {
		return fmt.Errorf("failed with error code: %d message: %s", resp.Code, resp.Msg)
	}
	return nil
}

// SendAuthenticatedHTTPRequest signs and sends an authenticated HTTP request
func (e *Exchange) SendAuthenticatedHTTPRequest(ctx context.Context, ep exchange.URL, method, path string, params url.Values, result any) error {
	creds, err := e.GetCredentials(ctx)
	if err != nil {
		return err
	}
	endpoint, err := e.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	if params == nil {
		params = url.Values{}
	}
	resp := &ResponseWrapper{
		Data: result,
	}
	newRequest := func() (*request.Item, error) {
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
		hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256, []byte(signaturePayload(params)), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}
		return &request.Item{
			Method:                 method,
			Path:                   endpoint + common.EncodeURLValues(path, params) + "&signature=" + hex.EncodeToString(hmacSigned),
			Headers:                map[string]string{"X-BX-APIKEY": creds.Key},
			Result:                 resp,
			Verbose:                e.Verbose,
			HTTPDebugging:          e.HTTPDebugging,
			HTTPRecording:          e.HTTPRecording,
			HTTPMockDataSliceLimit: e.HTTPMockDataSliceLimit,
		}, nil
	}
	if err := e.SendPayload(ctx, request.Unset, newRequest, request.AuthenticatedRequest); err != nil {
		return err
	}
	if resp.Code != 0 || resp.Msg != "" {
		return fmt.Errorf("failed with error code: %d message: %s", resp.Code, resp.Msg)
	}
	return nil
}

func signaturePayload(params url.Values) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for i, key := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(params.Get(key))
	}
	return sb.String()
}
