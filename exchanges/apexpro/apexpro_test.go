package apexpro

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/config"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
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

	err = ap.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// Implement tests for API endpoints below

func TestGetSystemTimeV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetSystemTimeV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAllConfigDataV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetAllConfigDataV3(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetMarketDepthV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetMarketDepthV3(context.Background(), "BTCUSDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// MarketDepthV3

func TestGetNewestTradingDataV3(t *testing.T) {
	t.Parallel()
	result, err := ap.GetNewestTradingDataV3(context.Background(), "BTC-USDT", 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}
