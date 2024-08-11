package apexpro

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	// Account ID: 603545650545558021
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

func TestGetSystemTimeV3(t *testing.T) {
	t.Parallel()
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

func TestGetAllSymbolsConfigDataV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetAllSymbolsConfigDataV1(context.Background())
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
	// ap.Verbose = true
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

func TestGenerateNonceV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GenerateNonceV1(context.Background(), starkKey, ethereumAddress, "9")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUsersData(t *testing.T) {
	t.Parallel()
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

func TestGetUsersDataV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUsersDataV1(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestEditUserData(t *testing.T) {
	t.Parallel()
	_, err := ap.EditUserDataV3(context.Background(), &EditUserDataParams{})
	require.ErrorIs(t, err, common.ErrNilPointer)

	// sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.EditUserDataV3(context.Background(), &EditUserDataParams{
		Email:                    "someone@thrasher.io",
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

func TestGetUserAccountDataV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountDataV1(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountBalance(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountBalance(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountBalanceV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserAccountBalanceV2(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserAccountBalanceV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetUserAccountBalanceV1(context.Background())
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

func TestGetUserTransferDataV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserTransferDataV1(context.Background(), currency.USDT, time.Time{}, time.Time{}, "", nil, 0, 10)
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

func TestGetUserWithdrawalListV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserWithdrawalListV1(context.Background(), "WITHDRAWAL", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFastAndCrossChainWithdrawalFees(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFastAndCrossChainWithdrawalFeesV2(context.Background(), 1, "1", currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
	_, err = ap.GetFastAndCrossChainWithdrawalFeesV2(context.Background(), 1, "", currency.BTC)
	require.ErrorIs(t, err, errChainIDMissing)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFastAndCrossChainWithdrawalFeesV2(context.Background(), 1.32, "1", currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFastAndCrossChainWithdrawalFeesV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFastAndCrossChainWithdrawalFeesV1(context.Background(), 1.32, "1")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAssetWithdrawalAndTransferLimit(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAssetWithdrawalAndTransferLimitV2(context.Background(), currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAssetWithdrawalAndTransferLimitV2(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAssetWithdrawalAndTransferLimitV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAssetWithdrawalAndTransferLimitV1(context.Background(), currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAssetWithdrawalAndTransferLimitV1(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserDepositWithdrawData(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserTransferData(context.Background(), 0, 10, "", "DEPOSIT", "", "", time.Now().Add(time.Hour*30), time.Now(), []string{"1"})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWithdrawalFees(t *testing.T) {
	t.Parallel()
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

func TestGetTradeHistoryV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetTradeHistoryV1(context.Background(), "BTC-USD", order.Sell.String(), "LIMIT", time.Time{}, time.Time{}, 0, 10)
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
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV3(context.Background(), "BTC-USDT", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPriceV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV2(context.Background(), "BTC-USDT", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPriceV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV1(context.Background(), "BTC-USDT", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	futuresTradablePair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		require.NoError(t, err)
		require.NotNil(t, ap.UserAccountDetail)
	}
	takerFeeRate := ap.UserAccountDetail.ContractAccount.TakerFeeRate.Float64()
	result, err := ap.CreateOrderV3(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Sell.String(),
		OrderType:       "LIMIT",
		Size:            123,
		Price:           1,
		LimitFee:        takerFeeRate * 123 * 1,
		ExpirationTime:  time.Now().Add(time.Hour * 240),
		TimeInForce:     "GTC",
		TriggerPrice:    0,
		TrailingPercent: 1,
		ClientOrderID:   2312312312,
		ReduceOnly:      true,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrder(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CancelPerpOrder(context.Background(), "123231")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrderByClientOrderID(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CancelPerpOrderByClientOrderID(context.Background(), "2312312")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelAllOpenOrdersV3(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.CancelAllOpenOrdersV3(context.Background(), []string{"BTC-USDC"})
	assert.NoError(t, err)
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
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOpenOrdersV2(context.Background(), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOpenOrdersV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOpenOrdersV1(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllOrderHistory(t *testing.T) {
	t.Parallel()
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

func TestGetAllOrderHistoryV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAllOrderHistoryV1(context.Background(), "BTC-USDT", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
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
	_, err := ap.getSingleOrderV2(context.Background(), "", "", "USDT", exchange.RestSpot)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)
	_, err = ap.getSingleOrderV2(context.Background(), "231232341", "", "", exchange.RestSpot)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
}

func TestGetSingleOrderByOrderIDV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByOrderIDV2(context.Background(), "231232341", currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSingleOrderByClientOrderIDV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetSingleOrderByClientOrderIDV2(context.Background(), "231232341", currency.EMPTYCODE)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
	_, err = ap.GetSingleOrderByClientOrderIDV2(context.Background(), "", currency.USDT)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByClientOrderIDV2(context.Background(), "231232341", currency.USDT)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSingleOrderByOrderIDV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByOrderIDV1(context.Background(), "231232341")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetSingleOrderByClientOrderIDV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetSingleOrderByClientOrderIDV1(context.Background(), "")
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetSingleOrderByClientOrderIDV1(context.Background(), "231232341")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetVerificationEmailLink(t *testing.T) {
	t.Parallel()
	err := ap.GetVerificationEmailLink(context.Background(), "", currency.USDC)
	require.ErrorIs(t, err, errUserIDRequired)
	err = ap.GetVerificationEmailLink(context.Background(), "123123", currency.USDC)
	assert.NoError(t, err)
}

func TestLinkDevice(t *testing.T) {
	t.Parallel()
	err := ap.LinkDevice(context.Background(), currency.EMPTYCODE, "1")
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
	err = ap.LinkDevice(context.Background(), currency.USDT, "")
	require.ErrorIs(t, err, errDeviceTypeIsRequired)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err = ap.LinkDevice(context.Background(), currency.USDT, "2")
	require.NoError(t, err)
}

func TestGetOrderClientOrderID(t *testing.T) {
	t.Parallel()
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

func TestGetFundingRateV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV1(context.Background(), "BTC-USDT", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
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

func TestGetUserHistorialProfitAndLossV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLossV1(context.Background(), "BTC-USDT", "LONG", time.Time{}, time.Time{}, 0, 100)
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

func TestGetYesterdaysPNLV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetYesterdaysPNLV1(context.Background())
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

func TestGetHistoricalAssetValueV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetHistoricalAssetValueV1(context.Background(), time.Now().Add(-time.Hour*50), time.Now())
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

func TestSetInitialMarginRateInfoV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfoV1(context.Background(), "BTC-USDT", 200)
	assert.NoError(t, err)
}

func TestSetInitialMarginRateInfoV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfoV2(context.Background(), "BTC-USDT", 200)
	assert.NoError(t, err)
}

func TestWithdrawAsset(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.WithdrawAsset(context.Background(), &AssetWithdrawalParams{
		Amount:           1,
		ClientWithdrawID: "123123",
		Timestamp:        time.Now(),
		EthereumAddress:  ethereumAddress,
		L2Key:            starkKey,
		ToChainID:        "3",
		L2SourceTokenID:  currency.USDC,
		L1TargetTokenID:  currency.USDC,
		IsFastWithdraw:   false,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUserWithdrawalV2(t *testing.T) {
	t.Parallel()
	// sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.UserWithdrawalV2(context.Background(), 1, "1231231", time.Now().Add(time.Hour*24), currency.USDC)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWithdrawalToAddressV2(t *testing.T) {
	t.Parallel()
	_, err := ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{})
	require.ErrorIs(t, err, order.ErrAmountBelowMin)
	_, err = ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{})
	require.ErrorIs(t, err, order.ErrClientOrderIDMustBeSet)
	_, err = ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{})
	require.ErrorIs(t, err, errExpirationTimeRequired)
	_, err = ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{})
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)
	_, err = ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{})
	require.ErrorIs(t, err, errEthereumAddressMissing)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.WithdrawalToAddressV2(context.Background(), &WithdrawalToAddressParams{
		Amount:          1,
		ClientID:        "12334",
		ExpirationTime:  time.Now().Add(time.Hour * 50),
		Asset:           currency.BTC,
		EthereumAddress: ethereumAddress,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWithdrawalToAddressV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.WithdrawalToAddressV1(context.Background(), &WithdrawalToAddressParams{
		Amount:          1,
		ClientID:        "12334",
		ExpirationTime:  time.Now().Add(time.Hour * 50),
		Asset:           currency.BTC,
		EthereumAddress: ethereumAddress,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestOrderCreationParamsFilter(t *testing.T) {
	t.Parallel()
	_, err := ap.orderCreationParamsFilter(context.Background(), nil)
	require.ErrorIs(t, err, order.ErrOrderDetailIsNil)
	_, err = ap.orderCreationParamsFilter(context.Background(), &CreateOrderParams{Side: order.Buy.String()})
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	futuresTradablePair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)
	arg := &CreateOrderParams{Symbol: futuresTradablePair}
	_, err = ap.orderCreationParamsFilter(context.Background(), arg)
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	arg.Side = order.Buy.String()
	_, err = ap.orderCreationParamsFilter(context.Background(), &CreateOrderParams{Symbol: futuresTradablePair, Side: order.Buy.String()})
	require.ErrorIs(t, err, order.ErrTypeIsInvalid)
	arg.OrderType = order.Limit.String()
	_, err = ap.orderCreationParamsFilter(context.Background(), arg)
	require.ErrorIs(t, err, order.ErrAmountBelowMin)
	arg.Size = 2
	_, err = ap.orderCreationParamsFilter(context.Background(), arg)
	require.ErrorIs(t, err, order.ErrPriceBelowMin)
	arg.Price = 123
	arg.LimitFee = -1
	_, err = ap.orderCreationParamsFilter(context.Background(), arg)
	require.ErrorIs(t, err, errLimitFeeRequired)
	arg.LimitFee = 0.003
	_, err = ap.orderCreationParamsFilter(context.Background(), arg)
	require.ErrorIs(t, err, errExpirationTimeRequired)
}

func TestCreateOrderV1(t *testing.T) {
	t.Parallel()
	futuresTradablePair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		require.NoError(t, err)
		require.NotNil(t, ap.UserAccountDetail)
	}
	takerFeeRate := ap.UserAccountDetail.ContractAccount.TakerFeeRate.Float64()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CreateOrderV1(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Sell.String(),
		OrderType:       "LIMIT",
		Size:            123,
		Price:           1,
		LimitFee:        takerFeeRate * 123 * 1,
		ExpirationTime:  time.Now().Add(time.Hour * 240),
		TimeInForce:     "GTC",
		TriggerPrice:    0,
		TrailingPercent: 1,
		ClientOrderID:   2312312312,
		ReduceOnly:      true,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateOrderV2(t *testing.T) {
	t.Parallel()
	futuresTradablePair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		require.NoError(t, err)
		require.NotNil(t, ap.UserAccountDetail)
	}
	takerFeeRate := ap.UserAccountDetail.ContractAccount.TakerFeeRate.Float64()
	// sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CreateOrderV2(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Sell.String(),
		OrderType:       "LIMIT",
		Size:            123,
		Price:           1,
		LimitFee:        takerFeeRate * 123 * 1,
		ExpirationTime:  time.Now().Add(time.Hour * 240),
		TimeInForce:     "GTC",
		TriggerPrice:    0,
		TrailingPercent: 1,
		ClientOrderID:   2312312312,
		ReduceOnly:      true,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFastWithdrawalV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.FastWithdrawalV2(context.Background(), &FastWithdrawalParams{
		Amount:       1,
		ClientID:     "123213",
		Expiration:   time.Now().Add(time.Hour * 45),
		Asset:        currency.USDC,
		ERC20Address: "0x0330eBB5e894720e6746070371F9Fd797BE9D074",
		ChainID:      "56",
		Fees:         0,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFastWithdrawalV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.FastWithdrawalV1(context.Background(), &FastWithdrawalParams{
		Amount:       1,
		ClientID:     "123213",
		Expiration:   time.Now().Add(time.Hour * 45),
		Asset:        currency.USDC,
		ERC20Address: "0x0330eBB5e894720e6746070371F9Fd797BE9D074",
		ChainID:      "56",
		Fees:         0,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}
