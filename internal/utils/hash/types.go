package hash

import "math/big"

// / First 64 bytes of the BLAKE2s input during group hash.
// / This is chosen to be some random string that we couldn't have anticipated when we designed
// / the algorithm, for rigidity purposes.
// / We deliberately use an ASCII hex string of 32 bytes here.
var GH_FIRST_BLOCK = []byte("096b36a5804bfacef1691e173c366a47ff5ba84a44f26ddd7e8d9f79d5b42df0")

// PedersenCfg represents a pedersen hash configuration options
type PedersenCfg struct {
	Comment        string        `json:"_comment"`
	FieldPrime     *big.Int      `json:"FIELD_PRIME"`
	FieldGen       int           `json:"FIELD_GEN"`
	EcOrder        *big.Int      `json:"EC_ORDER"`
	ALPHA          int           `json:"ALPHA"`
	BETA           *big.Int      `json:"BETA"`
	ConstantPoints [][2]*big.Int `json:"CONSTANT_POINTS"`
}

// Bn256RescueParams represents capacity, rate of hash, and other details of the rescue hash algorithm
type Bn256RescueParams struct {
	C              uint32
	R              uint32
	Rounds         uint32
	SecurityLevel  uint32
	RoundConstants []*big.Int
	MDSMatrix      []*big.Int
	SBox0          *PowerSBox
	SBox1          *QuinticSBox

	CustomGatesAllowed bool
}
