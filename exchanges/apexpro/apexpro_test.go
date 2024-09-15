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
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/margin"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey    = "1fcc2d47-6c5b-37e5-cf5d-ee07a2975c66"
	apiSecret = "cDrY2JYzwtAGpZPENY3FnSy5W80CHJYY-dtA7TFW"
	clientID  = "RvtoQ0zSFBTEq0Ll-tKn"

	// apiKey    = "3ee9f1f5-84e7-3b45-6b8d-e0f6249792a1"
	// apiSecret = "5EQDIF2x9p3o5Hf9xlCrmo4vrUGYoabhPI0U283X"
	// clientID  = "WriryhFSKb8SpOtmvpWY"

	starkKey            = "0x06c98993ca62f5e71dbe721f743045eff7475711b359681cd64364a60e677505"
	starkSecret         = "0x074bcbe7f64f95e8d3f1afda4c338775702b4d1db3651fc70bad95a160b7f9ae"
	starkKeyYCoordinate = "0x0207d57867e0820e0f7588339e8b7491ce1da964260044340e3fd27c718f2a91"

	// starkKey            = "0xf8c6635f9cfe85f46759dc2eebe71a45b765687e35dbe5e74e8bde347813ef"
	// starkSecret         = "0x607ba3969039f3e19006ff8f40629d20a7b7dac31d4019e0965fbf7c5c068a"
	// starkKeyYCoordinate = ""

	// starkKey            = "0x002474dee3cd13931e85b4e7bb4a501d32097192515a492289f4804046ace567"
	// starkSecret         = "0x064bbfda1ae95578713f23ad9dc4a19a0e8b2edc0efdc6819d389146b140f24b"
	// starkKeyYCoordinate = "0x01c67b1317067aba57ea69fa9f3e59f4e556e520ab81dedb39b58409b79e2373"

	ethereumAddress = "0x0330eBB5e894720e6746070371F9Fd797BE9D074"

	canManipulateRealOrders = true
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
	result, err := ap.GetMarketDepthV3(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV2(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV1(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetNewestTradingDataV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV3(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetNewestTradingDataV2(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV2(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}
func TestGetNewestTradingDataV1(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV1(context.Background(), "BTC-USDC", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV3(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV3(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetCandlestickChartDataV3(context.Background(), "BTC-USDC", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV2(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)

	result, err := ap.GetCandlestickChartDataV2(context.Background(), "BTC-USDC", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCandlestickChartDataV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetCandlestickChartDataV1(context.Background(), "", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetCandlestickChartDataV1(context.Background(), "BTC-USDC", kline.FiveMin, time.Time{}, time.Time{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTickerDataV3(t *testing.T) {
	t.Parallel()
	_, err := ap.GetTickerDataV3(context.Background(), "")
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)

	result, err := ap.GetTickerDataV3(context.Background(), "BTC-USDC")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTickerDataV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetTickerDataV2(context.Background(), "")
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)

	result, err := ap.GetTickerDataV2(context.Background(), "BTC-USDC")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingHistoryRate(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV3(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV3(context.Background(), "BTC-USDC", time.Time{}, time.Time{}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetFundingHistoryRateV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV2(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV2(context.Background(), "BTC-USDC", time.Time{}, time.Time{}, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestGetFundingHistoryRateV1(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingHistoryRateV1(context.Background(), "", time.Time{}, time.Time{}, 10, 0)
	require.ErrorIs(t, err, currency.ErrSymbolStringEmpty)
	result, err := ap.GetFundingHistoryRateV1(context.Background(), "BTC-USDC", time.Time{}, time.Time{}, 0, 0)
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
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
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

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
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
		Username:                 "Username",
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
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
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
	_, err = ap.GetWorstPriceV3(context.Background(), "BTC-USDC", "", 1)
	require.ErrorIs(t, err, order.ErrSideIsInvalid)
	_, err = ap.GetWorstPriceV3(context.Background(), "BTC-USDC", "SELL", 0)
	require.ErrorIs(t, err, order.ErrAmountBelowMin)
}

func TestGetWorstPriceV3(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV3(context.Background(), "BTC-USDC", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPriceV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV2(context.Background(), "BTC-USDC", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWorstPriceV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWorstPriceV1(context.Background(), "BTC-USDC", "SELL", 1)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	futuresTradablePair, err := currency.NewPairFromString("ETH-USDC")
	require.NoError(t, err)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ap.UserAccountDetail)
	}
	ap.Verbose = true
	result, err := ap.CreateOrderV3(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Buy.String(),
		OrderType:       "LIMIT",
		Size:            0.01,
		Price:           2250,
		TimeInForce:     "GOOD_TIL_CANCEL",
		TriggerPrice:    0,
		TrailingPercent: 0,
		ReduceOnly:      false,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrder(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CancelPerpOrder(context.Background(), 123231)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelPerpOrderByClientOrderID(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CancelPerpOrderByClientOrderID(context.Background(), 2312312)
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
	result, err := ap.GetAllOrderHistory(context.Background(), "BTC-USDC", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllOrderHistoryV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetAllOrderHistoryV2(context.Background(), currency.EMPTYCODE, "BTC-USDC", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAllOrderHistoryV2(context.Background(), currency.USDT, "BTC-USDC", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllOrderHistoryV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAllOrderHistoryV1(context.Background(), "BTC-USDC", "SELL", "MARKET", "OPEN", "HISTORY", time.Time{}, time.Time{}, 0, 10)
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
	_, err := ap.getSingleOrder(context.Background(), "", "", currency.USDC)
	require.ErrorIs(t, err, order.ErrOrderIDNotSet)
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

func TestGetOrderByClientOrderID(t *testing.T) {
	t.Parallel()
	_, err := ap.GetOrderByClientOrderID(context.Background(), "")
	require.ErrorIs(t, err, order.ErrClientOrderIDMustBeSet)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOrderByClientOrderID(context.Background(), "12343")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingRate(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV3(context.Background(), "BTC-USDC", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingRateV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV1(context.Background(), "BTC-USDC", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetFundingRateV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetFundingRateV2(context.Background(), currency.EMPTYCODE, "BTC-USDC", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetFundingRateV2(context.Background(), currency.USDT, "BTC-USDC", "LONG", "", time.Now().Add(-time.Hour*50), time.Now(), 10, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserHistorialProfitAndLoss(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLoss(context.Background(), "BTC-USDC", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserHistorialProfitAndLossV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLossV1(context.Background(), "BTC-USDC", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUserHistorialProfitAndLossV2(t *testing.T) {
	t.Parallel()
	_, err := ap.GetUserHistorialProfitAndLossV2(context.Background(), currency.EMPTYCODE, "BTC-USDC", "LONG", time.Time{}, time.Time{}, 0, 100)
	require.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserHistorialProfitAndLossV2(context.Background(), currency.USDT, "BTC-USDC", "LONG", time.Time{}, time.Time{}, 0, 100)
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
	err := ap.SetInitialMarginRateInfo(context.Background(), "BTC-USDC", 200)
	assert.NoError(t, err)
}

func TestSetInitialMarginRateInfoV1(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfoV1(context.Background(), "BTC-USDC", 200)
	assert.NoError(t, err)
}

func TestSetInitialMarginRateInfoV2(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.SetInitialMarginRateInfoV2(context.Background(), "BTC-USDC", 200)
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
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.UserWithdrawalV2(context.Background(), 1000, "1231231", time.Now().Add(time.Hour*24), currency.USDC)
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
		ClientOrderID:   "12334",
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
		ClientOrderID:   "12334",
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
	futuresTradablePair, err := currency.NewPairFromString("ETH-USDC")
	require.NoError(t, err)

	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		require.NoError(t, err)
		require.NotNil(t, ap.UserAccountDetail)
	}

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)

	ap.Verbose = true
	result, err := ap.CreateOrderV1(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Buy.String(),
		OrderType:       "LIMIT",
		Size:            0.01,
		Price:           2250,
		TimeInForce:     "GOOD_TIL_CANCEL",
		TriggerPrice:    0,
		TrailingPercent: 0,
		ReduceOnly:      false,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateOrderV2(t *testing.T) {
	t.Parallel()
	futuresTradablePair, err := currency.NewPairFromString("ETH-USDC")
	require.NoError(t, err)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(context.Background())
		require.NoError(t, err)
		require.NotNil(t, ap.UserAccountDetail)
	}
	ap.Verbose = true
	result, err := ap.CreateOrderV2(context.Background(), &CreateOrderParams{
		Symbol:          futuresTradablePair,
		Side:            order.Buy.String(),
		OrderType:       "LIMIT",
		Size:            0.01,
		Price:           2250,
		TimeInForce:     "GOOD_TIL_CANCEL",
		TriggerPrice:    0,
		TrailingPercent: 0,
		ReduceOnly:      false,
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

func TestUpdateOrderExecutionLimits(t *testing.T) {
	t.Parallel()
	pairs, err := ap.FetchTradablePairs(context.Background(), asset.Futures)
	assert.NoErrorf(t, err, "FetchTradablePairs should not error for %s", asset.Futures)
	assert.NotEmptyf(t, pairs, "Should get some pairs for %s", asset.Futures)

	err = ap.UpdateOrderExecutionLimits(context.Background(), asset.Futures)
	require.NoError(t, err)

	limits, err := ap.GetOrderExecutionLimits(asset.Futures, pairs[0])
	assert.NoErrorf(t, err, "GetOrderExecutionLimits should not error for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MinPrice, "MinPrice must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MaxPrice, "MaxPrice must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.PriceStepIncrementSize, "PriceStepIncrementSize must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MinimumBaseAmount, "MinimumBaseAmount must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MaximumBaseAmount, "MaximumBaseAmount must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.AmountStepIncrementSize, "AmountStepIncrementSize must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MarketMaxQty, "MarketMaxQty must be positive for %s pair %s", asset.Futures, pairs[0])
	assert.Positivef(t, limits.MaxTotalOrders, "MaxTotalOrders must be positive for %s pair %s", asset.Futures, pairs[0])
}

func TestIsPerpetualFutureCurrency(t *testing.T) {
	t.Parallel()
	is, err := ap.IsPerpetualFutureCurrency(asset.Futures, currency.NewPair(currency.BTC, currency.USDC))
	require.NoError(t, err)
	assert.True(t, is)
}

func TestGetFuturesContractDetails(t *testing.T) {
	t.Parallel()
	result, err := ap.GetFuturesContractDetails(context.Background(), asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetHistoricCandles(t *testing.T) {
	t.Parallel()
	pair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	result, err := ap.GetHistoricCandles(context.Background(), pair, asset.Futures, kline.OneMin, time.Now().Add(-time.Minute*3), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetHistoricCandlesExtended(t *testing.T) {
	t.Parallel()
	pair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	result, err := ap.GetHistoricCandlesExtended(context.Background(), pair, asset.Futures, kline.OneMin, time.Now().Add(-time.Minute*3), time.Now())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFetchTradablePairs(t *testing.T) {
	t.Parallel()
	result, err := ap.FetchTradablePairs(context.Background(), asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateTradablePairs(t *testing.T) {
	t.Parallel()
	err := ap.UpdateTradablePairs(context.Background(), true)
	assert.NoError(t, err)
}

func TestUpdateTicker(t *testing.T) {
	t.Parallel()
	pair, err := currency.NewPairFromString("BTC-USDC")
	require.NoError(t, err)

	result, err := ap.UpdateTicker(context.Background(), pair, asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateOrderbook(t *testing.T) {
	t.Parallel()
	pair, err := currency.NewPairFromString("BTCUSD")
	require.NoError(t, err)

	result, err := ap.UpdateOrderbook(context.Background(), pair, asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateAccountInfo(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.UpdateAccountInfo(context.Background(), asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAccountFundingHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetAccountFundingHistory(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetWithdrawalsHistory(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetWithdrawalsHistory(context.Background(), currency.USDC, asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetRecentTrades(t *testing.T) {
	t.Parallel()
	pair, err := currency.NewPairFromString("BTCUSD")
	require.NoError(t, err)

	result, err := ap.GetRecentTrades(context.Background(), pair, asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetServerTime(t *testing.T) {
	t.Parallel()
	result, err := ap.GetServerTime(context.Background(), asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCancelAllOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.CancelAllOrders(context.Background(), &order.Cancel{
		AssetType:  asset.Futures,
		MarginType: margin.Isolated,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOrderInfo(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetOrderInfo(context.Background(), "614463889001677573", currency.EMPTYPAIR, asset.Futures)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetActiveOrders(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetActiveOrders(context.Background(), &order.MultiOrderRequest{
		AssetType: asset.Futures,
		Type:      order.Limit,
		Side:      order.Buy,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWsConnect(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	err := ap.WsConnect()
	require.NoError(t, err)
}

var data = `{
    "amount_collateral": "4000000",
    "amount_fee": "4000",
    "amount_synthetic": "1000000",
    "asset_id_collateral": "0xa21edc9d9997b1b1956f542fe95922518a9e28ace11b7b2972a1974bf5971f",
    "asset_id_synthetic": "0x0",
    "expiration_timestamp": "1100000",
    "is_buying_synthetic": false,
    "nonce": "1001",
    "order_type": "LIMIT_ORDER_WITH_FEES",
    "position_id": "10000",
    "public_key": "0xf8c6635f9cfe85f46759dc2eebe71a45b765687e35dbe5e74e8bde347813ef",
    "signature": {
        "r": "0x07a15838aad9b20368dc4ba27613fd35ceec3b34be7a2cb913bca0fb06e98107",
        "s": "0x05007f40fddd9babae0c7362d3b4e9c152ed3fced7fe78435b302d825489298f"
    }
}`

// func TestRegistrationAndOnboarding(t *testing.T) {
// 	t.Parallel()
// 	// sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
// 	ap.Verbose = true
// 	result, err := ap.RegistrationAndOnboarding(context.Background(), starkKey, starkKeyYCoordinate, "0x0330eBB5e894720e6746070371F9Fd797BE9D074", "", "")
// 	require.NoError(t, err)
// 	assert.NotNil(t, result)
// }
