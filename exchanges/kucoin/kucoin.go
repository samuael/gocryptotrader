package kucoin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Kucoin is the overarching type across this package
type Kucoin struct {
	exchange.Base
}

const (
	kucoinAPIURL        = "https://api.kucoin.com/api"
	kucoinAPIVersion    = "1"
	kucoinAPIKeyVersion = "2"

	// Public endpoints
	kucoinGetSymbols             = "/v2/symbols"
	kucoinGetTicker              = "/v1/market/orderbook/level1"
	kucoinGetAllTickers          = "/v1/market/allTickers"
	kucoinGet24hrStats           = "/v1/market/stats"
	kucoinGetMarketList          = "/v1/markets"
	kucoinGetPartOrderbook20     = "/v1/market/orderbook/level2_20"
	kucoinGetPartOrderbook100    = "/v1/market/orderbook/level2_100"
	kucoinGetTradeHistory        = "/v1/market/histories"
	kucoinGetKlines              = "/v1/market/candles"
	kucoinGetCurrencies          = "/v1/currencies"
	kucoinGetCurrency            = "/v2/currencies/"
	kucoinGetFiatPrice           = "/v1/prices"
	kucoinGetMarkPrice           = "/v1/mark-price/%s/current"
	kucoinGetMarginConfiguration = "/v1/margin/config"
	kucoinGetServerTime          = "/v1/timestamp"
	kucoinGetServiceStatus       = "/v1/status"

	// Authenticated endpoints
	kucoinGetOrderbook         = "/v3/market/orderbook/level2"
	kucoinGetMarginAccount     = "/v1/margin/account"
	kucoinGetMarginRiskLimit   = "/v1/risk/limit/strategy"
	kucoinBorrowOrder          = "/v1/margin/borrow"
	kucoinGetOutstandingRecord = "/v1/margin/borrow/outstanding"
	kucoinGetRepaidRecord      = "/v1/margin/borrow/repaid"
	kucoinOneClickRepayment    = "/v1/margin/repay/all"
	kucoinRepaySingleOrder     = "/v1/margin/repay/single"
	kucoinPostLendOrder        = "/v1/margin/lend"
	kucoinCancelLendOrder      = "/v1/margin/lend/%s"
	kucoinSetAutoLend          = "/v1/margin/toggle-auto-lend"
	kucoinGetActiveOrder       = "/v1/margin/lend/active"
	kucoinGetLendHistory       = "/v1/margin/lend/done"
	kucoinGetUnsettleLendOrder = "/v1/margin/lend/trade/unsettled"
	kucoinGetSettleLendOrder   = "/v1/margin/lend/trade/settled"
	kucoinGetAccountLendRecord = "/v1/margin/lend/assets"
	kucoinGetLendingMarketData = "/v1/margin/market"
	kucoinGetMarginTradeData   = "/v1/margin/trade/last"

	kucoinGetIsolatedMarginPairConfig            = "/v1/isolated/symbols"
	kucoinGetIsolatedMarginAccountInfo           = "/v1/isolated/accounts"
	kucoinGetSingleIsolatedMarginAccountInfo     = "/v1/isolated/account/%s"
	kucoinInitiateIsolatedMarginBorrowing        = "/v1/isolated/borrow"
	kucoinGetIsolatedOutstandingRepaymentRecords = "/v1/isolated/borrow/outstanding"
	kucoinGetIsolatedMarginRepaymentRecords      = "/v1/isolated/borrow/repaid"
	kucoinInitiateIsolatedMarginQuickRepayment   = "/v1/isolated/repay/all"
	kucoinInitiateIsolatedMarginSingleRepayment  = "/v1/isolated/repay/single"

	kucoinPostOrder        = "/v1/orders"
	kucoinPostMarginOrder  = "/v1/margin/order"
	kucoinPostBulkOrder    = "/v1/orders/multi"
	kucoinOrderByID        = "/v1/orders/%s"             // used by CancelSingleOrder and GetOrderByID
	kucoinOrderByClientOID = "/v1/order/client-order/%s" // used by CancelOrderByClientOID and GetOrderByClientOID
	kucoinOrders           = "/v1/orders"                // used by CancelAllOpenOrders and GetOrders
	kucoinGetRecentOrders  = "/v1/limit/orders"

	kucoinGetFills       = "/v1/fills"
	kucoinGetRecentFills = "/v1/limit/fills"

	kucoinStopOrder                 = "/v1/stop-order"
	kucoinStopOrderByID             = "/v1/stop-order/%s"
	kucoinCancelAllStopOrder        = "/v1/stop-order/cancel"
	kucoinGetStopOrderByClientID    = "/v1/stop-order/queryOrderByClientOid"
	kucoinCancelStopOrderByClientID = "/v1/stop-order/cancelOrderByClientOid"

	// account
	kucoinAccount                        = "/v1/accounts"
	kucoinGetAccount                     = "/v1/accounts/%s"
	kucoinGetAccountLedgers              = "/v1/accounts/ledgers"
	kucoinGetSubAccountBalance           = "/v1/sub-accounts/%s"
	kucoinGetAggregatedSubAccountBalance = "/v1/sub-accounts"
	kucoinGetTransferableBalance         = "/v1/accounts/transferable"
	kucoinTransferMainToSubAccount       = "/v2/accounts/sub-transfer"
	kucoinInnerTransfer                  = "/v2/accounts/inner-transfer"

	// deposit
	kucoinCreateDepositAddress     = "/v1/deposit-addresses"
	kucoinGetDepositAddressV2      = "/v2/deposit-addresses"
	kucoinGetDepositAddressV1      = "/v1/deposit-addresses"
	kucoinGetDepositList           = "/v1/deposits"
	kucoinGetHistoricalDepositList = "/v1/hist-deposits"

	// withdrawal
	kucoinWithdrawal                  = "/v1/withdrawals"
	kucoinGetHistoricalWithdrawalList = "/v1/hist-withdrawals"
	kucoinGetWithdrawalQuotas         = "/v1/withdrawals/quotas"
	kucoinCancelWithdrawal            = "/v1/withdrawals/%s"

	kucoinBasicFee   = "/v1/base-fee"
	kucoinTradingFee = "/v1/trade-fees"
)

// GetSymbols gets pairs details on the exchange
func (ku *Kucoin) GetSymbols(ctx context.Context, ccy string) ([]SymbolInfo, error) {
	var params url.Values
	if ccy != "" {
		params.Set("market", ccy)
	}
	resp := []SymbolInfo{}
	return resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetSymbols, params), &resp)
}

// GetTicker gets pair ticker information
func (ku *Kucoin) GetTicker(ctx context.Context, pair string) (*Ticker, error) {
	resp := Ticker{}
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	return &resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetTicker, params), &resp)
}

