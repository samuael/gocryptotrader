package bingx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

func TestRateLimit_LimitStatic(t *testing.T) {
	t.Parallel()
	testTable := map[string]request.EndpointLimit{
		"spotV1CommonSymbolsEPL":           spotV1CommonSymbolsEPL,
		"spotV1MarketTradesEPL":            spotV1MarketTradesEPL,
		"spotV1MarketDepthEPL":             spotV1MarketDepthEPL,
		"spotV2MarketKlineEPL":             spotV2MarketKlineEPL,
		"spotV1Ticker24hrEPL":              spotV1Ticker24hrEPL,
		"spotV2MarketDepthEPL":             spotV2MarketDepthEPL,
		"spotV2TickerPriceEPL":             spotV2TickerPriceEPL,
		"spotV1TickerBookTickerEPL":        spotV1TickerBookTickerEPL,
		"marketHisV1KlineEPL":              marketHisV1KlineEPL,
		"marketHisV1TradeEPL":              marketHisV1TradeEPL,
		"spotV1TradeOrderEPL":              spotV1TradeOrderEPL,
		"spotV1TradeBatchOrdersEPL":        spotV1TradeBatchOrdersEPL,
		"spotV1TradeCancelEPL":             spotV1TradeCancelEPL,
		"spotV1TradeCancelOrdersEPL":       spotV1TradeCancelOrdersEPL,
		"spotV1TradeCancelOpenOrdersEPL":   spotV1TradeCancelOpenOrdersEPL,
		"spotV1TradeOrderCancelReplaceEPL": spotV1TradeOrderCancelReplaceEPL,
		"spotV1TradeQueryEPL":              spotV1TradeQueryEPL,
		"spotV1TradeOpenOrdersEPL":         spotV1TradeOpenOrdersEPL,
		"spotV1TradeHistoryOrdersEPL":      spotV1TradeHistoryOrdersEPL,
		"spotV1TradeMyTradesEPL":           spotV1TradeMyTradesEPL,
		"spotV1UserCommissionRateEPL":      spotV1UserCommissionRateEPL,
		"spotV1TradeCancelAllAfterEPL":     spotV1TradeCancelAllAfterEPL,
		"spotV1OcoOrderEPL":                spotV1OcoOrderEPL,
		"spotV1OcoCancelEPL":               spotV1OcoCancelEPL,
		"spotV1OcoOrderListEPL":            spotV1OcoOrderListEPL,
		"spotV1OcoOpenOrderListEPL":        spotV1OcoOpenOrderListEPL,
		"spotV1OcoHistoryOrderListEPL":     spotV1OcoHistoryOrderListEPL,
		"swapV2QuoteContractsEPL":          swapV2QuoteContractsEPL,
		"swapV2QuoteDepthEPL":              swapV2QuoteDepthEPL,
		"swapV2QuoteTradesEPL":             swapV2QuoteTradesEPL,
		"swapV2QuotePremiumIndexEPL":       swapV2QuotePremiumIndexEPL,
		"swapV2QuoteFundingRateEPL":        swapV2QuoteFundingRateEPL,
		"swapV3QuoteKlinesEPL":             swapV3QuoteKlinesEPL,
		"swapV2QuoteOpenInterestEPL":       swapV2QuoteOpenInterestEPL,
		"swapV2QuoteTickerEPL":             swapV2QuoteTickerEPL,
		"swapV1MarketHistoricalTradesEPL":  swapV1MarketHistoricalTradesEPL,
		"swapV2QuoteBookTickerEPL":         swapV2QuoteBookTickerEPL,
		"swapV1MarketMarkPriceKlinesEPL":   swapV1MarketMarkPriceKlinesEPL,
		"swapV1TickerPriceEPL":             swapV1TickerPriceEPL,
		"swapV1TradingRulesEPL":            swapV1TradingRulesEPL,
		"swapV2TradeOrderTestEPL":          swapV2TradeOrderTestEPL,
		"swapPlaceOrderEPL":                swapPlaceOrderEPL,
		"swapCancelOrderEPL":               swapCancelOrderEPL,
		"swapQueryOrderEPL":                swapQueryOrderEPL,
		"swapV1TradeAmendEPL":              swapV1TradeAmendEPL,
		"swapV2TradeBatchOrdersEPL":        swapV2TradeBatchOrdersEPL,
		"swapV2TradeCloseAllPositionsEPL":  swapV2TradeCloseAllPositionsEPL,
		"swapV2TradeAllOpenOrdersEPL":      swapV2TradeAllOpenOrdersEPL,
		"swapV2TradeOpenOrdersEPL":         swapV2TradeOpenOrdersEPL,
		"swapV2TradeOpenOrderEPL":          swapV2TradeOpenOrderEPL,
		"swapV2TradeMarginTypeEPL":         swapV2TradeMarginTypeEPL,
		"swapV2TradeLeverageEPL":           swapV2TradeLeverageEPL,
		"swapV2TradeForceOrdersEPL":        swapV2TradeForceOrdersEPL,
		"swapV2TradeAllOrdersEPL":          swapV2TradeAllOrdersEPL,
		"swapV2TradePositionMarginEPL":     swapV2TradePositionMarginEPL,
		"swapV2TradeAllFillOrdersEPL":      swapV2TradeAllFillOrdersEPL,
		"swapSetPositionModeEPL":           swapSetPositionModeEPL,
		"swapGetPositionModeEPL":           swapGetPositionModeEPL,
		"swapV1TradeCancelReplaceEPL":      swapV1TradeCancelReplaceEPL,
		"swapV1TradeBatchCancelReplaceEPL": swapV1TradeBatchCancelReplaceEPL,
		"swapV2TradeCancelAllAfterEPL":     swapV2TradeCancelAllAfterEPL,
		"swapV1TradeClosePositionEPL":      swapV1TradeClosePositionEPL,
		"swapV1TradeFullOrderEPL":          swapV1TradeFullOrderEPL,
		"swapV1MaintMarginRatioEPL":        swapV1MaintMarginRatioEPL,
		"swapV2TradeFillHistoryEPL":        swapV2TradeFillHistoryEPL,
		"swapV1TradePositionHistoryEPL":    swapV1TradePositionHistoryEPL,
		"swapV1PositionMarginHistoryEPL":   swapV1PositionMarginHistoryEPL,
		"swapV2TradeGetVstEPL":             swapV2TradeGetVstEPL,
		"swapV1TwapOrderEPL":               swapV1TwapOrderEPL,
		"swapV1TwapOpenOrdersEPL":          swapV1TwapOpenOrdersEPL,
		"swapV1TwapHistoryOrdersEPL":       swapV1TwapHistoryOrdersEPL,
		"swapV1TwapOrderDetailEPL":         swapV1TwapOrderDetailEPL,
		"swapV1TwapCancelOrderEPL":         swapV1TwapCancelOrderEPL,
		"swapV1TradeAssetModeEPL":          swapV1TradeAssetModeEPL,
		"swapV1TradeMultiAssetsRulesEPL":   swapV1TradeMultiAssetsRulesEPL,
		"swapV1UserMarginAssetsEPL":        swapV1UserMarginAssetsEPL,
		"swapV1TradeReverseEPL":            swapV1TradeReverseEPL,
		"swapV1TradeAutoAddMarginEPL":      swapV1TradeAutoAddMarginEPL,
		"coinMV1MarketContractsEPL":        coinMV1MarketContractsEPL,
		"coinMV1MarketPremiumIndexEPL":     coinMV1MarketPremiumIndexEPL,
		"coinMV1MarketOpenInterestEPL":     coinMV1MarketOpenInterestEPL,
		"coinMV1MarketKlinesEPL":           coinMV1MarketKlinesEPL,
		"coinMV1MarketDepthEPL":            coinMV1MarketDepthEPL,
		"coinMV1MarketTickerEPL":           coinMV1MarketTickerEPL,
		"coinMV1TradeOrderEPL":             coinMV1TradeOrderEPL,
		"coinMV1UserCommissionRateEPL":     coinMV1UserCommissionRateEPL,
		"coinMV1TradeLeverageEPL":          coinMV1TradeLeverageEPL,
		"coinMV1TradeAllOpenOrdersEPL":     coinMV1TradeAllOpenOrdersEPL,
		"coinMV1TradeCloseAllPositionsEPL": coinMV1TradeCloseAllPositionsEPL,
		"coinMV1UserPositionsEPL":          coinMV1UserPositionsEPL,
		"coinMV1UserBalanceEPL":            coinMV1UserBalanceEPL,
		"coinMV1TradeForceOrdersEPL":       coinMV1TradeForceOrdersEPL,
		"coinMV1TradeAllFillOrdersEPL":     coinMV1TradeAllFillOrdersEPL,
		"coinMV1TradeCancelOrderEPL":       coinMV1TradeCancelOrderEPL,
		"coinMV1TradeOpenOrdersEPL":        coinMV1TradeOpenOrdersEPL,
		"coinMV1TradeOrderDetailEPL":       coinMV1TradeOrderDetailEPL,
		"coinMV1TradeOrderHistoryEPL":      coinMV1TradeOrderHistoryEPL,
		"coinMV1TradeMarginTypeEPL":        coinMV1TradeMarginTypeEPL,
		"coinMV1TradePositionMarginEPL":    coinMV1TradePositionMarginEPL,
	}
	rl, err := request.New("coinMV1TradePositionMarginEPL", http.DefaultClient, request.WithLimiter(rateLimits))
	require.NoError(t, err)

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.NoError(t, rl.InitiateRateLimit(t.Context(), tt), "InitiateRateLimit must not error")
		})
	}
}
