package starkex

import "math/big"

type Signable interface {
	GetPedersenHash(pedersenHash func(...string) string) (string, error)
}

// GetPedersenHash implements the Signable interface and generates a pedersen hash of CreateOrderWithFeeParams
func (s *CreateOrderWithFeeParams) GetPedersenHash(pedersenHash func(...string) string) (string, error) {
	var assetIdSell, assetIdBuy, quantumsAmountSell, quantumsAmountBuy *big.Int
	if s.IsBuyingSynthetic {
		assetIdSell = s.AssetIdCollateral
		assetIdBuy = s.AssetIdSynthetic
		quantumsAmountSell = s.QuantumAmountCollateral
		quantumsAmountBuy = s.QuantumAmountSynthetic
	} else {
		assetIdSell = s.AssetIdSynthetic
		assetIdBuy = s.AssetIdCollateral
		quantumsAmountSell = s.QuantumAmountSynthetic
		quantumsAmountBuy = s.QuantumAmountCollateral
	}
	fee := s.QuantumAmountFee
	nonce := s.Nonce
	// part1
	part1 := big.NewInt(0).Set(quantumsAmountSell)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["quantums_amount"])
	part1.Add(part1, quantumsAmountBuy)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["quantums_amount"])
	part1.Add(part1, fee)
	part1.Lsh(part1, ORDER_FIELD_BIT_LENGTHS["nonce"])
	part1.Add(part1, nonce)
	// part2
	part2 := big.NewInt(ORDER_PREFIX)
	for i := 0; i < 3; i++ {
		part2.Lsh(part2, ORDER_FIELD_BIT_LENGTHS["position_id"])
		part2.Add(part2, s.PositionId)
	}
	part2.Lsh(part2, ORDER_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	part2.Add(part2, s.ExpirationEpochHours)
	part2.Lsh(part2, ORDER_PADDING_BITS)
	// pedersen hash
	assetHash := pedersenHash(pedersenHash(assetIdSell.String(), assetIdBuy.String()), s.AssetIdFee.String())
	part1Hash := pedersenHash(assetHash, part1.String())
	return pedersenHash(part1Hash, part2.String()), nil
}

// GetPedersenHash implements the Signable interface and generates a pedersen hash of WithdrawalToAddressParams
func (s *WithdrawalToAddressParams) GetPedersenHash(pedersenHash func(...string) string) (string, error) {
	// packed
	packed := big.NewInt(WITHDRAWAL_TO_ADDRESS_PREFIX)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["position_id"])
	packed.Add(packed, s.PositionID)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["nonce"])
	packed.Add(packed, s.Nonce)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["quantums_amount"])
	packed.Add(packed, s.Amount)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	packed.Add(packed, s.ExpirationEpochHours)
	packed.Lsh(packed, WITHDRAWAL_PADDING_BITS)
	// pedersen hash
	return pedersenHash(pedersenHash(s.AssetIDCollateral.String(), s.EthAddress.String()), packed.String()), nil
}

// GetPedersenHash implements the Signable interface and generates a pedersen hash of WithdrawalParams
func (s *WithdrawalParams) GetPedersenHash(pedersenHash func(...string) string) (string, error) {
	// packed
	packed := big.NewInt(WITHDRAWAL_PREFIX)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["position_id"])
	packed.Add(packed, s.PositionID)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["nonce"])
	packed.Add(packed, s.Nonce)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["quantums_amount"])
	packed.Add(packed, s.Amount)
	packed.Lsh(packed, WITHDRAWAL_FIELD_BIT_LENGTHS["expiration_epoch_hours"])
	packed.Add(packed, s.ExpirationEpochHours)
	packed.Lsh(packed, WITHDRAWAL_PADDING_BITS)
	// pedersen hash
	return pedersenHash(s.AssetIDCollateral.String(), packed.String()), nil
}