// GetAllTickers gets all trading pair ticker information including 24h volume
func (ku *Kucoin) GetAllTickers(ctx context.Context) ([]TickerInfo, error) {
	resp := struct {
		Time    kucoinTimeMilliSec `json:"time"`
		Tickers []TickerInfo       `json:"ticker"`
	}{}
	return resp.Tickers, ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetAllTickers, &resp)
}

// Get24hrStats get the statistics of the specified pair in the last 24 hours
func (ku *Kucoin) Get24hrStats(ctx context.Context, pair string) (*Stats24hrs, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	resp := Stats24hrs{}
	return &resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGet24hrStats, params), &resp)
}

// GetMarketList get the transaction currency for the entire trading market
func (ku *Kucoin) GetMarketList(ctx context.Context) ([]string, error) {
	resp := []string{}
	return resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetMarketList, &resp)
}

func processOB(ob [][2]string) ([]orderbook.Item, error) {
	o := make([]orderbook.Item, len(ob))
	for x := range ob {
		amount, err := strconv.ParseFloat(ob[x][1], 64)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(ob[x][0], 64)
		if err != nil {
			return nil, err
		}
		o[x] = orderbook.Item{
			Price:  price,
			Amount: amount,
		}
	}
	return o, nil
}

func constructOrderbook(o *orderbookResponse) (*Orderbook, error) {
	var (
		s   Orderbook
		err error
	)
	s.Bids, err = processOB(o.Bids)
	if err != nil {
		return nil, err
	}
	s.Asks, err = processOB(o.Asks)
	if err != nil {
		return nil, err
	}
	s.Time = o.Time.Time()
	return &s, err
}

// GetPartOrderbook20 gets orderbook for a specified pair with depth 20
func (ku *Kucoin) GetPartOrderbook20(ctx context.Context, pair string) (*Orderbook, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	var o orderbookResponse
	err := ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetPartOrderbook20, params), &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetPartOrderbook100 gets orderbook for a specified pair with depth 100
func (ku *Kucoin) GetPartOrderbook100(ctx context.Context, pair string) (*Orderbook, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	var o orderbookResponse
	err := ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetPartOrderbook100, params), &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetOrderbook gets full orderbook for a specified pair
func (ku *Kucoin) GetOrderbook(ctx context.Context, pair string) (*Orderbook, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	var o orderbookResponse
	err := ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetOrderbook, params), nil, &o)
	if err != nil {
		return nil, err
	}
	return constructOrderbook(&o)
}

// GetTradeHistory gets trade history of the specified pair
func (ku *Kucoin) GetTradeHistory(ctx context.Context, pair string) ([]Trade, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	resp := []Trade{}
	return resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetTradeHistory, params), &resp)
}

// GetKlines gets kline of the specified pair
func (ku *Kucoin) GetKlines(ctx context.Context, pair, period string, start, end time.Time) ([]Kline, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	var params url.Values
	params.Set("symbol", pair)
	if period == "" {
		return nil, errors.New("period can not be empty")
	}
	if !common.StringDataContains(validPeriods, period) {
		return nil, errors.New("invalid period")
	}
	params.Set("type", period)
	if !start.IsZero() {
		params.Set("startAt", strconv.FormatInt(start.Unix(), 10))
	}
	if !end.IsZero() {
		params.Set("endAt", strconv.FormatInt(end.Unix(), 10))
	}
	resp := [][7]string{}
	err := ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetKlines, params), &resp)
	if err != nil {
		return nil, err
	}
	klines := make([]Kline, len(resp))
	for i := range resp {
		t, err := strconv.ParseInt(resp[i][0], 10, 64)
		if err != nil {
			return nil, err
		}
		klines[i].StartTime = time.Unix(t, 0)
		klines[i].Open, err = strconv.ParseFloat(resp[i][1], 64)
		if err != nil {
			return nil, err
		}
		klines[i].Close, err = strconv.ParseFloat(resp[i][2], 64)
		if err != nil {
			return nil, err
		}
		klines[i].High, err = strconv.ParseFloat(resp[i][3], 64)
		if err != nil {
			return nil, err
		}
		klines[i].Low, err = strconv.ParseFloat(resp[i][4], 64)
		if err != nil {
			return nil, err
		}
		klines[i].Volume, err = strconv.ParseFloat(resp[i][5], 64)
		if err != nil {
			return nil, err
		}
		klines[i].Amount, err = strconv.ParseFloat(resp[i][6], 64)
		if err != nil {
			return nil, err
		}
	}
	return klines, nil
}

// GetCurrencies gets list of currencies
func (ku *Kucoin) GetCurrencies(ctx context.Context) ([]Currency, error) {
	resp := []Currency{}
	return resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetCurrencies, &resp)
}

// GetCurrencyDetail gets currency detail using currency code and chain information.
func (ku *Kucoin) GetCurrencyDetail(ctx context.Context, ccy, chain string) (*CurrencyDetail, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	if chain != "" {
		params.Set("chain", chain)
	}
	resp := CurrencyDetail{}
	return &resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetCurrency+strings.ToUpper(ccy), params), &resp)
}

// GetFiatPrice gets fiat prices of currencies, default base currency is USD
func (ku *Kucoin) GetFiatPrice(ctx context.Context, base, currencies string) (map[string]string, error) {
	var params url.Values
	if base != "" {
		params.Set("base", base)
	}
	if currencies != "" {
		params.Set("currencies", currencies)
	}
	resp := map[string]string{}
	return resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues(kucoinGetFiatPrice, params), &resp)
}

// GetMarkPrice gets index price of the specified pair
func (ku *Kucoin) GetMarkPrice(ctx context.Context, pair string) (*MarkPrice, error) {
	if pair == "" {
		return nil, currency.ErrCurrencyPairEmpty
	}
	resp := MarkPrice{}
	return &resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, fmt.Sprintf(kucoinGetMarkPrice, pair), &resp)
}

// GetMarginConfiguration gets configure info of the margin
func (ku *Kucoin) GetMarginConfiguration(ctx context.Context) (*MarginConfiguration, error) {
	resp := MarginConfiguration{}
	return &resp, ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetMarginConfiguration, &resp)
}

// GetMarginAccount gets configure info of the margin
func (ku *Kucoin) GetMarginAccount(ctx context.Context) (*MarginAccounts, error) {
	resp := MarginAccounts{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetMarginAccount, nil, &resp)
}

// GetMarginRiskLimit gets cross/isolated margin risk limit, default model is cross margin
func (ku *Kucoin) GetMarginRiskLimit(ctx context.Context, marginModel string) ([]MarginRiskLimit, error) {
	var params url.Values
	if marginModel != "" {
		params.Set("marginModel", marginModel)
	}
	resp := []MarginRiskLimit{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetMarginRiskLimit, params), nil, &resp)
}

