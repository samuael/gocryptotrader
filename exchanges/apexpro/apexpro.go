package apexpro

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/starkex"
	"github.com/thrasher-corp/gocryptotrader/types"
)

// Apexpro is the overarching type across this package
type Apexpro struct {
	exchange.Base

	// SymbolsConfig represents all symbols configuration.
	SymbolsConfig *AllSymbolsV1Config

	StarkConfig       *starkex.StarkConfig
	UserAccountDetail *UserAccountDetail
}

const (
	apexproAPIURL     = "https://pro.apex.exchange/api/"
	apexproTestAPIURL = "https://testnet.pro.apex.exchange/api/"

	apexProOmniAPIURL = "https://omni.apex.exchange/api/"

	// Public endpoints

	// Authenticated endpoints
)

var (
	errL2KeyMissing                  = errors.New("l2 Key is required")
	errEthereumAddressMissing        = errors.New("ethereum address is missing")
	errInvalidEthereumAddress        = errors.New("invalid ethereum address")
	errChainIDMissing                = errors.New("chain ID is missing")
	errOrderbookLevelIsRequired      = errors.New("orderbook level is required")
	errInvalidTimestamp              = errors.New("err invalid timestamp")
	errZeroKnowledgeAccountIDMissing = errors.New("zero knowledge account id is required")
	errSubAccountIDMissing           = errors.New("missing sub-account id")
	errUserNonceRequired             = errors.New("nonce is required")
	errInitialMarginRateRequired     = errors.New("initial margin rate required")
	errUserIDRequired                = errors.New("user ID is required")
	errDeviceTypeIsRequired          = errors.New("device type is required")
	errLimitFeeRequired              = errors.New("limit fee is required")
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
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v3/symbols", request.UnAuth, &resp)
}

// Apexpro retrieves all symbols and asset configurations from the V1 API.
func (ap *Apexpro) GetAllSymbolsConfigDataV1(ctx context.Context) (*AllSymbolsV1Config, error) {
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
	return resp, ap.SendHTTPRequest(ctx, exchange.RestSpot, "v2/symbols", request.UnAuth, &resp, true)
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

// GenerateNonceV3 generate and obtain nonce before registration. The nonce is used to assemble the signature field upon registration.
func (ap *Apexpro) GenerateNonceV3(ctx context.Context, l2Key, ethereumAddress, chainID string) (*NonceResponse, error) {
	return ap.generateNonce(ctx, l2Key, ethereumAddress, chainID, "v3/generate-nonce")
}

// GenerateNonceV2 before registering, generate and obtain a nonce. The nonce serves the purpose of assembling the signature field during the registration process.
func (ap *Apexpro) GenerateNonceV2(ctx context.Context, l2Key, ethereumAddress, chainID string) (*NonceResponse, error) {
	return ap.generateNonce(ctx, l2Key, ethereumAddress, chainID, "v2/generate-nonce")
}

// GenerateNonceV1 before registering, generate and obtain a nonce.
func (ap *Apexpro) GenerateNonceV1(ctx context.Context, l2Key, ethereumAddress, chainID string) (*NonceResponse, error) {
	return ap.generateNonce(ctx, l2Key, ethereumAddress, chainID, "v1/generate-nonce")
}

func (ap *Apexpro) generateNonce(ctx context.Context, l2Key, ethereumAddress, chainID, path string) (*NonceResponse, error) {
	if l2Key == "" {
		return nil, errL2KeyMissing
	}
	if ethereumAddress == "" {
		return nil, errEthereumAddressMissing
	}
	if chainID == "" {
		return nil, errChainIDMissing
	}
	params := url.Values{}
	params.Set("l2Key", l2Key)
	params.Set("ethAddress", ethereumAddress)
	params.Set("chainId", chainID)
	var resp *NonceResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, &resp, false)
}

// GetUsersDataV3 retrieves an account users information.
func (ap *Apexpro) GetUsersDataV3(ctx context.Context) (*UserData, error) {
	return ap.getUsersData(ctx, "v3/user")
}

// GetUsersDataV2 retrieves an account users information through the V2 API
func (ap *Apexpro) GetUsersDataV2(ctx context.Context) (*UserData, error) {
	return ap.getUsersData(ctx, "v2/user")
}

// GetUsersDataV1 retrieves an account users information through the V1 API
func (ap *Apexpro) GetUsersDataV1(ctx context.Context) (*UserData, error) {
	return ap.getUsersData(ctx, "v1/user")
}

func (ap *Apexpro) getUsersData(ctx context.Context, path string) (*UserData, error) {
	var resp *UserData
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.Unset, nil, &resp)
}

// EditUserDataV3 edits user's data.
func (ap *Apexpro) EditUserDataV3(ctx context.Context, arg *EditUserDataParams) (*UserDataResponse, error) {
	return ap.editUserData(ctx, arg, "v3/modify-user")
}

// EditUserDataV2 edits user's data through the V2 API.
func (ap *Apexpro) EditUserDataV2(ctx context.Context, arg *EditUserDataParams) (*UserDataResponse, error) {
	return ap.editUserData(ctx, arg, "v2/modify-user")
}

// EditUserDataV1 edits user's data through the V1 API.
func (ap *Apexpro) EditUserDataV1(ctx context.Context, arg *EditUserDataParams) (*UserDataResponse, error) {
	return ap.editUserData(ctx, arg, "v1/modify-user")
}

func (ap *Apexpro) editUserData(ctx context.Context, arg *EditUserDataParams, path string) (*UserDataResponse, error) {
	if arg == nil || *arg == (EditUserDataParams{}) {
		return nil, common.ErrNilPointer
	}
	params := url.Values{}
	if arg.UserData != "" {
		params.Set("userData", arg.UserData)
	}
	if arg.Email != "" {
		params.Set("email", arg.Email)
	}
	if arg.Username != "" {
		params.Set("username", arg.Username)
	}
	if arg.Country != "" {
		params.Set("country", arg.Country)
	}
	if arg.IsSharingUsername {
		params.Set("isSharingUsername", "true")
	}
	if arg.IsSharingAddress {
		params.Set("isSharingAddress", "true")
	}
	if arg.EmailNotifyGeneralEnable {
		params.Set("emailNotifyGeneralEnable", "true")
	}
	if arg.EmailNotifyTradingEnable {
		params.Set("emailNotifyTradingEnable", "true")
	}
	if arg.EmailNotifyAccountEnable {
		params.Set("emailNotifyAccountEnable", "true")
	}
	if arg.PopupNotifyTradingEnable {
		params.Set("popupNotifyTradingEnable", "true")
	}
	var resp *UserDataResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, &resp)
}

