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

	PedersenHashGenerators     []interface{}
	PedersenHashExp            []interface{}
	PedersenCircuitGenerators  []interface{}
	FixedBaseGenerators        []interface{}
	FixedBaseCircuitGenerators []interface{}
}

func NewAltJubJubBn256() *AltJubJubBn256 {
	edwardsD, _ := big.NewInt(0).SetString("12181644023421730124874158521699555681764249180949974110617291017600649128846", 0)
	montegomeryA := big.NewInt(168698)

	doubleMontgomery2A := new(big.Int).Mul(montegomeryA, montegomeryA)
	scale, _ := big.NewInt(0).SetString("6360561867910373094066688120553762416144456282423235903351243436111059670888", 0)

	return &AltJubJubBn256{
		MontgometryA:               montegomeryA,
		EdwardsD:                   edwardsD,
		Montgometry2A:              doubleMontgomery2A,
		Scale:                      scale,
		PedersenHashGenerators:     make([]interface{}, 0),
		PedersenHashExp:            make([]interface{}, 0),
		PedersenCircuitGenerators:  make([]interface{}, 0),
		FixedBaseGenerators:        make([]interface{}, 0),
		FixedBaseCircuitGenerators: make([]interface{}, 0),
	}
}

// const defaultJubjubConfigsPath = "internal/utils/hash/elliptic_curve_config/"

// Sign generates a jubjub signature from string
func (jb *JubJubConfig) Sign(privateKey, message string) (*big.Int, *big.Int, error) {
	return nil, nil, nil
}