// PostBorrowOrder used to post borrow order
func (ku *Kucoin) PostBorrowOrder(ctx context.Context, ccy, orderType, term string, size, maxRate float64) (*PostBorrowOrderResp, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if orderType == "" {
		return nil, errors.New("orderType can not be empty")
	}
	params["type"] = orderType
	if size == 0 {
		return nil, errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	if maxRate != 0 {
		params["maxRate"] = strconv.FormatFloat(maxRate, 'f', -1, 64)
	}
	if term != "" {
		params["term"] = term
	}
	resp := PostBorrowOrderResp{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinBorrowOrder, params, &resp)
}

// GetBorrowOrder gets borrow order information
func (ku *Kucoin) GetBorrowOrder(ctx context.Context, orderID string) (*BorrowOrder, error) {
	if orderID == "" {
		return nil, errors.New("empty orderID")
	}
	var params url.Values
	params.Set("orderId", orderID)
	resp := BorrowOrder{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinBorrowOrder, params), nil, &resp)
}

// GetOutstandingRecord gets outstanding record information
func (ku *Kucoin) GetOutstandingRecord(ctx context.Context, ccy string) ([]OutstandingRecord, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := []OutstandingRecord{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetOutstandingRecord, params), nil, &resp)
}

// GetRepaidRecord gets repaid record information
func (ku *Kucoin) GetRepaidRecord(ctx context.Context, ccy string) ([]RepaidRecord, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := []RepaidRecord{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetRepaidRecord, params), nil, &resp)
}

// OneClickRepayment used to compplete repayment in single go
func (ku *Kucoin) OneClickRepayment(ctx context.Context, ccy, sequence string, size float64) error {
	if ccy == "" {
		return currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if sequence == "" {
		return errors.New("sequence can not be empty")
	}
	params["sequence"] = sequence
	if size == 0 {
		return errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinOneClickRepayment, params, &struct{}{})
}

// SingleOrderRepayment used to repay single order
func (ku *Kucoin) SingleOrderRepayment(ctx context.Context, ccy, tradeID string, size float64) error {
	if ccy == "" {
		return currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if tradeID == "" {
		return errors.New("tradeId can not be empty")
	}
	params["tradeId"] = tradeID
	if size == 0 {
		return errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinRepaySingleOrder, params, &struct{}{})
}

// PostLendOrder used to create lend order
func (ku *Kucoin) PostLendOrder(ctx context.Context, ccy string, dailyIntRate, size float64, term int64) (string, error) {
	if ccy == "" {
		return "", currency.ErrCurrencyPairEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if dailyIntRate == 0 {
		return "", errors.New("dailyIntRate can not be zero")
	}
	params["dailyIntRate"] = strconv.FormatFloat(dailyIntRate, 'f', -1, 64)
	if size == 0 {
		return "", errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
	if term == 0 {
		return "", errors.New("term can not be zero")
	}
	params["term"] = strconv.FormatInt(term, 10)
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}
	return resp.OrderID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostLendOrder, params, &resp)
}

// CancelLendOrder used to cancel lend order
func (ku *Kucoin) CancelLendOrder(ctx context.Context, orderID string) error {
	resp := struct {
		Error
	}{}
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinCancelLendOrder, orderID), nil, &resp)
}

// SetAutoLend used to set up the automatic lending for a specified currency
func (ku *Kucoin) SetAutoLend(ctx context.Context, ccy string, dailyIntRate, retainSize float64, term int64, isEnable bool) error {
	if ccy == "" {
		return currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if dailyIntRate == 0 {
		return errors.New("dailyIntRate can not be zero")
	}
	params["dailyIntRate"] = strconv.FormatFloat(dailyIntRate, 'f', -1, 64)
	if retainSize == 0 {
		return errors.New("retainSize can not be zero")
	}
	params["retainSize"] = strconv.FormatFloat(retainSize, 'f', -1, 64)
	if term == 0 {
		return errors.New("term can not be zero")
	}
	params["term"] = strconv.FormatInt(term, 10)
	params["isEnable"] = isEnable
	resp := struct {
		Error
	}{}
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinSetAutoLend, params, &resp)
}

// GetActiveOrder gets active lend orders
func (ku *Kucoin) GetActiveOrder(ctx context.Context, ccy string) ([]LendOrder, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := struct {
		Data []LendOrder `json:"items"`
		Error
	}{}
	return resp.Data, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetActiveOrder, params), nil, &resp)
}

// GetLendHistory gets lend orders
func (ku *Kucoin) GetLendHistory(ctx context.Context, ccy string) ([]LendOrderHistory, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := struct {
		Data []LendOrderHistory `json:"items"`
		Error
	}{}
	return resp.Data, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetLendHistory, params), nil, &resp)
}

// GetUnsettleLendOrder gets outstanding lend order list
func (ku *Kucoin) GetUnsettleLendOrder(ctx context.Context, ccy string) ([]UnsettleLendOrder, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := struct {
		Data []UnsettleLendOrder `json:"items"`
		Error
	}{}
	return resp.Data, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetUnsettleLendOrder, params), nil, &resp)
}

// GetSettleLendOrder gets settle lend orders
func (ku *Kucoin) GetSettleLendOrder(ctx context.Context, ccy string) ([]SettleLendOrder, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := struct {
		Data []SettleLendOrder `json:"items"`
		Error
	}{}
	return resp.Data, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetSettleLendOrder, params), nil, &resp)
}

// GetAccountLendRecord get the lending history of the main account
func (ku *Kucoin) GetAccountLendRecord(ctx context.Context, ccy string) ([]LendRecord, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	resp := []LendRecord{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetAccountLendRecord, params), nil, &resp)
}

// GetLendingMarketData get the lending market data
func (ku *Kucoin) GetLendingMarketData(ctx context.Context, ccy string, term int64) ([]LendMarketData, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	if term != 0 {
		params.Set("term", strconv.FormatInt(term, 10))
	}
	resp := []LendMarketData{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetLendingMarketData, params), nil, &resp)
}

// GetMarginTradeData get the last 300 fills in the lending and borrowing market
func (ku *Kucoin) GetMarginTradeData(ctx context.Context, ccy string) ([]MarginTradeData, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	resp := []MarginTradeData{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetMarginTradeData, params), nil, &resp)
}

// GetIsolatedMarginPairConfig get the current isolated margin trading pair configuration
func (ku *Kucoin) GetIsolatedMarginPairConfig(ctx context.Context) ([]IsolatedMarginPairConfig, error) {
	resp := []IsolatedMarginPairConfig{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetIsolatedMarginPairConfig, nil, &resp)
}

// GetIsolatedMarginAccountInfo get all isolated margin accounts of the current user
func (ku *Kucoin) GetIsolatedMarginAccountInfo(ctx context.Context, balanceCurrency string) (*IsolatedMarginAccountInfo, error) {
	var params url.Values
	if balanceCurrency != "" {
		params.Set("balanceCurrency", balanceCurrency)
	}
	resp := IsolatedMarginAccountInfo{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedMarginAccountInfo, params), nil, &resp)
}

