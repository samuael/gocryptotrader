package apexpro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
)

const (
	apexProWebsocket        = "wss://qa-quote.omni.apex.exchange/realtime_public"
	apexProPrivateWebsocket = "wss://quote.omni.apex.exchange/realtime_private"

	chOrderbook   = "orderBook"
	chTrade       = "recentlyTrade"
	chTicker      = "instrumentInfo"
	chAllTickers  = "instrumentInfo.all"
	chCandlestick = "candle"

	// Authenticated websocket channels
)

var defaultChannels = []string{
	chOrderbook, chTrade, chTicker, chCandlestick, chAllTickers,
}

func generatePingMessage() ([]byte, error) {
	return json.Marshal(&WsMessage{
		Operation: "ping",
		Args:      []string{strconv.FormatInt(time.Now().UnixMilli(), 10)},
	})
}

// WsConnect creates a websocket connection
func (ap *Apexpro) WsConnect() error {
	if !ap.Websocket.IsEnabled() || !ap.IsEnabled() {
		return stream.ErrWebsocketNotEnabled
	}
	var dialer websocket.Dialer
	dialer.HandshakeTimeout = ap.Config.HTTPTimeout
	dialer.Proxy = http.ProxyFromEnvironment
	var err error
	err = ap.Websocket.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return fmt.Errorf("%v - Unable to connect to Websocket. Error: %s",
			ap.Name,
			err)
	}
	if ap.Websocket.CanUseAuthenticatedEndpoints() {
		err = ap.WsAuthConnect()
		if err != nil {
			ap.Websocket.SetCanUseAuthenticatedEndpoints(false)
		}
	}
	payload, err := generatePingMessage()
	if err != nil {
		return err
	}
	ap.Websocket.Conn.SetupPingHandler(stream.PingHandler{
		UseGorillaHandler: true,
		MessageType:       websocket.PongMessage,
		Message:           payload,
	})
	ap.Websocket.Wg.Add(1)
	go ap.wsReadData(ap.Websocket.Conn)
	subscriptions, err := ap.GenerateDefaultSubscriptions()
	if err != nil {
		return err
	}
	return ap.Subscribe(subscriptions)
	// return nil
}

// GenerateDefaultSubscriptions generates a default subscription list.
func (ap *Apexpro) GenerateDefaultSubscriptions() (subscription.List, error) {
	subscriptions := subscription.List{}
	// enabledPairs, err := ap.GetEnabledPairs(asset.Futures)
	// if err != nil {
	// 	return subscriptions, err
	// }
	enabledPairs := []currency.Pair{{Base: currency.BTC, Quote: currency.USDT}}
	for a := range defaultChannels {
		switch defaultChannels[a] {
		case chOrderbook:
			subscriptions = append(subscriptions, &subscription.Subscription{
				Channel: defaultChannels[a],
				Pairs:   enabledPairs,
				Levels:  200,
			})
		case chTrade, chTicker:
			subscriptions = append(subscriptions, &subscription.Subscription{
				Channel:  defaultChannels[a],
				Pairs:    enabledPairs,
				Interval: kline.HundredMilliseconds,
			})
		case chCandlestick:
			subscriptions = append(subscriptions, &subscription.Subscription{
				Channel:  defaultChannels[a],
				Pairs:    enabledPairs,
				Levels:   200,
				Interval: kline.FiveMin,
			})
		case chAllTickers:
			subscriptions = append(subscriptions, &subscription.Subscription{
				Channel: defaultChannels[a],
			})
		}
	}
	return subscriptions, nil
}

// WsAuthConnect creates a websocket connection and authenticates the private stream connection.
func (ap *Apexpro) WsAuthConnect() error {
	return nil
}

// Subscribe sends a websocket channel subscription.
func (ap *Apexpro) Subscribe(subscriptions subscription.List) error {
	payload, err := ap.handleSubscriptionPayload("subscribe", subscriptions)
	if err != nil {
		return err
	}
	err = ap.Websocket.Conn.SendJSONMessage(payload)
	if err != nil {
		return err
	}
	return ap.Websocket.AddSuccessfulSubscriptions(subscriptions...)
}

// Unsubscribe sends a websocket channel unsubscriptions.
func (ap *Apexpro) Unsubscribe(subscriptions subscription.List) error {
	payload, err := ap.handleSubscriptionPayload("unsubscribe", subscriptions)
	if err != nil {
		return err
	}
	return ap.Websocket.Conn.SendJSONMessage(payload)
}

