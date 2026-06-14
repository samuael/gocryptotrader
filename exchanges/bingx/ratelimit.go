package bingx

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Rate limit intervals.
const (
	oneSecond  = time.Second
	tenSeconds = time.Second * 10
)

const (
	spotV1CommonSymbolsEPL request.EndpointLimit = iota + 1
	spotV1MarketTradesEPL
	spotV1MarketDepthEPL
	spotV2MarketKlineEPL
	spotV1Ticker24hrEPL
	spotV2MarketDepthEPL
	spotV2TickerPriceEPL
	spotV1TickerBookTickerEPL
	marketHisV1KlineEPL
	marketHisV1TradeEPL
	spotV1TradeOrderEPL
	spotV1TradeBatchOrdersEPL
	spotV1TradeCancelEPL
	spotV1TradeCancelOrdersEPL
	spotV1TradeCancelOpenOrdersEPL
	spotV1TradeOrderCancelReplaceEPL
	spotV1TradeQueryEPL
	spotV1TradeOpenOrdersEPL
	spotV1TradeHistoryOrdersEPL
	spotV1TradeMyTradesEPL
	spotV1UserCommissionRateEPL
	spotV1TradeCancelAllAfterEPL
	spotV1OcoOrderEPL
	spotV1OcoCancelEPL
	spotV1OcoOrderListEPL
	spotV1OcoOpenOrderListEPL
	spotV1OcoHistoryOrderListEPL
	swapV2QuoteContractsEPL
	swapV2QuoteDepthEPL
	swapV2QuoteTradesEPL
	swapV2QuotePremiumIndexEPL
	swapV2QuoteFundingRateEPL
	swapV3QuoteKlinesEPL
	swapV2QuoteOpenInterestEPL
	swapV2QuoteTickerEPL
	swapV1MarketHistoricalTradesEPL
	swapV2QuoteBookTickerEPL
	swapV1MarketMarkPriceKlinesEPL
	swapV1TickerPriceEPL
	swapV1TradingRulesEPL
	swapV2TradeOrderTestEPL
	swapPlaceOrderEPL
	swapCancelOrderEPL
	swapQueryOrderEPL
	swapV1TradeAmendEPL
	swapV2TradeBatchOrdersEPL
	swapV2TradeCloseAllPositionsEPL
	swapV2TradeAllOpenOrdersEPL
	swapV2TradeOpenOrdersEPL
	swapV2TradeOpenOrderEPL
	swapV2TradeMarginTypeEPL
	swapV2TradeLeverageEPL
	swapV2TradeForceOrdersEPL
	swapV2TradeAllOrdersEPL
	swapV2TradePositionMarginEPL
	swapV2TradeAllFillOrdersEPL
	swapSetPositionModeEPL
	swapGetPositionModeEPL
	swapV1TradeCancelReplaceEPL
	swapV1TradeBatchCancelReplaceEPL
	swapV2TradeCancelAllAfterEPL
	swapV1TradeClosePositionEPL
	swapV1TradeFullOrderEPL
	swapV1MaintMarginRatioEPL
	swapV2TradeFillHistoryEPL
	swapV1TradePositionHistoryEPL
	swapV1PositionMarginHistoryEPL
	swapV2TradeGetVstEPL
	swapV1TwapOrderEPL
	swapV1TwapOpenOrdersEPL
	swapV1TwapHistoryOrdersEPL
	swapV1TwapOrderDetailEPL
	swapV1TwapCancelOrderEPL
	swapV1TradeAssetModeEPL
	swapV1TradeMultiAssetsRulesEPL
	swapV1UserMarginAssetsEPL
	swapV1TradeReverseEPL
	swapV1TradeAutoAddMarginEPL
	coinMV1MarketContractsEPL
	coinMV1MarketPremiumIndexEPL
	coinMV1MarketOpenInterestEPL
	coinMV1MarketKlinesEPL
	coinMV1MarketDepthEPL
	coinMV1MarketTickerEPL
	coinMV1TradeOrderEPL
	coinMV1UserCommissionRateEPL
	coinMV1TradeLeverageEPL
	coinMV1TradeAllOpenOrdersEPL
	coinMV1TradeCloseAllPositionsEPL
	coinMV1UserPositionsEPL
	coinMV1UserBalanceEPL
	coinMV1TradeForceOrdersEPL
	coinMV1TradeAllFillOrdersEPL
	coinMV1TradeCancelOrderEPL
	coinMV1TradeOpenOrdersEPL
	coinMV1TradeOrderDetailEPL
	coinMV1TradeOrderHistoryEPL
	coinMV1TradeMarginTypeEPL
	coinMV1TradePositionMarginEPL
)

