package zklink

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/bits-and-blooms/bitset"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/zklink/bn256/fr"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/zklink/helper"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20"
)

const PAD_MSG_BEFORE_HASH_BITS_LEN = 736

func RescueHashTransactionMsg(message *big.Int) (interface{}, error) {
	// msgBits := BytesIntoBits(message)
	if len(message.Bits()) <= PAD_MSG_BEFORE_HASH_BITS_LEN {
		return nil, errors.New("invalid message")
	}
	// bits := BytesIntoBits(message.Bytes())
	// bn256.Order
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

	roundConstants, err := func() ([]*fr.Element, error) {
		tag := []byte("Rescue_f")
		roundConstants := make([]*fr.Element, numRoundConstants)
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

			constantRepr := helper.NewElement().SetBytes(hashData)
			constant := constantRepr
			if constant.Cmp(helper.NewElement().SetZero()) != 0 {
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

	mdsMatrix, err := func() ([]*fr.Element, error) {
		// This tag is a first one in a sequence of b"ResMxxxx"
		// that produces MDS matrix without eigenvalues for rate = 2,
		// capacity = 1 variant over Bn254 curve
		tag := []byte("ResM0003")
		rng, err := func() (*chacha20.Cipher, error) {
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
			// aead, err := chacha20poly1305.NewX(byteSeed)
			// if err != nil {
			// 	log.Fatal("failed to create cipher:", err)
			// }
			// return aead, nil
			key := make([]byte, 32) // 256 bits for ChaCha20 key
			if _, err := rand.Read(key); err != nil {
				return nil, err
			}

			// Return a new ChaCha20 instance
			rng, err := chacha20.NewUnauthenticatedCipher(key, byteSeed)
			if err != nil {
				return nil, err
			}
			return rng, nil
		}()
		if err != nil {
			return nil, err
		}
		// var mdsMatrix helper.Matrix
		return generateMDSMatrix(stateWidth, rng)
		// return helper.GenMDS(int(stateWidth)), nil
	}()
	if err != nil {
		return nil, err
	}
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

// generateMDSMatrix ...
func generateMDSMatrix(t uint32, rng *chacha20.Cipher) ([]*fr.Element, error) {
	for {
		var x, y []*fr.Element
		for range t {
			x = append(x, genFr(rng))
			y = append(y, genFr(rng))
		}
		var invalid bool
		// quick and dirty check for uniqueness of x
		for i := range t {
			if invalid {
				continue
			}
			el := x[i]
			for _, other := range x[i+1:] {
				if el == other {
					invalid = true
					break
				}
			}
		}
		if invalid {
			continue
		}
		// quick and dirty check for uniqueness of y
		for i := range t {
			if invalid {
				continue
			}
			el := y[i]
			for _, other := range y[i+1:] {
				if el == other {
					invalid = true
					break
				}
			}
		}
		if invalid {
			continue
		}
		// quick and dirty check for uniqueness of x vs y
		for i := range t {
			if invalid {
				continue
			}
			el := x[i]
			for _, other := range y {
				if el == other {
					invalid = true
					break
				}
			}
		}

		if invalid {
			continue
		}

		// by previous checks we can be sure in uniqueness and perform subtractions easily
		mdsMatrix := make([]*fr.Element, t*t)
		for i, xi := range x {
			for j, yi := range y {
				placeInto := uint32(i)*t + uint32(j)
				element := new(fr.Element).Set(xi)
				element = element.Sub(element, yi)
				mdsMatrix[placeInto] = element
			}
		}

		// now we need to do the inverse
		// batch_inversion::<E>(&mut mds_matrix[..]);
		return BatchInvert(mdsMatrix), nil
		// return helper.GenMDS(int(t)), nil
	}
}

// BatchInvert returns a new slice with every element inverted.
// Uses Montgomery batch inversion trick
func BatchInvert(a []*fr.Element) []*fr.Element {
	res := make([]*fr.Element, len(a))
	if len(a) == 0 {
		return res
	}

	zeroes := bitset.New(uint(len(a)))
	accumulator := new(fr.Element).SetUint64(1)

	for i := 0; i < len(a); i++ {
		if a[i].IsZero() {
			zeroes.Set(uint(i))
			continue
		}
		res[i] = accumulator
		accumulator.Mul(accumulator, a[i])
	}

	accumulator.Inverse(accumulator)

	for i := len(a) - 1; i >= 0; i-- {
		if zeroes.Test(uint(i)) {
			continue
		}
		res[i].Mul(res[i], accumulator)
		accumulator.Mul(accumulator, a[i])
	}
	return res
}

// BatchInversion computes the inverses of elements in the slice using Montgomery's trick
func BatchInversion(v []*big.Int, modulus *big.Int) {
	if len(v) == 0 {
		return
	}
	// First pass: compute [a, ab, abc, ...]
	prod := make([]*big.Int, len(v))
	tmp := big.NewInt(1) // Assuming the one method initializes to multiplicative identity
	var zero = big.NewInt(0)
	j := 0

	for _, g := range v {
		if g.Cmp(zero) != 0 {
			tmp = new(big.Int).Mul(tmp, g)
			tmp = tmp.Mod(tmp, modulus)
			prod[j] = new(big.Int).Set(tmp)
			j++
		}
	}

	// Invert `tmp`
	tmp = tmp.ModInverse(tmp, modulus)

	// Second pass: iterate backwards to compute inverses
	for i := j - 1; i >= 0; i-- {
		g := v[i]
		if g.Cmp(zero) != 0 {
			newtmp := new(big.Int).Mul(tmp, g)
			newtmp = newtmp.Mod(newtmp, modulus)
			g.Mul(tmp, prod[i])
			g.Mod(g, modulus)
			tmp.Set(newtmp)
		}
	}
}

// genFr simulates generating a random number value using the provided ChaCha20 RNG.
func genFr(rng *chacha20.Cipher) *fr.Element {
	buf := make([]byte, 8)     // Assuming we need 64 bits for our random values
	rng.XORKeyStream(buf, buf) // Fill buf with random data
	value := binary.BigEndian.Uint64(buf)
	return helper.NewElement().SetBigInt(big.NewInt(0).SetUint64(value))
}

// PowerSBox ...
type PowerSBox struct {
	Power *big.Int
	Inv   uint64
}

// QuinticSBox ...
type QuinticSBox struct {
	Marker *big.Int
}
