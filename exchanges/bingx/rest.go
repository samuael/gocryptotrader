package bingx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
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

// GetSpotOrderbookDepth get the order book depth for a specified symbol. It provides the most recent buy and sell orders for the given trading pair.
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

func (e *Exchange) GetSpotKlineData(ctx context.Context, symbol currency.Pair, interval kline.Interval, startTime, endTime time.Time, timezone, limit uint64) ([]*SpotCandle, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	
	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatUint(limit, 10))
	}
	var resp []*SpotCandle
	return resp, e.SendHTTPRequest(ctx, exchange.RestSpot, "spot/v2/market/kline", params, &resp)
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
	if err = e.SendPayload(ctx, request.Unset, func() (*request.Item, error) {
		return item, nil
	}, request.UnauthenticatedRequest); err != nil {
		return err
	}
	if resp.Code != 0 || resp.Msg != "" {
		return fmt.Errorf("failed with error code: %d message: %s", resp.Code, resp.Msg)
	}
	return nil
}
