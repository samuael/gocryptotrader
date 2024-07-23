package apexpro

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Apexpro is the overarching type across this package
type Apexpro struct {
	exchange.Base
}

const (
	apexproAPIURL     = "https://pro.apex.exchange/api/"
	apexproTestAPIURL = "https://testnet.pro.apex.exchange/api/"

	// Public endpoints

	// Authenticated endpoints
)

var (
	errEthereumAddressMissing   = errors.New("ethereum address is missing")
	errL2KeyMissing             = errors.New("l2 Key is required")
	errOrderbookLevelIsRequired = errors.New("orderbook level is required")
)

// Start implementing public and private exchange API funcs below

// GetSystemTimeV3 retrieves V3 system time.
func (ap *Apexpro) GetSystemTimeV3(ctx context.Context) (time.Time, error) {
	return ap.getSystemTime(ctx, "v3/time")
}

// GetSystemTimeV2 retrieves V2 system time.
func (ap *Apexpro) GetSystemTimeV2(ctx context.Context) (time.Time, error) {
	return ap.getSystemTime(ctx, "v2/time")
}

// GetSystemTimeV2 retrieves V2 system time.
func (ap *Apexpro) GetSystemTimeV1(ctx context.Context) (time.Time, error) {
	return ap.getSystemTime(ctx, "v1/time")
}

func (ap *Apexpro) getSystemTime(ctx context.Context, path string) (time.Time, error) {
	resp := &struct {
		Time convert.ExchangeTime `json:"time"`
	}{}
	return resp.Time.Time(), ap.SendHTTPRequest(ctx, exchange.RestSpot, path, request.UnAuth, &resp)
}

// GetAllConfigDataV3 retrieves all symbols and asset configurations.
func (ap *Apexpro) GetAllConfigDataV3(ctx context.Context) (*AllSymbolsConfigs, error) {
	var resp *AllSymbolsConfigs
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v3/symbols", request.UnAuth, &resp, true)
}

// Apexpro retrieves all symbols and asset configurations from the V1 API.
func (ap *Apexpro) GetAllConfigDataV1(ctx context.Context) (*AllSymbolsV1Config, error) {
	var resp *AllSymbolsV1Config
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v1/symbols", request.UnAuth, &resp, true)
}

// GetMarketDepthV3 retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apexpro) GetMarketDepthV3(ctx context.Context, symbol string, limit int64) (*MarketDepthV3, error) {
	return ap.getMarketDepth(ctx, symbol, "v3/depth", limit)
}

// GetMarketDepthV2 retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apexpro) GetMarketDepthV2(ctx context.Context, symbol string, limit int64) (*MarketDepthV3, error) {
	return ap.getMarketDepth(ctx, symbol, "v2/depth", limit)
}

// GetMarketDepthV1 retrieve all active orderbook for one symbol, inclue all bids and asks.
func (ap *Apexpro) GetMarketDepthV1(ctx context.Context, symbol string, limit int64) (*MarketDepthV3, error) {
	return ap.getMarketDepth(ctx, symbol, "v1/depth", limit)
}

func (ap *Apexpro) getMarketDepth(ctx context.Context, symbol, path string, limit int64) (*MarketDepthV3, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)

	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp *MarketDepthV3
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

// GetNewestTradingDataV3 retrieve trading data.
func (ap *Apexpro) GetNewestTradingDataV3(ctx context.Context, symbol string, limit int64) ([]NewTradingData, error) {
	return ap.getNewestTradingData(ctx, symbol, "v3/trades", limit)
}

// GetNewestTradingDataV2 retrieve trading data.
func (ap *Apexpro) GetNewestTradingDataV2(ctx context.Context, symbol string, limit int64) ([]NewTradingData, error) {
	return ap.getNewestTradingData(ctx, symbol, "v2/trades", limit)
}

// GetNewestTradingDataV1 retrieve trading data.
func (ap *Apexpro) GetNewestTradingDataV1(ctx context.Context, symbol string, limit int64) ([]NewTradingData, error) {
	return ap.getNewestTradingData(ctx, symbol, "v1/trades", limit)
}

