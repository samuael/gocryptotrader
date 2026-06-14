package bingx

import (
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/margin"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
	testexch "github.com/thrasher-corp/gocryptotrader/internal/testing/exchange"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

var (
	e           *Exchange
	btcUSDT     currency.Pair
	coinMBTCUSD currency.Pair
)

func TestMain(m *testing.M) {
	e = new(Exchange)
	if err := testexch.Setup(e); err != nil {
		log.Fatal(err)
	}

	if apiKey != "" && apiSecret != "" {
		e.API.AuthenticatedSupport = true
		e.API.AuthenticatedWebsocketSupport = true
		e.SetCredentials(apiKey, apiSecret, "", "", "", "")
		e.Websocket.SetCanUseAuthenticatedEndpoints(true)
	}

	btcUSDT = currency.NewBTCUSDT()
	btcUSDT.Delimiter = currency.DashDelimiter
	coinMBTCUSD = currency.NewPair(currency.BTC, currency.USD)
	coinMBTCUSD.Delimiter = currency.DashDelimiter
	os.Exit(m.Run())
}

func TestGetSpotSymbols(t *testing.T) {
	t.Parallel()
	result, err := e.GetSpotSymbols(t.Context(), currency.EMPTYPAIR)
	require.NoError(t, err)
	assert.NotNil(t, result)

	result, err = e.GetSpotSymbols(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketTrades(t *testing.T) {
	t.Parallel()
	_, err := e.GetMarketTrades(t.Context(), currency.EMPTYPAIR, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetMarketTrades(t.Context(), btcUSDT, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOrderbookDepth(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotOrderbookDepth(t.Context(), currency.EMPTYPAIR, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSpotOrderbookDepth(t.Context(), btcUSDT, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotKlineData(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotKlineData(t.Context(), currency.EMPTYPAIR, kline.OneMin, time.Time{}, time.Time{}, 1, 100)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSpotKlineData(t.Context(), btcUSDT, kline.OneMin, time.Time{}, time.Time{}, 0, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGet24HrTickerPriceChange(t *testing.T) {
	t.Parallel()
	result, err := e.Get24HrTickerPriceChange(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOrderbookAggregation(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotOrderbookAggregation(t.Context(), currency.EMPTYPAIR, 20, "step0")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSpotOrderbookAggregation(t.Context(), btcUSDT, 20, "step0")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotSymbolPriceTicker(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotSymbolPriceTicker(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSpotSymbolPriceTicker(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotSymbolOrderbookTicker(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotSymbolOrderbookTicker(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSpotSymbolOrderbookTicker(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetHistoricalKlineData(t *testing.T) {
	t.Parallel()
	_, err := e.GetHistoricalKlineData(t.Context(), currency.EMPTYPAIR, kline.OneMin, time.Time{}, time.Time{}, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetHistoricalKlineData(t.Context(), btcUSDT, kline.OneMin, time.Time{}, time.Time{}, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOldSpotTrades(t *testing.T) {
	t.Parallel()
	_, err := e.GetOldSpotTrades(t.Context(), currency.EMPTYPAIR, 5, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetOldSpotTrades(t.Context(), btcUSDT, 5, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapContracts(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwapContracts(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetSwapOrderbookDepth(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapOrderbookDepth(t.Context(), currency.EMPTYPAIR, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSwapOrderbookDepth(t.Context(), btcUSDT, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapRecentTrades(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapRecentTrades(t.Context(), currency.EMPTYPAIR, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSwapRecentTrades(t.Context(), btcUSDT, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapMarkPriceAndFundingRate(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwapMarkPriceAndFundingRate(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	result, err = e.GetSwapMarkPriceAndFundingRate(t.Context(), currency.EMPTYPAIR)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetSwapFundingRateHistory(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwapFundingRateHistory(t.Context(), btcUSDT, time.Time{}, time.Time{}, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapKlineData(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapKlineData(t.Context(), currency.EMPTYPAIR, kline.OneHour, time.Time{}, time.Time{}, 0, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	_, err = e.GetSwapKlineData(t.Context(), btcUSDT, kline.Interval(0), time.Time{}, time.Time{}, 0, 5)
	require.ErrorIs(t, err, kline.ErrInvalidInterval)

	result, err := e.GetSwapKlineData(t.Context(), btcUSDT, kline.OneHour, time.Time{}, time.Time{}, 0, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapOpenInterest(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapOpenInterest(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSwapOpenInterest(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwap24HrTickerPriceChange(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwap24HrTickerPriceChange(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetSwapHistoricalTrades(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwapHistoricalTrades(t.Context(), btcUSDT, "", 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapSymbolOrderbookTicker(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapSymbolOrderbookTicker(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSwapSymbolOrderbookTicker(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapMarkPriceKlineData(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapMarkPriceKlineData(t.Context(), currency.EMPTYPAIR, kline.OneHour, time.Time{}, time.Time{}, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	_, err = e.GetSwapMarkPriceKlineData(t.Context(), btcUSDT, kline.Interval(0), time.Time{}, time.Time{}, 5)
	require.ErrorIs(t, err, kline.ErrInvalidInterval)

	result, err := e.GetSwapMarkPriceKlineData(t.Context(), btcUSDT, kline.OneHour, time.Time{}, time.Time{}, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapSymbolPriceTicker(t *testing.T) {
	t.Parallel()
	result, err := e.GetSwapSymbolPriceTicker(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetSwapTradingRules(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapTradingRules(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetSwapTradingRules(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMContracts(t *testing.T) {
	t.Parallel()
	result, err := e.GetCoinMContracts(t.Context(), currency.EMPTYPAIR)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	result, err = e.GetCoinMContracts(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetCoinMMarkPriceAndFundingRate(t *testing.T) {
	t.Parallel()
	result, err := e.GetCoinMMarkPriceAndFundingRate(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	result, err = e.GetCoinMMarkPriceAndFundingRate(t.Context(), currency.EMPTYPAIR)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetCoinMOpenInterest(t *testing.T) {
	t.Parallel()
	result, err := e.GetCoinMOpenInterest(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestGetCoinMKlineData(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMKlineData(t.Context(), currency.EMPTYPAIR, kline.OneHour, time.Time{}, time.Time{}, 0, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	_, err = e.GetCoinMKlineData(t.Context(), coinMBTCUSD, kline.Interval(0), time.Time{}, time.Time{}, 0, 5)
	require.ErrorIs(t, err, kline.ErrInvalidInterval)

	result, err := e.GetCoinMKlineData(t.Context(), coinMBTCUSD, kline.OneHour, time.Time{}, time.Time{}, 0, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMOrderbookDepth(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMOrderbookDepth(t.Context(), currency.EMPTYPAIR, 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	result, err := e.GetCoinMOrderbookDepth(t.Context(), coinMBTCUSD, 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinM24HrTickerPriceChange(t *testing.T) {
	t.Parallel()
	result, err := e.GetCoinM24HrTickerPriceChange(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestSignaturePayload(t *testing.T) {
	t.Parallel()
	params := url.Values{}
	params.Set("symbol", "BTC-USDT")
	params.Set("recvWindow", "0")
	params.Set("timestamp", "1696751141337")
	signingString := signaturePayload(params)
	assert.Equal(t, "recvWindow=0&symbol=BTC-USDT&timestamp=1696751141337", signingString, "signaturePayload should sort and concatenate parameters")

	hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256, []byte(signingString), []byte("SECRET_KEY"))
	require.NoError(t, err, "GetHMAC must not error")
	assert.Equal(t, "fe041f159118c90ac13eab4d32f9e2d75b80ca6fe17ca8acd290aba864753ce2", hex.EncodeToString(hmacSigned), "signature should match the documented HMAC-SHA256 value")
}

func TestTestSwapOrder(t *testing.T) {
	t.Parallel()
	_, err := e.TestSwapOrder(t.Context(), &PlaceSwapOrderRequest{})
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.TestSwapOrder(t.Context(), &PlaceSwapOrderRequest{Symbol: btcUSDT})
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)
	_, err = e.TestSwapOrder(t.Context(), &PlaceSwapOrderRequest{Symbol: btcUSDT, Type: "MARKET"})
	require.ErrorIs(t, err, order.ErrSideIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.TestSwapOrder(t.Context(), &PlaceSwapOrderRequest{Symbol: btcUSDT, Type: "MARKET", Side: order.Buy.String(), PositionSide: "LONG", Quantity: 1})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceSwapOrder(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceSwapOrder(t.Context(), &PlaceSwapOrderRequest{Symbol: btcUSDT})
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceSwapOrder(t.Context(), &PlaceSwapOrderRequest{Symbol: btcUSDT, Type: "LIMIT", Side: order.Buy.String(), PositionSide: "LONG", Price: 10000, Quantity: 1})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestModifySwapOrder(t *testing.T) {
	t.Parallel()
	_, err := e.ModifySwapOrder(t.Context(), currency.EMPTYPAIR, "123", "", 1)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.ModifySwapOrder(t.Context(), btcUSDT, "", "", 1)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)
	_, err = e.ModifySwapOrder(t.Context(), btcUSDT, "123", "", 0)
	require.ErrorIs(t, err, order.ErrAmountIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.ModifySwapOrder(t.Context(), btcUSDT, "123456", "", 2)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceMultipleSwapOrders(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceMultipleSwapOrders(t.Context(), nil)
	require.ErrorIs(t, err, errBatchOrdersRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceMultipleSwapOrders(t.Context(), []PlaceSwapOrderRequest{{Symbol: btcUSDT, Type: "MARKET", Side: order.Buy.String(), PositionSide: "LONG", Quantity: 1}})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCloseAllSwapPositions(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CloseAllSwapPositions(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelSwapOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CancelSwapOrder(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CancelSwapOrder(t.Context(), btcUSDT, 0, "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelSwapOrder(t.Context(), btcUSDT, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelMultipleSwapOrders(t *testing.T) {
	t.Parallel()
	_, err := e.CancelMultipleSwapOrders(t.Context(), currency.EMPTYPAIR, []int64{1}, nil)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CancelMultipleSwapOrders(t.Context(), btcUSDT, nil, nil)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelMultipleSwapOrders(t.Context(), btcUSDT, []int64{123456}, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelAllSwapOpenOrders(t *testing.T) {
	t.Parallel()
	_, err := e.CancelAllSwapOpenOrders(t.Context(), btcUSDT, "")
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelAllSwapOpenOrders(t.Context(), btcUSDT, "LIMIT")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapOpenOrders(t.Context(), btcUSDT, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapPendingOrderStatus(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapPendingOrderStatus(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapPendingOrderStatus(t.Context(), btcUSDT, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapOrderDetails(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapOrderDetails(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapOrderDetails(t.Context(), btcUSDT, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapMarginType(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapMarginType(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapMarginType(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestChangeSwapMarginType(t *testing.T) {
	t.Parallel()
	err := e.ChangeSwapMarginType(t.Context(), currency.EMPTYPAIR, "CROSSED")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	err = e.ChangeSwapMarginType(t.Context(), btcUSDT, "")
	require.ErrorIs(t, err, margin.ErrInvalidMarginType)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	err = e.ChangeSwapMarginType(t.Context(), btcUSDT, "CROSSED")
	require.NoError(t, err)
}

func TestGetSwapLeverage(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapLeverage(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapLeverage(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetSwapLeverage(t *testing.T) {
	t.Parallel()
	_, err := e.SetSwapLeverage(t.Context(), currency.EMPTYPAIR, "LONG", 5)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.SetSwapLeverage(t.Context(), btcUSDT, "", 5)
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = e.SetSwapLeverage(t.Context(), btcUSDT, "LONG", 0)
	require.ErrorIs(t, err, errLeverageRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetSwapLeverage(t.Context(), btcUSDT, "LONG", 5)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapForceOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapForceOrders(t.Context(), btcUSDT, "", "", time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapOrderHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapOrderHistory(t.Context(), btcUSDT, "", 0, time.Time{}, time.Time{}, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestModifyIsolatedPositionMargin(t *testing.T) {
	t.Parallel()
	_, err := e.ModifyIsolatedPositionMargin(t.Context(), currency.EMPTYPAIR, 1, "1", "LONG", 0)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.ModifyIsolatedPositionMargin(t.Context(), btcUSDT, 0, "1", "LONG", 0)
	require.ErrorIs(t, err, errMarginAmountRequired)
	_, err = e.ModifyIsolatedPositionMargin(t.Context(), btcUSDT, 1, "", "LONG", 0)
	require.ErrorIs(t, err, errAdjustmentTypeRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.ModifyIsolatedPositionMargin(t.Context(), btcUSDT, 3, "1", "LONG", 0)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapHistoricalTransactionOrders(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapHistoricalTransactionOrders(t.Context(), btcUSDT, "", "", 0, time.Now().Add(-time.Hour), time.Now())
	require.ErrorIs(t, err, errTradingUnitRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapHistoricalTransactionOrders(t.Context(), btcUSDT, "", "COIN", 0, time.Now().Add(-time.Hour*24), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetSwapPositionMode(t *testing.T) {
	t.Parallel()
	_, err := e.SetSwapPositionMode(t.Context(), "")
	require.ErrorIs(t, err, errPositionModeRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetSwapPositionMode(t.Context(), "true")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapPositionMode(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapPositionMode(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelReplaceSwapOrder(t *testing.T) {
	t.Parallel()
	arg := &PlaceSwapOrderRequest{Symbol: btcUSDT, Type: "LIMIT", Side: order.Buy.String(), PositionSide: "LONG", Price: 10000, Quantity: 1}
	_, err := e.CancelReplaceSwapOrder(t.Context(), "", "", 1, "", arg)
	require.ErrorIs(t, err, errCancelReplaceModeEmpty)
	_, err = e.CancelReplaceSwapOrder(t.Context(), "STOP_ON_FAILURE", "", 0, "", arg)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelReplaceSwapOrder(t.Context(), "STOP_ON_FAILURE", "", 123456, "", arg)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBatchCancelReplaceSwapOrders(t *testing.T) {
	t.Parallel()
	_, err := e.BatchCancelReplaceSwapOrders(t.Context(), "")
	require.ErrorIs(t, err, errBatchOrdersRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.BatchCancelReplaceSwapOrders(t.Context(), `[{"cancelReplaceMode":"STOP_ON_FAILURE","cancelOrderId":123456,"symbol":"BTC-USDT","type":"LIMIT","side":"BUY","positionSide":"LONG","price":10000,"quantity":1}]`)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetSwapCancelAllAfter(t *testing.T) {
	t.Parallel()
	_, err := e.SetSwapCancelAllAfter(t.Context(), "", 10)
	require.ErrorIs(t, err, errCountdownTypeRequired)
	_, err = e.SetSwapCancelAllAfter(t.Context(), "ACTIVATE", 5)
	require.ErrorIs(t, err, errCountdownTimeoutInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetSwapCancelAllAfter(t.Context(), "ACTIVATE", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCloseSwapPositionByID(t *testing.T) {
	t.Parallel()
	_, err := e.CloseSwapPositionByID(t.Context(), "")
	require.ErrorIs(t, err, errPositionIDRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CloseSwapPositionByID(t.Context(), "1769649551460794368")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllSwapOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetAllSwapOrders(t.Context(), btcUSDT, 0, time.Time{}, time.Time{}, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapMaintMarginRatio(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapMaintMarginRatio(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapMaintMarginRatio(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapTransactionDetails(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapTransactionDetails(t.Context(), btcUSDT, "", 0, 0, 1, 50, time.Now().Add(-time.Hour*24), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapPositionHistory(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapPositionHistory(t.Context(), currency.EMPTYPAIR, "", 0, time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapPositionHistory(t.Context(), btcUSDT, "", 0, time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSwapIsolatedMarginChangeHistory(t *testing.T) {
	t.Parallel()
	_, err := e.GetSwapIsolatedMarginChangeHistory(t.Context(), currency.EMPTYPAIR, "1", time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.GetSwapIsolatedMarginChangeHistory(t.Context(), btcUSDT, "", time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.ErrorIs(t, err, errPositionIDRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSwapIsolatedMarginChangeHistory(t.Context(), btcUSDT, "1847596444958068736", time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestApplyVST(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.ApplyVST(t.Context(), "0", 500000)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceTWAPOrder(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceTWAPOrder(t.Context(), &PlaceTWAPOrderRequest{})
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.PlaceTWAPOrder(t.Context(), &PlaceTWAPOrderRequest{Symbol: btcUSDT})
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = e.PlaceTWAPOrder(t.Context(), &PlaceTWAPOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), PositionSide: "LONG"})
	require.ErrorIs(t, err, order.ErrUnknownPriceType)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceTWAPOrder(t.Context(), &PlaceTWAPOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), PositionSide: "LONG", PriceType: "constant", PriceVariance: "2000", TriggerPrice: "68000", Interval: 8, AmountPerOrder: "0.1", TotalAmount: "0.5"})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTWAPOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetTWAPOpenOrders(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTWAPHistoricalOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetTWAPHistoricalOrders(t.Context(), btcUSDT, time.Now().Add(-time.Hour*24), time.Now(), 1, 50)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTWAPOrderDetails(t *testing.T) {
	t.Parallel()
	_, err := e.GetTWAPOrderDetails(t.Context(), "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetTWAPOrderDetails(t.Context(), "123456")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelTWAPOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CancelTWAPOrder(t.Context(), "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelTWAPOrder(t.Context(), "123456")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSwitchMultiAssetsMode(t *testing.T) {
	t.Parallel()
	_, err := e.SwitchMultiAssetsMode(t.Context(), "")
	require.ErrorIs(t, err, errAssetModeRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SwitchMultiAssetsMode(t.Context(), "multiAssetsMode")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMultiAssetsMode(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetMultiAssetsMode(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMultiAssetsRules(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetMultiAssetsRules(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMultiAssetsMargin(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetMultiAssetsMargin(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestOneClickReversePosition(t *testing.T) {
	t.Parallel()
	_, err := e.OneClickReversePosition(t.Context(), "", btcUSDT, "", "")
	require.ErrorIs(t, err, errReverseTypeRequired)
	_, err = e.OneClickReversePosition(t.Context(), "Reverse", currency.EMPTYPAIR, "", "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.OneClickReversePosition(t.Context(), "Reverse", btcUSDT, "", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetHedgeModeAutoAddMargin(t *testing.T) {
	t.Parallel()
	_, err := e.SetHedgeModeAutoAddMargin(t.Context(), currency.EMPTYPAIR, 1, "true", 3)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.SetHedgeModeAutoAddMargin(t.Context(), btcUSDT, 0, "true", 3)
	require.ErrorIs(t, err, errPositionIDRequired)
	_, err = e.SetHedgeModeAutoAddMargin(t.Context(), btcUSDT, 1, "", 3)
	require.ErrorIs(t, err, errFunctionSwitchRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetHedgeModeAutoAddMargin(t.Context(), btcUSDT, 1, "true", 3)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceSpotOrder(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceSpotOrder(t.Context(), &PlaceSpotOrderRequest{})
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.PlaceSpotOrder(t.Context(), &PlaceSpotOrderRequest{Symbol: btcUSDT})
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = e.PlaceSpotOrder(t.Context(), &PlaceSpotOrderRequest{Symbol: btcUSDT, Side: order.Buy.String()})
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceSpotOrder(t.Context(), &PlaceSpotOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), OrderType: "MARKET", Quantity: 0.001})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceMultipleSpotOrders(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceMultipleSpotOrders(t.Context(), nil, false)
	require.ErrorIs(t, err, errBatchOrdersRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceMultipleSpotOrders(t.Context(), []PlaceSpotOrderRequest{{Symbol: btcUSDT, Side: order.Buy.String(), OrderType: "LIMIT", Quantity: 0.001, Price: 10000}}, false)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelSpotOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CancelSpotOrder(t.Context(), currency.EMPTYPAIR, 1, "", "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CancelSpotOrder(t.Context(), btcUSDT, 0, "", "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelSpotOrder(t.Context(), btcUSDT, 123456, "", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelMultipleSpotOrders(t *testing.T) {
	t.Parallel()
	_, err := e.CancelMultipleSpotOrders(t.Context(), currency.EMPTYPAIR, 0, "1,2", "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CancelMultipleSpotOrders(t.Context(), btcUSDT, 0, "", "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelMultipleSpotOrders(t.Context(), btcUSDT, 0, "123456,234567", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelAllSpotOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelAllSpotOpenOrders(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelReplaceSpotOrder(t *testing.T) {
	t.Parallel()
	arg := &PlaceSpotOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), OrderType: "LIMIT", Quantity: 0.001, Price: 10000}
	_, err := e.CancelReplaceSpotOrder(t.Context(), "", 1, "", "", arg)
	require.ErrorIs(t, err, errCancelReplaceModeEmpty)
	_, err = e.CancelReplaceSpotOrder(t.Context(), "STOP_ON_FAILURE", 0, "", "", arg)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelReplaceSpotOrder(t.Context(), "STOP_ON_FAILURE", 123456, "", "", arg)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOrderDetails(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotOrderDetails(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.GetSpotOrderDetails(t.Context(), btcUSDT, 0, "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOrderDetails(t.Context(), btcUSDT, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOpenOrders(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOrderHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOrderHistory(t.Context(), btcUSDT, 0, "", "", time.Time{}, time.Time{}, 1, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotTransactionDetails(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotTransactionDetails(t.Context(), btcUSDT, 0, time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotCommissionRate(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotCommissionRate(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotCommissionRate(t.Context(), btcUSDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetSpotCancelAllAfter(t *testing.T) {
	t.Parallel()
	_, err := e.SetSpotCancelAllAfter(t.Context(), "", 10)
	require.ErrorIs(t, err, errCountdownTypeRequired)
	_, err = e.SetSpotCancelAllAfter(t.Context(), "ACTIVATE", 5)
	require.ErrorIs(t, err, errCountdownTimeoutInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetSpotCancelAllAfter(t.Context(), "ACTIVATE", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateSpotOCOOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{})
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT})
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT, Side: order.Buy.String()})
	require.ErrorIs(t, err, order.ErrAmountIsInvalid)
	_, err = e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), Quantity: 0.001})
	require.ErrorIs(t, err, order.ErrPriceMustBeSetIfLimitOrder)
	_, err = e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), Quantity: 0.001, LimitPrice: 48000})
	require.ErrorIs(t, err, order.ErrPriceMustBeSetIfLimitOrder)
	_, err = e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), Quantity: 0.001, LimitPrice: 48000, OrderPrice: 88000})
	require.ErrorIs(t, err, errTriggerPriceRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CreateSpotOCOOrder(t.Context(), &CreateSpotOCOOrderRequest{Symbol: btcUSDT, Side: order.Buy.String(), Quantity: 0.001, LimitPrice: 48000, OrderPrice: 88000, TriggerPrice: 87000})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelSpotOCOOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CancelSpotOCOOrder(t.Context(), "", "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelSpotOCOOrder(t.Context(), "1827980248763858944", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOCOOrderList(t *testing.T) {
	t.Parallel()
	_, err := e.GetSpotOCOOrderList(t.Context(), "", "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOCOOrderList(t.Context(), "1827968196914479104", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOpenOCOOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOpenOCOOrders(t.Context(), 1, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSpotOCOHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetSpotOCOHistory(t.Context(), time.Now().Add(-time.Hour*24), time.Now(), 1, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPlaceCoinMOrder(t *testing.T) {
	t.Parallel()
	_, err := e.PlaceCoinMOrder(t.Context(), &PlaceCoinMOrderRequest{})
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.PlaceCoinMOrder(t.Context(), &PlaceCoinMOrderRequest{Symbol: coinMBTCUSD})
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)
	_, err = e.PlaceCoinMOrder(t.Context(), &PlaceCoinMOrderRequest{Symbol: coinMBTCUSD, Type: "MARKET"})
	require.ErrorIs(t, err, order.ErrSideIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.PlaceCoinMOrder(t.Context(), &PlaceCoinMOrderRequest{Symbol: coinMBTCUSD, Type: "MARKET", Side: order.Buy.String(), PositionSide: "LONG", Quantity: 1})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMCommissionRate(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMCommissionRate(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMLeverage(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMLeverage(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMLeverage(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetCoinMLeverage(t *testing.T) {
	t.Parallel()
	_, err := e.SetCoinMLeverage(t.Context(), currency.EMPTYPAIR, "LONG", 4)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.SetCoinMLeverage(t.Context(), coinMBTCUSD, "", 4)
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = e.SetCoinMLeverage(t.Context(), coinMBTCUSD, "LONG", 0)
	require.ErrorIs(t, err, errLeverageRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.SetCoinMLeverage(t.Context(), coinMBTCUSD, "LONG", 4)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelAllCoinMOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelAllCoinMOrders(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCloseAllCoinMPositions(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CloseAllCoinMPositions(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMPositions(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMPositions(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMAccountAssets(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMAccountAssets(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMForceOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMForceOrders(t.Context(), coinMBTCUSD, "", time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMOrderTradeDetail(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMOrderTradeDetail(t.Context(), "", 1, 100)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMOrderTradeDetail(t.Context(), "1796163365782945792", 1, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelCoinMOrder(t *testing.T) {
	t.Parallel()
	_, err := e.CancelCoinMOrder(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.CancelCoinMOrder(t.Context(), coinMBTCUSD, 0, "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	result, err := e.CancelCoinMOrder(t.Context(), coinMBTCUSD, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMOpenOrders(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMOrderDetail(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMOrderDetail(t.Context(), currency.EMPTYPAIR, 1, "")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	_, err = e.GetCoinMOrderDetail(t.Context(), coinMBTCUSD, 0, "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMOrderDetail(t.Context(), coinMBTCUSD, 123456, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMOrderHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMOrderHistory(t.Context(), coinMBTCUSD, 0, time.Time{}, time.Time{}, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCoinMMarginType(t *testing.T) {
	t.Parallel()
	_, err := e.GetCoinMMarginType(t.Context(), currency.EMPTYPAIR)
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e)
	result, err := e.GetCoinMMarginType(t.Context(), coinMBTCUSD)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetCoinMMarginType(t *testing.T) {
	t.Parallel()
	err := e.SetCoinMMarginType(t.Context(), currency.EMPTYPAIR, "ISOLATED")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	err = e.SetCoinMMarginType(t.Context(), coinMBTCUSD, "")
	require.ErrorIs(t, err, margin.ErrInvalidMarginType)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	err = e.SetCoinMMarginType(t.Context(), coinMBTCUSD, "ISOLATED")
	require.NoError(t, err)
}

func TestAdjustCoinMIsolatedMargin(t *testing.T) {
	t.Parallel()
	err := e.AdjustCoinMIsolatedMargin(t.Context(), currency.EMPTYPAIR, 0.01, 1, "LONG")
	require.ErrorIs(t, err, currency.ErrCurrencyPairEmpty)
	err = e.AdjustCoinMIsolatedMargin(t.Context(), coinMBTCUSD, 0, 1, "LONG")
	require.ErrorIs(t, err, errMarginAmountRequired)
	err = e.AdjustCoinMIsolatedMargin(t.Context(), coinMBTCUSD, 0.01, 0, "LONG")
	require.ErrorIs(t, err, errAdjustmentTypeRequired)
	err = e.AdjustCoinMIsolatedMargin(t.Context(), coinMBTCUSD, 0.01, 1, "")
	require.ErrorIs(t, err, order.ErrSideIsInvalid)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, e, canManipulateRealOrders)
	err = e.AdjustCoinMIsolatedMargin(t.Context(), coinMBTCUSD, 0.01, 1, "LONG")
	require.NoError(t, err)
}
