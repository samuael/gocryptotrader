package starkex

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestECDSASignature(t *testing.T) {
	t.Parallel()
	magHash, ok := big.NewInt(0).SetString("0x011049f4032190ec4b5a9420cc77006d13a260df46bfcacf60a53f447a5a925d", 0)
	require.True(t, ok)

	publicX, ok := big.NewInt(0).SetString("0xf8c6635f9cfe85f46759dc2eebe71a45b765687e35dbe5e74e8bde347813ef", 0)
	require.True(t, ok)
	// publicSecret, ok := big.NewInt(0).SetString("0x607ba3969039f3e19006ff8f40629d20a7b7dac31d4019e0965fbf7c5c068a", 0)
	// require.True(t, ok)

	// publicY, ok := big.NewInt(0).SetString("0x0207d57867e0820e0f7588339e8b7491ce1da964260044340e3fd27c718f2a91", 0)
	// require.True(t, ok)

	sfg, err := NewStarkExConfig("apexpro")
	require.NoError(t, err)
	require.NotNil(t, sfg)

	publicY := sfg.GetYCoordinate(publicX)
	require.True(t, ok)

	r, ok := big.NewInt(0).SetString("0x07a15838aad9b20368dc4ba27613fd35ceec3b34be7a2cb913bca0fb06e98107", 0)
	require.True(t, ok)
	s, ok := big.NewInt(0).SetString("0x05007f40fddd9babae0c7362d3b4e9c152ed3fced7fe78435b302d825489298f", 0)
	require.True(t, ok)

	// r, s, err := sfg.SignECDSA(magHash, publicSecret)
	// require.NoError(t, err)

	// publicKeyYCoordinate, ok := big.NewInt(0).SetString(starkPublicKeyYCoordinate, 0)
	// if !ok {
	// 	publicKeyYCoordinate = sfg.GetYCoordinate(publicKey)
	// 	if publicKeyYCoordinate.Cmp(big.NewInt(0)) == 0 {
	// 		return "", fmt.Errorf("%w, invalid stark public key x coordinat", ErrInvalidPublicKey)
	// 	}
	// }
	ok = sfg.Verify(magHash, r, s, [2]*big.Int{publicX, publicY})
	require.True(t, ok,
		ErrFailedToGenerateSignature)
}