// GetUserAccountDataV3 get an account for a user by id. Using the client, the id will be generated with client information and an Ethereum address.
func (ap *Apexpro) GetUserAccountDataV3(ctx context.Context) (*UserAccountDetail, error) {
	return ap.getUserAccountData(ctx, "v3/account")
}

// GetUserAccountDataV2 get a user account detail throught the V2 API.
func (ap *Apexpro) GetUserAccountDataV2(ctx context.Context) (*UserAccountDetail, error) {
	return ap.getUserAccountData(ctx, "v2/account")
}

func (ap *Apexpro) getUserAccountData(ctx context.Context, path string) (*UserAccountDetail, error) {
	var resp *UserAccountDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, nil, &resp)
}

// GetUserAccountDataV1 get an account for a user by id
func (ap *Apexpro) GetUserAccountDataV1(ctx context.Context) (*UserAccountDetailV1, error) {
	var resp *UserAccountDetailV1
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v1/account", request.UnAuth, nil, &resp)
}

// GetUserAccountBalance retrieves user account balance information.
func (ap *Apexpro) GetUserAccountBalance(ctx context.Context) (*UserAccountBalanceResponse, error) {
	var resp *UserAccountBalanceResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/account-balance", request.UnAuth, nil, &resp)
}

// GetUserAccountBalance retrieves user account balance information through the V2 API.
func (ap *Apexpro) GetUserAccountBalanceV2(ctx context.Context) (*UserAccountBalanceV2Response, error) {
	var resp *UserAccountBalanceV2Response
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v2/account-balance", request.UnAuth, nil, &resp)
}

// GetUserAccountBalanceV1 retrive user account baance
func (ap *Apexpro) GetUserAccountBalanceV1(ctx context.Context) (*UserAccountBalanceResponse, error) {
	var resp *UserAccountBalanceResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v1/account-balance", request.UnAuth, nil, &resp)
}

// UserWithdrawalsV2

// GetUserTransferDataV2 retrieves user's asset transfer information.
func (ap *Apexpro) GetUserTransferDataV2(ctx context.Context, ccy currency.Code, startTime, endTime time.Time, transferType string, chainIDs []string, limit, page int64) (*UserWithdrawalsV2, error) {
	return ap.getUserTransferData(ctx, ccy, startTime, endTime, transferType, "v2/transfers", chainIDs, limit, page)
}

// GetUserTransferDataV1 retrieves user's deposit data.
func (ap *Apexpro) GetUserTransferDataV1(ctx context.Context, ccy currency.Code, startTime, endTime time.Time, transferType string, chainIDs []string, limit, page int64) (*UserWithdrawalsV2, error) {
	return ap.getUserTransferData(ctx, ccy, startTime, endTime, transferType, "v1/transfers", chainIDs, limit, page)
}

