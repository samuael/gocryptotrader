package paxos

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Paxos is the overarching type across this package
type Paxos struct {
	exchange.Base
}

const (
	paxosAPIURL     = "https://api.paxos.com/"
	paxosTestAPIURL = "https://api.sandbox.paxos.com/"
	paxosAPIVersion = "v2/"

	// Public endpoints

	// Authenticated endpoints
)

// Start implementing public and private exchange API funcs below

type timeSpecifier uint8

const (
	lessthanTime = iota
	lessthanEqualTo
	equalTo
	greaterthanEqualTo
	greaterthan
)

// GetProfiles return the associated Profiles for the current Account.
// The paginated results default to the maximum limit of 1,000 Profiles, unless otherwise specified with the limit parameter. Every paginated response contains a next_page field until the last page is reached. Pass the next_page value into the page_cursor field of the next request to retrieve the next page of results.
func (pa *Paxos) GetProfiles(ctx context.Context, createdAt time.Time, timeSpec timeSpecifier, limit int64, fetchingOrder, orderBy, pageCursor string) (*Profiles, error) {
	var resp *Profiles
	return resp, nil
}

// GetMarkets retrieves the set of current available markets for trading with details.
func (pa *Paxos) GetMarkets(ctx context.Context) (*MarketInstruments, error) {
	var resp *MarketInstruments
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets", &resp)
}

// GetOrderbook retrieves the full list of bids and asks of the order book at individual price levels with resting quantities per level.
func (pa *Paxos) GetOrderbook(ctx context.Context, market string) (*Orderbook, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *Orderbook
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/order-book", &resp)
}

// GetRecentExecutions retrieves the list of 2000 most recent executions by all users to occur in the order book.
func (pa *Paxos) GetRecentExecutions(ctx context.Context, market string) (*RecentExecutions, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *RecentExecutions
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/recent-executions", &resp)
}

// GetTicker retrieves order book statistics of the exchange over the last 24 hours and from midnight UTC until current time.
func (pa *Paxos) GetTicker(ctx context.Context, market string) (*TckerDetail, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *TckerDetail
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/ticker", &resp)
}

// -------------------------------------------------   Pricing Endpoints  ---------------------------------------

// GetPrices retrieve current prices, as well as 24 hour prior (yesterday) prices, for the specified markets.
func (pa *Paxos) GetPrices(ctx context.Context, markets currency.Pairs) (*MarketPricing, error) {
	if len(markets) == 0 {
		return nil, currency.ErrCurrencyPairsEmpty
	}
	marketsByte, err := json.Marshal(markets.Strings())
	if err != nil {
		return nil, err
	}
	var resp *MarketPricing
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "all-markets/prices?markets="+string(marketsByte), &resp)
}

// GetPriceTickers retrieves order book statistics of the exchange for all markets over the last 24 hours and from midnight UTC until current time.
func (pa *Paxos) GetPriceTickers(ctx context.Context) (*TckerDetail, error) {
	var resp *TckerDetail
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "all-markets/ticker", &resp)
}

// GetHistoricalPrices retrieves a set of average prices at a certain increment of time for the requested market.
// This endpoint is suitable for retrieving historical average price trends where price precision, execution, and other detailed information is not required.
// func (pa *Paxos) GetHistoricalPrices(ctx context.Context, )

// SendHTTPRequest sends an http request to a desired
func (pa *Paxos) SendHTTPRequest(ctx context.Context, ep exchange.URL, requestPath string, result interface{}) (err error) {
	endpoint, err := pa.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	requestType := request.AuthType(request.UnauthenticatedRequest)
	newRequest := func() (*request.Item, error) {
		path := endpoint + paxosAPIVersion + requestPath
		println("Path: ", path)
		headers := make(map[string]string)
		headers["Content-Type"] = "application/json"
		return &request.Item{
			Method:        http.MethodGet,
			Path:          path,
			Headers:       headers,
			Result:        &result,
			Verbose:       pa.Verbose,
			HTTPDebugging: pa.HTTPDebugging,
			HTTPRecording: pa.HTTPRecording,
		}, nil
	}
	return pa.SendPayload(ctx, request.Unset, newRequest, requestType)
}
