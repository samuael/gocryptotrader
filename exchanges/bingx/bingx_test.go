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
