package zklink

import "math/big"

// NewZkLinkSigner creates a new zklink signer instance.
func NewZkLinkSigner() *ZkLinkSigner {
	return &ZkLinkSigner{}
}

// Sign generates a new zkline signature instance from signable interfaces.
func (sr *ZkLinkSigner) Sign(arg Signable) (*ZkLinkSignature, error) {
	return sr.SignMusig(arg.GetBytes())
}

// SignMusig uses musig Schnorr signature scheme.
// It is impossible to restore signer for signature, that is why we provide public key of the signer
// along with signature.
func (sg *ZkLinkSigner) SignMusig(msg *big.Int) (*ZkLinkSignature, error) {
	return nil, nil
}

// func (sg *ZkLinkSigner) SignMusig(msg []uint8) (*ZkLinkSignature, error) {