func (ap *Apexpro) handleSubscriptionPayload(operation string, subscriptions subscription.List) (*WsMessage, error) {
	susbcriptionPayload := &WsMessage{
		Operation: operation,
		Args:      []string{},
	}
	pairFormat, err := ap.GetPairFormat(asset.Futures, true)
	if err != nil {
		return nil, err
	}
	for s := range subscriptions {
		subscriptions[s].Pairs = subscriptions[s].Pairs.Format(pairFormat)
		switch subscriptions[s].Channel {
		case chOrderbook:
			if subscriptions[s].Levels == 0 {
				return nil, errOrderbookLevelIsRequired
			}
			for p := range subscriptions[s].Pairs {
				susbcriptionPayload.Args = append(susbcriptionPayload.Args, subscriptions[s].Channel+strconv.Itoa(subscriptions[s].Levels)+".H."+subscriptions[s].Pairs[p].String())
			}
		case chTrade, chTicker:
			for p := range subscriptions[s].Pairs {
				susbcriptionPayload.Args = append(susbcriptionPayload.Args, subscriptions[s].Channel+".H."+subscriptions[s].Pairs[p].String())
			}
		case chCandlestick:
			if subscriptions[s].Interval == kline.Interval(0) {
				return nil, kline.ErrInvalidInterval
			}
			intervalString, err := intervalToString(subscriptions[s].Interval)
			if err != nil {
				return nil, err
			}
			for p := range subscriptions[s].Pairs {
				susbcriptionPayload.Args = append(susbcriptionPayload.Args, subscriptions[s].Channel+"."+intervalString+"."+subscriptions[s].Pairs[p].String())
			}
		case chAllTickers:
			susbcriptionPayload.Args = append(susbcriptionPayload.Args, subscriptions[s].Channel)
		}
	}
	return susbcriptionPayload, nil
}

func (ap *Apexpro) wsReadData(conn stream.Connection) {
	defer ap.Websocket.Wg.Done()
	for {
		response := conn.ReadMessage()
		if response.Raw == nil {
			return
		}
		err := ap.wsHandleData(response.Raw)
		if err != nil {
			ap.Websocket.DataHandler <- err
		}
	}
}

func (ap *Apexpro) wsHandleData(respRaw []byte) error {
	var response WsMessage
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	switch response.Operation {
	case "pong":
	case chOrderbook:
		return ap.processOrderbook(respRaw)
	case chTrade:
		return ap.processTrades(respRaw)
	case chTicker:
		return ap.processTickerData(respRaw)
	case chCandlestick:
		return ap.processCandlestickData(respRaw)
	case chAllTickers:
		return ap.processAllTickers(respRaw)
	}
	return nil
}

func (ap *Apexpro) processOrderbook(respRaw []byte) error {
	var resp *WsDepth
	var cp currency.Pair
	err := json.Unmarshal(respRaw, &resp)
	if err != nil {
		return err
	}
	cp, err = currency.NewPairFromString(resp.Data.Symbol)
	if err != nil {
		return err
	}
	asks := make(orderbook.Tranches, len(resp.Data.Asks))
	for a := range resp.Data.Asks {
		asks[a].Price = resp.Data.Asks[a][0].Float64()
		asks[a].Amount = resp.Data.Asks[a][1].Float64()
	}
	bids := make(orderbook.Tranches, len(resp.Data.Bids))
	for b := range resp.Data.Bids {
		bids[b].Price = resp.Data.Bids[b][0].Float64()
		bids[b].Amount = resp.Data.Bids[b][1].Float64()
	}
	if resp.Type == "delta" {
		return ap.Websocket.Orderbook.Update(&orderbook.Update{
			Bids:       bids,
			Asks:       asks,
			Pair:       cp,
			UpdateID:   resp.Data.UpdateID,
			UpdateTime: resp.Timestamp.Time(),
			Asset:      asset.Futures,
		})
	}
	return ap.Websocket.Orderbook.LoadSnapshot(&orderbook.Base{
		Pair:            cp,
		Asset:           asset.Spot,
		Exchange:        ap.Name,
		LastUpdateID:    resp.Data.UpdateID,
		VerifyOrderbook: ap.CanVerifyOrderbook,
		LastUpdated:     time.Now(),
		Asks:            asks,
		Bids:            bids,
	})
}

