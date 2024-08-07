package starkex

import (
	"crypto/sha256"
	"math/big"
)

// func  initMsg() error {
// 	exp, err := time.Parse("2006-01-02T15:04:05.000Z", s.param.Expiration)
// 	if err != nil {
// 		return err
// 	}
// 	QuantumAmount, err := decimal.NewFromString(s.param.HumanAmount)
// 	if err != nil {
// 		return err
// 	}
// 	s.msg.QuantumAmount = QuantumAmount.Mul(resolutionUsdc).BigInt()
// 	s.msg.PositionId = big.NewInt(s.param.PositionId)
// 	s.msg.Nonce = NonceByClientId(s.param.ClientId)
// 	s.msg.ExpirationEpochHours = big.NewInt(int64(math.Ceil(float64(exp.Unix()) / float64(ONE_HOUR_IN_SECONDS))))
// 	return nil
// }

// NonceByClientId generate nonce by clientId
func NonceByClientId(clientId string) *big.Int {
	h := sha256.New()
	h.Write([]byte(clientId))

	a := new(big.Int)
	a.SetBytes(h.Sum(nil))
	res := a.Mod(a, big.NewInt(1<<32))
	return res
}

// func signWithdraw(param *WithdrawalToAddressParams, resolution, collateralAssetID string) (string, error) {
// 	quantumAmount := decimal.New(int64(param.Amount), 1).Mul(resolution).BigInt()
// 	nonce := big.NewInt(int64(param.Nonce))
// 	expHrs := big.NewInt(int64(math.Ceil(float64(param.ExpirationTimestamp) / float64(3600))))

// 	net := COLLATERAL_ASSET_ID_BY_NETWORK_ID[param.NetworkID]
// 	if net == nil {
// 		return "", errors.New(fmt.Sprintf("invalid network_id: %v", param.NetworkID))
// 	}
// 	// packed
// 	packed := big.NewInt(6)
// 	packed.Lsh(packed, 64)
// 	packed.Add(packed, big.NewInt(int64(param.PositionID)))
// 	packed.Lsh(packed, 32)
// 	packed.Add(packed, nonce)
// 	packed.Lsh(packed, 63)
// 	packed.Add(packed, quantumAmount)
// 	packed.Lsh(packed, 32)
// 	packed.Add(packed, expHrs)
// 	packed.Lsh(packed, 49)
// 	// pedersen hash
// 	cfg, err := hash.LoadPedersenConfig("internal/utils/hash/pedersen_config/apexpro.json")
// 	if err != nil {
// 		return "", err
// 	}
// 	return cfg.PedersenHash(net.String(), packed.String()), nil
// }
