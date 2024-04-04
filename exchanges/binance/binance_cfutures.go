package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/types"
)

const (

	// Auth

	cfuturesLimit              = "LIMIT"
	cfuturesMarket             = "MARKET"
	cfuturesStop               = "STOP"
	cfuturesTakeProfit         = "TAKE_PROFIT"
	cfuturesStopMarket         = "STOP_MARKET"
	cfuturesTakeProfitMarket   = "TAKE_PROFIT_MARKET"
	cfuturesTrailingStopMarket = "TRAILING_STOP_MARKET"
)

// FuturesExchangeInfo stores CoinMarginedFutures, data
func (b *Binance) FuturesExchangeInfo(ctx context.Context) (*CExchangeInfo, error) {
	var resp *CExchangeInfo
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/exchangeInfo?", cFuturesDefaultRate, &resp)
}

// GetFuturesOrderbook gets orderbook data for CoinMarginedFutures,
func (b *Binance) GetFuturesOrderbook(ctx context.Context, symbol currency.Pair, limit int64) (*OrderBook, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("symbol", symbolValue)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}

	rateBudget := cFuturesDefaultRate
	switch {
	case limit == 5, limit == 10, limit == 20, limit == 50:
		rateBudget = cFuturesOrderbook50Rate
	case limit >= 100 && limit < 500:
		rateBudget = cFuturesOrderbook100Rate
	case limit >= 500 && limit < 1000:
		rateBudget = cFuturesOrderbook500Rate
	case limit == 1000:
		rateBudget = cFuturesOrderbook1000Rate
	}
	var data OrderbookData
	err = b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/depth?"+params.Encode(), rateBudget, &data)
	if err != nil {
		return nil, err
	}

	resp := &OrderBook{
		Symbol: symbol.String(),
		Asks:   make([]OrderbookItem, len(data.Asks)),
		Bids:   make([]OrderbookItem, len(data.Bids)),
	}
	for x := range data.Asks {
		resp.Asks[x] = OrderbookItem{
			Price:    data.Asks[x][0].Float64(),
			Quantity: data.Asks[x][1].Float64(),
		}
	}

	for y := range data.Bids {
		resp.Bids[y] = OrderbookItem{
			Price:    data.Bids[y][0].Float64(),
			Quantity: data.Bids[y][1].Float64(),
		}
	}
	return resp, nil
}

