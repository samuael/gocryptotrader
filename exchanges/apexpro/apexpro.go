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
	"github.com/thrasher-corp/gocryptotrader/types"
)

// Apexpro is the overarching type across this package
type Apexpro struct {
	exchange.Base
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
	errChainIDMissing                = errors.New("chain ID is missing")
	errOrderbookLevelIsRequired      = errors.New("orderbook level is required")
	errInvalidTimestamp              = errors.New("err invalid timestamp")
	errZeroKnowledgeAccountIDMissing = errors.New("zero knowledge account id is required")
	errSubAccountIDMissing           = errors.New("missing sub-account id")
	errUserNonceRequired             = errors.New("nonce is required")
	errInitialMarginRateRequired     = errors.New("initial margin rate required")
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

// GenerateNonceV3 generate and obtain nonce before registration. The nonce is used to assemble the signature field upon registration.
func (ap *Apexpro) GenerateNonceV3(ctx context.Context, l2Key, ethereumAddress, chainID string) (*NonceResponse, error) {
	return ap.generateNonce(ctx, l2Key, ethereumAddress, chainID, "v3/generate-nonce")
}

// GenerateNonceV2 before registering, generate and obtain a nonce. The nonce serves the purpose of assembling the signature field during the registration process.
func (ap *Apexpro) GenerateNonceV2(ctx context.Context, l2Key, ethereumAddress, chainID string) (*NonceResponse, error) {
	return ap.generateNonce(ctx, l2Key, ethereumAddress, chainID, "v2/generate-nonce")
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

// GetUserAccountDataV2 get an user account detail by use id throught the V2 API.
func (ap *Apexpro) GetUserAccountDataV2(ctx context.Context) (*UserAccountDetail, error) {
	return ap.getUserAccountData(ctx, "v2/account")
}

func (ap *Apexpro) getUserAccountData(ctx context.Context, path string) (*UserAccountDetail, error) {
	var resp *UserAccountDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, nil, &resp)
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

// GetUserTransferDataV2 retrieves user's asset transfer information.
func (ap *Apexpro) GetUserTransferDataV2(ctx context.Context, ccy currency.Code, startTime, endTime time.Time, transferType string, chainIDs []string, limit, page int64) (*UserWithdrawalsV2, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v2/transfers", request.UnAuth, params, &resp)
}

// GetUserWithdrawalListV2 retrieves asset withdrawal list.
func (ap *Apexpro) GetUserWithdrawalListV2(ctx context.Context, transferType string, startTime, endTime time.Time, page, limit int64) (*WithdrawalsV2, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v2/withdraw-list", request.UnAuth, params, &resp)
}

// GetFastAndCrossChainWithdrawalFees retrieves fee information of fast and cross-chain withdrawal transactions.
func (ap *Apexpro) GetFastAndCrossChainWithdrawalFees(ctx context.Context, amount float64, chainID string, token currency.Code) (*FastAndCrossChainWithdrawalFees, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	if chainID == "" {
		return nil, errChainIDMissing
	}
	params := url.Values{}
	params.Set("token", token.String())
	if amount > 0 {
		params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	}
	if chainID != "" {
		params.Set("chainId", chainID)
	}
	var resp *FastAndCrossChainWithdrawalFees
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v2/uncommon-withdraw-fee", request.UnAuth, params, &resp)
}

// GetAssetWithdrawalAndTransferLimit represents an asset withdrawal and transfer limit per interval.
func (ap *Apexpro) GetAssetWithdrawalAndTransferLimit(ctx context.Context, currencyID currency.Code) (*TransferAndWithdrawalLimit, error) {
	if currencyID.IsEmpty() {
		return nil, fmt.Errorf("%w, currencyID is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("currencyId", currencyID.String())
	var resp *TransferAndWithdrawalLimit
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, "v2/transfer-limit", request.UnAuth, params, &resp)
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
	return ap.getTradeHistory(ctx, symbol, side, orderType, "", "v3/fills", startTime, endTime, page, limit)
}

// GetTradeHistoryV2 retrieves trade fills history through the v2 API
func (ap *Apexpro) GetTradeHistoryV2(ctx context.Context, symbol, side, orderType string, token currency.Code, startTime, endTime time.Time, page, limit int64) (*TradeHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getTradeHistory(ctx, symbol, side, orderType, token.String(), "v2/fills", startTime, endTime, page, limit)
}

func (ap *Apexpro) getTradeHistory(ctx context.Context, symbol, side, orderType, token, path string, startTime, endTime time.Time, page, limit int64) (*TradeHistory, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetWorstPriceV3 retrieves the worst market price from orderbook
func (ap *Apexpro) GetWorstPriceV3(ctx context.Context, symbol, side string, amount float64) (*SymbolWorstPrice, error) {
	return ap.getWorstPrice(ctx, symbol, side, "v3/get-worst-price", amount)
}

// GetWorstPriceV2 retrieves the worst market price from orderbook
func (ap *Apexpro) GetWorstPriceV2(ctx context.Context, symbol, side string, amount float64) (*SymbolWorstPrice, error) {
	return ap.getWorstPrice(ctx, symbol, side, "v2/get-worst-price", amount)
}

func (ap *Apexpro) getWorstPrice(ctx context.Context, symbol, side, path string, amount float64) (*SymbolWorstPrice, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// CancelPerpOrder cancels a perpetual contract order cancellation.
func (ap *Apexpro) CancelPerpOrder(ctx context.Context, id string) (types.Number, error) {
	if id == "" {
		return 0, order.ErrOrderIDNotSet
	}
	params := url.Values{}
	params.Set("id", id)
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, "v3/delete-order", request.UnAuth, params, &resp)
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v2/open-orders", request.UnAuth, params, &resp)
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
	return ap.getSingleOrderV2(ctx, orderID, "v2/get-order", token)
}

// GetSingleOrderByClientOrderIDV2 retrieves a single order detail by client supplied order ID through the V2 API
func (ap *Apexpro) GetSingleOrderByClientOrderIDV2(ctx context.Context, orderID string, token currency.Code) (*OrderDetail, error) {
	return ap.getSingleOrderV2(ctx, orderID, "v2/order-by-client-order-id", token)
}

func (ap *Apexpro) getSingleOrderV2(ctx context.Context, orderID, path string, token currency.Code) (*OrderDetail, error) {
	if orderID == "" {
		return nil, order.ErrOrderIDNotSet
	}
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("id", orderID)
	params.Set("token", token.String())
	var resp *OrderDetail
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, http.MethodGet, path, request.UnAuth, params, &resp)
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
	return ap.getFundingRate(ctx, symbol, side, status, "", "v3/funding", startTime, endTime, limit, page)
}

// GetFundingRateV2 retrieves a funding rate infor for a contract.
func (ap *Apexpro) GetFundingRateV2(ctx context.Context, token currency.Code, symbol, side, status string, startTime, endTime time.Time, limit, page int64) (*FundingRateResponse, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getFundingRate(ctx, symbol, side, status, token.String(), "v2/funding", startTime, endTime, limit, page)
}

func (ap *Apexpro) getFundingRate(ctx context.Context, symbol, side, status, token, path string, startTime, endTime time.Time, limit, page int64) (*FundingRateResponse, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetUserHistorialProfitAndLoss retrieves a profit and loss history of order positions
func (ap *Apexpro) GetUserHistorialProfitAndLoss(ctx context.Context, symbol, positionType string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
	return ap.getUserHistorialProfitAndLoss(ctx, symbol, positionType, "", "v3/historical-pnl", startTime, endTime, page, limit)
}

// GetUserHistorialProfitAndLossV2 retrieves a profit and loss history of order positions.
func (ap *Apexpro) GetUserHistorialProfitAndLossV2(ctx context.Context, token currency.Code, symbol, positionType string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getUserHistorialProfitAndLoss(ctx, symbol, positionType, token.String(), "v2/historical-pnl", startTime, endTime, page, limit)
}

func (ap *Apexpro) getUserHistorialProfitAndLoss(ctx context.Context, symbol, positionType, token, path string, startTime, endTime time.Time, page, limit int64) (*PNLHistory, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// GetYesterdaysPNL retrieves yesterdays profit and loss(PNL)
func (ap *Apexpro) GetYesterdaysPNL(ctx context.Context) (types.Number, error) {
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/yesterday-pnl", request.UnAuth, nil, &resp)
}

// GetYesterdaysPNLV2 retrieves yesterdays profit and loss(PNL)
func (ap *Apexpro) GetYesterdaysPNLV2(ctx context.Context, token currency.Code) (types.Number, error) {
	if token.IsEmpty() {
		return 0, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	params := url.Values{}
	params.Set("token", token.String())
	var resp types.Number
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v2/yesterday-pnl", request.UnAuth, params, &resp)
}

// GetHistoricalAssetValue retrieves a historical asset value
func (ap *Apexpro) GetHistoricalAssetValue(ctx context.Context, startTime, endTime time.Time) (*AssetValueHistory, error) {
	return ap.getHistoricalAssetValue(ctx, "", "v3/history-value", startTime, endTime)
}

// GetHistoricalAssetValueV2 retrieves a historical asset value
func (ap *Apexpro) GetHistoricalAssetValueV2(ctx context.Context, token currency.Code, startTime, endTime time.Time) (*AssetValueHistory, error) {
	if token.IsEmpty() {
		return nil, fmt.Errorf("%w, token is required", currency.ErrCurrencyCodeEmpty)
	}
	return ap.getHistoricalAssetValue(ctx, token.String(), "v2/history-value", startTime, endTime)
}

func (ap *Apexpro) getHistoricalAssetValue(ctx context.Context, token, path string, startTime, endTime time.Time) (*AssetValueHistory, error) {
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
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, path, request.UnAuth, params, &resp)
}

// SetInitialMarginRateInfo sets an initial margin rate
func (ap *Apexpro) SetInitialMarginRateInfo(ctx context.Context, symbol string, initialMarginRate float64) error {
	return ap.setInitialMarginRateInfo(ctx, symbol, "v3/set-initial-margin-rate", initialMarginRate)
}

// SetInitialMarginRateInfoV2 sets an initial margin rate
func (ap *Apexpro) SetInitialMarginRateInfoV2(ctx context.Context, symbol string, initialMarginRate float64) error {
	return ap.setInitialMarginRateInfo(ctx, symbol, "v2/set-initial-margin-rate", initialMarginRate)
}

func (ap *Apexpro) setInitialMarginRateInfo(ctx context.Context, symbol, path string, initialMarginRate float64) error {
	if symbol == "" {
		return currency.ErrSymbolStringEmpty
	}
	if initialMarginRate <= 0 {
		return errInitialMarginRateRequired
	}
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("initialMarginRate", strconv.FormatFloat(initialMarginRate, 'f', -1, 64))
	return ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodPost, path, request.UnAuth, params, nil)
}

// WithdrawAsset posts an asset withdrawal
func (ap *Apexpro) WithdrawAsset(ctx context.Context, arg *AssetWithdrawalParams) (*WithdrawalResponse, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return nil, err
	}
	params := url.Values{}
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
	if arg.ZKAccountID == "" {
		return nil, errZeroKnowledgeAccountIDMissing
	}
	if arg.SubAccountID == "" {
		return nil, errSubAccountIDMissing
	}
	if arg.L2Key == "" && creds.L2Key == "" {
		return nil, errL2KeyMissing
	} else if arg.L2Key == "" {
		arg.L2Key = creds.L2Key
	}
	if arg.ToChainID == "" {
		return nil, fmt.Errorf("%w, toChainID is required", errChainIDMissing)
	}
	if arg.Nonce == 0 {
		return nil, errUserNonceRequired
	}
	if arg.L2SourceTokenID != "" {
		params.Set("l2SourceTokenId", arg.L2SourceTokenID)
	}
	if arg.L1TargetTokenID != "" {
		params.Set("l1TargetTokenId", arg.L1TargetTokenID)
	}
	if arg.Fee != 0 {
		params.Set("fee", strconv.FormatFloat(arg.Fee, 'f', -1, 64))
	}
	params.Set("isFastWithdraw", strconv.FormatBool(arg.IsFastWithdraw))
	params.Set("nonce", strconv.FormatInt(arg.Nonce, 10))
	params.Set("toChainId", arg.ToChainID)
	params.Set("l2Key", arg.L2Key)
	params.Set("subAccountId", arg.SubAccountID)
	params.Set("zkAccountId", arg.ZKAccountID)
	params.Set("zkAccountId", arg.ZKAccountID)
	params.Set("ethAddress", arg.EthereumAddress)
	params.Set("timestamp", strconv.FormatInt(arg.Timestamp.UnixMilli(), 10))
	params.Set("clientWithdrawId", arg.ClientWithdrawID)
	params.Set("amount", strconv.FormatFloat(arg.Amount, 'f', -1, 64))

	// TODO: generate signature and fill in the parameters

	var resp *WithdrawalResponse
	return resp, ap.SendAuthenticatedHTTPRequest(ctx, exchange.RestFutures, http.MethodGet, "v3/withdraw-fee", request.UnAuth, params, &resp)
}

// ----------------------------------------------------- Private V2 Endpoints --------------------------------------------------------------------------------

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