func (ap *Apexpro) processTrades(respRaw []byte) error {
	var resp *WsTrade
	err := json.Unmarshal(respRaw, &resp)
	if err != nil {
		return err
	}
	saveTradeData := ap.IsSaveTradeDataEnabled()
	if !saveTradeData &&
		!ap.IsTradeFeedEnabled() {
		return nil
	}
	trades := make([]trade.Data, len(resp.Data))
	for a := range resp.Data {
		cp, err := currency.NewPairFromString(resp.Data[a].Symbol)
		if err != nil {
			return err
		}
		trades[a] = trade.Data{
			CurrencyPair: cp,
			Timestamp:    resp.Data[a].Timestamp.Time(),
			Price:        resp.Data[a].Price.Float64(),
			Amount:       resp.Data[a].Volume.Float64(),
			Exchange:     ap.Name,
			AssetType:    asset.Futures,
			TID:          resp.Data[a].OrderID,
		}
	}
	return ap.Websocket.Trade.Update(saveTradeData, trades...)
}

func (ap *Apexpro) processTickerData(respRaw []byte) error {
	var resp *WsTicker
	err := json.Unmarshal(respRaw, &resp)
	if err != nil {
		return err
	}
	cp, err := currency.NewPairFromString(resp.Data.Symbol)
	if err != nil {
		return err
	}
	ap.Websocket.DataHandler <- ticker.Price{
		Last:         resp.Data.LastPrice.Float64(),
		High:         resp.Data.HighPrice24H.Float64(),
		Low:          resp.Data.LowPrice24H.Float64(),
		Volume:       resp.Data.Volume24H.Float64(),
		OpenInterest: resp.Data.OpenInterest.Float64(),
		MarkPrice:    resp.Data.OraclePrice.Float64(),
		IndexPrice:   resp.Data.IndexPrice.Float64(),
		Pair:         cp,
		ExchangeName: ap.Name,
		AssetType:    asset.Futures,
	}
	return nil
}

func (ap *Apexpro) processCandlestickData(respRaw []byte) error {
	var resp *WsCandlesticks
	err := json.Unmarshal(respRaw, &resp)
	if err != nil {
		return err
	}
	for a := range resp.Data {
		pair, err := currency.NewPairFromString(resp.Data[a].Symbol)
		if err != nil {
			return err
		}
		ap.Websocket.DataHandler <- stream.KlineData{
			Timestamp:  resp.Timestamp.Time(),
			Pair:       pair,
			AssetType:  asset.Futures,
			Exchange:   ap.Name,
			StartTime:  resp.Data[a].Start.Time(),
			Interval:   resp.Data[a].Interval,
			OpenPrice:  resp.Data[a].Open.Float64(),
			ClosePrice: resp.Data[a].Close.Float64(),
			HighPrice:  resp.Data[a].High.Float64(),
			LowPrice:   resp.Data[a].Low.Float64(),
			Volume:     resp.Data[a].Volume.Float64(),
		}
	}
	return nil
}

func (ap *Apexpro) processAllTickers(respRaw []byte) error {
	var resp *WsSymbolsTickerInformaton
	err := json.Unmarshal(respRaw, &resp)
	if err != nil {
		return err
	}
	tickerData := make([]ticker.Price, len(resp.Data))
	for a := range resp.Data {
		pair, err := currency.NewPairFromString(resp.Data[a].Symbol)
		if err != nil {
			return err
		}
		tickerData[a] = ticker.Price{
			Last:         resp.Data[a].LastPrice.Float64(),
			High:         resp.Data[a].Highest24Hr.Float64(),
			Low:          resp.Data[a].Lowest24Hr.Float64(),
			Volume:       resp.Data[a].Volume24Hr.Float64(),
			Open:         resp.Data[a].OpeningPrice.Float64(),
			OpenInterest: resp.Data[a].OpenInterest.Float64(),
			MarkPrice:    resp.Data[a].MarkPrice.Float64(),
			IndexPrice:   resp.Data[a].IndexPrice.Float64(),
			Pair:         pair,
			ExchangeName: ap.Name,
			AssetType:    asset.Futures,
			LastUpdated:  resp.Timestamp.Time(),
		}
	}
	ap.Websocket.DataHandler <- tickerData
	return nil
}