func (ap *Apexpro) getUserTransferData(ctx context.Context, ccy currency.Code, startTime, endTime time.Time, transferType, path string, chainIDs []string, limit, page int64) (*UserWithdrawalsV2, error) {
	if ccy.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := url.Values{}
	params.Set("currencyId", ccy.String())
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if len(chainIDs) > 0 {
		params.Set("chainIds", strings.Join(chainIDs, ","))
	}
	if transferType != "" {
		params.Set("transferType", transferType)
	}
	var resp *UserWithdrawalsV2
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetUserWithdrawalListV2 retrieves asset withdrawal list.
func (ap *Apexpro) GetUserWithdrawalListV2(ctx context.Context, transferType string, startTime, endTime time.Time, page, limit int64) (*WithdrawalsV2, error) {
	return ap.getUserWithdrawalList(ctx, transferType, "v2/withdraw-list", startTime, endTime, page, limit)
}

// GetUserWithdrawalListV1 returns the user withdrawal list.
func (ap *Apexpro) GetUserWithdrawalListV1(ctx context.Context, transferType string, startTime, endTime time.Time, page, limit int64) (*WithdrawalsV2, error) {
	return ap.getUserWithdrawalList(ctx, transferType, "v1/withdraw-list", startTime, endTime, page, limit)
}

func (ap *Apexpro) getUserWithdrawalList(ctx context.Context, transferType, path string, startTime, endTime time.Time, page, limit int64) (*WithdrawalsV2, error) {
	params := url.Values{}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if transferType != "" {
		params.Set("transferType", transferType)
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	var resp *WithdrawalsV2
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetFastAndCrossChainWithdrawalFees retrieves fee information of fast and cross-chain withdrawal transactions.
func (ap *Apexpro) GetFastAndCrossChainWithdrawalFeesV2(ctx context.Context, amount float64, chainID string, token currency.Code) (*FastAndCrossChainWithdrawalFees, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	if chainID == "" {
		return nil, errChainIDMissing
	}
	return ap.getFastAndCrossChainWithdrawalFees(ctx, amount, chainID, "v2/uncommon-withdraw-fee", token.String())
}

// GetFastAndCrossChainWithdrawalFeesV1 retrieves fee information of fast and cross-chain withdrawals.
func (ap *Apexpro) GetFastAndCrossChainWithdrawalFeesV1(ctx context.Context, amount float64, chainID string) (*FastAndCrossChainWithdrawalFees, error) {
	return ap.getFastAndCrossChainWithdrawalFees(ctx, amount, chainID, "v1/uncommon-withdraw-fee", "")
}

func (ap *Apexpro) getFastAndCrossChainWithdrawalFees(ctx context.Context, amount float64, chainID, path, token string) (*FastAndCrossChainWithdrawalFees, error) {
	params := url.Values{}
	params.Set("token", token)
	if amount > 0 {
		params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	}
	if chainID != "" {
		params.Set("chainId", chainID)
	}
	var resp *FastAndCrossChainWithdrawalFees
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetAssetWithdrawalAndTransferLimitV2 retrieves an asset withdrawal and transfer limit per interval.
func (ap *Apexpro) GetAssetWithdrawalAndTransferLimitV2(ctx context.Context, currencyID currency.Code) (*TransferAndWithdrawalLimit, error) {
	return ap.getAssetWithdrawalAndTransferLimit(ctx, currencyID, "v2/transfer-limit")
}

// GetAssetWithdrawalAndTransferLimitV1 retrieves an asset withdrawal and transfer limit per interval.
func (ap *Apexpro) GetAssetWithdrawalAndTransferLimitV1(ctx context.Context, currencyID currency.Code) (*TransferAndWithdrawalLimit, error) {
	return ap.getAssetWithdrawalAndTransferLimit(ctx, currencyID, "v1/transfer-limit")
}

func (ap *Apexpro) getAssetWithdrawalAndTransferLimit(ctx context.Context, currencyID currency.Code, path string) (*TransferAndWithdrawalLimit, error) {
	if currencyID.IsEmpty() {
		return nil, fmt.Errorf("%w, currencyID is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("currencyId", currencyID.String())
	var resp *TransferAndWithdrawalLimit
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetUserTransferData retrieves user's asset transfer data.
// Direction: possible values are 'NEXT' and 'PREVIOUS'
// TransfersType: possible values are 'DEPOSIT', 'WITHDRAW' ,'FAST_WITHDRAW' ,'OMNI_TO_PERP' for spot account -> contract account,'OMNI_FROM_PERP' for spot account <- contract account,'AFFILIATE_REBATE' affliate rebate,'REFERRAL_REBATE' for referral rebate,'BROKER_REBATE' for broker rebate
func (ap *Apexpro) GetUserTransferData(ctx context.Context, id, limit int64, tokenID, transferType, subAccountID, direction string, startAt, endAt time.Time, chainIDs []string) (*UserWithdrawals, error) {
	if startAt.IsZero() {
		return nil, fmt.Errorf("%w, startTime is required", errInvalidTimestamp)
	}
	if endAt.IsZero() {
		return nil, fmt.Errorf("%w, endTime is required", errInvalidTimestamp)
	}
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if id != 0 {
		params.Set("id", strconv.FormatInt(id, 10))
	}
	if transferType != "" {
		params.Set("transferType", transferType)
	}
	if tokenID != "" {
		params.Set("tokenId", tokenID)
	}
	if subAccountID != "" {
		params.Set("subAccountId", subAccountID)
	}
	if direction != "" {
		params.Set("direction", direction)
	}
	if len(chainIDs) > 0 {
		params.Add("chainIds", "1")
	}
	params.Set("endTimeExclusive", strconv.FormatInt(endAt.UnixMilli(), 10))
	params.Set("beginTimeInclusive", strconv.FormatInt(startAt.UnixMilli(), 10))
	var resp *UserWithdrawals
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v3/transfers", request.UnAuth, params, &resp)
}

// GetWithdrawalFees retrieves list of withdrawal fees.
// the withdrawal need zkAvailableAmount >= withdrawAmount
// the fast withdrawal needzkAvailableAmount >= withdrawAmount && fastPoolAvailableAmount>= withdrawAmount
func (ap *Apexpro) GetWithdrawalFees(ctx context.Context, amount float64, chainIDs []string, tokenID int64) (*WithdrawalFeeInfos, error) {
	params := url.Values{}
	if amount != 0 {
		params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	}
	if len(chainIDs) > 0 {
		for a := range chainIDs {
			params.Set("chainId", chainIDs[a])
		}
	}
	if tokenID != 0 {
		params.Set("tokenId", strconv.FormatInt(tokenID, 10))
	}
	var resp *WithdrawalFeeInfos
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v3/withdraw-fee", request.UnAuth, params, &resp)
}

// GetContractAccountTransferLimits retrieves a transfer limit of a contract.
func (ap *Apexpro) GetContractAccountTransferLimits(ctx context.Context, ccy currency.Code) (*ContractTransferLimit, error) {
	if ccy.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := url.Values{}
	params.Set("token", ccy.String())
	var resp *ContractTransferLimit
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v3/contract-transfer-limit", request.UnAuth, params, &resp)
}

// GetTradeHistory retrieves trade fills history
func (ap *Apexpro) GetTradeHistory(ctx context.Context, symbol, side, orderType string, startTime, endTime time.Time, page, limit int64) (*TradeHistory, error) {
	return ap.getTradeHistory(ctx, symbol, side, orderType, "", "v3/fills", startTime, endTime, page, limit, exchange.RestFutures)
}

// GetTradeHistoryV2 retrieves trade fills history through the v2 API
func (ap *Apexpro) GetTradeHistoryV2(ctx context.Context, symbol, side, orderType string, token currency.Code, startTime, endTime time.Time, page, limit int64) (*TradeHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getTradeHistory(ctx, symbol, side, orderType, token.String(), "v2/fills", startTime, endTime, page, limit, exchange.RestSpot)
}

// GetTradeHistoryV1 retrieves trade fills history through the v1 API
func (ap *Apexpro) GetTradeHistoryV1(ctx context.Context, symbol, side, orderType string, startTime, endTime time.Time, page, limit int64) (*TradeHistory, error) {
	return ap.getTradeHistory(ctx, symbol, side, orderType, "", "v1/fills", startTime, endTime, page, limit, exchange.RestSpot)
}

func (ap *Apexpro) getTradeHistory(ctx context.Context, symbol, side, orderType, token, path string, startTime, endTime time.Time, page, limit int64, ePath exchange.URL) (*TradeHistory, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if token != "" {
		params.Set("token", token)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("orderType", orderType)
	}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp *TradeHistory
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetWorstPriceV3 retrieves the worst market price from orderbook
func (ap *Apexpro) GetWorstPriceV3(ctx context.Context, symbol, side string, amount float64) (*SymbolWorstPrice, error) {
	return ap.getWorstPrice(ctx, symbol, side, "v3/get-worst-price", amount, exchange.RestFutures)
}

// GetWorstPriceV2 retrieves the worst market price from orderbook
func (ap *Apexpro) GetWorstPriceV2(ctx context.Context, symbol, side string, amount float64) (*SymbolWorstPrice, error) {
	return ap.getWorstPrice(ctx, symbol, side, "v2/get-worst-price", amount, exchange.RestSpot)
}

// GetWorstPriceV1 retrieves the worst market price from orderbook
func (ap *Apexpro) GetWorstPriceV1(ctx context.Context, symbol, side string, amount float64) (*SymbolWorstPrice, error) {
	return ap.getWorstPrice(ctx, symbol, side, "v1/get-worst-price", amount, exchange.RestSpot)
}

func (ap *Apexpro) orderCreationParamsFilter(ctx context.Context, arg *CreateOrderParams) (url.Values, error) {
	if arg == nil || *arg == (CreateOrderParams{}) {
		return nil, order.ErrOrderDetailIsNil
	}
	if arg.Symbol.IsEmpty() {
		return nil, currency.ErrSymbolStringEmpty
	}
	if arg.Side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if arg.OrderType == "" {
		return nil, order.ErrTypeIsInvalid
	}
	if arg.Size <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	if arg.Price <= 0 {
		return nil, order.ErrPriceBelowMin
	}
	if arg.LimitFee < 0 {
		return nil, errLimitFeeRequired
	}
	if arg.ExpirationTime.IsZero() {
		return nil, errExpirationTimeRequired
	}
	signature, err := ap.ProcessOrderSignature(ctx, arg)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("symbol", arg.Symbol.String())
	params.Set("side", arg.Side)
	params.Set("type", arg.OrderType)
	params.Set("size", strconv.FormatFloat(arg.Size, 'f', -1, 64))
	params.Set("price", strconv.FormatFloat(arg.Price, 'f', -1, 64))
	params.Set("limitFee", strconv.FormatFloat(arg.LimitFee, 'f', -1, 64))
	params.Set("expiration", strconv.FormatInt(arg.ExpirationTime.UnixMilli(), 10))
	params.Set("signature", signature)
	if arg.TimeInForce != "" {
		params.Set("timeInForce", arg.TimeInForce)
	}
	if arg.TriggerPrice > 0 {
		params.Set("triggerPrice", strconv.FormatFloat(arg.TriggerPrice, 'f', -1, 64))
	}
	if arg.TrailingPercent > 0 {
		params.Set("trailingPercent", strconv.FormatFloat(arg.TrailingPercent, 'f', -1, 64))
	}
	if arg.ClientOrderID != 0 {
		params.Set("clientOrderId", strconv.FormatInt(arg.ClientOrderID, 10))
	}
	return params, nil
}

// CreateOrderV3 creates a new order.
func (ap *Apexpro) CreateOrderV3(ctx context.Context, arg *CreateOrderParams) (*OrderDetail, error) {
	params, err := ap.orderCreationParamsFilter(ctx, arg)
	if err != nil {
		return nil, err
	}
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "v3/order", request.UnAuth, params, &resp)
}

// CreateOrderV2 creates a new order through the v2 API
func (ap *Apexpro) CreateOrderV2(ctx context.Context, arg *CreateOrderParams) (*OrderDetail, error) {
	params, err := ap.orderCreationParamsFilter(ctx, arg)
	if err != nil {
		return nil, err
	}
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v2/create-order", request.UnAuth, params, &resp)
}

// CreateOrderV1 creates a new order through the v2 API
func (ap *Apexpro) CreateOrderV1(ctx context.Context, arg *CreateOrderParams) (*OrderDetail, error) {
	params, err := ap.orderCreationParamsFilter(ctx, arg)
	if err != nil {
		return nil, err
	}
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v1/create-order", request.UnAuth, params, &resp)
}

// FastWithdrawalV2 withdraws an asset
func (ap *Apexpro) FastWithdrawalV1(ctx context.Context, arg *FastWithdrawalParams) (*WithdrawalResponse, error) {
	return ap.fastWithdrawal(ctx, arg, "v1/fast-withdraw")
}

// FastWithdrawalV2 withdraws an asset
func (ap *Apexpro) FastWithdrawalV2(ctx context.Context, arg *FastWithdrawalParams) (*WithdrawalResponse, error) {
	return ap.fastWithdrawal(ctx, arg, "v2/fast-withdraw")
}

func (ap *Apexpro) fastWithdrawal(ctx context.Context, arg *FastWithdrawalParams, path string) (*WithdrawalResponse, error) {
	if arg == nil || *arg == (FastWithdrawalParams{}) {
		return nil, common.ErrNilPointer
	}
	if arg.Amount <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	if arg.ClientID != "" {
		return nil, order.ErrClientOrderIDMustBeSet
	}
	if arg.Expiration.IsZero() {
		return nil, errExpirationTimeRequired
	}
	if arg.Asset.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	if arg.ERC20Address == "" {
		return nil, errEthereumAddressMissing
	}
	if arg.Fees <= 0 {
		return nil, errLimitFeeRequired
	}
	if arg.ChainID == "" {
		return nil, errChainIDMissing
	}
	signature, err := ap.ProcessConditionalTransfer(ctx, arg)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("amount", strconv.FormatFloat(arg.Amount, 'f', -1, 64))
	params.Set("clientId", arg.ClientID)
	params.Set("expiration", strconv.FormatInt(arg.Expiration.UnixMilli(), 10))
	params.Set("asset", arg.Asset.String())
	params.Set("erc20Address", arg.ERC20Address)
	params.Set("fee", strconv.FormatFloat(arg.Fees, 'f', -1, 64))
	params.Set("chainId", arg.ChainID)
	if arg.IPAccountID != "" {
		params.Set("lpAccountId", arg.IPAccountID)
	}
	params.Set("signature", signature)
	var resp *WithdrawalResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, &resp)
}

func (ap *Apexpro) getWorstPrice(ctx context.Context, symbol, side, path string, amount float64, ePath exchange.URL) (*SymbolWorstPrice, error) {
	if symbol == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	if side == "" {
		return nil, order.ErrSideIsInvalid
	}
	if amount <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	params := url.Values{}
	params.Set("size", strconv.FormatFloat(amount, 'f', -1, 64))
	params.Set("side", side)
	params.Set("symbol", symbol)
	var resp *SymbolWorstPrice
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// CancelPerpOrder cancels a perpetual contract order cancellation.
func (ap *Apexpro) CancelPerpOrder(ctx context.Context, orderID string) (types.Number, error) {
	return ap.cancelOrderByID(ctx, orderID, "v3/delete-order")
}

// CancelPerpOrderByClientOrderID cancels a perpetual contract order by client order ID.
func (ap *Apexpro) CancelPerpOrderByClientOrderID(ctx context.Context, clientOrderID string) (types.Number, error) {
	return ap.cancelOrderByID(ctx, clientOrderID, "v3/delete-client-order-id")
}

func (ap *Apexpro) cancelOrderByID(ctx context.Context, id, path string) (types.Number, error) {
	if id == "" {
		return 0, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("id", id)
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, &resp)
}

// CancelAllOpenOrdersV3 cancels all open orders
func (ap *Apexpro) CancelAllOpenOrdersV3(ctx context.Context, symbols []string) error {
	params := url.Values{}
	if len(symbols) > 0 {
		params.Set("symbol", strings.Join(symbols, ","))
	}
	return ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v3/delete-open-orders", request.UnAuth, params, nil)
}

// CancelPerpOrderV2 cancels a perpetual contract futures order.
func (ap *Apexpro) CancelPerpOrderV2(ctx context.Context, orderID string, token currency.Code) (types.Number, error) {
	if orderID == "" {
		return 0, order.ErrOrderIDNotSet
	}
	if token.IsEmpty() {
		return 0, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("id", orderID)
	params.Set("token", token.String())
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v2/delete-order", request.UnAuth, params, &resp)
}

// GetOpenOrders retrieves an active orders
func (ap *Apexpro) GetOpenOrders(ctx context.Context) ([]OrderDetail, error) {
	var resp []OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/open-orders", request.UnAuth, nil, &resp)
}

// GetOpenOrders retrieves an active orders
func (ap *Apexpro) GetOpenOrdersV2(ctx context.Context, token currency.Code) ([]OrderDetail, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("token", token.String())
	var resp []OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v2/open-orders", request.UnAuth, params, &resp)
}

// GetOpenOrdersV1 retrieves an active orders
func (ap *Apexpro) GetOpenOrdersV1(ctx context.Context) ([]OrderDetail, error) {
	var resp []OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v1/open-orders", request.UnAuth, nil, &resp)
}

// GetAllOrderHistory retrieves all order history
// possible ordersKind are "ACTIVE","CONDITION", and "HISTORY"
func (ap *Apexpro) GetAllOrderHistory(ctx context.Context, symbol, side, orderType, orderStatus, ordersKind string, startTime, endTime time.Time, page, limit int64) (*OrderHistoryResponse, error) {
	return ap.getAllOrderHistory(ctx, symbol, side, orderType, orderStatus, ordersKind, "", "v3/history-orders", startTime, endTime, page, limit)
}

// GetAllOrderHistoryV2 retrieves all order history
// possible ordersKind are "ACTIVE","CONDITION", and "HISTORY"
func (ap *Apexpro) GetAllOrderHistoryV2(ctx context.Context, token currency.Code, symbol, side, orderType, orderStatus, ordersKind string, startTime, endTime time.Time, page, limit int64) (*OrderHistoryResponse, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getAllOrderHistory(ctx, symbol, side, orderType, orderStatus, ordersKind, token.String(), "v2/history-orders", startTime, endTime, page, limit)
}

// GetAllOrderHistoryV1 retrieves all order history
// possible ordersKind are "ACTIVE","CONDITION", and "HISTORY"
func (ap *Apexpro) GetAllOrderHistoryV1(ctx context.Context, symbol, side, orderType, orderStatus, ordersKind string, startTime, endTime time.Time, page, limit int64) (*OrderHistoryResponse, error) {
	return ap.getAllOrderHistory(ctx, symbol, side, orderType, orderStatus, ordersKind, "", "v1/history-orders", startTime, endTime, page, limit)
}

func (ap *Apexpro) getAllOrderHistory(ctx context.Context, symbol, side, orderType, orderStatus, ordersKind, token, path string, startTime, endTime time.Time, page, limit int64) (*OrderHistoryResponse, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if ordersKind != "" {
		params.Set("orderType", ordersKind)
	}
	if side != "" {
		params.Set("side", side)
	}
	if orderType != "" {
		params.Set("type", orderType)
	}
	if orderStatus != "" {
		params.Set("status", orderStatus)
	}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
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
	var resp *OrderHistoryResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetOrderID retrieves a single order by ID.
func (ap *Apexpro) GetOrderID(ctx context.Context, orderID string) (*OrderDetail, error) {
	if orderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	return ap.getOrderID(ctx, orderID, "v3/order")
}

// // GetSingleOrderByOrderIDV2 retrieves a single order detail by ID through the V2 API
func (ap *Apexpro) GetSingleOrderByOrderIDV2(ctx context.Context, orderID string, token currency.Code) (*OrderDetail, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getSingleOrderV2(ctx, orderID, "v2/get-order", token.String(), exchange.RestSpot)
}

// GetSingleOrderByClientOrderIDV2 retrieves a single order detail by client supplied order ID through the V2 API
func (ap *Apexpro) GetSingleOrderByClientOrderIDV2(ctx context.Context, orderID string, token currency.Code) (*OrderDetail, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getSingleOrderV2(ctx, orderID, "v2/order-by-client-order-id", token.String(), exchange.RestSpot)
}

// GetSingleOrderByOrderIDV1 retrieves a single order detail by ID through the V1 API
func (ap *Apexpro) GetSingleOrderByOrderIDV1(ctx context.Context, orderID string) (*OrderDetail, error) {
	return ap.getSingleOrderV2(ctx, orderID, "v1/get-order", "", exchange.RestSpot)
}

// GetSingleOrderByClientOrderIDV1 retrieves a single order detail by client supplied order ID through the V1 API
func (ap *Apexpro) GetSingleOrderByClientOrderIDV1(ctx context.Context, orderID string) (*OrderDetail, error) {
	return ap.getSingleOrderV2(ctx, orderID, "v1/order-by-client-order-id", "", exchange.RestSpot)
}

func (ap *Apexpro) getSingleOrderV2(ctx context.Context, orderID, path, token string, ePath exchange.URL) (*OrderDetail, error) {
	if orderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("id", orderID)
	params.Set("token", token)
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetVerificationEmailLink retrieves a link to the verification email
func (ap *Apexpro) GetVerificationEmailLink(ctx context.Context, userID string, token currency.Code) error {
	params := url.Values{}
	if userID == "" {
		return errUserIDRequired
	}
	params.Set("userId", userID)
	if !token.IsEmpty() {
		params.Set("token", token.String())
	}
	return ap.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues("v1/verify-email", params), request.UnAuth, nil)
}

// LinkDevice bind a device to an account.
// possible device type values: 1 (ios_firebase), 2 (android_firebase)
func (ap *Apexpro) LinkDevice(ctx context.Context, deviceToken currency.Code, deviceType string) error {
	if deviceToken.IsEmpty() {
		return fmt.Errorf("%w, device token is required", currency.ErrCurrencyCodeEmpty)
	}
	if deviceType == "" {
		return errDeviceTypeIsRequired
	}
	params := url.Values{}
	params.Set("deviceToken", deviceType)
	params.Set("deviceType", deviceType)
	return ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodPost, "v1/bind-device", request.UnAuth, params, nil)
}

// GetOrderClientOrderID retrieves a single order by client order ID.
func (ap *Apexpro) GetOrderClientOrderID(ctx context.Context, clientOrderID string) (*OrderDetail, error) {
	if clientOrderID == "" {
		return nil, order.ErrClientOrderIDMustBeSet
	}
	return ap.getOrderID(ctx, clientOrderID, "v3/order-by-client-order-id")
}

func (ap *Apexpro) getOrderID(ctx context.Context, id, path string) (*OrderDetail, error) {
	params := url.Values{}
	params.Set("id", id)
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetFundingRateV3 retrieves a funding rate information.
func (ap *Apexpro) GetFundingRateV3(ctx context.Context, symbol, side, status string, startTime, endTime time.Time, limit, page int64) (*FundingRateResponse, error) {
	return ap.getFundingRate(ctx, symbol, side, status, "", "v3/funding", startTime, endTime, limit, page, exchange.RestFutures)
}

// GetFundingRateV2 retrieves a funding rate infor for a contract.
func (ap *Apexpro) GetFundingRateV2(ctx context.Context, token currency.Code, symbol, side, status string, startTime, endTime time.Time, limit, page int64) (*FundingRateResponse, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getFundingRate(ctx, symbol, side, status, token.String(), "v2/funding", startTime, endTime, limit, page, exchange.RestSpot)
}

// GetFundingRateV1 retrieves a funding rate information.
func (ap *Apexpro) GetFundingRateV1(ctx context.Context, symbol, side, status string, startTime, endTime time.Time, limit, page int64) (*FundingRateResponse, error) {
	return ap.getFundingRate(ctx, symbol, side, status, "", "v1/funding", startTime, endTime, limit, page, exchange.RestSpot)
}

func (ap *Apexpro) getFundingRate(ctx context.Context, symbol, side, status, token, path string, startTime, endTime time.Time, limit, page int64, ePath exchange.URL) (*FundingRateResponse, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if side != "" {
		params.Set("side", side)
	}
	if status != "" {
		params.Set("status", status)
	}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp *FundingRateResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetUserHistorialProfitAndLoss retrieves a profit and loss history of order positions
func (ap *Apexpro) GetUserHistorialProfitAndLoss(ctx context.Context, symbol, positionType string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
	return ap.getUserHistorialProfitAndLoss(ctx, symbol, positionType, "", "v3/historical-pnl", startTime, endTime, page, limit, exchange.RestFutures)
}

// GetUserHistorialProfitAndLossV1 retrieves a profit and loss history of order positions through the V1 API endpoint.
func (ap *Apexpro) GetUserHistorialProfitAndLossV1(ctx context.Context, symbol, positionType string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
	return ap.getUserHistorialProfitAndLoss(ctx, symbol, positionType, "", "v1/historical-pnl", startTime, endTime, page, limit, exchange.RestSpot)
}

// GetUserHistorialProfitAndLossV2 retrieves a profit and loss history of order positions.
func (ap *Apexpro) GetUserHistorialProfitAndLossV2(ctx context.Context, token currency.Code, symbol, positionType string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getUserHistorialProfitAndLoss(ctx, symbol, positionType, token.String(), "v2/historical-pnl", startTime, endTime, page, limit, exchange.RestSpot)
}

func (ap *Apexpro) getUserHistorialProfitAndLoss(ctx context.Context, symbol, positionType, token, path string, startTime, endTime time.Time, page, limit int64, ePath exchange.URL) (*PNLHistory, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if symbol != "" {
		params.Set("symbol", symbol)
	}
	if positionType != "" {
		params.Set("type", positionType)
	}
	if !startTime.IsZero() {
		params.Set("beginTimeInclusive", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTimeExclusive", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	if page > 0 {
		params.Set("page", strconv.FormatInt(page, 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	var resp *PNLHistory
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetYesterdaysPNL retrieves yesterdays profit and loss(PNL)
func (ap *Apexpro) GetYesterdaysPNL(ctx context.Context) (types.Number, error) {
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/yesterday-pnl", request.UnAuth, nil, &resp)
}

// GetYesterdaysPNLV1 retrieves yesterdays profit and loss(PNL)
func (ap *Apexpro) GetYesterdaysPNLV1(ctx context.Context) (types.Number, error) {
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v1/yesterday-pnl", request.UnAuth, nil, &resp)
}

// GetYesterdaysPNLV2 retrieves yesterdays profit and loss(PNL)
func (ap *Apexpro) GetYesterdaysPNLV2(ctx context.Context, token currency.Code) (types.Number, error) {
	if token.IsEmpty() {
		return 0, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("token", token.String())
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v2/yesterday-pnl", request.UnAuth, params, &resp)
}

// GetHistoricalAssetValue retrieves a historical asset value
func (ap *Apexpro) GetHistoricalAssetValue(ctx context.Context, startTime, endTime time.Time) (*AssetValueHistory, error) {
	return ap.getHistoricalAssetValue(ctx, "", "v3/history-value", startTime, endTime, exchange.RestFutures)
}

// GetHistoricalAssetValueV1 retrieves a historical asset value through the V1 APi endpoints.
func (ap *Apexpro) GetHistoricalAssetValueV1(ctx context.Context, startTime, endTime time.Time) (*AssetValueHistory, error) {
	return ap.getHistoricalAssetValue(ctx, "", "v1/history-value", startTime, endTime, exchange.RestSpot)
}

// GetHistoricalAssetValueV2 retrieves a historical asset value
func (ap *Apexpro) GetHistoricalAssetValueV2(ctx context.Context, token currency.Code, startTime, endTime time.Time) (*AssetValueHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getHistoricalAssetValue(ctx, token.String(), "v2/history-value", startTime, endTime, exchange.RestSpot)
}

func (ap *Apexpro) getHistoricalAssetValue(ctx context.Context, token, path string, startTime, endTime time.Time, ePath exchange.URL) (*AssetValueHistory, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if !startTime.IsZero() {
		params.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	}
	if !endTime.IsZero() {
		params.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}
	var resp *AssetValueHistory
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodGet, path, request.UnAuth, params, &resp)
}

// SetInitialMarginRateInfo sets an initial margin rate
func (ap *Apexpro) SetInitialMarginRateInfo(ctx context.Context, symbol string, initialMarginRate float64) error {
	return ap.setInitialMarginRateInfo(ctx, symbol, "v3/set-initial-margin-rate", initialMarginRate, exchange.RestFutures)
}

// SetInitialMarginRateInfoV1 sets an initial margin rate
func (ap *Apexpro) SetInitialMarginRateInfoV1(ctx context.Context, symbol string, initialMarginRate float64) error {
	return ap.setInitialMarginRateInfo(ctx, symbol, "v1/set-initial-margin-rate", initialMarginRate, exchange.RestSpot)
}

// SetInitialMarginRateInfoV2 sets an initial margin rate
func (ap *Apexpro) SetInitialMarginRateInfoV2(ctx context.Context, symbol string, initialMarginRate float64) error {
	return ap.setInitialMarginRateInfo(ctx, symbol, "v2/set-initial-margin-rate", initialMarginRate, exchange.RestSpot)
}

func (ap *Apexpro) setInitialMarginRateInfo(ctx context.Context, symbol, path string, initialMarginRate float64, ePath exchange.URL) error {
	if symbol == "" {
		return currency.ErrSymbolStringEmpty
	}
	if initialMarginRate <= 0 {
		return errInitialMarginRateRequired
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("initialMarginRate", strconv.FormatFloat(initialMarginRate, 'f', -1, 64))
	return ap.SendAuthenticatedHTTPRequest(ctx, ePath, http.MethodPost, path, request.UnAuth, params, nil)
}

// WithdrawAsset posts an asset withdrawal
func (ap *Apexpro) WithdrawAsset(ctx context.Context, arg *AssetWithdrawalParams) (*WithdrawalResponse, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return nil, err
	}
	if arg.Amount <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	if arg.ClientWithdrawID == "" {
		return nil, order.ErrClientOrderIDMustBeSet
	}
	if arg.Timestamp.IsZero() {
		return nil, errInvalidTimestamp
	}
	if arg.EthereumAddress == "" && creds.SubAccount == "" {
		return nil, errEthereumAddressMissing
	} else if arg.EthereumAddress == "" {
		arg.EthereumAddress = creds.SubAccount
	}
	if arg.L2Key == "" && creds.L2Key == "" {
		return nil, errL2KeyMissing
	} else if arg.L2Key == "" {
		arg.L2Key = creds.L2Key
	}
	if arg.ToChainID == "" {
		return nil, fmt.Errorf("%w, toChainID is required", errChainIDMissing)
	}
	if arg.L2SourceTokenID.IsEmpty() {
		return nil, fmt.Errorf("%w, l2SourceTokenId is required", currency.ErrCurrencyCodeEmpty)
	}
	if arg.L1TargetTokenID.IsEmpty() {
		return nil, fmt.Errorf("%w, l1TargetTokenId is required", currency.ErrCurrencyCodeEmpty)
	}
	if arg.Nonce == "" {
		arg.Nonce = strconv.FormatInt(ap.Websocket.Conn.GenerateMessageID(true), 10)
	}
	params := url.Values{}
	params.Set("amount", strconv.FormatFloat(arg.Amount, 'f', -1, 64))
	params.Set("clientWithdrawId", arg.ClientWithdrawID)
	params.Set("timestamp", strconv.FormatInt(arg.Timestamp.UnixMilli(), 10))
	params.Set("ethAddress", arg.EthereumAddress)
	params.Set("subAccountId", arg.SubAccountID)
	params.Set("l2Key", arg.L2Key)
	params.Set("toChainId", arg.ToChainID)
	params.Set("l2SourceTokenId", arg.L2SourceTokenID.String())
	params.Set("l1TargetTokenId", arg.L1TargetTokenID.String())
	if arg.Fee != 0 {
		params.Set("fee", strconv.FormatFloat(arg.Fee, 'f', -1, 64))
	}
	params.Set("isFastWithdraw", strconv.FormatBool(arg.IsFastWithdraw))
	params.Set("nonce", arg.Nonce)
	signature, err := ap.ProcessWithdrawalToAddressSignatureV3(ctx, arg)
	if err != nil {
		return nil, err
	}
	if arg.ZKAccountID == "" {
		return nil, errZeroKnowledgeAccountIDMissing
	}
	params.Set("zkAccountId", arg.ZKAccountID)
	params.Set("signature", signature)
	var resp *WithdrawalResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/withdrawal", request.UnAuth, params, &resp)
}

// ----------------------------------------------------- Private V2 Endpoints --------------------------------------------------------------------------------

// UserWithdrawalV2 withdraws an asset
func (ap *Apexpro) UserWithdrawalV2(ctx context.Context, amount float64, clientID string, expiration time.Time, asset currency.Code) (*WithdrawalResponse, error) {
	if amount <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	if clientID == "" {
		return nil, errUserNonceRequired
	}
	if expiration.IsZero() {
		return nil, errExpirationTimeRequired
	}
	if asset.IsEmpty() {
		return nil, currency.ErrCurrencyCodeEmpty
	}
	params := url.Values{}
	signature, err := ap.ProcessWithdrawalSignature(ctx, &WithdrawalParams{
		Amount:         amount,
		ClientID:       strconv.FormatInt(ap.Websocket.Conn.GenerateMessageID(true), 10),
		ExpirationTime: expiration,
		Asset:          asset,
	})
	if err != nil {
		return nil, err
	}
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	params.Set("clientId", clientID)
	params.Set("expiration", strconv.FormatInt(expiration.UnixMilli(), 10))
	params.Set("asset", asset.String())
	params.Set("signature", signature)
	var resp *WithdrawalResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v2/create-withdrawal", request.UnAuth, params, &resp)
}

// WithdrawalToAddressV1 withdraws as asset to an ethereum address
func (ap *Apexpro) WithdrawalToAddressV1(ctx context.Context, arg *WithdrawalToAddressParams) (*WithdrawalResponse, error) {
	return ap.withdrawalToAddress(ctx, arg, "v1/create-withdrawal-to-address")
}

// WithdrawalToAddressV2 withdraws as asset to an ethereum address
func (ap *Apexpro) WithdrawalToAddressV2(ctx context.Context, arg *WithdrawalToAddressParams) (*WithdrawalResponse, error) {
	return ap.withdrawalToAddress(ctx, arg, "v2/create-withdrawal-to-address")
}

func (ap *Apexpro) withdrawalToAddress(ctx context.Context, arg *WithdrawalToAddressParams, path string) (*WithdrawalResponse, error) {
	if arg == nil || *arg == (WithdrawalToAddressParams{}) {
		return nil, common.ErrNilPointer
	}
	if arg.Amount <= 0 {
		return nil, order.ErrAmountBelowMin
	}
	if arg.ClientID == "" {
		return nil, order.ErrClientOrderIDMustBeSet
	}
	if arg.ExpirationTime.IsZero() {
		return nil, errExpirationTimeRequired
	}
	if arg.Asset.IsEmpty() {
		return nil, fmt.Errorf("%w, asset is required", currency.ErrCurrencyCodeEmpty)
	}
	if arg.EthereumAddress == "" {
		return nil, errEthereumAddressMissing
	}
	signature, err := ap.ProcessWithdrawalToAddressSignature(ctx, arg)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("amount", strconv.FormatFloat(arg.Amount, 'f', -1, 64))
	params.Set("clientId", arg.ClientID)
	params.Set("expiration", strconv.FormatInt(arg.ExpirationTime.UnixMilli(), 10))
	params.Set("asset", arg.Asset.String())
	params.Set("ethAddress", arg.EthereumAddress)
	params.Set("signature", signature)
	var resp *WithdrawalResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, &resp)
}

// CrossChainWithdrawals withdraws an asset through different chains
// func (ap *Apexpro) CrossChainWithdrawals(ctx context.Context, arg *FastWithdrawalParams) (*WithdrawalResponse, error) {

// 	params := url.Values{}
// 	var resp *WithdrawalResponse
// 	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v1/cross-chain-withdraw", request.UnAuth, params, &resp)
// }

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
	return ap.SendPayload(ctx, f, func() (*request.Item, error) {
		return &request.Item{
			Method:        http.MethodGet,
			Path:          endpointPath + path,
			Result:        response,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording,
		}, nil
	}, request.UnauthenticatedRequest)
}

// SendAuthenticatedHTTPRequest sends an authenticated HTTP request.
func (ap *Apexpro) SendAuthenticatedHTTPRequest(ctx context.Context, ePath exchange.URL, method, path string, f request.EndpointLimit, params url.Values, result interface{}, onboarding ...bool) error {
	creds, err := ap.GetCredentials(ctx)
	if err != nil {
		return err
	}
	endpointPath, err := ap.API.Endpoints.GetURL(ePath)
	if err != nil {
		return err
	}
	response := &UserResponse{
		Data: result,
	}
	var dataString string
	if method == http.MethodGet {
		path = common.EncodeURLValues(path, params)
	} else {
		dataString = params.Encode()
	}
	err = ap.SendPayload(ctx, f, func() (*request.Item, error) {
		timestamp := time.Now().UnixMilli()
		message := strconv.FormatInt(timestamp, 10) + method + ("/api/" + path) + dataString
		encodedSecret := base64.StdEncoding.EncodeToString([]byte(creds.Secret))
		var hmacSigned []byte
		hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256,
			[]byte(message),
			[]byte(encodedSecret))
		if err != nil {
			return nil, err
		}
		headers := make(map[string]string)
		headers["APEX-API-KEY"] = creds.Key
		headers["APEX-SIGNATURE"] = base64.StdEncoding.EncodeToString(hmacSigned)
		headers["APEX-TIMESTAMP"] = strconv.FormatInt(timestamp, 10)
		headers["APEX-PASSPHRASE"] = creds.ClientID
		if len(onboarding) > 0 && onboarding[0] {
			if creds.SubAccount == "" {
				return nil, errEthereumAddressMissing
			}
			headers = make(map[string]string)
			headers["APEX-SIGNATURE"] = base64.StdEncoding.EncodeToString(hmacSigned)
			headers["APEX-ETHEREUM-ADDRESS"] = creds.SubAccount
		} else if len(onboarding) > 0 {
			headers = make(map[string]string)
		}
		reqItem := &request.Item{
			Method:        method,
			Path:          endpointPath + path,
			Headers:       headers,
			Result:        response,
			Verbose:       ap.Verbose,
			HTTPDebugging: ap.HTTPDebugging,
			HTTPRecording: ap.HTTPRecording,
		}
		if dataString != "" {
			reqItem.Body = bytes.NewBuffer([]byte(dataString))
		}
		return reqItem, nil
	}, request.AuthenticatedRequest)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("code: %d msg: %q", response.Code, response.Message)
	}
	return nil
}
