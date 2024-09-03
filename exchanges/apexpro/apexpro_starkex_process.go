package apexpro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

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
			break
		}
	}
	if contractDetail == nil {
		return "", fmt.Errorf("%w, contract: %s", errContractNotFound, arg.Symbol.String())
	}
	syntheticAssetID, ok := big.NewInt(0).SetString(contractDetail.StarkExSyntheticAssetID, 0)
	if !ok {
		return "", fmt.Errorf("%w, syntheticAssetId: %s", errInvalidAssetID, contractDetail.StarkExSyntheticAssetID)
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	takerFeeRate := -1.
	for k := range ap.UserAccountDetail.Accounts {
		if ap.UserAccountDetail.Accounts[k].Token == contractDetail.SettleCurrencyID {
			takerFeeRate = ap.UserAccountDetail.Accounts[k].TakerFeeRate.Float64()
			break
		}
	}
	if takerFeeRate == -1. {
		return "", fmt.Errorf("%w, account with a settlement "+contractDetail.SettleCurrencyID+" is missing", errLimitFeeRequired)
	}
	arg.LimitFee = takerFeeRate * arg.Size * arg.Price
	var collateralAsset *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == contractDetail.SettleCurrencyID {
			collateralAsset = &ap.SymbolsConfig.Data.Currency[c]
			break
		}
	}
	if collateralAsset == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	collateralAssetID, ok := big.NewInt(0).SetString(collateralAsset.StarkExAssetID, 0)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, collateralAsset.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 10)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	syntheticResolution, err := decimal.NewFromString(contractDetail.StarkExResolution)
	if err != nil {
		return "", err
	}
	collateralResolution, err := decimal.NewFromString(collateralAsset.StarkExResolution)
	if err != nil {
		return "", err
	}
	arg.Side = strings.ToUpper(arg.Side)
	isBuy := arg.Side == "BUY"
	var quantumsAmountCollateral decimal.Decimal
	if isBuy {
		quantumsAmountCollateral = size.Mul(price).Mul(collateralResolution).RoundUp(0)
	} else {
		quantumsAmountCollateral = size.Mul(price).Mul(collateralResolution).RoundDown(0)
	}
	quantumsAmountSynthetic := size.Mul(syntheticResolution)
	limitFeeRounded := decimal.NewFromFloat(takerFeeRate)
	// arg.ClientOrderID = "10001"
	arg.ClientOrderID = randomClientID()
	// expEpoch, _ := big.NewInt(0).SetString("1100000", 0)
	expEpoch := big.NewInt(int64(math.Ceil(float64(time.Now().Add(time.Hour*24*28).UnixMilli()) / float64(3600*1000))))
	arg.ExpirationTime = expEpoch.Int64() * 3600 * 1000
	newArg := &starkex.CreateOrderWithFeeParams{
		OrderType:               "LIMIT_ORDER_WITH_FEES",
		AssetIDSynthetic:        syntheticAssetID,
		AssetIDCollateral:       collateralAssetID,
		AssetIDFee:              collateralAssetID,
		QuantumAmountSynthetic:  quantumsAmountSynthetic.BigInt(),
		QuantumAmountCollateral: quantumsAmountCollateral.BigInt(),
		QuantumAmountFee:        limitFeeRounded.Mul(quantumsAmountCollateral).RoundUp(0).BigInt(),
		IsBuyingSynthetic:       isBuy,
		PositionID:              positionID,
		Nonce:                   nonceFromClientID(arg.ClientOrderID),
		ExpirationEpochHours:    expEpoch,
	}

	val, _ := json.Marshal(newArg)
	println(string(val))

	return ap.StarkConfig.Sign(newArg, creds.L2Secret, creds.L2Key, creds.L2KeyYCoordinate)
}