var rateLimits = request.RateLimitDefinitions{
	spotV1CommonSymbolsEPL:           request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV1MarketTradesEPL:            request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV1MarketDepthEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV2MarketKlineEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV1Ticker24hrEPL:              request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV2MarketDepthEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV2TickerPriceEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV1TickerBookTickerEPL:        request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	marketHisV1KlineEPL:              request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	marketHisV1TradeEPL:              request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	spotV1TradeOrderEPL:              request.NewRateLimitWithWeight(oneSecond, 10, 1),
	spotV1TradeBatchOrdersEPL:        request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1TradeCancelEPL:             request.NewRateLimitWithWeight(oneSecond, 5, 1),
	spotV1TradeCancelOrdersEPL:       request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1TradeCancelOpenOrdersEPL:   request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1TradeOrderCancelReplaceEPL: request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1TradeQueryEPL:              request.NewRateLimitWithWeight(oneSecond, 10, 1),
	spotV1TradeOpenOrdersEPL:         request.NewRateLimitWithWeight(oneSecond, 10, 1),
	spotV1TradeHistoryOrdersEPL:      request.NewRateLimitWithWeight(oneSecond, 10, 1),
	spotV1TradeMyTradesEPL:           request.NewRateLimitWithWeight(oneSecond, 5, 1),
	spotV1UserCommissionRateEPL:      request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1TradeCancelAllAfterEPL:     request.NewRateLimitWithWeight(oneSecond, 1, 1),
	spotV1OcoOrderEPL:                request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1OcoCancelEPL:               request.NewRateLimitWithWeight(oneSecond, 2, 1),
	spotV1OcoOrderListEPL:            request.NewRateLimitWithWeight(oneSecond, 5, 1),
	spotV1OcoOpenOrderListEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	spotV1OcoHistoryOrderListEPL:     request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2QuoteContractsEPL:          request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteDepthEPL:              request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteTradesEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuotePremiumIndexEPL:       request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteFundingRateEPL:        request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV3QuoteKlinesEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteOpenInterestEPL:       request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteTickerEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV1MarketHistoricalTradesEPL:  request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV2QuoteBookTickerEPL:         request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV1MarketMarkPriceKlinesEPL:   request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV1TickerPriceEPL:             request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	swapV1TradingRulesEPL:            request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeOrderTestEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapPlaceOrderEPL:                request.NewRateLimitWithWeight(oneSecond, 10, 1),
	swapCancelOrderEPL:               request.NewRateLimitWithWeight(oneSecond, 10, 1),
	swapQueryOrderEPL:                request.NewRateLimitWithWeight(oneSecond, 30, 1),
	swapV1TradeAmendEPL:              request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeBatchOrdersEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeCloseAllPositionsEPL:  request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeAllOpenOrdersEPL:      request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeOpenOrdersEPL:         request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeOpenOrderEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeMarginTypeEPL:         request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV2TradeLeverageEPL:           request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeForceOrdersEPL:        request.NewRateLimitWithWeight(oneSecond, 10, 1),
	swapV2TradeAllOrdersEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradePositionMarginEPL:     request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV2TradeAllFillOrdersEPL:      request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapSetPositionModeEPL:           request.NewRateLimitWithWeight(oneSecond, 4, 1),
	swapGetPositionModeEPL:           request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TradeCancelReplaceEPL:      request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1TradeBatchCancelReplaceEPL: request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV2TradeCancelAllAfterEPL:     request.NewRateLimitWithWeight(oneSecond, 1, 1),
	swapV1TradeClosePositionEPL:      request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1TradeFullOrderEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1MaintMarginRatioEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeFillHistoryEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1TradePositionHistoryEPL:    request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1PositionMarginHistoryEPL:   request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV2TradeGetVstEPL:             request.NewRateLimitWithWeight(oneSecond, 5, 1),
	swapV1TwapOrderEPL:               request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TwapOpenOrdersEPL:          request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TwapHistoryOrdersEPL:       request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TwapOrderDetailEPL:         request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TwapCancelOrderEPL:         request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TradeAssetModeEPL:          request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TradeMultiAssetsRulesEPL:   request.NewRateLimitWithWeight(tenSeconds, 300, 1),
	swapV1UserMarginAssetsEPL:        request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TradeReverseEPL:            request.NewRateLimitWithWeight(oneSecond, 2, 1),
	swapV1TradeAutoAddMarginEPL:      request.NewRateLimitWithWeight(oneSecond, 2, 1),
	coinMV1MarketContractsEPL:        request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1MarketPremiumIndexEPL:     request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1MarketOpenInterestEPL:     request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1MarketKlinesEPL:           request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1MarketDepthEPL:            request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1MarketTickerEPL:           request.NewRateLimitWithWeight(tenSeconds, 500, 1),
	coinMV1TradeOrderEPL:             request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1UserCommissionRateEPL:     request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeLeverageEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeAllOpenOrdersEPL:     request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeCloseAllPositionsEPL: request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1UserPositionsEPL:          request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1UserBalanceEPL:            request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeForceOrdersEPL:       request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeAllFillOrdersEPL:     request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeCancelOrderEPL:       request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeOpenOrdersEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeOrderDetailEPL:       request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeOrderHistoryEPL:      request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradeMarginTypeEPL:        request.NewRateLimitWithWeight(oneSecond, 5, 1),
	coinMV1TradePositionMarginEPL:    request.NewRateLimitWithWeight(oneSecond, 5, 1),
}