// GetSingleIsolatedMarginAccountInfo get single isolated margin accounts of the current user
func (ku *Kucoin) GetSingleIsolatedMarginAccountInfo(ctx context.Context, symbol string) (*AssetInfo, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not be empty")
	}
	resp := AssetInfo{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetSingleIsolatedMarginAccountInfo, symbol), nil, &resp)
}

// InitiateIsolateMarginBorrowing initiates isolated margin borrowing
func (ku *Kucoin) InitiateIsolateMarginBorrowing(ctx context.Context, symbol, ccy, borrowStrategy, period string, size, maxRate int64) (*IsolatedMarginBorrowing, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not be empty")
	}
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	params["currency"] = ccy
	if borrowStrategy == "" {
		return nil, errors.New("borrowStrategy can not be empty")
	}
	params["borrowStrategy"] = borrowStrategy
	if size == 0 {
		return nil, errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)

	if period != "" {
		params["period"] = period
	}
	if maxRate == 0 {
		params["maxRate"] = strconv.FormatInt(maxRate, 10)
	}
	resp := IsolatedMarginBorrowing{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginBorrowing, params, &resp)
}

// GetIsolatedOutstandingRepaymentRecords get the outstanding repayment records of isolated margin positions
func (ku *Kucoin) GetIsolatedOutstandingRepaymentRecords(ctx context.Context, symbol, ccy string, pageSize, currentPage int64) ([]OutstandingRepaymentRecord, error) {
	var params url.Values
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	resp := struct {
		CurrentPage int64                        `json:"currentPage"`
		PageSize    int64                        `json:"pageSize"`
		TotalNum    int64                        `json:"totalNum"`
		TotalPage   int64                        `json:"totalPage"`
		Items       []OutstandingRepaymentRecord `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedOutstandingRepaymentRecords, params), nil, &resp)
}

// GetIsolatedMarginRepaymentRecords get the repayment records of isolated margin positions
func (ku *Kucoin) GetIsolatedMarginRepaymentRecords(ctx context.Context, symbol, ccy string, pageSize, currentPage int64) ([]CompletedRepaymentRecord, error) {
	var params url.Values
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	resp := struct {
		CurrentPage int64                      `json:"currentPage"`
		PageSize    int64                      `json:"pageSize"`
		TotalNum    int64                      `json:"totalNum"`
		TotalPage   int64                      `json:"totalPage"`
		Items       []CompletedRepaymentRecord `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetIsolatedMarginRepaymentRecords, params), nil, &resp)
}

// InitiateIsolatedMarginQuickRepayment is used to initiate quick repayment for isolated margin accounts
func (ku *Kucoin) InitiateIsolatedMarginQuickRepayment(ctx context.Context, symbol, ccy, seqStrategy string, size int64) error {
	if symbol == "" {
		return currency.ErrCurrencyPairEmpty
	}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	if ccy == "" {
		return currency.ErrCurrencyCodeEmpty
	}
	params["currency"] = ccy
	if seqStrategy == "" {
		return errors.New("seqStrategy can not be empty")
	}
	params["seqStrategy"] = seqStrategy
	if size == 0 {
		return errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)
	resp := struct {
		Error
	}{}
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginQuickRepayment, params, &resp)
}

// InitiateIsolatedMarginSingleRepayment is used to initiate quick repayment for single margin accounts
func (ku *Kucoin) InitiateIsolatedMarginSingleRepayment(ctx context.Context, symbol, ccy, loanID string, size int64) error {
	if symbol == "" {
		return currency.ErrCurrencyPairEmpty
	}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	if ccy == "" {
		return currency.ErrCurrencyCodeEmpty
	}
	params["currency"] = ccy
	if loanID == "" {
		return errors.New("loanId can not be empty")
	}
	params["loanId"] = loanID
	if size == 0 {
		return errors.New("size can not be zero")
	}
	params["size"] = strconv.FormatInt(size, 10)
	resp := struct {
		Error
	}{}
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInitiateIsolatedMarginSingleRepayment, params, &resp)
}

// GetCurrentServerTime gets the server time
func (ku *Kucoin) GetCurrentServerTime(ctx context.Context) (time.Time, error) {
	resp := struct {
		Timestamp kucoinTimeMilliSec `json:"data"`
		Error
	}{}
	return resp.Timestamp.Time(), ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetServerTime, &resp)
}

// GetServiceStatus gets the service status
func (ku *Kucoin) GetServiceStatus(ctx context.Context) (status, message string, err error) {
	resp := struct {
		Status  string `json:"status"`
		Message string `json:"msg"`
	}{}
	return resp.Status, resp.Message, ku.SendHTTPRequest(ctx, exchange.RestSpot, kucoinGetServiceStatus, &resp)
}

// PostOrder used to place two types of orders: limit and market
// Note: use this only for SPOT trades
func (ku *Kucoin) PostOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, selfTradePrevention, timeInForce string, size, price, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg bool) (string, error) {
	params := make(map[string]interface{})
	if clientOID == "" {
		return "", errors.New("clientOid can not be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return "", errors.New("side can not be empty")
	}
	params["side"] = side
	if symbol == "" {
		return "", fmt.Errorf("%w, empty symbol", currency.ErrCurrencyPairEmpty)
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if selfTradePrevention != "" {
		params["stp"] = selfTradePrevention
	}
	switch orderType {
	case "limit", "":
		if price <= 0 {
			return "", errors.New("price can not be empty")
		}
		params["price"] = price
		if size <= 0 {
			return "", errors.New("size can not be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	case "market":
		switch {
		case size > 0:
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		case funds > 0:
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		default:
			return "", errors.New("atleast one required among size and funds")
		}
	default:
		return "", errors.New("invalid orderType")
	}
	if orderType != "" {
		params["type"] = orderType
	}
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}
	return resp.OrderID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostOrder, params, &resp)
}

// PostMarginOrder used to place two types of margin orders: limit and market
func (ku *Kucoin) PostMarginOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, selfTradePrevention, marginModel, timeInForce string, price, size, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg, autoBorrow bool) (*PostMarginOrderResp, error) {
	params := make(map[string]interface{})
	if clientOID == "" {
		return nil, errors.New("clientOid can not be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return nil, errors.New("side can not be empty")
	}
	params["side"] = side
	if symbol == "" {
		return nil, fmt.Errorf("%w, empty symbol", currency.ErrCurrencyPairEmpty)
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if selfTradePrevention != "" {
		params["stp"] = selfTradePrevention
	}
	if marginModel != "" {
		params["marginMode"] = marginModel
	}
	params["autoBorrow"] = autoBorrow
	switch orderType {
	case "limit", "":
		if price <= 0 {
			return nil, errors.New("price can not be empty")
		}
		params["price"] = price
		if size <= 0 {
			return nil, errors.New("size can not be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	case "market":
		switch {
		case size > 0:
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		case funds > 0:
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		default:
			return nil, errors.New("atleast one required among size and funds")
		}
	default:
		return nil, errors.New("invalid orderType")
	}
	if orderType != "" {
		params["type"] = orderType
	}
	resp := struct {
		PostMarginOrderResp
		Error
	}{}
	return &resp.PostMarginOrderResp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostMarginOrder, params, &resp)
}

// PostBulkOrder used to place 5 orders at the same time. The order type must be a limit order of the same symbol
// Note: it supports only SPOT trades
// Note: To check if order was posted successfully, check status field in response
func (ku *Kucoin) PostBulkOrder(ctx context.Context, symbol string, orderList []OrderRequest) ([]PostBulkOrderResp, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not be empty")
	}
	for i := range orderList {
		if orderList[i].ClientOID == "" {
			return nil, errors.New("clientOid can not be empty")
		}
		if orderList[i].Side == "" {
			return nil, errors.New("side can not be empty")
		}
		if orderList[i].Price <= 0 {
			return nil, errors.New("price must be positive")
		}
		if orderList[i].Size <= 0 {
			return nil, errors.New("size must be positive")
		}
	}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	params["orderList"] = orderList
	resp := []PostBulkOrderResp{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinPostBulkOrder, params, &resp)
}

// CancelSingleOrder used to cancel single order previously placed
func (ku *Kucoin) CancelSingleOrder(ctx context.Context, orderID string) ([]string, error) {
	if orderID == "" {
		return nil, errors.New("orderID can not be empty")
	}
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}
	return resp.CancelledOrderIDs, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinOrderByID, orderID), nil, &resp)
}

