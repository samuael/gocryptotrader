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
	part2Hash := pedersenHash(part1Hash, part2.String())
	return part2Hash, nil
}
