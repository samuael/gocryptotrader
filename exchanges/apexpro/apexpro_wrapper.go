package apexpro

import (
	"context"
	"fmt"
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

// SetDefaults sets the basic defaults for Apexpro
func (ap *Apexpro) SetDefaults() {
	ap.Name = "Apexpro"
	ap.Enabled = true
	ap.Verbose = false
	ap.API.CredentialsValidator.RequiresKey = true
	ap.API.CredentialsValidator.RequiresSecret = true

	requestFmt := &currency.PairFormat{Uppercase: true, Delimiter: ""}
	configFmt := &currency.PairFormat{Uppercase: true, Delimiter: "-"}
	err := ap.StoreAssetPairFormat(asset.Futures, currency.PairStore{
		RequestFormat: requestFmt,
		ConfigFormat:  configFmt,
	})
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}

	ap.Features = exchange.Features{
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
	ap.Requester, err = request.New(ap.Name,
		common.NewHTTPClientWithTimeout(exchange.DefaultHTTPTimeout))
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}

	ap.API.Endpoints = ap.NewEndpoints()
	err = ap.API.Endpoints.SetDefaultEndpoints(map[exchange.URL]string{
		exchange.RestSpotSupplementary:      apexproAPIURL,
		exchange.RestSpot:                   apexproTestAPIURL,
		exchange.WebsocketSpot:              apexProWebsocket,
		exchange.WebsocketSpotSupplementary: apexProPrivateWebsocket,
	})
	if err != nil {
		log.Errorln(log.ExchangeSys, err)
	}
	ap.Websocket = stream.NewWebsocket()
	ap.WebsocketResponseMaxLimit = exchange.DefaultWebsocketResponseMaxLimit
	ap.WebsocketResponseCheckTimeout = exchange.DefaultWebsocketResponseCheckTimeout
	ap.WebsocketOrderbookBufferLimit = exchange.DefaultWebsocketOrderbookBufferLimit
}

// Setup takes in the supplied exchange configuration details and sets params
func (ap *Apexpro) Setup(exch *config.Exchange) error {
	err := exch.Validate()
	if err != nil {
		return err
	}
	if !exch.Enabled {
		ap.SetEnabled(false)
		return nil
	}
	err = ap.SetupDefaults(exch)
	if err != nil {
		return err
	}
	wsRunningEndpoint, err := ap.API.Endpoints.GetURL(exchange.WebsocketSpot)
	if err != nil {
		return err
	}

	err = ap.Websocket.Setup(
		&stream.WebsocketSetup{
			ExchangeConfig:        exch,
			DefaultURL:            apexProWebsocket,
			RunningURL:            wsRunningEndpoint,
			Connector:             ap.WsConnect,
			Subscriber:            ap.Subscribe,
			Unsubscriber:          ap.Unsubscribe,
			GenerateSubscriptions: ap.GenerateDefaultSubscriptions,
			Features:              &ap.Features.Supports.WebsocketCapabilities,
		})
	if err != nil {
		return err
	}
	err = ap.Websocket.SetupNewConnection(stream.ConnectionSetup{
		URL:                  apexProWebsocket,
		ResponseCheckTimeout: exch.WebsocketResponseCheckTimeout,
		ResponseMaxLimit:     exch.WebsocketResponseMaxLimit,
	})
	if err != nil {
		return err
	}

	return ap.Websocket.SetupNewConnection(stream.ConnectionSetup{
		URL:                  apexProPrivateWebsocket,
		ResponseCheckTimeout: exch.WebsocketResponseCheckTimeout,
		ResponseMaxLimit:     exch.WebsocketResponseMaxLimit,
		Authenticated:        true,
	})
}

// FetchTradablePairs returns a list of the exchanges tradable pairs
func (ap *Apexpro) FetchTradablePairs(ctx context.Context, a asset.Item) (currency.Pairs, error) {
	if !ap.SupportsAsset(a) {
		return nil, fmt.Errorf("%w %v", asset.ErrNotSupported, a)
	}
	configs, err := ap.GetAllConfigDataV3(ctx)
	if err != nil {
		return nil, err
	}
	tradablePairs := make(currency.Pairs, 0, len((configs.ContractConfig.PerpetualContract)))
	for a := range configs.ContractConfig.PerpetualContract {
		if !configs.ContractConfig.PerpetualContract[a].EnableTrade {
			continue
		}
		cp, err := currency.NewPairFromString(configs.ContractConfig.PerpetualContract[a].Symbol)
		if err != nil {
			return nil, err
		}
		tradablePairs = append(tradablePairs, cp)
	}
	return tradablePairs, nil
}

// UpdateTradablePairs updates the exchanges available pairs and stores
// them in the exchanges config
func (ap *Apexpro) UpdateTradablePairs(ctx context.Context, forceUpdate bool) error {
	pairs, err := ap.FetchTradablePairs(ctx, asset.Futures)
	if err != nil {
		return err
	}
	return ap.UpdatePairs(pairs, asset.Futures, true, forceUpdate)
}

