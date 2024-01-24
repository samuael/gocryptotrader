package paxos

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = "a455062e-47c6-4756-89e9-9aefbacdb57c"
	apiSecret               = "J3IJL32UJ7KT56IRJ6VGPATKHQ"
	canManipulateRealOrders = false
)

var pa = &Paxos{}

func TestMain(m *testing.M) {
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Paxos")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret
	pa.SetDefaults()
	if apiKey != "" && apiSecret != "" {
		exchCfg.API.AuthenticatedSupport = true
		exchCfg.API.AuthenticatedWebsocketSupport = true
	}

	err = pa.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// Ensures that this exchange package is compatible with IBotExchange
func TestInterface(t *testing.T) {
	var e exchange.IBotExchange
	if e = new(Paxos); e == nil {
		t.Fatal("unable to allocate exchange")
	}
}

// Implement tests for API endpoints below

func TestGetMarkets(t *testing.T) {
	t.Parallel()
	result, err := pa.GetMarkets(context.Background())
	if err != nil {
		t.Fatal(err)
	} else {
		for a := range result.Markets {
			println(result.Markets[a].Market)
		}
	}
}

func TestGetOrderbook(t *testing.T) {
	t.Parallel()
	_, err := pa.GetOrderbook(context.Background(), "BCHUSD")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRecentExecutions(t *testing.T) {
	t.Parallel()
	_, err := pa.GetRecentExecutions(context.Background(), "BCHUSD")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTicker(t *testing.T) {
	t.Parallel()
	_, err := pa.GetTicker(context.Background(), "BCHUSD")
	if err != nil {
		t.Error(err)
	}
}

func TestGetPrices(t *testing.T) {
	t.Parallel()
	_, err := pa.GetPrices(context.Background(), currency.Pairs{{Base: currency.BTC, Quote: currency.USD}, {Base: currency.ETH, Quote: currency.USD}})
	if err != nil {
		t.Error(err)
	}
}

func TestGetTickers(t *testing.T) {
	t.Parallel()
	_, err := pa.GetPriceTickers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetProfiles(t *testing.T) {
	t.Parallel()
	// sharedtestvalues.SkipTestIfCredentialsUnset(t, pa)
	_, err := pa.GetProfiles(context.Background(), time.Now().Add(-time.Hour*24*5), lessthanEqualTo, 20, "DESC", "", "")
	if err != nil {
		t.Fatal(err)
	}
}
