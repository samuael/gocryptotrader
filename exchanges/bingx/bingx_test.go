package bingx

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/currency"
	testexch "github.com/thrasher-corp/gocryptotrader/internal/testing/exchange"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
)

var e *Exchange
var btcUSDT currency.Pair

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
	os.Exit(m.Run())
}

func TestGetSpotSymbols(t *testing.T) {
	t.Parallel()
	e.Verbose = true
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
