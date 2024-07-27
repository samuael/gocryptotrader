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
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GenerateNonce(context.Background(), starkKey, ethereumAddress, "9")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestRegistrationAndOnboarding(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.RegistrationAndOnboarding(context.Background(), starkKey, ethereumAddress, "", "")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetUsersData(t *testing.T) {
	t.Parallel()
	ap.Verbose = true
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUsersData(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestEditUserData(t *testing.T) {
	t.Parallel()
	_, err := ap.EditUserData(context.Background(), &EditUserDataParams{})
	require.ErrorIs(t, err, common.ErrNilPointer)

	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap, canManipulateRealOrders)
	result, err := ap.EditUserData(context.Background(), &EditUserDataParams{
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
	result, err := ap.GetUserAccountData(context.Background())
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

func TestGetUserDepositWithdrawData(t *testing.T) {
	t.Parallel()
	sharedtestvalues.SkipTestIfCredentialsUnset(t, ap)
	result, err := ap.GetUserTransferData(context.Background(), "", "", "WITHDRAW", "", "", time.Time{}, time.Time{}, []string{}, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}
