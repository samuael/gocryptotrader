package kucoin

import (
	"encoding/json"
	"time"

	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
)

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsPositionStatus) UnmarshalJSON(data []byte) error {
	type Alias WsPositionStatus
	chil := &struct {
		*Alias
		TimestampMS int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.TimestampMS = time.UnixMilli(chil.TimestampMS)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsDebtRatioChange) UnmarshalJSON(data []byte) error {
	type Alias WsDebtRatioChange
	chil := &struct {
		*Alias
		TimestampMS int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.TimestampMS)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsMarginTradeOrderEntersEvent) UnmarshalJSON(data []byte) error {
	type Alias WsMarginTradeOrderEntersEvent
	chil := &struct {
		*Alias
		TimestampNS int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(0, chil.TimestampNS)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsMarginTradeOrderDoneEvent) UnmarshalJSON(data []byte) error {
	type Alias WsMarginTradeOrderDoneEvent
	chil := &struct {
		*Alias
		TimestampNS int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(0, chil.TimestampNS)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsStopOrder) UnmarshalJSON(data []byte) error {
	type Alias WsStopOrder
	chil := &struct {
		*Alias
		CreatedAt   int64 `json:"createdAt"`
		TimestampNS int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(0, chil.TimestampNS)
	a.CreatedAt = time.Unix(0, chil.CreatedAt)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsMarginFundingBook) UnmarshalJSON(data []byte) error {
	type Alias WsMarginFundingBook
	chil := &struct {
		*Alias
		TimestampNS int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(0, chil.TimestampNS)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsTradeOrder) UnmarshalJSON(data []byte) error {
	type Alias WsTradeOrder
	chil := &struct {
		*Alias
		OrderTime int64 `json:"orderTime"`
		Timestamp int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(chil.Timestamp/1e3, chil.Timestamp%1e3)
	a.OrderTime = time.Unix(chil.OrderTime/1e3, chil.OrderTime%1e3)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsAccountBalance) UnmarshalJSON(data []byte) error {
	type Alias WsAccountBalance
	chil := &struct {
		*Alias
		Time int64 `json:"time,string"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Time = time.UnixMilli(chil.Time)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesTicker) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesTicker
	chil := &struct {
		*Alias
		FilledTime int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.FilledTime = time.Unix(0, chil.FilledTime)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesOrderbokInfo) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesOrderbokInfo
	chil := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesExecutionData) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesExecutionData
	chil := &struct {
		*Alias
		Time int64 `json:"time"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Time = time.Unix(0, chil.Time)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesFundingBegin) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesFundingBegin
	chil := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesTransactionStatisticsTimeEvent) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesTransactionStatisticsTimeEvent
	chil := &struct {
		*Alias
		SnapshotTime int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.SnapshotTime = time.Unix(0, chil.SnapshotTime)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesTradeOrder) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesTradeOrder
	chil := &struct {
		*Alias
		Timestamp int64 `json:"ts"`
		OrderTime int64 `json:"orderTime"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.Unix(0, chil.Timestamp)
	a.OrderTime = time.Unix(0, chil.OrderTime)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsStopOrderLifecycleEvent) UnmarshalJSON(data []byte) error {
	type Alias WsStopOrderLifecycleEvent
	chil := &struct {
		*Alias
		CreatedAt int64 `json:"createdAt"`
		Timestamp int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.CreatedAt = time.UnixMilli(chil.CreatedAt)
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesOrderMarginEvent) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesOrderMarginEvent
	chil := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesAvailableBalance) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesAvailableBalance
	chil := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesWithdrawalAmountAndTransferOutAmountEvent) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesWithdrawalAmountAndTransferOutAmountEvent
	chil := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Timestamp = time.UnixMilli(chil.Timestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesPosition) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesPosition
	chil := &struct {
		*Alias
		OpeningTimestamp int64 `json:"openingTimestamp"` // Open time
		CurrentTimestamp int64 `json:"currentTimestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.OpeningTimestamp = time.UnixMilli(chil.OpeningTimestamp)
	a.CurrentTimestamp = time.UnixMilli(chil.CurrentTimestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesMarkPricePositionChanges) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesMarkPricePositionChanges
	chil := &struct {
		*Alias
		CurrentTimestamp int64 `json:"currentTimestamp"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.CurrentTimestamp = time.UnixMilli(chil.CurrentTimestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsFuturesPositionFundingSettlement) UnmarshalJSON(data []byte) error {
	type Alias WsFuturesPositionFundingSettlement
	chil := &struct {
		*Alias
		FundingTime      int64 `json:"fundingTime"`
		CurrentTimestamp int64 `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.FundingTime = time.UnixMilli(chil.FundingTime)
	a.CurrentTimestamp = time.Unix(0, chil.CurrentTimestamp)
	return nil
}

// UnmarshalJSON deserialises the JSON info, including the timestamp
func (a *WsOrderbookLevel5) UnmarshalJSON(data []byte) error {
	type Alias WsOrderbookLevel5
	chil := &struct {
		*Alias
		Asks      [][2]float64 `json:"asks"`
		Bids      [][2]float64 `json:"bids"`
		Timestamp int64        `json:"ts"`
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &chil); err != nil {
		return err
	}
	a.Asks = make([]orderbook.Item, len(chil.Asks))
	for x := range chil.Asks {
		a.Asks[x] = orderbook.Item{
			Price:  chil.Asks[x][0],
			Amount: chil.Asks[x][1],
		}
	}
	a.Bids = make([]orderbook.Item, len(chil.Bids))
	for x := range chil.Bids {
		a.Bids[x] = orderbook.Item{
			Price:  chil.Bids[x][0],
			Amount: chil.Bids[x][1],
		}
	}
	a.Timestamp = time.Unix(0, chil.Timestamp)
	return nil
}