package starkex

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// FactToCondition Generate the condition, signed as part of a conditional transfer.
func FactToCondition(factRegistryAddress, fact string) *big.Int {
	data := strings.TrimPrefix(factRegistryAddress, "0x") + fact
	hexBytes, _ := hex.DecodeString(data)
	hash := crypto.Keccak256Hash(hexBytes)
	fst := hash.Big()
	fst.And(fst, BitMask250)
	return fst
}
