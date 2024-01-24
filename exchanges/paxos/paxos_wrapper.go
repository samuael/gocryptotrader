package paxos

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/account"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/deposit"
	"github.com/thrasher-corp/gocryptotrader/exchanges/fundingrate"
	"github.com/thrasher-corp/gocryptotrader/exchanges/futures"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/protocol"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
	"github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/gocryptotrader/portfolio/withdraw"
)

// GetDefaultConfig returns a default exchange config
func (pa *Paxos) GetDefaultConfig(ctx context.Context) (*config.Exchange, error) {
	pa.SetDefaults()
	exchCfg, err := pa.GetStandardConfig()
	if err != nil {
		return nil, err
	}

	err = pa.SetupDefaults(exchCfg)
	if err != nil {
		return nil, err
	}

	if pa.Features.Supports.RESTCapabilities.AutoPairUpdates {
		err := pa.UpdateTradablePairs(ctx, true)
		if err != nil {
			return nil, err
		}
	}
	return exchCfg, nil
}

// SetDefaults sets the basic defaults for Paxos
func (pa *Paxos) SetDefaults() {
	pa.Name = "Paxos"
	pa.Enabled = true
	pa.Verbose = true
	pa.API.CredentialsValidator.RequiresKey = true
	pa.API.CredentialsValidator.RequiresSecret = true
	requestFmt := &currency.PairFormat{ /*Set pair request formatting details here for e.g.*/ Uppercase: true, Delimiter: ":"}
	// Config format denotes what the pair as a string will be, when saved to
	// the config.json file.
	configFmt := &currency.PairFormat{ /*Set pair request formatting details here*/ }
	err := pa.SetGlobalPairsManager(requestFmt, configFmt /*multiple assets can be set here using the asset package ie asset.Spot*/)
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}

	fmt1 := currency.PairStore{
		RequestFormat: &currency.PairFormat{Uppercase: true},
		ConfigFormat:  &currency.PairFormat{Uppercase: true},
	}

	fmt2 := currency.PairStore{
		RequestFormat: &currency.PairFormat{Uppercase: true},
		ConfigFormat:  &currency.PairFormat{Uppercase: true, Delimiter: ":"},
	}

	err = pa.StoreAssetPairFormat(asset.Spot, fmt1)
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}
	err = pa.StoreAssetPairFormat(asset.Margin, fmt2)
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}

	// Fill out the capabilities/features that the exchange supports
	pa.Features = exchange.Features{
		Supports: exchange.FeaturesSupported{
			REST:      true,
			Websocket: true,
			RESTCapabilities: protocol.Features{
				TickerFetching:    true,
				OrderbookFetching: true,
			},
			WebsocketCapabilities: protocol.Features{
				TickerFetching:    true,
				OrderbookFetching: true,
			},
			WithdrawPermissions: exchange.AutoWithdrawCrypto |
				exchange.AutoWithdrawFiat,
		},
		Enabled: exchange.FeaturesEnabled{
			AutoPairUpdates: true,
		},
	}
	// NOTE: SET THE EXCHANGES RATE LIMIT HERE
	pa.Requester, err = request.New(pa.Name,
		common.NewHTTPClientWithTimeout(exchange.DefaultHTTPTimeout))
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}

	// NOTE: SET THE URLs HERE
	pa.API.Endpoints = pa.NewEndpoints()
	pa.API.Endpoints.SetDefaultEndpoints(map[exchange.URL]string{
		exchange.RestSpot:              paxosAPIURL,
		exchange.RestSpotSupplementary: paxosTestAPIURL,
		// exchange.WebsocketSpotSupplementary:
	})
	pa.Websocket = stream.New()
	pa.WebsocketResponseMaxLimit = exchange.DefaultWebsocketResponseMaxLimit
	pa.WebsocketResponseCheckTimeout = exchange.DefaultWebsocketResponseCheckTimeout
	pa.WebsocketOrderbookBufferLimit = exchange.DefaultWebsocketOrderbookBufferLimit
}

// Setup takes in the supplied exchange configuration details and sets params
func (pa *Paxos) Setup(exch *config.Exchange) error {
	err := exch.Validate()
	if err != nil {
		return err
	}
	if !exch.Enabled {
		pa.SetEnabled(false)
		return nil
	}
	err = pa.SetupDefaults(exch)
	if err != nil {
		return err
	}

	/*
		wsRunningEndpoint, err := pa.API.Endpoints.GetURL(exchange.WebsocketSpot)
		if err != nil {
			return err
		}

		// If websocket is supported, please fill out the following

		err = pa.Websocket.Setup(
			&stream.WebsocketSetup{
				ExchangeConfig:  exch,
				DefaultURL:      paxosWSAPIURL,
				RunningURL:      wsRunningEndpoint,
				Connector:       pa.WsConnect,
				Subscriber:      pa.Subscribe,
				UnSubscriber:    pa.Unsubscribe,
				Features:        &pa.Features.Supports.WebsocketCapabilities,
			})
		if err != nil {
			return err
		}

		pa.WebsocketConn = &stream.WebsocketConnection{
			ExchangeName:         pa.Name,
			URL:                  pa.Websocket.GetWebsocketURL(),
			ProxyURL:             pa.Websocket.GetProxyAddress(),
			Verbose:              pa.Verbose,
			ResponseCheckTimeout: exch.WebsocketResponseCheckTimeout,
			ResponseMaxLimit:     exch.WebsocketResponseMaxLimit,
		}
	*/
	return nil
}

