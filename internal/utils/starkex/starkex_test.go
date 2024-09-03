package starkex

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const MOCK_PUBLIC_KEY = "0x3b865a18323b8d147a12c556bfb1d502516c325b1477a23ba6c77af31f020fd"
const MOCK_PRIVATE_KEY = "0x58c7d5a90b1776bde86ebac077e053ed85b0f7164f53b080304a531947f46e3"

func TestECDSASignature(t *testing.T) {
	t.Parallel()
	magHash, ok := big.NewInt(0).SetString("0x011049f4032190ec4b5a9420cc77006d13a260df46bfcacf60a53f447a5a925d", 0)
	require.True(t, ok)

	publicX, ok := big.NewInt(0).SetString(MOCK_PUBLIC_KEY, 0)
	require.True(t, ok)

	publicSecret, ok := big.NewInt(0).SetString(MOCK_PRIVATE_KEY, 0)
	require.True(t, ok)

	sfg, err := NewStarkExConfig("apexpro")
	require.NoError(t, err)
	require.NotNil(t, sfg)

	r, s, err := sfg.SignECDSA(magHash, publicSecret)
	require.NoError(t, err)

	publicY := sfg.GetYCoordinate(publicX)
	ok = sfg.Verify(magHash, r, s, [2]*big.Int{publicX, publicY})
	require.True(t, ok,
		ErrFailedToGenerateSignature)
}

func TestOrderSign(t *testing.T) {
	t.Parallel()
	sfg, err := NewStarkExConfig("apexpro")
	require.NoError(t, err)
	require.NotNil(t, sfg)

	syntheticAssetID, ok := big.NewInt(0).SetString("344400637343183300222065759427231744", 10)
	require.True(t, ok)

	collateralAssetID, ok := big.NewInt(0).SetString("1147032829293317481173155891309375254605214077236177772270270553197624560221", 10)
	require.True(t, ok)

	arg := &CreateOrderWithFeeParams{
		OrderType:               "LIMIT_ORDER_WITH_FEES",
		AssetIDSynthetic:        syntheticAssetID,
		AssetIDCollateral:       collateralAssetID,
		AssetIDFee:              collateralAssetID,
		QuantumAmountSynthetic:  big.NewInt(100000000),
		QuantumAmountCollateral: big.NewInt(200000000),
		QuantumAmountFee:        big.NewInt(100000),
		IsBuyingSynthetic:       false,
		PositionID:              big.NewInt(603545650545558021),
		Nonce:                   big.NewInt(3762202436),
		ExpirationEpochHours:    big.NewInt(479941),
	}
	signature, err := sfg.Sign(arg, MOCK_PRIVATE_KEY, MOCK_PUBLIC_KEY, "")
	require.NoError(t, err)
	assert.NotEmpty(t, signature)
}

func TestGetYCoordinate(t *testing.T) {
	t.Parallel()
	sfg, err := NewStarkExConfig("apexpro")
	require.NoError(t, err)
	require.NotNil(t, sfg)

	publicX, ok := big.NewInt(0).SetString(MOCK_PUBLIC_KEY, 0)
	require.True(t, ok)

	result := sfg.GetYCoordinate(publicX)
	require.NotNil(t, result)
}
