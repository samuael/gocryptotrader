package apexpro

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/currency"
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

// ProcessOrderSignature processes order request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessOrderSignature(ctx context.Context, arg *CreateOrderParams) (string, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return "", err
	}
	if creds.L2Secret == "" {
		return "", starkex.ErrInvalidPrivateKey
	}
	if arg.ExpirationTime.IsZero() {
		return "", errExpirationTimeRequired
	}
	price := decimal.NewFromFloat(arg.Price)
	size := decimal.NewFromFloat(arg.Size)

	// check if the all symbols config is loaded, if not load
	if ap.SymbolsConfig == nil {
		ap.SymbolsConfig, err = ap.GetAllSymbolsConfigDataV1(ctx)
		if err != nil {
			return "", err
		}
	}

	var contractDetail *PerpetualContractDetail
	for a := range ap.SymbolsConfig.Data.PerpetualContract {
		if ap.SymbolsConfig.Data.PerpetualContract[a].Symbol == arg.Symbol.String() {
			contractDetail = &ap.SymbolsConfig.Data.PerpetualContract[a]
			if !contractDetail.EnableTrade {
				return "", currency.ErrPairNotEnabled
			}
		}
	}
	if contractDetail == nil {
		return "", fmt.Errorf("%w, contract: %s", errContractNotFound, arg.Symbol.String())
	}
	contractDetail.StarkExSyntheticAssetID = strings.TrimPrefix(contractDetail.StarkExSyntheticAssetID, "0x")
	syntheticAssetID, ok := big.NewInt(0).SetString(contractDetail.StarkExSyntheticAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, syntheticAssetId: %s", errInvalidAssetID, contractDetail.StarkExSyntheticAssetID)
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return "", err
		}
	}
	var currencyInfo *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == contractDetail.SettleCurrencyID {
			currencyInfo = &ap.SymbolsConfig.Data.Currency[c]
			break
		}
	}
	if currencyInfo == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	currencyInfo.StarkExAssetID = strings.TrimPrefix(currencyInfo.StarkExAssetID, "0x")
	assetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.ID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	resolution, err := decimal.NewFromString(contractDetail.StarkExResolution)
	if err != nil {
		return "", err
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
	return ap.StarkConfig.Sign(&starkex.CreateOrderWithFeeParams{
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
	}, creds.L2Secret)
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

// ProcessWithdrawalToAddressSignature processes withdrawal to specified ethereum address request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessWithdrawalToAddressSignature(ctx context.Context, arg *AssetWithdrawalParams) (string, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return "", err
	}
	var currencyInfo *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == arg.L1TargetTokenID.String() {
			currencyInfo = &ap.SymbolsConfig.Data.Currency[c]
			break
		}
	}
	println("arg.L1TargetTokenID.String(): ", arg.L1TargetTokenID.String())
	if currencyInfo == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return "", err
		}
	}
	if arg.ZKAccountID == "" {
		arg.ZKAccountID = ap.UserAccountDetail.SpotAccount.ZkAccountID
	}
	currencyInfo.StarkExAssetID = strings.TrimPrefix(currencyInfo.StarkExAssetID, "0x")
	collateralAssetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	if arg.EthereumAddress == "" {
		return "", errEthereumAddressMissing
	}
	arg.EthereumAddress = strings.TrimPrefix(arg.EthereumAddress, "0x")
	ethereumAddress, ok := big.NewInt(0).SetString(arg.EthereumAddress, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidEthereumAddress, arg.EthereumAddress)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.ID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	amount := decimal.NewFromFloat(arg.Amount)
	return ap.StarkConfig.Sign(&starkex.WithdrawalToAddressParams{
		AssetIDCollateral:    collateralAssetID,
		EthAddress:           ethereumAddress,
		PositionID:           positionID,
		Amount:               amount.Mul(resolution).BigInt(),
		Nonce:                NonceByClientId(arg.Nonce),
		ExpirationEpochHours: big.NewInt(int64(math.Ceil(float64(arg.Timestamp.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS),
	}, creds.L2Secret)
}

// ProcessWithdrawalSignature processes withdrawal request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessWithdrawalSignature(ctx context.Context, arg *WithdrawalParams) (string, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return "", err
	}
	var currencyInfo *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == arg.Asset.String() {
			currencyInfo = &ap.SymbolsConfig.Data.Currency[c]
			break
		}
	}
	if currencyInfo == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return "", err
		}
	}
	currencyInfo.StarkExAssetID = strings.TrimPrefix(currencyInfo.StarkExAssetID, "0x")
	collateralAssetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.ID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV3(context.Background())
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	amount := decimal.NewFromFloat(arg.Amount)
	return ap.StarkConfig.Sign(&starkex.WithdrawalParams{
		AssetIDCollateral:    collateralAssetID,
		PositionID:           positionID,
		Amount:               amount.Mul(resolution).BigInt(),
		Nonce:                NonceByClientId(arg.ClientID),
		ExpirationEpochHours: big.NewInt(int64(math.Ceil(float64(arg.ExpirationTime.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS),
	}, creds.L2Secret)
}