// GetFuturesPublicTrades gets recent public trades for CoinMarginedFutures,
func (b *Binance) GetFuturesPublicTrades(ctx context.Context, symbol currency.Pair, limit int64) ([]FuturesPublicTradesData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []FuturesPublicTradesData
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/trades?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetFuturesHistoricalTrades gets historical public trades for CoinMarginedFutures,
func (b *Binance) GetFuturesHistoricalTrades(ctx context.Context, symbol currency.Pair, fromID string, limit int64) ([]UPublicTradesData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if fromID != "" {
		params.Set("fromID", fromID)
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []UPublicTradesData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/historicalTrades", params, cFuturesHistoricalTradesRate, &resp)
}

// GetPastPublicTrades gets past public trades for CoinMarginedFutures,
func (b *Binance) GetPastPublicTrades(ctx context.Context, symbol currency.Pair, limit, fromID int64) ([]FuturesPublicTradesData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if fromID != 0 {
		params.Set("fromID", strconv.FormatInt(fromID, 10))
	}
	var resp []FuturesPublicTradesData
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/trades?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetFuturesAggregatedTradesList gets aggregated trades list for CoinMarginedFutures,
func (b *Binance) GetFuturesAggregatedTradesList(ctx context.Context, symbol currency.Pair, fromID, limit int64, startTime, endTime time.Time) ([]AggregatedTrade, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if fromID != 0 {
		params.Set("fromID", strconv.FormatInt(fromID, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []AggregatedTrade
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/aggTrades?"+params.Encode(), cFuturesHistoricalTradesRate, &resp)
}

// GetIndexAndMarkPrice gets index and mark prices  for CoinMarginedFutures,
func (b *Binance) GetIndexAndMarkPrice(ctx context.Context, symbol, pair string) ([]IndexMarkPrice, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []IndexMarkPrice
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/premiumIndex?"+params.Encode(), cFuturesIndexMarkPriceRate, &resp)
}

// GetFundingRateInfo returns extra details about funding rates
func (b *Binance) GetFundingRateInfo(ctx context.Context) ([]FundingRateInfoResponse, error) {
	var resp []FundingRateInfoResponse
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/fundingInfo", uFuturesDefaultRate, &resp)
}

// GetFuturesKlineData gets futures kline data for CoinMarginedFutures,
func (b *Binance) GetFuturesKlineData(ctx context.Context, symbol currency.Pair, interval string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	if !common.StringDataCompare(validFuturesIntervals, interval) {
		return nil, errors.New("invalid interval parsed")
	}
	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("interval", interval)
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}

	var data [][12]types.Number
	rateBudget := getKlineRateBudget(limit)
	err := b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/klines?"+params.Encode(), rateBudget, &data)
	if err != nil {
		return nil, err
	}

	resp := make([]FuturesCandleStick, len(data))
	for x := range data {
		resp[x] = FuturesCandleStick{
			OpenTime:                time.UnixMilli(data[x][0].Int64()),
			Open:                    data[x][1].Float64(),
			High:                    data[x][2].Float64(),
			Low:                     data[x][3].Float64(),
			Close:                   data[x][4].Float64(),
			Volume:                  data[x][5].Float64(),
			CloseTime:               time.UnixMilli(data[x][6].Int64()),
			BaseAssetVolume:         data[x][7].Float64(),
			NumberOfTrades:          data[x][8].Int64(),
			TakerBuyVolume:          data[x][9].Float64(),
			TakerBuyBaseAssetVolume: data[x][10].Float64(),
		}
	}
	return resp, nil
}

// GetContinuousKlineData gets continuous kline data
func (b *Binance) GetContinuousKlineData(ctx context.Context, pair, contractType, interval string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	if pair == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	if !common.StringDataCompare(validContractType, contractType) {
		return nil, errors.New("invalid contractType")
	}
	if !common.StringDataCompare(validFuturesIntervals, interval) {
		return nil, errors.New("invalid interval parsed")
	}
	params := url.Values{}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	params.Set("pair", pair)
	params.Set("contractType", contractType)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("interval", interval)
	rateBudget := getKlineRateBudget(limit)
	var data [][12]types.Number
	err := b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/continuousKlines?"+params.Encode(), rateBudget, &data)
	if err != nil {
		return nil, err
	}
	resp := make([]FuturesCandleStick, len(data))
	for x := range data {
		resp[x] = FuturesCandleStick{
			OpenTime:                time.UnixMilli(data[x][0].Int64()),
			Open:                    data[x][1].Float64(),
			High:                    data[x][2].Float64(),
			Low:                     data[x][3].Float64(),
			Close:                   data[x][4].Float64(),
			Volume:                  data[x][5].Float64(),
			CloseTime:               time.UnixMilli(data[x][6].Int64()),
			BaseAssetVolume:         data[x][7].Float64(),
			NumberOfTrades:          data[x][8].Int64(),
			TakerBuyVolume:          data[x][9].Float64(),
			TakerBuyBaseAssetVolume: data[x][10].Float64(),
		}
	}
	return resp, nil
}

// GetIndexPriceKlines gets continuous kline data
func (b *Binance) GetIndexPriceKlines(ctx context.Context, pair, interval string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	if !common.StringDataCompare(validFuturesIntervals, interval) {
		return nil, errors.New("invalid interval parsed")
	}
	params := url.Values{}
	params.Set("pair", pair)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("interval", interval)
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	rateBudget := getKlineRateBudget(limit)
	var data [][12]types.Number
	err := b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/indexPriceKlines?"+params.Encode(), rateBudget, &data)
	if err != nil {
		return nil, err
	}
	resp := make([]FuturesCandleStick, len(data))
	for x := range data {
		resp[x] = FuturesCandleStick{
			OpenTime:  time.UnixMilli(data[x][0].Int64()),
			Open:      data[x][1].Float64(),
			High:      data[x][2].Float64(),
			Low:       data[x][3].Float64(),
			Close:     data[x][4].Float64(),
			CloseTime: time.UnixMilli(data[x][6].Int64()),
		}
	}
	return resp, nil
}

// GetMarkPriceKline gets mark price kline data
func (b *Binance) GetMarkPriceKline(ctx context.Context, symbol currency.Pair, interval string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	return b.getKline(ctx, symbol, interval, "/dapi/v1/markPriceKlines", limit, startTime, endTime)
}

// GetPremiumIndexKlineData premium index kline bars of a symbol.
// Klines are uniquely identified by their open time.
func (b *Binance) GetPremiumIndexKlineData(ctx context.Context, symbol currency.Pair, interval string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	return b.getKline(ctx, symbol, interval, "/dapi/v1/premiumIndexKlines", limit, startTime, endTime)
}

func (b *Binance) getKline(ctx context.Context, symbol currency.Pair, interval, path string, limit int64, startTime, endTime time.Time) ([]FuturesCandleStick, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	if !common.StringDataCompare(validFuturesIntervals, interval) {
		return nil, fmt.Errorf("%w, interval = %s", kline.ErrInvalidInterval, interval)
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("interval", interval)
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var data [][12]types.Number
	rateBudget := getKlineRateBudget(limit)
	err = b.SendHTTPRequest(ctx, exchange.RestCoinMargined, common.EncodeURLValues(path, params), rateBudget, &data)
	if err != nil {
		return nil, err
	}
	resp := make([]FuturesCandleStick, len(data))
	for x := range data {
		resp[x] = FuturesCandleStick{
			OpenTime:  time.UnixMilli(data[x][0].Int64()),
			Open:      data[x][1].Float64(),
			High:      data[x][2].Float64(),
			Low:       data[x][3].Float64(),
			Close:     data[x][4].Float64(),
			CloseTime: time.UnixMilli(data[x][6].Int64()),
		}
	}
	return resp, nil
}

func getKlineRateBudget(limit int64) request.EndpointLimit {
	rateBudget := cFuturesDefaultRate
	switch {
	case limit > 0 && limit < 100:
		rateBudget = cFuturesKline100Rate
	case limit >= 100 && limit < 500:
		rateBudget = cFuturesKline500Rate
	case limit >= 500 && limit < 1000:
		rateBudget = cFuturesKline1000Rate
	case limit >= 1000:
		rateBudget = cFuturesKlineMaxRate
	}
	return rateBudget
}

// GetFuturesSwapTickerChangeStats gets 24hr ticker change stats for CoinMarginedFutures,
func (b *Binance) GetFuturesSwapTickerChangeStats(ctx context.Context, symbol currency.Pair, pair string) ([]PriceChangeStats, error) {
	params := url.Values{}
	rateLimit := cFuturesTickerPriceHistoryRate
	if !symbol.IsEmpty() {
		rateLimit = cFuturesDefaultRate
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []PriceChangeStats
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/ticker/24hr?"+params.Encode(), rateLimit, &resp)
}

// FuturesGetFundingHistory gets funding history for CoinMarginedFutures,
func (b *Binance) FuturesGetFundingHistory(ctx context.Context, symbol currency.Pair, limit int64, startTime, endTime time.Time) ([]FundingRateHistory, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []FundingRateHistory
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/fundingRate?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetFuturesSymbolPriceTicker gets price ticker for symbol
func (b *Binance) GetFuturesSymbolPriceTicker(ctx context.Context, symbol currency.Pair, pair string) ([]SymbolPriceTicker, error) {
	params := url.Values{}
	rateLimit := cFuturesOrderbookTickerAllRate
	if !symbol.IsEmpty() {
		rateLimit = cFuturesDefaultRate
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []SymbolPriceTicker
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/ticker/price?"+params.Encode(), rateLimit, &resp)
}

// GetFuturesOrderbookTicker gets orderbook ticker for symbol
func (b *Binance) GetFuturesOrderbookTicker(ctx context.Context, symbol currency.Pair, pair string) ([]SymbolOrderBookTicker, error) {
	params := url.Values{}
	rateLimit := cFuturesOrderbookTickerAllRate
	if !symbol.IsEmpty() {
		rateLimit = cFuturesDefaultRate
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []SymbolOrderBookTicker
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/ticker/bookTicker?"+params.Encode(), rateLimit, &resp)
}

// GetCFuturesIndexPriceConstituents retrieved index price constituents detail
func (b *Binance) GetCFuturesIndexPriceConstituents(ctx context.Context, symbol currency.Pair) (*CFuturesIndexPriceConstituents, error) {
	if symbol.IsEmpty() {
		return nil, currency.ErrSymbolStringEmpty
	}
	params := url.Values{}
	params.Set("symbol", symbol.String())
	var resp *CFuturesIndexPriceConstituents
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, common.EncodeURLValues("/dapi/v1/constituents", params), cFuturesDefaultRate, &resp)
}

// OpenInterest gets open interest data for a symbol
func (b *Binance) OpenInterest(ctx context.Context, symbol currency.Pair) (*OpenInterestData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	var resp *OpenInterestData
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/dapi/v1/openInterest?symbol="+symbolValue, cFuturesDefaultRate, &resp)
}

// CFuturesQuarterlyContractSettlementPrice retrieves coin margined futures quarterly contract settlement price
func (b *Binance) CFuturesQuarterlyContractSettlementPrice(ctx context.Context, pair currency.Pair) ([]SettlementPrice, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var resp []SettlementPrice
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/delivery-price?pair="+pair.String(), cFuturesDefaultRate, &resp)
}

// GetOpenInterestStats gets open interest stats for a symbol
func (b *Binance) GetOpenInterestStats(ctx context.Context, pair, contractType, period string, limit int64, startTime, endTime time.Time) ([]OpenInterestStats, error) {
	if !common.StringDataCompare(validContractType, contractType) {
		return nil, errors.New("invalid contractType")
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period")
	}
	params := url.Values{}
	params.Set("contractType", contractType)
	params.Set("period", period)
	if pair != "" {
		params.Set("pair", pair)
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []OpenInterestStats
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/openInterestHist?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetTraderFuturesAccountRatio gets a traders futures account long/short ratio
func (b *Binance) GetTraderFuturesAccountRatio(ctx context.Context, pair currency.Pair, period string, limit int64, startTime, endTime time.Time) ([]TopTraderAccountRatio, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period")
	}
	params := url.Values{}
	params.Set("pair", pair.String())
	params.Set("period", period)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []TopTraderAccountRatio
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/topLongShortAccountRatio?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetTraderFuturesPositionsRatio gets a traders futures positions' long/short ratio
func (b *Binance) GetTraderFuturesPositionsRatio(ctx context.Context, pair currency.Pair, period string, limit int64, startTime, endTime time.Time) ([]TopTraderPositionRatio, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period")
	}
	params := url.Values{}
	params.Set("pair", pair.String())
	params.Set("period", period)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []TopTraderPositionRatio
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/topLongShortPositionRatio?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetMarketRatio gets global long/short ratio
func (b *Binance) GetMarketRatio(ctx context.Context, pair currency.Pair, period string, limit int64, startTime, endTime time.Time) ([]TopTraderPositionRatio, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period")
	}
	params := url.Values{}
	params.Set("pair", pair.String())
	params.Set("period", period)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []TopTraderPositionRatio
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/globalLongShortAccountRatio?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetFuturesTakerVolume gets futures taker buy/sell volumes
func (b *Binance) GetFuturesTakerVolume(ctx context.Context, pair currency.Pair, contractType, period string, limit int64, startTime, endTime time.Time) ([]TakerBuySellVolume, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if !common.StringDataCompare(validContractType, contractType) {
		return nil, errors.New("invalid contractType")
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period parsed")
	}
	params := url.Values{}
	params.Set("pair", pair.String())
	params.Set("contractType", contractType)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("period", period)
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []TakerBuySellVolume
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/takerBuySellVol?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// GetFuturesBasisData gets futures basis data
func (b *Binance) GetFuturesBasisData(ctx context.Context, pair currency.Pair, contractType, period string, limit int64, startTime, endTime time.Time) ([]FuturesBasisData, error) {
	if pair.IsEmpty() {
		return nil, currency.ErrCurrencyPairEmpty
	}
	if !common.StringDataCompare(validContractType, contractType) {
		return nil, errors.New("invalid contractType")
	}
	if !common.StringDataCompare(validFuturesIntervals, period) {
		return nil, errors.New("invalid period parsed")
	}
	params := url.Values{}
	params.Set("pair", pair.String())
	params.Set("contractType", contractType)
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Set("period", period)
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []FuturesBasisData
	return resp, b.SendHTTPRequest(ctx, exchange.RestCoinMargined, "/futures/data/basis?"+params.Encode(), cFuturesDefaultRate, &resp)
}

// FuturesNewOrder sends a new futures order to the exchange
func (b *Binance) FuturesNewOrder(ctx context.Context, x *FuturesNewOrderRequest) (*FuturesOrderPlaceData, error) {
	symbolValue, err := b.FormatSymbol(x.Symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	params.Set("side", x.Side)
	if x.PositionSide != "" {
		if !common.StringDataCompare(validPositionSide, x.PositionSide) {
			return nil, errors.New("invalid positionSide")
		}
		params.Set("positionSide", x.PositionSide)
	}
	params.Set("type", x.OrderType)
	if string(x.TimeInForce) != "" {
		params.Set("timeInForce", string(x.TimeInForce))
	}
	if x.ReduceOnly {
		params.Set("reduceOnly", "true")
	}
	if x.NewClientOrderID != "" {
		params.Set("newClientOrderID", x.NewClientOrderID)
	}
	if x.ClosePosition != "" {
		params.Set("closePosition", x.ClosePosition)
	}
	if x.WorkingType != "" {
		if !common.StringDataCompare(validWorkingType, x.WorkingType) {
			return nil, errors.New("invalid workingType")
		}
		params.Set("workingType", x.WorkingType)
	}
	if x.NewOrderRespType != "" {
		if !common.StringDataCompare(validNewOrderRespType, x.NewOrderRespType) {
			return nil, errors.New("invalid newOrderRespType")
		}
		params.Set("newOrderRespType", x.NewOrderRespType)
	}
	if x.Quantity != 0 {
		params.Set("quantity", strconv.FormatFloat(x.Quantity, 'f', -1, 64))
	}
	if x.Price != 0 {
		params.Set("price", strconv.FormatFloat(x.Price, 'f', -1, 64))
	}
	if x.StopPrice != 0 {
		params.Set("stopPrice", strconv.FormatFloat(x.StopPrice, 'f', -1, 64))
	}
	if x.ActivationPrice != 0 {
		params.Set("activationPrice", strconv.FormatFloat(x.ActivationPrice, 'f', -1, 64))
	}
	if x.CallbackRate != 0 {
		params.Set("callbackRate", strconv.FormatFloat(x.CallbackRate, 'f', -1, 64))
	}
	if x.PriceProtect {
		params.Set("priceProtect", "TRUE")
	}
	var resp *FuturesOrderPlaceData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/order", params, cFuturesOrdersDefaultRate, &resp)
}

// FuturesBatchOrder sends a batch order request
func (b *Binance) FuturesBatchOrder(ctx context.Context, data []PlaceBatchOrderData) ([]FuturesOrderPlaceData, error) {
	for x := range data {
		unformattedPair, err := currency.NewPairFromString(data[x].Symbol)
		if err != nil {
			return nil, err
		}
		formattedPair, err := b.FormatExchangeCurrency(unformattedPair, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		data[x].Symbol = formattedPair.String()
		if data[x].PositionSide != "" {
			if !common.StringDataCompare(validPositionSide, data[x].PositionSide) {
				return nil, errors.New("invalid positionSide")
			}
		}
		if data[x].WorkingType != "" {
			if !common.StringDataCompare(validWorkingType, data[x].WorkingType) {
				return nil, errors.New("invalid workingType")
			}
		}
		if data[x].NewOrderRespType != "" {
			if !common.StringDataCompare(validNewOrderRespType, data[x].NewOrderRespType) {
				return nil, errors.New("invalid newOrderRespType")
			}
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("batchOrders", string(jsonData))
	var resp []FuturesOrderPlaceData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/batchOrders", params, cFuturesBatchOrdersRate, &resp)
}

// FuturesBatchCancelOrders sends a batch request to cancel orders
func (b *Binance) FuturesBatchCancelOrders(ctx context.Context, symbol currency.Pair, orderList, origClientOrderIDList []string) ([]BatchCancelOrderData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if len(orderList) != 0 {
		jsonOrderList, err := json.Marshal(orderList)
		if err != nil {
			return nil, err
		}
		params.Set("orderIdList", string(jsonOrderList))
	}
	if len(origClientOrderIDList) != 0 {
		jsonCliOrdIDList, err := json.Marshal(origClientOrderIDList)
		if err != nil {
			return nil, err
		}
		params.Set("origClientOrderIdList", string(jsonCliOrdIDList))
	}
	var resp []BatchCancelOrderData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodDelete, "/dapi/v1/batchOrders", params, cFuturesOrdersDefaultRate, &resp)
}

// FuturesGetOrderData gets futures order data
func (b *Binance) FuturesGetOrderData(ctx context.Context, symbol currency.Pair, orderID, origClientOrderID string) (*FuturesOrderGetData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	var resp *FuturesOrderGetData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/order", params, cFuturesOrdersDefaultRate, &resp)
}

// FuturesCancelOrder cancels a futures order
func (b *Binance) FuturesCancelOrder(ctx context.Context, symbol currency.Pair, orderID, origClientOrderID string) (*FuturesOrderGetData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	var resp *FuturesOrderGetData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodDelete, "/dapi/v1/order", params, cFuturesOrdersDefaultRate, &resp)
}

// FuturesCancelAllOpenOrders cancels a futures order
func (b *Binance) FuturesCancelAllOpenOrders(ctx context.Context, symbol currency.Pair) (*GenericAuthResponse, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	var resp *GenericAuthResponse
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodDelete, "/dapi/v1/allOpenOrders", params, cFuturesOrdersDefaultRate, &resp)
}

// AutoCancelAllOpenOrders cancels all open futures orders
// countdownTime 1000 = 1s, example - to cancel all orders after 30s (countdownTime: 30000)
func (b *Binance) AutoCancelAllOpenOrders(ctx context.Context, symbol currency.Pair, countdownTime int64) (*AutoCancelAllOrdersData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	params.Set("countdownTime", strconv.FormatInt(countdownTime, 10))
	var resp *AutoCancelAllOrdersData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/countdownCancelAll", params, cFuturesCancelAllOrdersRate, &resp)
}

// FuturesOpenOrderData gets open order data for CoinMarginedFutures
func (b *Binance) FuturesOpenOrderData(ctx context.Context, symbol currency.Pair, orderID, origClientOrderID string) (*FuturesOrderGetData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	var resp *FuturesOrderGetData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/openOrder", params, cFuturesOrdersDefaultRate, &resp)
}

// GetFuturesAllOpenOrders gets all open orders data for CoinMarginedFutures,
func (b *Binance) GetFuturesAllOpenOrders(ctx context.Context, symbol currency.Pair, pair string) ([]FuturesOrderData, error) {
	var p string
	var err error
	rateLimit := cFuturesGetAllOpenOrdersRate
	params := url.Values{}
	if !symbol.IsEmpty() {
		rateLimit = cFuturesOrdersDefaultRate
		p, err = b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", p)
	} else {
		// extend the receive window when all currencies to prevent "recvwindow" error
		params.Set("recvWindow", "10000")
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []FuturesOrderData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/openOrders", params, rateLimit, &resp)
}

// GetAllFuturesOrders gets all orders active cancelled or filled
func (b *Binance) GetAllFuturesOrders(ctx context.Context, symbol, pair currency.Pair, startTime, endTime time.Time, orderID, limit int64) ([]FuturesOrderData, error) {
	params := url.Values{}
	rateLimit := cFuturesPairOrdersRate
	if !symbol.IsEmpty() {
		rateLimit = cFuturesSymbolOrdersRate
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if !pair.IsEmpty() {
		params.Set("pair", pair.String())
	}
	if orderID != 0 {
		params.Set("orderID", strconv.FormatInt(orderID, 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []FuturesOrderData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/allOrders", params, rateLimit, &resp)
}

// GetFuturesAccountBalance gets account balance data for CoinMarginedFutures, account
func (b *Binance) GetFuturesAccountBalance(ctx context.Context) ([]FuturesAccountBalanceData, error) {
	var resp []FuturesAccountBalanceData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/balance", nil, cFuturesDefaultRate, &resp)
}

// GetFuturesAccountInfo gets account info data for CoinMarginedFutures, account
func (b *Binance) GetFuturesAccountInfo(ctx context.Context) (*FuturesAccountInformation, error) {
	var resp *FuturesAccountInformation
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/account", nil, cFuturesAccountInformationRate, &resp)
}

// FuturesChangeInitialLeverage changes initial leverage for the account
func (b *Binance) FuturesChangeInitialLeverage(ctx context.Context, symbol currency.Pair, leverage float64) (*FuturesLeverageData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	if leverage < 1 || leverage > 125 {
		return nil, errors.New("invalid leverage")
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	params.Set("leverage", strconv.FormatFloat(leverage, 'f', -1, 64))
	var resp *FuturesLeverageData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/leverage", params, cFuturesDefaultRate, &resp)
}

// FuturesChangeMarginType changes margin type
func (b *Binance) FuturesChangeMarginType(ctx context.Context, symbol currency.Pair, marginType string) (*GenericAuthResponse, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	if !common.StringDataCompare(validMarginType, marginType) {
		return nil, errors.New("invalid marginType")
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	params.Set("marginType", marginType)
	var resp *GenericAuthResponse
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/marginType", params, cFuturesDefaultRate, &resp)
}

// ModifyIsolatedPositionMargin changes margin for an isolated position
func (b *Binance) ModifyIsolatedPositionMargin(ctx context.Context, symbol currency.Pair, positionSide, changeType string, amount float64) (*FuturesMarginUpdatedResponse, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	cType, ok := validMarginChange[changeType]
	if !ok {
		return nil, errors.New("invalid changeType")
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	if positionSide != "" {
		if !common.StringDataCompare(validPositionSide, positionSide) {
			return nil, errors.New("invalid positionSide")
		}
		params.Set("positionSide", positionSide)
	}
	params.Set("type", strconv.FormatInt(cType, 10))
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	var resp *FuturesMarginUpdatedResponse
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodPost, "/dapi/v1/positionMargin", params, cFuturesDefaultRate, &resp)
}

// FuturesMarginChangeHistory gets past margin changes for positions
func (b *Binance) FuturesMarginChangeHistory(ctx context.Context, symbol currency.Pair, changeType string, startTime, endTime time.Time, limit int64) ([]GetPositionMarginChangeHistoryData, error) {
	symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
	if err != nil {
		return nil, err
	}
	cType, ok := validMarginChange[changeType]
	if !ok {
		return nil, errors.New("invalid changeType")
	}
	params := url.Values{}
	params.Set("symbol", symbolValue)
	params.Set("type", strconv.FormatInt(cType, 10))
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []GetPositionMarginChangeHistoryData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/positionMargin/history", params, cFuturesDefaultRate, &resp)
}

// FuturesPositionsInfo gets futures positions info
// "pair" for coinmarginedfutures in GCT terms is the pair base
// eg ADAUSD_PERP the "pair" parameter is ADAUSD
func (b *Binance) FuturesPositionsInfo(ctx context.Context, marginAsset, pair string) ([]FuturesPositionInformation, error) {
	params := url.Values{}
	if marginAsset != "" {
		params.Set("marginAsset", marginAsset)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []FuturesPositionInformation
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/positionRisk", params, cFuturesDefaultRate, &resp)
}

// FuturesTradeHistory gets trade history for CoinMarginedFutures, account
func (b *Binance) FuturesTradeHistory(ctx context.Context, symbol currency.Pair, pair string, startTime, endTime time.Time, limit, fromID int64) ([]FuturesAccountTradeList, error) {
	params := url.Values{}
	rateLimit := cFuturesPairOrdersRate
	if !symbol.IsEmpty() {
		rateLimit = cFuturesSymbolOrdersRate
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if pair != "" {
		params.Set("pair", pair)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if fromID != 0 {
		params.Set("fromId", strconv.FormatInt(fromID, 10))
	}
	var resp []FuturesAccountTradeList
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/userTrades", params, rateLimit, &resp)
}

// FuturesIncomeHistory gets income history for CoinMarginedFutures,
func (b *Binance) FuturesIncomeHistory(ctx context.Context, symbol currency.Pair, incomeType string, startTime, endTime time.Time, limit int64) ([]FuturesIncomeHistoryData, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if incomeType != "" {
		if !common.StringDataCompare(validIncomeType, incomeType) {
			return nil, fmt.Errorf("invalid incomeType: %v", incomeType)
		}
		params.Set("incomeType", incomeType)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp []FuturesIncomeHistoryData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/income", params, cFuturesIncomeHistoryRate, &resp)
}

// FuturesNotionalBracket gets futures notional bracket
func (b *Binance) FuturesNotionalBracket(ctx context.Context, pair string) ([]NotionalBracketData, error) {
	params := url.Values{}
	if pair != "" {
		params.Set("pair", pair)
	}
	var resp []NotionalBracketData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/leverageBracket", params, cFuturesDefaultRate, &resp)
}

// FuturesForceOrders gets futures forced orders
func (b *Binance) FuturesForceOrders(ctx context.Context, symbol currency.Pair, autoCloseType string, startTime, endTime time.Time) ([]ForcedOrdersData, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	if autoCloseType != "" {
		if !common.StringDataCompare(validAutoCloseTypes, autoCloseType) {
			return nil, errors.New("invalid autoCloseType")
		}
		params.Set("autoCloseType", autoCloseType)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return nil, common.ErrStartAfterEnd
		}
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp []ForcedOrdersData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/forceOrders", params, cFuturesDefaultRate, &resp)
}

// FuturesPositionsADLEstimate estimates ADL on positions
func (b *Binance) FuturesPositionsADLEstimate(ctx context.Context, symbol currency.Pair) ([]ADLEstimateData, error) {
	params := url.Values{}
	if !symbol.IsEmpty() {
		symbolValue, err := b.FormatSymbol(symbol, asset.CoinMarginedFutures)
		if err != nil {
			return nil, err
		}
		params.Set("symbol", symbolValue)
	}
	var resp []ADLEstimateData
	return resp, b.SendAuthHTTPRequest(ctx, exchange.RestCoinMargined, http.MethodGet, "/dapi/v1/adlQuantile", params, cFuturesAccountInformationRate, &resp)
}

// FetchCoinMarginExchangeLimits fetches coin margined order execution limits
func (b *Binance) FetchCoinMarginExchangeLimits(ctx context.Context) ([]order.MinMaxLevel, error) {
	coinFutures, err := b.FuturesExchangeInfo(ctx)
	if err != nil {
		return nil, err
	}
	limits := make([]order.MinMaxLevel, 0, len(coinFutures.Symbols))
	for x := range coinFutures.Symbols {
		symbol := strings.Split(coinFutures.Symbols[x].Symbol, currency.UnderscoreDelimiter)
		var cp currency.Pair
		cp, err = currency.NewPairFromStrings(symbol[0], symbol[1])
		if err != nil {
			return nil, err
		}
		if len(coinFutures.Symbols[x].Filters) < 6 {
			continue
		}
		limits = append(limits, order.MinMaxLevel{
			Pair:                    cp,
			Asset:                   asset.CoinMarginedFutures,
			MinPrice:                coinFutures.Symbols[x].Filters[0].MinPrice,
			MaxPrice:                coinFutures.Symbols[x].Filters[0].MaxPrice,
			PriceStepIncrementSize:  coinFutures.Symbols[x].Filters[0].TickSize,
			MaximumBaseAmount:       coinFutures.Symbols[x].Filters[1].MaxQty,
			MinimumBaseAmount:       coinFutures.Symbols[x].Filters[1].MinQty,
			AmountStepIncrementSize: coinFutures.Symbols[x].Filters[1].StepSize,
			MarketMinQty:            coinFutures.Symbols[x].Filters[2].MinQty,
			MarketMaxQty:            coinFutures.Symbols[x].Filters[2].MaxQty,
			MarketStepIncrementSize: coinFutures.Symbols[x].Filters[2].StepSize,
			MaxTotalOrders:          coinFutures.Symbols[x].Filters[3].Limit,
			MaxAlgoOrders:           coinFutures.Symbols[x].Filters[4].Limit,
			MultiplierUp:            coinFutures.Symbols[x].Filters[5].MultiplierUp,
			MultiplierDown:          coinFutures.Symbols[x].Filters[5].MultiplierDown,
			MultiplierDecimal:       coinFutures.Symbols[x].Filters[5].MultiplierDecimal,
		})
	}
	return limits, nil
}
