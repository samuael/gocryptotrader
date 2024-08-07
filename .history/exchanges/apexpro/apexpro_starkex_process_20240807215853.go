package apexpro

import (
	"context"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/starkex"
)

var (
	errExpirationTimeRequired         = errors.New("expiration time is required")
	errContractNotFound               = errors.New("contract not found")
	errSettlementCurrencyInfoNotFound = errors.New("settlement currency information not found")
	errInvalidAssetID                 = errors.New("invalid asset ID provided")
	errInvalidPositionIDMissing       = errors.New("invalid position or account ID")
)

const ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS = 24 * 7 // Seven days.
const NONCE_UPPER_BOUND_EXCLUSIVE = 1 << 32            // 1 << ORDER_FIELD_BIT_LENGTHS['nonce']

// ProcessOrder processes order request parameter and generates a CreateOrderWithFeeParams, which can be used later for generating signature for
func (ap *Apexpro) ProcessOrder(ctx context.Context, arg *CreateOrderParams) (*starkex.CreateOrderWithFeeParams, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return nil, err
	}
	if creds.L2Secret == "" {
		return nil, starkex.ErrInvalidPrivateKey
	}
	if arg.ExpirationTime.IsZero() {
		return nil, errExpirationTimeRequired
	}
	price := decimal.NewFromFloat(arg.Price)
	size := decimal.NewFromFloat(arg.Size)

	format, err := ap.GetPairFormat(asset.Futures, true)
	if err != nil {
		return nil, err
	}
	arg.Symbol = arg.Symbol.Format(format)

	// check if the all symbols config is loaded, if not load
	if ap.SymbolsConfig == nil {
		ap.SymbolsConfig, err = ap.GetAllConfigDataV1(ctx)
		if err != nil {
			return nil, err
		}
	}

	var contractDetail *PerpetualContractDetail
	for a := range ap.SymbolsConfig.Data.PerpetualContract {
		if ap.SymbolsConfig.Data.PerpetualContract[a].Symbol == arg.Symbol.String() {
			contractDetail = &ap.SymbolsConfig.Data.PerpetualContract[a]
			if contractDetail.EnableTrade {
				return nil, currency.ErrPairNotEnabled
			}
		}
	}
	if contractDetail == nil {
		return nil, errContractNotFound
	}
	syntheticAssetID, ok := big.NewInt(0).SetString(contractDetail.StarkExSyntheticAssetID, 16)
	if !ok {
		return nil, errInvalidAssetID
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return nil, err
		}
	}
	var currencyInfo *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == contractDetail.SettleCurrencyID {
			currencyInfo = &ap.SymbolsConfig.Data.Currency[c]
		}
	}
	if currencyInfo == nil {
		return nil, errSettlementCurrencyInfoNotFound
	}
	assetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return nil, errInvalidAssetID
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.ID, 0)
	if !ok {
		return nil, errInvalidPositionIDMissing
	}
	resolution, err := decimal.NewFromString(contractDetail.StarkExResolution)
	if err != nil {
		return nil, err
	}
	var quantumsAmountSynthetic = decimal.NewFromFloat(0)
	arg.Side = strings.ToUpper(arg.Side)
	isBuy := arg.Side == "BUY"
	if isBuy {
		quantumsAmountSynthetic = size.Mul(price).Mul(resolution).RoundUp(0)
	} else {
		quantumsAmountSynthetic = size.Mul(price).Mul(resolution).RoundDown(0)
	}
	limitFeeRounded := decimal.NewFromFloat(arg.LimitFee)
	if arg.Nonce == "" {
		arg.Nonce = strconv.FormatInt(ap.Websocket.Conn.GenerateMessageID(true), 10)
	}
	return &starkex.CreateOrderWithFeeParams{
		OrderType:               "LIMIT_ORDER_WITH_FEES",
		AssetIdSynthetic:        syntheticAssetID,
		AssetIdCollateral:       assetID,
		AssetIdFee:              assetID,
		QuantumAmountSynthetic:  size.Mul(resolution).BigInt(),
		QuantumAmountCollateral: quantumsAmountSynthetic.BigInt(),
		QuantumAmountFee:        limitFeeRounded.Mul(quantumsAmountSynthetic).RoundUp(0).BigInt(),
		IsBuyingSynthetic:       isBuy,
		PositionId:              positionID,
		Nonce:                   NonceByClientId(arg.Nonce),
		ExpirationEpochHours:    big.NewInt(int64(math.Ceil(float64(arg.ExpirationTime.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS),
	}, nil
}

// NonceByClientId generate nonce by clientId
func NonceByClientId(clientId string) *big.Int {
	h := sha256.New()
	h.Write([]byte(clientId))

	a := new(big.Int)
	a.SetBytes(h.Sum(nil))
	res := a.Mod(a, big.NewInt(NONCE_UPPER_BOUND_EXCLUSIVE))
	return res
}