func (ap *Apexpro) getNewestTradingData(ctx context.Context, symbol, path string, limit int64) ([]NewTradingData, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []NewTradingData
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

func intervalToString(interval kline.Interval) (string, error) {
	intervalToStringMap := map[kline.Interval]string{
		kline.OneMin: "1", kline.FiveMin: "5", kline.FifteenMin: "15", kline.ThirtyMin: "30", kline.OneHour: "60", kline.TwoHour: "120", kline.FourHour: "240", kline.SixHour: "360", kline.SevenHour: "720", kline.OneDay: "D", kline.OneMonth: "M", kline.OneWeek: "W"}
	intervalString, okay := intervalToStringMap[interval]
	if !okay {
		return "", kline.ErrUnsupportedInterval
	}
	return intervalString, nil
}

// GetCandlestickChartDataV3 etrieves all candlestick chart data.
// Candlestick chart time indicators: Numbers represent minutes, D for Days, M for Month and W for Week — 1 5 15 30 60 120 240 360 720 "D" "M" "W"
func (ap *Apexpro) GetCandlestickChartDataV3(ctx context.Context, symbol string, interval kline.Interval, startTime, endTime time.Time, limit int64) (map[string][]CandlestickData, error) {
	return ap.getCandlestickChartData(ctx, symbol, "v3/klines", interval, startTime, endTime, limit)
}

// GetCandlestickChartDataV2 retrieves v2 all candlestick chart data.
func (ap *Apexpro) GetCandlestickChartDataV2(ctx context.Context, symbol string, interval kline.Interval, startTime, endTime time.Time, limit int64) (map[string][]CandlestickData, error) {
	return ap.getCandlestickChartData(ctx, symbol, "v2/klines", interval, startTime, endTime, limit)
}

// GetCandlestickChartDataV2 retrieves v2 all candlestick chart data.
func (ap *Apexpro) GetCandlestickChartDataV1(ctx context.Context, symbol string, interval kline.Interval, startTime, endTime time.Time, limit int64) (map[string][]CandlestickData, error) {
	return ap.getCandlestickChartData(ctx, symbol, "v1/klines", interval, startTime, endTime, limit)
}

func (ap *Apexpro) getCandlestickChartData(ctx context.Context, symbol, path string, interval kline.Interval, startTime, endTime time.Time, limit int64) (map[string][]CandlestickData, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	if interval != kline.Interval(0) {
		intervalString, err := intervalToString(interval)
		if err != nil {
			return nil, err
		}
		params.Set("interval", intervalString)
	}
	if !startTime.IsZero() {
		params.Set("start", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("end", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp map[string][]CandlestickData
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

// GetTickerDataV3 get the latest data on symbol tickers.
func (ap *Apexpro) GetTickerDataV3(ctx context.Context, symbol string) ([]TickerData, error) {
	return ap.getTickerData(ctx, symbol, "v3/ticker")
}

// GetTickerDataV2 get the latest data on symbol tickers.
func (ap *Apexpro) GetTickerDataV2(ctx context.Context, symbol string) ([]TickerData, error) {
	return ap.getTickerData(ctx, symbol, "v2/ticker")
}

func (ap *Apexpro) getTickerData(ctx context.Context, symbol, path string) ([]TickerData, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	var resp []TickerData
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

// GetFundingHistoryRateV3 retrieves a funding history rate.
func (ap *Apexpro) GetFundingHistoryRateV3(ctx context.Context, symbol string, beginTime, endTime time.Time, page, limit int64) (*FundingRateHistory, error) {
	return ap.getFundingHistoryRate(ctx, symbol, "v3/history-funding", beginTime, endTime, page, limit)
}

// GetFundingHistoryRateV2 retrieves a funding history rate.
func (ap *Apexpro) GetFundingHistoryRateV2(ctx context.Context, symbol string, beginTime, endTime time.Time, page, limit int64) (*FundingRateHistory, error) {
	return ap.getFundingHistoryRate(ctx, symbol, "v2/history-funding", beginTime, endTime, page, limit)
}

// GetFundingHistoryRateV1 retrieves a funding history rate.
func (ap *Apexpro) GetFundingHistoryRateV1(ctx context.Context, symbol string, beginTime, endTime time.Time, page, limit int64) (*FundingRateHistory, error) {
	return ap.getFundingHistoryRate(ctx, symbol, "v2/history-funding", beginTime, endTime, page, limit)
}

func (ap *Apexpro) getFundingHistoryRate(ctx context.Context, symbol, path string, beginTime, endTime time.Time, page, limit int64) (*FundingRateHistory, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	if !beginTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(beginTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	var resp *FundingRateHistory
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

// GetAllConfigDataV2 retrieves USDC and USDT config
func (ap *Apexpro) GetAllConfigDataV2(ctx context.Context) (*V2ConfigData, error) {
	var resp *V2ConfigData
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v2/symbols", request.UnAuth, &resp)
}

// GetCheckIfUserExistsV2 checks existence of a persion using ther ethereum Address
func (ap *Apexpro) GetCheckIfUserExistsV2(ctx context.Context, ethAddress string) (bool, error) {
	return ap.getCheckIfUserExists(ctx, ethAddress, "v2/check-user-exist")
}

// GetCheckIfUserExistsV1 checks existence of a persion using ther ethereum Address
func (ap *Apexpro) GetCheckIfUserExistsV1(ctx context.Context, ethAddress string) (bool, error) {
	return ap.getCheckIfUserExists(ctx, ethAddress, "v1/check-user-exist")
}

func (ap *Apexpro) getCheckIfUserExists(ctx context.Context, ethAddress, path string) (bool, error) {
	if ethAddress == "" {
		return false, errEthereumAddressMissing
	}
	params := url.Values{}
	params.Set("ethAddress", ethAddress)
	var resp bool
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(path, params), request.UnAuth, &resp)
}

// ----------------------------------------------------------------     Authenticated Endpoints ----------------------------------------------------------------

// GenerateNonce generate and obtain nonce before registration. The nonce is used to assemble the signature field upon registration.
func (ap *Apexpro) GenerateNonce(ctx context.Context, l2Key, ethAddress, chainID string) (*NonceResponse, error) {
	if l2Key == "" {
		return nil, errL2KeyMissing
	}
	params := url.Values{}
	params.Set("l2Key", l2Key)
	var resp *NonceResponse
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues("v3/generate-nonce", params), request.UnAuth, &resp)
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