// ProcessWithdrawalToAddressSignatureV3 processes withdrawal to specified ethereum address request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessWithdrawalToAddressSignatureV3(ctx context.Context, arg *AssetWithdrawalParams) (string, error) {
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
	if currencyInfo == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	if arg.ZKAccountID == "" {
		arg.ZKAccountID = ap.UserAccountDetail.ID
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
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	amount := decimal.NewFromFloat(arg.Amount)
	// expEpoch := big.NewInt(int64(math.Ceil(float64(time.Now().Add(time.Hour*24*28).UnixMilli()) / float64(3600*1000))))
	// arg.E = expEpoch.Int64() * 3600 * 1000
	return ap.StarkConfig.Sign(&starkex.WithdrawalToAddressParams{
		AssetIDCollateral:    collateralAssetID,
		EthAddress:           ethereumAddress,
		PositionID:           positionID,
		Amount:               amount.Mul(resolution).BigInt(),
		Nonce:                nonceFromClientID(arg.Nonce),
		ExpirationEpochHours: big.NewInt(int64(math.Ceil(float64(arg.Timestamp.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS),
	}, creds.L2Secret, creds.L2Key, creds.L2KeyYCoordinate)
}

// ProcessWithdrawalToAddressSignature processes withdrawal to specified ethereum address request parameter and generates a starkEx signature for V1 and V2 api endpoints
func (ap *Apexpro) ProcessWithdrawalToAddressSignature(ctx context.Context, arg *WithdrawalToAddressParams) (string, error) {
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
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	collateralAssetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	if arg.EthereumAddress == "" {
		return "", errEthereumAddressMissing
	}
	ethereumAddress, ok := big.NewInt(0).SetString(arg.EthereumAddress, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidEthereumAddress, arg.EthereumAddress)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	expEpoch := big.NewInt(int64(math.Ceil(float64(arg.ExpirationTime.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS)
	arg.ExpEpoch = expEpoch.Int64()
	amount := decimal.NewFromFloat(arg.Amount)
	return ap.StarkConfig.Sign(&starkex.WithdrawalToAddressParams{
		AssetIDCollateral:    collateralAssetID,
		EthAddress:           ethereumAddress,
		PositionID:           positionID,
		Amount:               amount.Mul(resolution).BigInt(),
		Nonce:                nonceFromClientID(arg.ClientID),
		ExpirationEpochHours: expEpoch,
	}, creds.L2Secret, "", "")
}

// ProcessWithdrawalSignature processes withdrawal request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessWithdrawalSignature(ctx context.Context, arg *WithdrawalParams) (string, error) {
	creds, err := ap.GetCredentials(context.Background())
	if err != nil {
		return "", err
	}
	var collateralInfo *V1CurrencyConfig
	for c := range ap.SymbolsConfig.Data.Currency {
		if ap.SymbolsConfig.Data.Currency[c].ID == arg.Asset.String() {
			collateralInfo = &ap.SymbolsConfig.Data.Currency[c]
			break
		}
	}
	if collateralInfo == nil {
		return "", errSettlementCurrencyInfoNotFound
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	collateralAssetID, ok := big.NewInt(0).SetString(collateralInfo.StarkExAssetID, 0)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, collateralInfo.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	collateralResolution, err := decimal.NewFromString(collateralInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	arg.ClientID = randomClientID()
	amount := decimal.NewFromFloat(arg.Amount)
	expEpoch := big.NewInt(int64(math.Ceil(float64(arg.ExpirationTime.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS)
	arg.ExpEpoch = expEpoch.Int64()
	return ap.StarkConfig.Sign(&starkex.WithdrawalParams{
		AssetIDCollateral:    collateralAssetID,
		PositionID:           positionID,
		Amount:               amount.Mul(collateralResolution).BigInt(),
		Nonce:                nonceFromClientID(arg.ClientID),
		ExpirationEpochHours: big.NewInt(int64(math.Ceil(float64(arg.ExpirationTime.Unix()) / float64(3600)))),
	}, creds.L2Secret, "", "")
}

// ProcessTransferSignature processes withdrawal request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessTransferSignature(ctx context.Context, arg *FastWithdrawalParams) (string, error) {
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
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	collateralAssetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	arg.ClientID = randomClientID()
	amount := decimal.NewFromFloat(arg.Amount)
	return ap.StarkConfig.Sign(&starkex.TransferParams{
		AssetID:              collateralAssetID,
		AssetIDFee:           big.NewInt(0),
		SenderPositionID:     positionID,
		QuantumsAmount:       amount.Mul(resolution).BigInt(),
		Nonce:                nonceFromClientID(arg.ClientID),
		ExpirationEpochHours: big.NewInt(int64(math.Ceil(float64(arg.Expiration.Unix())/float64(3600))) + ORDER_SIGNATURE_EXPIRATION_BUFFER_HOURS),
	}, creds.L2Secret, "", "")
}

// ProcessConditionalTransfer processes conditional transfer request parameter and generates a starkEx signature
func (ap *Apexpro) ProcessConditionalTransfer(ctx context.Context, arg *FastWithdrawalParams) (string, error) {
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
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	// currencyInfo.StarkExAssetID = strings.TrimPrefix(currencyInfo.StarkExAssetID, "0x")
	collateralAssetID, ok := big.NewInt(0).SetString(currencyInfo.StarkExAssetID, 16)
	if !ok {
		return "", fmt.Errorf("%w, assetId: %s", errInvalidAssetID, currencyInfo.StarkExAssetID)
	}
	positionID, ok := big.NewInt(0).SetString(ap.UserAccountDetail.PositionID, 0)
	if !ok {
		return "", errInvalidPositionIDMissing
	}
	if ap.UserAccountDetail == nil {
		ap.UserAccountDetail, err = ap.GetUserAccountDataV2(ctx)
		if err != nil {
			return "", err
		}
	}
	resolution, err := decimal.NewFromString(currencyInfo.StarkExResolution)
	if err != nil {
		return "", err
	}
	amount := decimal.NewFromFloat(arg.Amount)
	arg.ClientID = randomClientID()
	return ap.StarkConfig.Sign(&starkex.ConditionalTransferParams{
		AssetID:          collateralAssetID,
		AssetIDFee:       big.NewInt(0),
		SenderPositionID: positionID,
		QuantumsAmount:   amount.Mul(resolution).BigInt(),
		Nonce:            nonceFromClientID(arg.ClientID),
	}, creds.L2Secret, "", "")
}