// UpdateTicker updates and returns the ticker for a currency pair
func (ap *Apexpro) UpdateTicker(ctx context.Context, p currency.Pair, assetType asset.Item) (*ticker.Price, error) {
	// NOTE: EXAMPLE FOR GETTING TICKER PRICE
	/*
		tickerPrice := new(ticker.Price)
		tick, err := ap.GetTicker(p.String())
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
		err = ticker.ProcessTicker(ap.Name, tickerPrice, assetType)
		if err != nil {
			return tickerPrice, err
		}
	*/
	return ticker.GetTicker(ap.Name, p, assetType)
}

// UpdateTickers updates all currency pairs of a given asset type
func (ap *Apexpro) UpdateTickers(ctx context.Context, assetType asset.Item) error {
	// NOTE: EXAMPLE FOR GETTING TICKER PRICE
	/*
			tick, err := ap.GetTickers()
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
func (ap *Apexpro) FetchTicker(ctx context.Context, p currency.Pair, assetType asset.Item) (*ticker.Price, error) {
	tickerNew, err := ticker.GetTicker(ap.Name, p, assetType)
	if err != nil {
		return ap.UpdateTicker(ctx, p, assetType)
	}
	return tickerNew, nil
}

// FetchOrderbook returns orderbook base on the currency pair
func (ap *Apexpro) FetchOrderbook(ctx context.Context, pair currency.Pair, assetType asset.Item) (*orderbook.Base, error) {
	ob, err := orderbook.Get(ap.Name, pair, assetType)
	if err != nil {
		return ap.UpdateOrderbook(ctx, pair, assetType)
	}
	return ob, nil
}

// UpdateOrderbook updates and returns the orderbook for a currency pair
func (ap *Apexpro) UpdateOrderbook(ctx context.Context, pair currency.Pair, assetType asset.Item) (*orderbook.Base, error) {
	book := &orderbook.Base{
		Exchange:        ap.Name,
		Pair:            pair,
		Asset:           assetType,
		VerifyOrderbook: ap.CanVerifyOrderbook,
	}

	// NOTE: UPDATE ORDERBOOK EXAMPLE
	/*
		orderbookNew, err := ap.GetOrderBook(exchange.FormatExchangeCurrency(ap.Name, p).String(), 1000)
		if err != nil {
			return book, err
		}

		book.Bids = make([]orderbook.Tranche, len(orderbookNew.Bids))
		for x := range orderbookNew.Bids {
			book.Bids[x] = orderbook.Tranche{
				Amount: orderbookNew.Bids[x].Quantity,
				Price: orderbookNew.Bids[x].Price,
			}
		}

		book.Asks = make([]orderbook.Tranche, len(orderbookNew.Asks))
		for x := range orderbookNew.Asks {
			book.Asks[x] = orderbook.Tranche{
				Amount: orderBookNew.Asks[x].Quantity,
				Price: orderBookNew.Asks[x].Price,
			}
		}
	*/

	err := book.Process()
	if err != nil {
		return book, err
	}

	return orderbook.Get(ap.Name, pair, assetType)
}

// UpdateAccountInfo retrieves balances for all enabled currencies
func (ap *Apexpro) UpdateAccountInfo(ctx context.Context, assetType asset.Item) (account.Holdings, error) {
	// If fetching requires more than one asset type please set
	// HasAssetTypeAccountSegregation to true in RESTCapabilities above.
	return account.Holdings{}, common.ErrNotYetImplemented
}

// FetchAccountInfo retrieves balances for all enabled currencies
func (ap *Apexpro) FetchAccountInfo(ctx context.Context, assetType asset.Item) (account.Holdings, error) {
	// Example implementation below:
	// 	creds, err := ap.GetCredentials(ctx)
	// 	if err != nil {
	// 		return account.Holdings{}, err
	// 	}
	// 	acc, err := account.GetHoldings(ap.Name, creds, assetType)
	// 	if err != nil {
	// 		return ap.UpdateAccountInfo(ctx, assetType)
	// 	}
	// 	return acc, nil
	return account.Holdings{}, common.ErrNotYetImplemented
}

