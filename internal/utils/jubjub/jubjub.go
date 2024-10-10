package jubjub

import "math/big"

// JubJubConfig represents configurations parameters on twisted edwards curves
type JubJubConfig struct {
	PrimeField   *big.Int   `json:"p"`
	ACoefficient *big.Int   `json:"a"`
	DCoefficient *big.Int   `json:"d"`
	Generator    []*big.Int `json:"g"`
	Order        *big.Int   `json:"n"`
	Coefactor    *big.Int   `json:"c"`
}

type AltJubJubBn256 struct {
	EdwardsD      *big.Int
	MontgometryA  *big.Int
	Montgometry2A *big.Int
	Scale         *big.Int

	// PedersenHashGenerators
}

// const defaultJubjubConfigsPath = "internal/utils/hash/elliptic_curve_config/"

// Sign generates a jubjub signature from string
func (jb *JubJubConfig) Sign(privateKey, message string) (*big.Int, *big.Int, error) {
	return nil, nil, nil
}