// CancelOrderByClientOID used to cancel order via the clientOid
func (ku *Kucoin) CancelOrderByClientOID(ctx context.Context, orderID string) (cancelledOrderID, clientOrderID string, err error) {
	resp := struct {
		CancelledOrderID string `json:"cancelledOrderId"`
		ClientOID        string `json:"clientOid"`
		Error
	}{}
	return resp.CancelledOrderID, resp.ClientOID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinOrderByClientOID, orderID), nil, &resp)
}

// CancelAllOpenOrders used to cancel all order based upon the parameters passed
func (ku *Kucoin) CancelAllOpenOrders(ctx context.Context, symbol, tradeType string) ([]string, error) {
	var params url.Values
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}
	return resp.CancelledOrderIDs, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinOrders, params), nil, &resp)
}

// GetOrders gets the user order list
func (ku *Kucoin) GetOrders(ctx context.Context, status, symbol, side, orderType, tradeType string, startAt, endAt time.Time) ([]OrderDetail, error) {
	var params url.Values
	if status != "" {
		params.Set("status", status)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	resp := struct {
		CurrentPage int64         `json:"currentPage"`
		PageSize    int64         `json:"pageSize"`
		TotalNum    int64         `json:"totalNum"`
		TotalPage   int64         `json:"totalPage"`
		Items       []OrderDetail `json:"items"`
	}{}
	err := ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinOrders, params), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetRecentOrders get orders in the last 24 hours.
func (ku *Kucoin) GetRecentOrders(ctx context.Context) ([]OrderDetail, error) {
	resp := []OrderDetail{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetRecentOrders, nil, &resp)
}

// GetOrderByID get a single order info by order ID
func (ku *Kucoin) GetOrderByID(ctx context.Context, orderID string) (*OrderDetail, error) {
	if orderID == "" {
		return nil, errors.New("orderID can not be empty")
	}
	resp := OrderDetail{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinOrderByID, orderID), nil, &resp)
}

// GetOrderByClientSuppliedOrderID get a single order info by client order ID
func (ku *Kucoin) GetOrderByClientSuppliedOrderID(ctx context.Context, clientOID string) (*OrderDetail, error) {
	if clientOID == "" {
		return nil, errors.New("client order ID can not be empty")
	}
	resp := OrderDetail{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinOrderByClientOID, clientOID), nil, &resp)
}

// GetFills get fills
func (ku *Kucoin) GetFills(ctx context.Context, orderID, symbol, side, orderType, tradeType string, startAt, endAt time.Time) ([]Fill, error) {
	var params url.Values
	if orderID != "" {
		params.Set("orderId", orderID)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	resp := []Fill{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetFills, params), nil, &resp)
}

// GetRecentFills get a list of 1000 fills in last 24 hours
func (ku *Kucoin) GetRecentFills(ctx context.Context) ([]Fill, error) {
	resp := []Fill{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetRecentFills, nil, &resp)
}

// PostStopOrder used to place two types of stop orders: limit and market
func (ku *Kucoin) PostStopOrder(ctx context.Context, clientOID, side, symbol, orderType, remark, stop, price, stopPrice, stp, tradeType, timeInForce string, size, cancelAfter, visibleSize, funds float64, postOnly, hidden, iceberg bool) (string, error) {
	params := make(map[string]interface{})
	if clientOID == "" {
		return "", errors.New("clientOid can not be empty")
	}
	params["clientOid"] = clientOID
	if side == "" {
		return "", errors.New("side can not be empty")
	}
	params["side"] = side
	if symbol == "" {
		return "", fmt.Errorf("%w, empty symbol", currency.ErrCurrencyPairEmpty)
	}
	params["symbol"] = symbol
	if remark != "" {
		params["remark"] = remark
	}
	if stop != "" {
		params["stop"] = stop
		if stopPrice == "" {
			return "", errors.New("stopPrice can not be empty when stop is set")
		}
		params["stopPrice"] = stopPrice
	}
	if stp != "" {
		params["stp"] = stp
	}
	if tradeType != "" {
		params["tradeType"] = tradeType
	}
	switch orderType {
	case "limit", "":
		if price == "" {
			return "", errors.New("price can not be empty")
		}
		params["price"] = price
		if size <= 0 {
			return "", errors.New("size can not be zero or negative")
		}
		params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		if timeInForce != "" {
			params["timeInForce"] = timeInForce
		}
		if cancelAfter > 0 && timeInForce == "GTT" {
			params["cancelAfter"] = strconv.FormatFloat(cancelAfter, 'f', -1, 64)
		}
		params["postOnly"] = postOnly
		params["hidden"] = hidden
		params["iceberg"] = iceberg
		if visibleSize > 0 {
			params["visibleSize"] = strconv.FormatFloat(visibleSize, 'f', -1, 64)
		}
	case "market":
		switch {
		case size > 0:
			params["size"] = strconv.FormatFloat(size, 'f', -1, 64)
		case funds > 0:
			params["funds"] = strconv.FormatFloat(funds, 'f', -1, 64)
		default:
			return "", errors.New("atleast one required among size and funds")
		}
	default:
		return "", errors.New("invalid orderType")
	}
	if orderType != "" {
		params["type"] = orderType
	}
	resp := struct {
		OrderID string `json:"orderId"`
		Error
	}{}
	return resp.OrderID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinStopOrder, params, &resp)
}

