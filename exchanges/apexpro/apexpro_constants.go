package apexpro

import "math/big"

var one = big.NewInt(1)

// BIT_MASK_250 (2 ** 250) - 1
var BIT_MASK_250 = big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(250), nil), one)
