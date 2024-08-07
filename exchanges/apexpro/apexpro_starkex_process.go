package apexpro

import (
	"context"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/starkex"
)

var (
	errExpirationTimeRequired = errors.New("expiration time is required")
	errContractNotFound       = errors.New("contract not found")
)

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
	var positionID string // contract position
	var assetID string
	syntheticAssetID = big.NewFromString(contractDetail.StarkExSyntheticAssetID)
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return nil, err
		}
	}
	positionID = ap.UserAccountDetail.ID
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
	argNew := &starkex.CreateOrderWithFeeParams{
		// OrderType
		// AssetIdSynthetic
		// AssetIdCollateral
		// AssetIdFee
		// QuantumAmountSynthetic
		// QuantumAmountCollateral
		// QuantumAmountFee
		// IsBuyingSynthetic
		// PositionId
		// Nonce
		// ExpirationEpochHours
	}
}
