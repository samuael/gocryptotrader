package zklink

import (
	"errors"
	"math/big"
)

var (
	errSignatureFailed      = errors.New("signature failed")
	errInvalidPrivateKey    = errors.New("invalid private key")
	errInvalidSeed          = errors.New("invalid seed")
	errInvalidPublicKey     = errors.New("invalid public key")
	errInvalidPublicKeyHash = errors.New("invalid public key hash")

	errInvalidEthSigner     = errors.New("invalid eth signer")
	errMissingEthPrivateKey = errors.New("Ethereum private key required to perform an operation")
	errMissingEthSigner     = errors.New("EthereumSigner required to perform an operation")
	errSigningFailed        = errors.New("Signing failed: {0}")
	errUnlockingFailed      = errors.New("Unlocking failed")
	errInvalidRawTx         = errors.New("Decode raw transaction failed")
	errEip712Failed         = errors.New("Eip712 error")
	errNoSigningKey         = errors.New("signing key is not set in account")
	errDefineAddress        = errors.New("address determination error")
	errRecoverAddress       = errors.New("recover address from signature failed: {0}")
	errLengthMismatched     = errors.New("signature length mismatch")
	errCryptoError          = errors.New("Crypto Error")
	errInvalidETHSignature  = errors.New("invalid eth signature string")
)

const PACKED_POINT_SIZE = 32
const SIGNATURE_SIZE = 96

const NEW_PUBKEY_HASH_BYTES_LEN = 20
const NEW_PUBKEY_HASH_WIDTH = NEW_PUBKEY_HASH_BYTES_LEN * 8

type ZkLinkSigner struct {
}

type ZkLinkSignature struct {
}

// ContractBuilder holds a contract builder parameters
type ContractBuilder struct {
	AccountID    *big.Int
	SubAccountID *big.Int
	SlotID       *big.Int
	Nonce        *big.Int
	PairID       *big.Int
	Size         *big.Int
	Price        *big.Int
	Direction    bool
	TakerFeeRate *big.Int
	MakerFeeRate *big.Int
	HasSubsidy   bool
}

// WithdrawBuilder holds an asset withdrawal builder parameters
type WithdrawBuilder struct {
	AccountID    *big.Int
	SubAccountID *big.Int
	ToChainID    *big.Int
	ToAddress    *big.Int
	// L2SourceToken TokenId
	// L1TargetToken TokenId
	Amount *big.Int
	// DataHash         *H256
	Fee              *big.Int
	Nonce            *big.Int
	WithdrawFeeRatio *big.Int
	WithdrawToL1     bool
	Timestamp        *big.Int
}

// TransferBuilder holds an asset transfer builder parameters for zklink signature
type TransferBuilder struct {
	AccountID        *big.Int
	ToAddress        *big.Int
	FromSubAccountID *big.Int
	ToSubAccountID   *big.Int
	Token            *big.Int
	Amount           *big.Int
	Fee              *big.Int
	Nonce            *big.Int
	Timestamp        *big.Int
}

// Bn256RescueParams
type Bn256RescueParams struct {
	C             uint64
	R             uint64
	Rounds        uint64
	SecurityLevel uint64

	CustomGatesAllowed bool
}