// GetFundingHistory returns funding history, deposits and
// withdrawals
func (ap *Apexpro) GetAccountFundingHistory(ctx context.Context) ([]exchange.FundingHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// GetWithdrawalsHistory returns previous withdrawals data
func (ap *Apexpro) GetWithdrawalsHistory(ctx context.Context, c currency.Code, a asset.Item) ([]exchange.WithdrawalHistory, error) {
	return nil, common.ErrNotYetImplemented
}

// GetRecentTrades returns the most recent trades for a currency and asset
func (ap *Apexpro) GetRecentTrades(ctx context.Context, p currency.Pair, assetType asset.Item) ([]trade.Data, error) {
	return nil, common.ErrNotYetImplemented
}

// GetHistoricTrades returns historic trade data within the timeframe provided
func (ap *Apexpro) GetHistoricTrades(ctx context.Context, p currency.Pair, assetType asset.Item, timestampStart, timestampEnd time.Time) ([]trade.Data, error) {
	return nil, common.ErrNotYetImplemented
}

// GetServerTime returns the current exchange server time.
func (ap *Apexpro) GetServerTime(ctx context.Context, a asset.Item) (time.Time, error) {
	return time.Time{}, common.ErrNotYetImplemented
}

// SubmitOrder submits a new order
func (ap *Apexpro) SubmitOrder(ctx context.Context, s *order.Submit) (*order.SubmitResponse, error) {
	if err := s.Validate(ap.GetTradingRequirements()); err != nil {
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
func (ap *Apexpro) ModifyOrder(ctx context.Context, action *order.Modify) (*order.ModifyResponse, error) {
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
func (ap *Apexpro) CancelOrder(ctx context.Context, ord *order.Cancel) error {
	// if err := ord.Validate(ord.StandardCancel()); err != nil {
	//	 return err
	// }
	return common.ErrNotYetImplemented
}

// CancelBatchOrders cancels orders by their corresponding ID numbers
func (ap *Apexpro) CancelBatchOrders(ctx context.Context, orders []order.Cancel) (*order.CancelBatchResponse, error) {
	return nil, common.ErrNotYetImplemented
}

// CancelAllOrders cancels all orders associated with a currency pair
func (ap *Apexpro) CancelAllOrders(ctx context.Context, orderCancellation *order.Cancel) (order.CancelAllResponse, error) {
	// if err := orderCancellation.Validate(); err != nil {
	//	 return err
	// }
	return order.CancelAllResponse{}, common.ErrNotYetImplemented
}

// GetOrderInfo returns order information based on order ID
func (ap *Apexpro) GetOrderInfo(ctx context.Context, orderID string, pair currency.Pair, assetType asset.Item) (*order.Detail, error) {
	return nil, common.ErrNotYetImplemented
}

// GetDepositAddress returns a deposit address for a specified currency
func (ap *Apexpro) GetDepositAddress(ctx context.Context, c currency.Code, accountID string, chain string) (*deposit.Address, error) {
	return nil, common.ErrNotYetImplemented
}

// WithdrawCryptocurrencyFunds returns a withdrawal ID when a withdrawal is
// submitted
func (ap *Apexpro) WithdrawCryptocurrencyFunds(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// WithdrawFiatFunds returns a withdrawal ID when a withdrawal is
// submitted
func (ap *Apexpro) WithdrawFiatFunds(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// WithdrawFiatFundsToInternationalBank returns a withdrawal ID when a withdrawal is
// submitted
func (ap *Apexpro) WithdrawFiatFundsToInternationalBank(ctx context.Context, withdrawRequest *withdraw.Request) (*withdraw.ExchangeResponse, error) {
	// if err := withdrawRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetActiveOrders retrieves any orders that are active/open
func (ap *Apexpro) GetActiveOrders(ctx context.Context, getOrdersRequest *order.MultiOrderRequest) (order.FilteredOrders, error) {
	// if err := getOrdersRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetOrderHistory retrieves account order information
// Can Limit response to specific order status
func (ap *Apexpro) GetOrderHistory(ctx context.Context, getOrdersRequest *order.MultiOrderRequest) (order.FilteredOrders, error) {
	// if err := getOrdersRequest.Validate(); err != nil {
	//	return nil, err
	// }
	return nil, common.ErrNotYetImplemented
}

// GetFeeByType returns an estimate of fee based on the type of transaction
func (ap *Apexpro) GetFeeByType(ctx context.Context, feeBuilder *exchange.FeeBuilder) (float64, error) {
	return 0, common.ErrNotYetImplemented
}

// ValidateAPICredentials validates current credentials used for wrapper
func (ap *Apexpro) ValidateAPICredentials(ctx context.Context, assetType asset.Item) error {
	_, err := ap.UpdateAccountInfo(ctx, assetType)
	return ap.CheckTransientError(err)
}

// GetHistoricCandles returns candles between a time period for a set time interval
func (ap *Apexpro) GetHistoricCandles(ctx context.Context, pair currency.Pair, a asset.Item, interval kline.Interval, start, end time.Time) (*kline.Item, error) {
	return nil, common.ErrNotYetImplemented
}

// GetHistoricCandlesExtended returns candles between a time period for a set time interval
func (ap *Apexpro) GetHistoricCandlesExtended(ctx context.Context, pair currency.Pair, a asset.Item, interval kline.Interval, start, end time.Time) (*kline.Item, error) {
	return nil, common.ErrNotYetImplemented
}

// GetFuturesContractDetails returns all contracts from the exchange by asset type
func (ap *Apexpro) GetFuturesContractDetails(context.Context, asset.Item) ([]futures.Contract, error) {
	return nil, common.ErrNotYetImplemented
}

// GetLatestFundingRates returns the latest funding rates data
func (ap *Apexpro) GetLatestFundingRates(_ context.Context, _ *fundingrate.LatestRateRequest) ([]fundingrate.LatestRateResponse, error) {
	return nil, common.ErrNotYetImplemented
}

// UpdateOrderExecutionLimits updates order execution limits
func (ap *Apexpro) UpdateOrderExecutionLimits(_ context.Context, _ asset.Item) error {
	return common.ErrNotYetImplemented
}