// CancelStopOrder used to cancel single stop order previously placed
func (ku *Kucoin) CancelStopOrder(ctx context.Context, orderID string) ([]string, error) {
	if orderID == "" {
		return nil, errors.New("orderID can not be empty")
	}
	resp := struct {
		Data []string `json:"cancelledOrderIds"`
		Error
	}{}
	return resp.Data, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinStopOrderByID, orderID), nil, &resp)
}

// CancelAllStopOrder used to cancel all order based upon the parameters passed
func (ku *Kucoin) CancelAllStopOrder(ctx context.Context, symbol, tradeType, orderIDs string) ([]string, error) {
	var params url.Values
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if orderIDs != "" {
		params.Set("orderIds", orderIDs)
	}
	resp := struct {
		CancelledOrderIDs []string `json:"cancelledOrderIds"`
		Error
	}{}
	return resp.CancelledOrderIDs, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinCancelAllStopOrder, params), nil, &resp)
}

// GetStopOrder used to cancel single stop order previously placed
func (ku *Kucoin) GetStopOrder(ctx context.Context, orderID string) (*StopOrder, error) {
	if orderID == "" {
		return nil, errors.New("orderID can not be empty")
	}
	resp := struct {
		StopOrder
		Error
	}{}
	return &resp.StopOrder, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinStopOrderByID, orderID), nil, &resp)
}

// GetAllStopOrder get all current untriggered stop orders
func (ku *Kucoin) GetAllStopOrder(ctx context.Context, symbol, side, orderType, tradeType, orderIDs string, startAt, endAt time.Time, currentPage, pageSize int64) ([]StopOrder, error) {
	var params url.Values
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if tradeType != "" {
		params.Set("tradeType", tradeType)
	}
	if orderIDs != "" {
		params.Set("orderIds", orderIDs)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.Unix(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.Unix(), 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	resp := struct {
		CurrentPage int64       `json:"currentPage"`
		PageSize    int64       `json:"pageSize"`
		TotalNum    int64       `json:"totalNum"`
		TotalPage   int64       `json:"totalPage"`
		Items       []StopOrder `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinStopOrder, params), nil, &resp)
}

// GetStopOrderByClientID get a stop order information via the clientOID
func (ku *Kucoin) GetStopOrderByClientID(ctx context.Context, symbol, clientOID string) ([]StopOrder, error) {
	if clientOID == "" {
		return nil, errors.New("clientOID can not be empty")
	}
	var params url.Values
	params.Set("clientOid", clientOID)
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	resp := []StopOrder{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetStopOrderByClientID, params), nil, &resp)
}

// CancelStopOrderByClientID used to cancel a stop order via the clientOID.
func (ku *Kucoin) CancelStopOrderByClientID(ctx context.Context, symbol, clientOID string) (cancelledOrderID, clientOrderID string, err error) {
	if clientOID == "" {
		return "", "", errors.New("clientOID can not be empty")
	}
	var params url.Values
	params.Set("clientOid", clientOID)
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	resp := struct {
		CancelledOrderID string `json:"cancelledOrderId"`
		ClientOID        string `json:"clientOid"`
		Error
	}{}
	return resp.CancelledOrderID, resp.ClientOID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, common.EncodeURLValues(kucoinCancelStopOrderByClientID, params), nil, &resp)
}

// CreateAccount creates an account
func (ku *Kucoin) CreateAccount(ctx context.Context, ccy, accountType string) (string, error) {
	if accountType == "" {
		return "", errors.New("accountType can not be empty")
	}
	params := make(map[string]interface{})
	params["type"] = accountType
	if ccy == "" {
		return "", currency.ErrCurrencyPairEmpty
	}
	params["currency"] = ccy
	resp := struct {
		ID string `json:"id"`
		Error
	}{}
	return resp.ID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinAccount, params, &resp)
}

// GetAllAccounts get all accounts
func (ku *Kucoin) GetAllAccounts(ctx context.Context, ccy, accountType string) ([]AccountInfo, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if accountType != "" {
		params.Set("type", accountType)
	}
	resp := []AccountInfo{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinAccount, params), nil, &resp)
}

// GetAccount get information of single account
func (ku *Kucoin) GetAccount(ctx context.Context, accountID string) (*AccountInfo, error) {
	resp := AccountInfo{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetAccount, accountID), nil, &resp)
}

// GetAccountLedgers get the history of deposit/withdrawal of all accounts, supporting inquiry of various currencies
func (ku *Kucoin) GetAccountLedgers(ctx context.Context, ccy, direction, bizType string, startAt, endAt time.Time) ([]LedgerInfo, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if direction != "" {
		params.Set("direction", direction)
	}
	if bizType != "" {
		params.Set("bizType", bizType)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	resp := struct {
		CurrentPage int64        `json:"currentPage"`
		PageSize    int64        `json:"pageSize"`
		TotalNum    int64        `json:"totalNum"`
		TotalPage   int64        `json:"totalPage"`
		Items       []LedgerInfo `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetAccountLedgers, params), nil, &resp)
}

// GetSubAccountBalance get account info of a sub-user specified by the subUserID
func (ku *Kucoin) GetSubAccountBalance(ctx context.Context, subUserID string) (*SubAccountInfo, error) {
	resp := SubAccountInfo{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, fmt.Sprintf(kucoinGetSubAccountBalance, subUserID), nil, &resp)
}

// GetAggregatedSubAccountBalance get the account info of all sub-users
func (ku *Kucoin) GetAggregatedSubAccountBalance(ctx context.Context) ([]SubAccountInfo, error) {
	resp := []SubAccountInfo{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, kucoinGetAggregatedSubAccountBalance, nil, &resp)
}

// GetTransferableBalance get the transferable balance of a specified account
func (ku *Kucoin) GetTransferableBalance(ctx context.Context, ccy, accountType, tag string) (*TransferableBalanceInfo, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	if accountType == "" {
		return nil, errors.New("accountType can not be empty")
	}
	params.Set("type", accountType)
	if tag != "" {
		params.Set("tag", tag)
	}
	resp := TransferableBalanceInfo{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetTransferableBalance, params), nil, &resp)
}

