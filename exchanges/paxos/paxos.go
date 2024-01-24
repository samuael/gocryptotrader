package paxos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
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
	params := url.Values{}
	switch timeSpec {
	case lessthanTime:
		params.Set("created_at.lt", createdAt.Format(time.RFC3339))
	case lessthanEqualTo:
		params.Set("created_at.lte", createdAt.Format(time.RFC3339))
	case equalTo:
		params.Set("created_at.eq", createdAt.Format(time.RFC3339))
	case greaterthanEqualTo:
		params.Set("created_at.gte", createdAt.Format(time.RFC3339))
	case greaterthan:
		params.Set("created_at.gt", createdAt.Format(time.RFC3339))
	}
	if limit > 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if fetchingOrder != "" {
		params.Set("order", fetchingOrder)
	}
	if orderBy != "" {
		params.Set("order_by", orderBy)
	}
	if pageCursor != "" {
		params.Set("page_cursor", pageCursor)
	}
	var resp *Profiles
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, common.EncodeURLValues("/profiles", params), &resp, true)
}

// GetMarkets retrieves the set of current available markets for trading with details.
func (pa *Paxos) GetMarkets(ctx context.Context) (*MarketInstruments, error) {
	var resp *MarketInstruments
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets", &resp, false)
}

// GetOrderbook retrieves the full list of bids and asks of the order book at individual price levels with resting quantities per level.
func (pa *Paxos) GetOrderbook(ctx context.Context, market string) (*Orderbook, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *Orderbook
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/order-book", &resp, false)
}

// GetRecentExecutions retrieves the list of 2000 most recent executions by all users to occur in the order book.
func (pa *Paxos) GetRecentExecutions(ctx context.Context, market string) (*RecentExecutions, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *RecentExecutions
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/recent-executions", &resp, false)
}

// GetTicker retrieves order book statistics of the exchange over the last 24 hours and from midnight UTC until current time.
func (pa *Paxos) GetTicker(ctx context.Context, market string) (*TckerDetail, error) {
	if market == "" {
		return nil, currency.ErrSymbolStringEmpty
	}
	var resp *TckerDetail
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "markets/"+market+"/ticker", &resp, false)
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
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "all-markets/prices?markets="+string(marketsByte), &resp, false)
}

// GetPriceTickers retrieves order book statistics of the exchange for all markets over the last 24 hours and from midnight UTC until current time.
func (pa *Paxos) GetPriceTickers(ctx context.Context) (*TckerDetail, error) {
	var resp *TckerDetail
	return resp, pa.SendHTTPRequest(ctx, exchange.RestSpot, "all-markets/ticker", &resp, false)
}

// GetHistoricalPrices retrieves a set of average prices at a certain increment of time for the requested market.
// This endpoint is suitable for retrieving historical average price trends where price precision, execution, and other detailed information is not required.
// func (pa *Paxos) GetHistoricalPrices(ctx context.Context ){
// return "", nil
// }

var permissions = map[string]string{}

// SendHTTPRequest sends an http request to a desired
func (pa *Paxos) SendHTTPRequest(ctx context.Context, ep exchange.URL, requestPath string, result interface{}, authenticated bool) (err error) {
	endpoint, err := pa.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	requestType := request.AuthType(request.UnauthenticatedRequest)
	var accessToken string
	if authenticated {
		// switch requestPath {
		// case "/profiles":
		accessToken, err = pa.AuthenticateEndpoint(ctx, []string{"funding:read_profile"})
		if err != nil {
			return err
		}
		// case "":
		// }
	}
	newRequest := func() (*request.Item, error) {
		path := endpoint + paxosAPIVersion + requestPath
		headers := make(map[string]string)
		headers["Content-Type"] = "application/json"
		if authenticated {
			headers["Authorization"] = "Bearer " + accessToken
		}
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
	err = pa.SendPayload(ctx, request.Unset, newRequest, requestType)
	if err != nil {
		if authenticated {
			return fmt.Errorf("%w %v", request.ErrAuthRequestFailed, err)
		}
		return err
	}
	return nil
}

// AuthenticateEndpoint authenticates and generate a token.
func (pa *Paxos) AuthenticateEndpoint(ctx context.Context, scops []string) (string, error) {
	creds, err := pa.GetCredentials(ctx)
	if err != nil {
		return "", err
	}
	var result string
	newRequest := func() (*request.Item, error) {
		path := "https://oauth.paxos.com/oauth2/token"
		headers := make(map[string]string)
		headers["grant_type"] = "client_credentials"
		headers["client_id"] = creds.ClientID
		headers["client_secret"] = creds.Secret
		headers["scope"] = strings.Join(scops, " ")
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
	return result, pa.SendPayload(ctx, request.Unset, newRequest, request.AuthenticatedRequest)
}
