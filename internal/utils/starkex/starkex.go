package starkex

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	path "github.com/thrasher-corp/gocryptotrader/internal/testing/utils"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/hash"
	math_utils "github.com/thrasher-corp/gocryptotrader/internal/utils/math"
)

// Error declarations.
var (
	ErrInvalidPrivateKey  = errors.New("invalid private key")
	ErrInvalidHashPayload = errors.New("invalid hash payload")
)

// StarkConfig represents a stark configuration
type StarkConfig struct {
	PedersenCfg *hash.PedersenCfg
}

var zero = big.NewInt(0)
var one = big.NewInt(1)

const defaultPedersenConfigsPath = "internal/utils/hash/pedersen_config/"

// NewStarkExConfig returns a stark configuration given the exchange name
func NewStarkExConfig(exchangeName string) (*StarkConfig, error) {
	rootPath, err := path.RootPathFromCWD()
	if err != nil {
		return nil, err
	}
	pedersenConfig, err := hash.LoadPedersenConfig(rootPath + "/" + defaultPedersenConfigsPath + strings.ToLower(exchangeName) + ".json")
	if err != nil {
		return nil, err
	}
	return &StarkConfig{
		PedersenCfg: pedersenConfig,
	}, nil
}

// Sign generates a signature out using the users private key and signable order params.
func (sfg *StarkConfig) Sign(sgn Signable, starkPrivateKey string) (string, error) {
	println(sfg.PedersenCfg.FieldPrime.String())
	println(sfg.PedersenCfg.EcOrder.String())
	println(sfg.PedersenCfg.BETA.String())
	pHash, err := sgn.GetPedersenHash(sfg.PedersenCfg.PedersenHash)
	if err != nil {
		return pHash, err
	}
	r, s, err := sfg.ECDSAHash(pHash, starkPrivateKey)
	if err != nil {
		return "", err
	}
	// Concatenate the hex strings
	return math_utils.IntToHex32(r) + math_utils.IntToHex32(s), nil
}

// ECDSAHash generates an ECDSA signature given the private key and pedersen signed hash
func (sfg *StarkConfig) ECDSAHash(message, starkPrivateKey string) (*big.Int, *big.Int, error) {
	priKey, okay := big.NewInt(0).SetString(starkPrivateKey, 0)
	if !okay {
		return nil, nil, fmt.Errorf("%w, %v", ErrInvalidPrivateKey, starkPrivateKey)
	}
	msgHash, okay := new(big.Int).SetString(message, 10)
	if !okay {
		return nil, nil, ErrInvalidHashPayload
	}
	seed := 0
	EcGen := sfg.PedersenCfg.ConstantPoints[1]
	alpha := sfg.PedersenCfg.ALPHA
	nBit := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(251), nil)
	for {
		k := math_utils.GenerateKRfc6979(msgHash, priKey, sfg.PedersenCfg.EcOrder, seed)
		//	Update seed for next iteration in case the value of k is bad.
		if seed == 0 {
			seed = 1
		} else {
			seed += 1
		}
		// Cannot fail because 0 < k < EC_ORDER and EC_ORDER is prime.
		x := math_utils.ECMult(k, EcGen, alpha, sfg.PedersenCfg.FieldPrime)[0]
		// !(1 <= x < 2 ** N_ELEMENT_BITS_ECDSA)
		if !(x.Cmp(one) > 0 && x.Cmp(nBit) < 0) {
			continue
		}
		// msg_hash + r * priv_key
		x1 := big.NewInt(0).Add(msgHash, big.NewInt(0).Mul(x, priKey))
		// (msg_hash + r * priv_key) % EC_ORDER == 0
		if big.NewInt(0).Mod(x1, sfg.PedersenCfg.EcOrder).Cmp(zero) == 0 {
			continue
		}
		// w = div_mod(k, msg_hash + r * priv_key, EC_ORDER)
		w := math_utils.DivMod(k, x1, sfg.PedersenCfg.EcOrder)
		// not (1 <= w < 2 ** N_ELEMENT_BITS_ECDSA)
		if !(w.Cmp(one) > 0 && w.Cmp(nBit) < 0) {
			continue
		}
		s1 := math_utils.DivMod(one, w, sfg.PedersenCfg.EcOrder)
		return x, s1, nil
	}
}
