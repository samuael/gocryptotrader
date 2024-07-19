package apexpro

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Apexpro is the overarching type across this package
type Apexpro struct {
	exchange.Base
}

const (
	apexproAPIURL     = "https://pro.apex.exchange/api/"
	apexproAPIVersion = "v1/"
	apexproTestAPIURL = "https://testnet.pro.apex.exchange/api/"

	// Public endpoints

	// Authenticated endpoints
)

// Start implementing public and private exchange API funcs below

// GetSystemTimeV3 retrieves V3 system time.
func (ap *Apexpro) GetSystemTimeV3(ctx context.Context) (time.Time, error) {
	resp := &struct {
		Time convert.ExchangeTime `json:"time"`
	}{}
	return resp.Time.Time(), ap.SendHTTPRequest(ctx, exchange.RestSpot, "v3/time", request.UnAuth, &resp)
}

// GetAllConfigDataV3 retrieves all symbols and asset configurations.
func (ap *Apexpro) GetAllConfigDataV3(ctx context.Context) (*AllSymbolsConfigs, error) {
	var resp *AllSymbolsConfigs
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v3/symbols", request.UnAuth, &resp, true)
}

// GetMarketDepthV3 retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apexpro) GetMarketDepthV3(ctx context.Context, symbol string, limit int64) ([]MarketDepthV3, error) {
	var resp []MarketDepthV3
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v3/depth", request.UnAuth, &resp)
}

// GetNewestTradingDataV3 retrieve trading data.
func (ap *Apexpro) GetNewestTradingDataV3(ctx context.Context, symbol string, limit int64) ([]NewTradingData, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []NewTradingData
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues("v3/trades", params), request.UnAuth, &resp)
}

// SendHTTPRequest sends an unauthenticated request
func (ap *Apexpro) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, f request.EndpointLimit, result interface{}, useAsItIs ...bool) error {
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	var response interface{}
	if len(useAsItIs) > 0 && useAsItIs[0] {
		response = result
	} else {
		response = &struct {
			Data interface{} `json:"data"`
		}{
			Data: result,
		}
	}
	item := &request.Item{
		Method:        http.MethodGet,
		Path:          endpointPath + path,
		Result:        response,
		Verbose:       ap.Verbose,
		HTTPDebugging: ap.HTTPDebugging,
		HTTPRecording: ap.HTTPRecording}

	return ap.SendPayload(ctx, f, func() (*request.Item, error) {
		return item, nil
	}, request.UnauthenticatedRequest)
}