// Start starts the Paxos go routine
func (pa *Paxos) Start(ctx context.Context, wg *sync.WaitGroup) error {
	if wg == nil {
		return fmt.Errorf("%T %w", wg, common.ErrNilPointer)
	}
	wg.Add(1)
	go func() {
		pa.Run(ctx)
		wg.Done()
	}()
	return nil
}

// Run implements the Paxos wrapper
func (pa *Paxos) Run(ctx context.Context) {
	if pa.Verbose {
		log.Debugf(log.ExchangeSys,
			"%s Websocket: %s.",
			pa.Name,
			common.IsEnabled(pa.Websocket.IsEnabled()))
		pa.PrintEnabledPairs()
	}

	if !pa.GetEnabledFeatures().AutoPairUpdates {
		return
	}

	err := pa.UpdateTradablePairs(ctx, false)
	if err != nil {
		log.Errorf(log.ExchangeSys,
			"%s failed to update tradable pairs. Err: %s",
			pa.Name,
			err)
	}
}

// FetchTradablePairs returns a list of the exchanges tradable pairs
func (pa *Paxos) FetchTradablePairs(ctx context.Context, a asset.Item) (currency.Pairs, error) {
	// Implement fetching the exchange available pairs if supported
	return nil, nil
}

// UpdateTradablePairs updates the exchanges available pairs and stores
// them in the exchanges config
func (pa *Paxos) UpdateTradablePairs(ctx context.Context, forceUpdate bool) error {
	assetTypes := pa.GetAssetTypes(false)
	for x := range assetTypes {
		pairs, err := pa.FetchTradablePairs(ctx, assetTypes[x])
		if err != nil {
			return err
		}

		err = pa.UpdatePairs(pairs, assetTypes[x], false, forceUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateTicker updates and returns the ticker for a currency pair
func (pa *Paxos) UpdateTicker(ctx context.Context, p currency.Pair, assetType asset.Item) (*ticker.Price, error) {
	// NOTE: EXAMPLE FOR GETTING TICKER PRICE
	/*
		tickerPrice := new(ticker.Price)
		tick, err := pa.GetTicker(p.String())
		if err != nil {
			return tickerPrice, err
		}
		tickerPrice = &ticker.Price{
			High:    tick.High,
			Low:     tick.Low,
			Bid:     tick.Bid,
			Ask:     tick.Ask,
			Open:    tick.Open,
			Close:   tick.Close,
			Pair:    p,
		}
		err = ticker.ProcessTicker(pa.Name, tickerPrice, assetType)
		if err != nil {
			return tickerPrice, err
		}
	*/
	return ticker.GetTicker(pa.Name, p, assetType)
}

// UpdateTickers updates all currency pairs of a given asset type
func (pa *Paxos) UpdateTickers(ctx context.Context, assetType asset.Item) error {
	// NOTE: EXAMPLE FOR GETTING TICKER PRICE
	/*
			tick, err := pa.GetTickers()
			if err != nil {
				return err
			}
		    for y := range tick {
		        cp, err := currency.NewPairFromString(tick[y].Symbol)
		        if err != nil {
		            return err
		        }
		        err = ticker.ProcessTicker(&ticker.Price{
		            Last:         tick[y].LastPrice,
		            High:         tick[y].HighPrice,
		            Low:          tick[y].LowPrice,
		            Bid:          tick[y].BidPrice,
		            Ask:          tick[y].AskPrice,
		            Volume:       tick[y].Volume,
		            QuoteVolume:  tick[y].QuoteVolume,
		            Open:         tick[y].OpenPrice,
		            Close:        tick[y].PrevClosePrice,
		            Pair:         cp,
		            ExchangeName: b.Name,
		            AssetType:    assetType,
		        })
		        if err != nil {
		            return err
		        }
		    }
	*/
	return nil
}

// FetchTicker returns the ticker for a currency pair
func (pa *Paxos) FetchTicker(ctx context.Context, p currency.Pair, assetType asset.Item) (*ticker.Price, error) {
	tickerNew, err := ticker.GetTicker(pa.Name, p, assetType)
	if err != nil {
		return pa.UpdateTicker(ctx, p, assetType)
	}
	return tickerNew, nil
}

// FetchOrderbook returns orderbook base on the currency pair
func (pa *Paxos) FetchOrderbook(ctx context.Context, pair currency.Pair, assetType asset.Item) (*orderbook.Base, error) {
	ob, err := orderbook.Get(pa.Name, pair, assetType)
	if err != nil {
		return pa.UpdateOrderbook(ctx, pair, assetType)
	}
	return ob, nil
}

// UpdateOrderbook updates and returns the orderbook for a currency pair
func (pa *Paxos) UpdateOrderbook(ctx context.Context, pair currency.Pair, assetType asset.Item) (*orderbook.Base, error) {
	book := &orderbook.Base{
		Exchange:        pa.Name,
		Pair:            pair,
		Asset:           assetType,
		VerifyOrderbook: pa.CanVerifyOrderbook,
	}

	// NOTE: UPDATE ORDERBOOK EXAMPLE
	/*
		orderbookNew, err := pa.GetOrderBook(exchange.FormatExchangeCurrency(pa.Name, p).String(), 1000)
		if err != nil {
			return book, err
		}

		book.Bids = make([]orderbook.Item, len(orderbookNew.Bids))
		for x := range orderbookNew.Bids {
			book.Bids[x] = orderbook.Item{
				Amount: orderbookNew.Bids[x].Quantity,
				Price: orderbookNew.Bids[x].Price,
			}
		}

		book.Asks = make([]orderbook.Item, len(orderbookNew.Asks))
		for x := range orderbookNew.Asks {
			book.Asks[x] = orderbook.Item{
				Amount: orderBookNew.Asks[x].Quantity,
				Price: orderBookNew.Asks[x].Price,
			}
		}
	*/

	err := book.Process()
	if err != nil {
		return book, err
	}

	return orderbook.Get(pa.Name, pair, assetType)
}

// UpdateAccountInfo retrieves balances for all enabled currencies
func (pa *Paxos) UpdateAccountInfo(ctx context.Context, assetType asset.Item) (account.Holdings, error) {
	// If fetching requires more than one asset type please set
	// HasAssetTypeAccountSegregation to true in RESTCapabilities above.
	return account.Holdings{}, common.ErrNotYetImplemented
}

// FetchAccountInfo retrieves balances for all enabled currencies
func (pa *Paxos) FetchAccountInfo(ctx context.Context, assetType asset.Item) (account.Holdings, error) {
	// Example implementation below:
	// 	creds, err := pa.GetCredentials(ctx)
	// 	if err != nil {
	// 		return account.Holdings{}, err
	// 	}
	// 	acc, err := account.GetHoldings(pa.Name, creds, assetType)
	// 	if err != nil {
	// 		return pa.UpdateAccountInfo(ctx, assetType)
	// 	}
	// 	return acc, nil
	return account.Holdings{}, common.ErrNotYetImplemented
}

// GetFundingHistory returns funding history, deposits and
// withdrawals
func (pa *Paxos) GetAccountFundingHistory(ctx context.Context) ([]exchange.FundingHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// GetWithdrawalsHistory returns previous withdrawals data
func (pa *Paxos) GetWithdrawalsHistory(ctx context.Context, c currency.Code, a asset.Item) ([]exchange.WithdrawalHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// GetRecentTrades returns the most recent trades for a currency and asset
func (pa *Paxos) GetRecentTrades(ctx context.Context, p currency.Pair, assetType asset.Item) ([]trade.Data, error) {
	return nil, common.ErrNotYetImplemented
}

// GetHistoricTrades returns historic trade data within the timeframe provided
func (pa *Paxos) GetHistoricTrades(ctx context.Context, p currency.Pair, assetType asset.Item, timestampStart, timestampEnd time.Time) ([]trade.Data, error) {
	return nil, common.ErrNotYetImplemented
}

// GetServerTime returns the current exchange server time.
func (pa *Paxos) GetServerTime(ctx context.Context, a asset.Item) (time.Time, error) {
	return time.Time{}, common.ErrNotYetImplemented
}

// SubmitOrder submits a new order
func (pa *Paxos) SubmitOrder(ctx context.Context, s *order.Submit) (*order.SubmitResponse, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	// When an order has been submitted you can use this helpful constructor to
	// return. Please add any additional order details to the
	// order.SubmitResponse if you think they are applicable.
	// resp, err := s.DeriveSubmitResponse( /*newOrderID*/)
	// if err != nil {
	// 	return nil, nil
	// }
	// resp.Date = exampleTime // e.g. If this is supplied by the exchanges API.
	// return resp, nil
	return nil, common.ErrNotYetImplemented
}

// ModifyOrder will allow of changing orderbook placement and limit to
// market conversion
func (pa *Paxos) ModifyOrder(ctx context.Context, action *order.Modify) (*order.ModifyResponse, error) {
	if err := action.Validate(); err != nil {
		return nil, err
	}
	// When an order has been modified you can use this helpful constructor to
	// return. Please add any additional order details to the
	// order.ModifyResponse if you think they are applicable.
	// resp, err := action.DeriveModifyResponse()
	// if err != nil {
	// 	return nil, nil
	// }
	// resp.OrderID = maybeANewOrderID // e.g. If this is supplied by the exchanges API.
	return nil, common.ErrNotYetImplemented
}

// CancelOrder cancels an order by its corresponding ID number
func (pa *Paxos) CancelOrder(ctx context.Context, ord *order.Cancel) error {
	// if err := ord.Validate(ord.StandardCancel()); err != nil {
	//	 return err
	// }
	return common.ErrNotYetImplemented
}

// CancelBatchOrders cancels orders by their corresponding ID numbers
func (pa *Paxos) CancelBatchOrders(ctx context.Context, orders []order.Cancel) (*order.CancelBatchResponse, error) {
	return nil, common.ErrNotYetImplemented
}

// CancelAllOrders cancels all orders associated with a currency pair
func (pa *Paxos) CancelAllOrders(ctx context.Context, orderCancellation *order.Cancel) (order.CancelAllResponse, error) {
	// if err := orderCancellation.Validate(); err != nil {
	//	 return err
	// }
	return order.CancelAllResponse{}, common.ErrNotYetImplemented
}

// GetOrderInfo returns order information based on order ID
func (pa *Paxos) GetOrderInfo(ctx context.Context, orderID string, pair currency.Pair, assetType asset.Item) (*order.Detail, error) {
	return nil, common.ErrNotYetImplemented
}

// GetDepositAddress returns a deposit address for a specified currency
func (pa *Paxos) GetDepositAddress(ctx context.Context, c currency.Code, accountID string, chain string) (*deposit.Address, error) {
	return nil, common.ErrNotYetImplemented
}

// WithdrawCryptocurrencyFunds returns a withdrawal ID when a withdrawal is
// submitted
func (pa *Paxos) WithdrawCryptocurrencyFunds(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// WithdrawFiatFunds returns a withdrawal ID when a withdrawal is
// submitted
func (pa *Paxos) WithdrawFiatFunds(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// WithdrawFiatFundsToInternationalBank returns a withdrawal ID when a withdrawal is
// submitted
func (pa *Paxos) WithdrawFiatFundsToInternationalBank(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetActiveOrders retrieves any orders that are active/open
func (pa *Paxos) GetActiveOrders(ctx context.Context, getOrdersRequest *order.MultiOrderRequest) (order.FilteredOrders, error) {
	// if err := getOrdersRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetOrderHistory retrieves account order information
// Can Limit response to specific order status
func (pa *Paxos) GetOrderHistory(ctx context.Context, getOrdersRequest *order.MultiOrderRequest) (order.FilteredOrders, error) {
	// if err := getOrdersRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetFeeByType returns an estimate of fee based on the type of transaction
func (pa *Paxos) GetFeeByType(ctx context.Context, feeBuilder *exchange.FeeBuilder) (float64, error) {
	return 0, common.ErrNotYetImplemented
}

// ValidateAPICredentials validates current credentials used for wrapper
func (pa *Paxos) ValidateAPICredentials(ctx context.Context, assetType asset.Item) error {
	_, err := pa.UpdateAccountInfo(ctx, assetType)
	return pa.CheckTransientError(err)
}

// GetHistoricCandles returns candles between a time period for a set time interval
func (pa *Paxos) GetHistoricCandles(ctx context.Context, pair currency.Pair, a asset.Item, interval kline.Interval, start, end time.Time) (*kline.Item, error) {
	return nil, common.ErrNotYetImplemented
}

// GetHistoricCandlesExtended returns candles between a time period for a set time interval
func (pa *Paxos) GetHistoricCandlesExtended(ctx context.Context, pair currency.Pair, a asset.Item, interval kline.Interval, start, end time.Time) (*kline.Item, error) {
	return nil, common.ErrNotYetImplemented
}

// GetFuturesContractDetails returns all contracts from the exchange by asset type
func (pa *Paxos) GetFuturesContractDetails(context.Context, asset.Item) ([]futures.Contract, error) {
	return nil, common.ErrNotYetImplemented
}

// GetLatestFundingRates returns the latest funding rates data
func (pa *Paxos) GetLatestFundingRates(_ context.Context, _ *fundingrate.LatestRateRequest) ([]fundingrate.LatestRateResponse, error) {
	return nil, common.ErrNotYetImplemented
}

// UpdateOrderExecutionLimits updates order execution limits
func (pa *Paxos) UpdateOrderExecutionLimits(_ context.Context, _ asset.Item) error {
	return common.ErrNotYetImplemented
}