// TransferMainToSubAccount used to transfer funds from main account to sub-account
func (ku *Kucoin) TransferMainToSubAccount(ctx context.Context, clientOID, ccy, amount, direction, accountType, subAccountType, subUserID string) (string, error) {
	if clientOID == "" {
		return "", errors.New("clientOID can not be empty")
	}
	params := make(map[string]interface{})
	params["clientOid"] = clientOID
	if ccy == "" {
		return "", currency.ErrCurrencyPairEmpty
	}
	params["currency"] = ccy
	if amount == "" {
		return "", errors.New("amount can not be empty")
	}
	params["amount"] = amount
	if direction == "" {
		return "", errors.New("direction can not be empty")
	}
	params["direction"] = direction
	if accountType != "" {
		params["accountType"] = accountType
	}
	if subAccountType != "" {
		params["subAccountType"] = subAccountType
	}
	if subUserID == "" {
		return "", errors.New("subUserID can not be empty")
	}
	params["subUserId"] = subUserID
	resp := struct {
		OrderID string `json:"orderId"`
	}{}
	return resp.OrderID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinTransferMainToSubAccount, params, &resp)
}

// MakeInnerTransfer used to transfer funds between accounts internally
func (ku *Kucoin) MakeInnerTransfer(ctx context.Context, clientOID, ccy, from, to, amount, fromTag, toTag string) (string, error) {
	if clientOID == "" {
		return "", errors.New("clientOID can not be empty")
	}
	params := make(map[string]interface{})
	params["clientOid"] = clientOID
	if ccy == "" {
		return "", currency.ErrCurrencyPairEmpty
	}
	params["currency"] = ccy
	if amount == "" {
		return "", errors.New("amount can not be empty")
	}
	params["amount"] = amount
	if from == "" {
		return "", errors.New("from can not be empty")
	}
	params["from"] = from
	if to == "" {
		return "", errors.New("to can not be empty")
	}
	params["to"] = to
	if fromTag != "" {
		params["fromTag"] = fromTag
	}
	if toTag != "" {
		params["toTag"] = toTag
	}
	resp := struct {
		OrderID string `json:"orderId"`
	}{}
	return resp.OrderID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinInnerTransfer, params, &resp)
}

// CreateDepositAddress create a deposit address for a currency you intend to deposit
func (ku *Kucoin) CreateDepositAddress(ctx context.Context, ccy, chain string) (*DepositAddress, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if chain != "" {
		params["chain"] = chain
	}
	resp := DepositAddress{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinCreateDepositAddress, params, &resp)
}

// GetDepositAddressV2 get all deposit addresses for the currency you intend to deposit
func (ku *Kucoin) GetDepositAddressV2(ctx context.Context, ccy string) ([]DepositAddress, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	resp := []DepositAddress{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositAddressV2, params), nil, &resp)
}

// GetDepositAddressV1 get a deposit address for the currency you intend to deposit
func (ku *Kucoin) GetDepositAddressV1(ctx context.Context, ccy, chain string) ([]DepositAddress, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	if chain != "" {
		params.Set("chain", chain)
	}
	resp := []DepositAddress{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositAddressV1, params), nil, &resp)
}

// GetDepositList get deposit list items and sorted to show the latest first
func (ku *Kucoin) GetDepositList(ctx context.Context, ccy, status string, startAt, endAt time.Time) ([]Deposit, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	resp := struct {
		CurrentPage int64     `json:"currentPage"`
		PageSize    int64     `json:"pageSize"`
		TotalNum    int64     `json:"totalNum"`
		TotalPage   int64     `json:"totalPage"`
		Items       []Deposit `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetDepositList, params), nil, &resp)
}

// GetHistoricalDepositList get historical deposit list items
func (ku *Kucoin) GetHistoricalDepositList(ctx context.Context, ccy, status string, startAt, endAt time.Time) ([]HistoricalDepositWithdrawal, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	resp := struct {
		CurrentPage int64                         `json:"currentPage"`
		PageSize    int64                         `json:"pageSize"`
		TotalNum    int64                         `json:"totalNum"`
		TotalPage   int64                         `json:"totalPage"`
		Items       []HistoricalDepositWithdrawal `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetHistoricalDepositList, params), nil, &resp)
}

