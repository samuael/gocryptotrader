package apexpro

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey    = ""
	apiSecret = ""
	clientID  = ""

	starkKey            = ""
	starkSecret         = ""
	starkKeyYCoordinate = ""

	ethereumAddress = ""

	canManipulateRealOrders = false
)

var ap = &Apexpro{}

func TestMain(m *testing.M) {
	ap.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Apexpro")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret
	exchCfg.API.Credentials.ClientID = clientID

	exchCfg.API.Credentials.L2Key = starkKey
	exchCfg.API.Credentials.L2Secret = starkSecret
	exchCfg.API.Credentials.L2KeyYCoordinate = starkKeyYCoordinate
	exchCfg.API.Credentials.Subaccount = ethereumAddress

	err = ap.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := ap.UpdateTradablePairs(context.Background(), true); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func (ap *Apexpro) isL2CredentialsProvided() bool {
	_, err := ap.GetCredentials(context.Background())
	if err != nil {
		return false
	}
	return true
}

// Implement tests for API endpoints below

func TestGetSystemTimeV3(t *testing.T) {
	t.Parallel()

	enabledPairs, err := ap.GetEnabledPairs(asset.Futures)
	require.NoError(t, err)
	println(strings.Join(enabledPairs.Strings(), ","))
	result, err := ap.GetSystemTimeV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSystemTimeV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetSystemTimeV2(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSystemTimeV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetSystemTimeV1(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllConfigDataV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetAllConfigDataV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllConfigDataV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetAllConfigDataV1(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV3(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV2(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV1(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetNewestTradingDataV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV3(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetNewestTradingDataV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV2(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}
func TestGetNewestTradingDataV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV1(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV3(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV3(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetCandlestickChartDataV3(context.Background(), "BTCUSDT", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV2(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetCandlestickChartDataV2(context.Background(), "BTCUSDT", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV1(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetCandlestickChartDataV1(context.Background(), "BTCUSDT", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTickerDataV3(t *testing.T) {
	t.Parallel()
	_, err := ap.GetTickerDataV3(context.Background(), "")
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetTickerDataV3(context.Background(), "BTCUSDT")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTickerDataV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetTickerDataV2(context.Background(), "")
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetTickerDataV2(context.Background(), "BTCUSDT")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingHistoryRate(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV3(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV3(context.Background(), "BTCUSDT", time.Time{}, time.Time{}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetFundingHistoryRateV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV2(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV2(context.Background(), "BTCUSDT", time.Time{}, time.Time{}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetFundingHistoryRateV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV1(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV1(context.Background(), "BTCUSDT", time.Time{}, time.Time{}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetAllConfigDataV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetAllConfigDataV2(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetCheckIfUserExistsV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetCheckIfUserExistsV2(context.Background(), "0x0330eBB5e894720e6746070371F9Fd797BE9D074")
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetCheckIfUserExistsV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetCheckIfUserExistsV1(context.Background(), "0x0330eBB5e894720e6746070371F9Fd797BE9D074")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWsConnect(t *testing.T) {
	t.Parallel()
	err := ap.WsConnect()
	require.NoError(t, err)
	time.Sleep(time.Second * 23)
}

func TestGenerateNonce(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GenerateNonceV3(context.Background(), starkKey, ethereumAddress, "9")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGenerateNonceV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GenerateNonceV2(context.Background(), starkKey, ethereumAddress, "9")
	require.NoError(t, err)
	assert.NotNil(t, result)
}
func TestGetUsersData(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUsersDataV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUsersDataV2GetUsersDataV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUsersDataV2(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestEditUserData(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	_, err := ap.EditUserDataV3(context.Background(), &EditUserDataParams{})
	require.ErrorIs(t, err, common.ErrNilPointer)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.EditUserDataV3(context.Background(), &EditUserDataParams{
		Email:                    "samuaeladnew@gmail.com",
		UserData:                 "",
		Username:                 "Thrasher",
		IsSharingUsername:        true,
		Country:                  "Ethiopia",
		EmailNotifyGeneralEnable: true,
		EmailNotifyTradingEnable: true,
		EmailNotifyAccountEnable: true,
		PopupNotifyTradingEnable: true,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}
func TestEditUserDataV2(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	_, err := ap.EditUserDataV2(context.Background(), &EditUserDataParams{})
	require.ErrorIs(t, err, common.ErrNilPointer)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.EditUserDataV2(context.Background(), &EditUserDataParams{
		Email:                    "samuaeladnew@gmail.com",
		UserData:                 "",
		Username:                 "Thrasher",
		IsSharingUsername:        true,
		Country:                  "Ethiopia",
		EmailNotifyGeneralEnable: true,
		EmailNotifyTradingEnable: true,
		EmailNotifyAccountEnable: true,
		PopupNotifyTradingEnable: true,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountData(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountDataV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountDataV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountDataV2(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountBalance(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountBalance(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountBalanceV2(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountBalanceV2(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserTransferDataV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserTransferDataV2(context.Background(), currency.USDT, time.Now().Add(-time.Hour*50), time.Now(), "DEPOSIT", nil, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserWithdrawalListV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserWithdrawalListV2(context.Background(), "WITHDRAWAL", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFastAndCrossChainWithdrawalFees(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	_, err := ap.GetFastAndCrossChainWithdrawalFees(context.Background(), 1, "1", currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
	_, err = ap.GetFastAndCrossChainWithdrawalFees(context.Background(), 1, "", currency.BTC)
	require.ErrorIs(t, err, errChainIDMissing)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFastAndCrossChainWithdrawalFees(context.Background(), 1.32, "1", currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAssetWithdrawalAndTransferLimit(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAssetWithdrawalAndTransferLimit(context.Background(), currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAssetWithdrawalAndTransferLimit(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserDepositWithdrawData(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserTransferData(context.Background(), 0, 10, "", "DEPOSIT", "", "", time.Now().Add(time.Hour*30), time.Now(), []string{"1"})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWithdrawalFees(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWithdrawalFees(context.Background(), 12, []string{"1"}, 140)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetContractAccountTransferLimits(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetContractAccountTransferLimits(context.Background(), currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTradeHistory(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetTradeHistory(context.Background(), "BTC-USD", order.Sell.String(), "LIMIT", time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTradeHistoryV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetTradeHistoryV2(context.Background(), "BTC-USD", order.Sell.String(), "LIMIT", currency.USDT, time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPrice(t *testing.T) {
	t.Parallel()
	_, err := ap.GetWorstPriceV3(context.Background(), "", "SELL", 1)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	_, err = ap.GetWorstPriceV3(context.Background(), "BTC-USDT", "", 1)
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = ap.GetWorstPriceV3(context.Background(), "BTC-USDT", "SELL", 0)
	require.ErrorIs(t, err, order.ErrAmountBelowMin)
}

func TestGetWorstPriceV3(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV3(context.Background(), "BTC-USDT", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPriceV2(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV2(context.Background(), "BTC-USDT", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrder(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.CancelPerpOrder(context.Background(), "123231")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrderV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.CancelPerpOrderV2(context.Background(), "123231", currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOpenOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOpenOrders(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOpenOrdersV2(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOpenOrdersV2(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllOrderHistory(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAllOrderHistory(context.Background(), "BTC-USDT", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllOrderHistoryV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAllOrderHistoryV2(context.Background(), currency.EMPTYCODE, "BTC-USDT", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAllOrderHistoryV2(context.Background(), currency.USDT, "BTC-USDT", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOrderID(t *testing.T) {
	t.Parallel()
	_, err := ap.GetOrderID(context.Background(), "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOrderID(context.Background(), "12343")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSingleOrderV2(t *testing.T) {
	t.Parallel()
	_, err := ap.getSingleOrderV2(context.Background(), "", "", currency.USDT)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)
	_, err = ap.getSingleOrderV2(context.Background(), "231232341", "", currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
}

func TestGetSingleOrderByOrderIDV2(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByOrderIDV2(context.Background(), "231232341", currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSingleOrderByClientOrderIDV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByClientOrderIDV2(context.Background(), "231232341", currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOrderClientOrderID(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	_, err := ap.GetOrderClientOrderID(context.Background(), "")
	require.ErrorIs(t, err, order.ErrClientOrderIDMustBeSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOrderClientOrderID(context.Background(), "12343")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingRate(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV3(context.Background(), "BTC-USDT", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingRateV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingRateV2(context.Background(), currency.EMPTYCODE, "BTC-USDT", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV2(context.Background(), currency.USDT, "BTC-USDT", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserHistorialProfitAndLoss(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLoss(context.Background(), "BTC-USDT", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserHistorialProfitAndLossV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetUserHistorialProfitAndLossV2(context.Background(), currency.EMPTYCODE, "BTC-USDT", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLossV2(context.Background(), currency.USDT, "BTC-USDT", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetYesterdaysPNL(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetYesterdaysPNL(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetYesterdaysPNLV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetYesterdaysPNLV2(context.Background(), currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetYesterdaysPNLV2(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetHistoricalAssetValue(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetHistoricalAssetValue(context.Background(), time.Now().Add(-time.Hour*50), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetHistoricalAssetValueV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetHistoricalAssetValueV2(context.Background(), currency.EMPTYCODE, time.Now().Add(-time.Hour*50), time.Now())
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetHistoricalAssetValueV2(context.Background(), currency.USDC, time.Now().Add(-time.Hour*50), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSetInitialMarginRateInfo(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfo(context.Background(), "BTC-USDT", 200)
	assert.NoError(t, err)
}

func TestSetInitialMarginRateInfoV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfoV2(context.Background(), "BTC-USDT", 200)
	assert.NoError(t, err)
}
