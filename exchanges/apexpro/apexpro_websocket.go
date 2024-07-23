package apexpro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
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
	return nil
}

// GenerateDefaultSubscriptions generates a default subscription list.
func (ap *Apexpro) GenerateDefaultSubscriptions() (subscription.List, error) {
	subscriptions := subscription.List{}
	enabledPairs, err := ap.GetEnabledPairs(asset.Futures)
	if err != nil {
		return subscriptions, err
	}
	for a := range defaultChannels {
		switch defaultChannels[a] {
		case chOrderbook:
			subscriptions = append(subscriptions, &subscription.Subscription{
				Channel: defaultChannels[a],
				Pairs:   enabledPairs,
				Levels:  100,
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
				Levels:   100,
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
	case chTrade:
	case chTicker:
	case chCandlestick:
	case chAllTickers:
	}
	return nil
}
