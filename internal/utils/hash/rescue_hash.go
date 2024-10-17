package hash

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"
)

const PAD_MSG_BEFORE_HASH_BITS_LEN = 736

func RescueHashTransactionMsg(message *big.Int) (interface{}, error) {
	// msgBits := BytesIntoBits(message)
	if len(message.Bits()) <= PAD_MSG_BEFORE_HASH_BITS_LEN {
		return nil, errors.New("invalid message")
	}
	// bits := BytesIntoBits(message.Bytes())
	return nil, nil
}

func BytesIntoBits(message []byte) []bool {
	msg := make([]bool, len(message)*8)
	for _, b := range message {
		x := b
		for i := range 8 {
			msg = append(msg, (x&(1<<i)) != 0)
		}
	}
	return msg
}

// PackBitsIntoBytes packas slice of bool to slice of bytes
func PackBitsIntoBytes(bits []bool) []byte {
	if len(bits)%8 > 0 {
		bits = append(make([]bool, 8-(len(bits)%8)), bits...)
	}
	bytesMsg := make([]byte, int64(float64(len(bits)/8)))
	start := 0
	for i := 0; i < len(bits); i += 8 {
		var val uint8
		for _, b := range bits[start : i+8] {
			if b {
				val |= 1
			}
		}
		start += 8
		bytesMsg = append(bytesMsg, byte(val))
	}
	return bytesMsg
}

func RescueHashFr(input []bool) *big.Int {
	// packed := ComputeMultiPacking(input)
	// spotgeOutput := RescueHash()
	return nil
}

func ComputeMultiPacking(input []bool) []*big.Int {
	if len(input)%8 > 0 {
		input = append(make([]bool, 8-(len(input)%8)), input...)
	}
	result := []*big.Int{}
	start := 0
	for i := 0; i < len(input); i += 8 {
		cur := big.NewInt(0)
		coeff := big.NewInt(1)

		sample := input[start : i+8]

		for _, b := range sample {
			if b {
				cur.Add(cur, coeff)
			}
			coeff.Mul(coeff, coeff)
		}
		result = append(result, cur)
		start += 8
	}
	return result
}

func (b *Bn256RescueParams) NewCheck2Into1() (*Bn256RescueParams, error) {
	c := uint32(1)
	r := uint32(2)
	round := uint32(22)
	securityLevel := uint32(12)
	return b.NewForParams(c, r, round, securityLevel)
}

func (b *Bn256RescueParams) NewForParams(c, r, rounds, securityLevel uint32) (*Bn256RescueParams, error) {
	stateWidth := c + r
	numRoundConstants := int((1 + rounds*2) * stateWidth)

	roundConstants, err := func() ([]*big.Int, error) {
		tag := []byte("Rescue_f")
		roundConstants := make([]*big.Int, numRoundConstants)
		nonce := uint32(0)
		var nonceBytes []byte

		for {
			binary.BigEndian.PutUint32(nonceBytes[:], nonce)
			hasher, err := blake2s.New256(nil)
			if err != nil {
				return nil, err
			}
			// nonceBytes = []byte()
			hasher.Write(tag)
			hasher.Write(GH_FIRST_BLOCK)
			hasher.Write(nonceBytes)

			hashData := hasher.Sum(nil)
			if len(hashData) != 32 {
				return nil, fmt.Errorf("expecting a hash length of 32 bytes, got a hash with length %d", len(hashData))
			}

			constantRepr := big.NewInt(0).SetBytes(hashData)
			constant := constantRepr
			if constant.Cmp(big.NewInt(0)) != 0 {
				roundConstants = append(roundConstants, constant)
			}

			if len(roundConstants) == numRoundConstants {
				break
			}
			nonce += 1
		}
		return roundConstants, nil
	}()
	if err != nil {
		return nil, err
	}

	mdsMatrix, err := func() ([]*big.Int, error) {
		// This tag is a first one in a sequence of b"ResMxxxx"
		// that produces MDS matrix without eigenvalues for rate = 2,
		// capacity = 1 variant over Bn254 curve
		tag := []byte("ResM0003")
		rng, err := func() (cipher.AEAD, error) {
			hasher, err := blake2s.New256(nil)
			if err != nil {
				return nil, err
			}
			// nonceBytes = []byte()
			hasher.Write(tag)
			hasher.Write(GH_FIRST_BLOCK)
			hashData := hasher.Sum(nil)
			if len(hashData) != 32 {
				return nil, fmt.Errorf("expecting a hash length of 32 bytes, got a hash with length %d", len(hashData))
			}
			var seed [8]uint32
			if len(hashData) < 32 {
				return nil, errors.New("digest is not large enough")
			}
			// Populate the seed array from the byte slice
			for i := 0; i < 8; i++ {
				seed[i] = binary.BigEndian.Uint32(hashData[i*4 : i*4+4])
			}

			// Convert uint32 array to a byte slice for the seed
			var byteSeed []byte
			buf := make([]byte, 4)
			for _, v := range seed {
				binary.BigEndian.PutUint32(buf, v)
				byteSeed = append(byteSeed, buf...)
			}

			// Create a ChaCha20-Poly1305 cipher
			aead, err := chacha20poly1305.NewX(byteSeed)
			if err != nil {
				log.Fatal("failed to create cipher:", err)
			}
			return aead, nil
		}()
		if err != nil {
			return nil, err
		}
		// return rng, nil
		//
		// TODO: ...
		// generatemdsMatrix()
		println(rng)
		return nil, nil
	}()
	return &Bn256RescueParams{
		C:                  c,
		R:                  r,
		Rounds:             rounds,
		SecurityLevel:      securityLevel,
		RoundConstants:     roundConstants,
		MDSMatrix:          mdsMatrix,
		CustomGatesAllowed: false,
	}, nil
}

type PowerSBox struct {
	Power *big.Int
	Inv   uint64
}

type QuinticSBox struct {
	Marker *big.Int
}
