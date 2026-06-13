package bingx

import (
	"context"
	"encoding/hex"
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
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Exchange implements exchange.IBotExchange and contains additional specific API methods for interacting with Bingx
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

// GetCoinMContracts retrieves coin-margined (Coin-M) futures contracts and their trading rules.
func (e *Exchange) GetCoinMContracts(ctx context.Context, symbol currency.Pair) ([]*CoinMContract, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMContract
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/contracts", params, &resp)
}

// GetCoinMMarkPriceAndFundingRate retrieves the latest mark price, index price and funding rate for coin-margined futures symbols.
func (e *Exchange) GetCoinMMarkPriceAndFundingRate(ctx context.Context, symbol currency.Pair) ([]*CoinMMarkPriceFundingRate, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMMarkPriceFundingRate
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/premiumIndex", params, &resp)
}

// GetCoinMOpenInterest retrieves the open interest of coin-margined futures symbols.
func (e *Exchange) GetCoinMOpenInterest(ctx context.Context, symbol currency.Pair) ([]*CoinMOpenInterest, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMOpenInterest
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/openInterest", params, &resp)
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
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/klines", params, &resp)
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
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/depth", params, &resp)
}

// GetCoinM24HrTickerPriceChange retrieves the 24-hour rolling window price change statistics for coin-margined futures symbols.
func (e *Exchange) GetCoinM24HrTickerPriceChange(ctx context.Context, symbol currency.Pair) ([]*CoinMTicker24Hr, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		params.Set("symbol", symbol.String())
	}
	var resp []*CoinMTicker24Hr
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "cswap/v1/market/ticker", params, &resp)
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