// GetWithdrawalList get withdrawal list items
func (ku *Kucoin) GetWithdrawalList(ctx context.Context, ccy, status string, startAt, endAt time.Time) ([]Withdrawal, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	resp := struct {
		CurrentPage int64        `json:"currentPage"`
		PageSize    int64        `json:"pageSize"`
		TotalNum    int64        `json:"totalNum"`
		TotalPage   int64        `json:"totalPage"`
		Items       []Withdrawal `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinWithdrawal, params), nil, &resp)
}

// GetHistoricalWithdrawalList get historical withdrawal list items
func (ku *Kucoin) GetHistoricalWithdrawalList(ctx context.Context, ccy, status string, startAt, endAt time.Time, currentPage, pageSize int64) ([]HistoricalDepositWithdrawal, error) {
	var params url.Values
	if ccy != "" {
		params.Set("currency", ccy)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startAt.IsZero() {
		params.Set("startAt", strconv.FormatInt(startAt.UnixMilli(), 10))
	}
	if !endAt.IsZero() {
		params.Set("endAt", strconv.FormatInt(endAt.UnixMilli(), 10))
	}
	if currentPage != 0 {
		params.Set("currentPage", strconv.FormatInt(currentPage, 10))
	}
	if pageSize != 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	resp := struct {
		CurrentPage int64                         `json:"currentPage"`
		PageSize    int64                         `json:"pageSize"`
		TotalNum    int64                         `json:"totalNum"`
		TotalPage   int64                         `json:"totalPage"`
		Items       []HistoricalDepositWithdrawal `json:"items"`
	}{}
	return resp.Items, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetHistoricalWithdrawalList, params), nil, &resp)
}

// GetWithdrawalQuotas get withdrawal quota details
func (ku *Kucoin) GetWithdrawalQuotas(ctx context.Context, ccy, chain string) (*WithdrawalQuota, error) {
	if ccy == "" {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	var params url.Values
	params.Set("currency", ccy)
	if chain != "" {
		params.Set("chain", chain)
	}
	resp := WithdrawalQuota{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinGetWithdrawalQuotas, params), nil, &resp)
}

// ApplyWithdrawal create a withdrawal request
func (ku *Kucoin) ApplyWithdrawal(ctx context.Context, ccy, address, memo, remark, chain, feeDeductType string, isInner bool, amount float64) (string, error) {
	if ccy == "" {
		return "", currency.ErrCurrencyPairEmpty
	}
	params := make(map[string]interface{})
	params["currency"] = ccy
	if address == "" {
		return "", errors.New("address can not be empty")
	}
	params["address"] = address
	if amount == 0 {
		return "", errors.New("amount can not be empty")
	}
	params["amount"] = amount
	if memo != "" {
		params["memo"] = memo
	}
	params["isInner"] = isInner
	if remark != "" {
		params["remark"] = remark
	}
	if chain != "" {
		params["chain"] = chain
	}
	if feeDeductType != "" {
		params["feeDeductType"] = feeDeductType
	}
	resp := struct {
		WithdrawalID string `json:"withdrawalId"`
		Error
	}{}
	return resp.WithdrawalID, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, kucoinWithdrawal, params, &resp)
}

// CancelWithdrawal used to cancel a withdrawal request
func (ku *Kucoin) CancelWithdrawal(ctx context.Context, withdrawalID string) error {
	return ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodDelete, fmt.Sprintf(kucoinCancelWithdrawal, withdrawalID), nil, &struct{}{})
}

// GetBasicFee get basic fee rate of users
func (ku *Kucoin) GetBasicFee(ctx context.Context, currencyType string) (*Fees, error) {
	var params url.Values
	if currencyType != "" {
		params.Set("currencyType", currencyType)
	}
	resp := Fees{}
	return &resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinBasicFee, params), nil, &resp)
}

// GetTradingFee get fee rate of trading pairs
func (ku *Kucoin) GetTradingFee(ctx context.Context, symbols string) ([]Fees, error) {
	var params url.Values
	if symbols != "" {
		params.Set("symbols", symbols)
	}
	resp := []Fees{}
	return resp, ku.SendAuthHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, common.EncodeURLValues(kucoinTradingFee, params), nil, &resp)
}

// SendHTTPRequest sends an unauthenticated HTTP request
func (ku *Kucoin) SendHTTPRequest(ctx context.Context, ePath exchange.URL, path string, result interface{}) error {
	value := reflect.ValueOf(result)
	if value.Kind() != reflect.Pointer {
		return errInvalidResultInterface
	}
	var resp UnmarshalTo
	_, okay := result.(UnmarshalTo)
	if okay {
		resp = result.(UnmarshalTo)
	} else {
		resp = &Response{Data: result}
	}
	endpointPath, err := ku.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	err = ku.SendPayload(ctx, request.Unset, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          endpointPath + path,
			Result:        resp,
			Verbose:       ku.Verbose,
			HTTPDebugging: ku.HTTPDebugging,
			HTTPRecording: ku.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return resp.GetError()
}

// SendAuthHTTPRequest sends an authenticated HTTP request
// Request parameters are added to path variable for GET and DELETE request and for other requests its passed in params variable
func (ku *Kucoin) SendAuthHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, params map[string]interface{}, result interface{}) error {
	value := reflect.ValueOf(result)
	if value.Kind() != reflect.Pointer {
		return errInvalidResultInterface
	}
	creds, err := ku.GetCredentials(ctx)
	if err != nil {
		return err
	}
	var resp UnmarshalTo
	_, okay := result.(UnmarshalTo)
	if okay {
		resp = result.(UnmarshalTo)
	} else {
		resp = &Response{Data: result}
	}
	endpointPath, err := ku.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(result)
	if val.IsNil() || val.Kind() != reflect.Pointer {
		return fmt.Errorf("%w receiver has to be non-nil pointer", errInvalidResponseReciever)
	}
	err = ku.SendPayload(ctx, request.Unset, func() (*request.Item, error) {
		var (
			body    io.Reader
			payload []byte
		)
		if len(params) != 0 {
			payload, err = json.Marshal(params)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(payload)
		}
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		var signHash, passPhraseHash []byte
		signHash, err = crypto.GetHMAC(crypto.HashSHA256, []byte(timeStamp+method+"/api"+path+string(payload)), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}
		passPhraseHash, err = crypto.GetHMAC(crypto.HashSHA256, []byte(creds.PEMKey), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}
		headers := map[string]string{
			"KC-API-KEY":         creds.Key,
			"KC-API-SIGN":        crypto.Base64Encode(signHash),
			"KC-API-TIMESTAMP":   timeStamp,
			"KC-API-PASSPHRASE":  crypto.Base64Encode(passPhraseHash),
			"KC-API-KEY-VERSION": kucoinAPIKeyVersion,
			"Content-Type":       "application/json",
		}
		return &request.Item{
			Method:        method,
			Path:          endpointPath + path,
			Headers:       headers,
			Body:          body,
			Result:        &resp,
			AuthRequest:   true,
			Verbose:       ku.Verbose,
			HTTPDebugging: ku.HTTPDebugging,
			HTTPRecording: ku.HTTPRecording}, nil
	})
	if err != nil {
		return err
	}
	return resp.GetError()
}

func (ku *Kucoin) intervalToString(interval kline.Interval) (string, error) {
	switch interval {
	case kline.OneMin:
		return "1min", nil
	case kline.ThreeMin:
		return "3min", nil
	case kline.FiveMin:
		return "5min", nil
	case kline.FifteenMin:
		return "15min", nil
	case kline.ThirtyMin:
		return "30min", nil
	case kline.OneHour:
		return "1hour", nil
	case kline.FourHour:
		return "4hour", nil
	case kline.SixHour:
		return "6hour", nil
	case kline.EightHour:
		return "8hour", nil
	case kline.TwelveHour:
		return "12hour", nil
	case kline.OneDay:
		return "1day", nil
	case kline.OneWeek:
		return "1week", nil
	default:
		return "", kline.ErrUnsupportedInterval
	}
}

func (ku *Kucoin) stringToInterval(interval string) (kline.Interval, error) {
	switch interval {
	case "1min":
		return kline.OneMin, nil
	case "3min":
		return kline.ThreeMin, nil
	case "5min":
		return kline.FiveMin, nil
	case "15min":
		return kline.FifteenMin, nil
	case "30min":
		return kline.ThirtyMin, nil
	case "1hour":
		return kline.OneHour, nil
	case "4hour":
		return kline.FourHour, nil
	case "6hour":
		return kline.SixHour, nil
	case "8hour":
		return kline.EightHour, nil
	case "12hour":
		return kline.TwelveHour, nil
	case "1day":
		return kline.OneDay, nil
	case "1week":
		return kline.OneWeek, nil
	default:
		return 0, kline.ErrUnsupportedInterval
	}
}

func (ku *Kucoin) stringToOrderStatus(status string) (order.Status, error) {
	switch status {
	case "match":
		return order.Filled, nil
	case "open":
		return order.Open, nil
	case "done":
		return order.Closed, nil
	default:
		return order.StringToOrderStatus(status)
	}
}

func (ku *Kucoin) accountTypeToString(a asset.Item) string {
	switch a {
	case asset.Spot:
		return "trade"
	case asset.Margin:
		return "margin"
	case asset.Empty:
		return ""
	default:
		return "main"
	}
}

func (ku *Kucoin) accountToTradeTypeString(a asset.Item, marginMode string) string {
	switch a {
	case asset.Spot:
		return "TRADE"
	case asset.Margin:
		if strings.EqualFold(marginMode, "isolated") {
			return "MARGIN_ISOLATED_TRADE"
		}
		return "MARGIN_TRADE"
	default:
		return ""
	}
}