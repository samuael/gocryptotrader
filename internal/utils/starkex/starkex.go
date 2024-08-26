package starkex

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	path "github.com/thrasher-corp/gocryptotrader/internal/testing/utils"
	"github.com/thrasher-corp/gocryptotrader/internal/utils/hash"
	math_utils "github.com/thrasher-corp/gocryptotrader/internal/utils/math"
)

// Error declarations.
var (
	ErrInvalidPrivateKey  = errors.New("invalid private key")
	ErrInvalidHashPayload = errors.New("invalid hash payload")
)

// StarkConfig represents a stark configuration
type StarkConfig struct {
	*elliptic.CurveParams
	EcGenX           *big.Int
	EcGenY           *big.Int
	MinusShiftPointX *big.Int
	MinusShiftPointY *big.Int
	Max              *big.Int
	Alpha            *big.Int
	ConstantPoints   [][2]*big.Int
	PedersenHash     func(...string) string
}

var one = big.NewInt(1)

const defaultPedersenConfigsPath = "internal/utils/hash/pedersen_config/"

// NewStarkExConfig returns a stark configuration given the exchange name
func NewStarkExConfig(exchangeName string) (*StarkConfig, error) {
	rootPath, err := path.RootPathFromCWD()
	if err != nil {
		return nil, err
	}
	pedersenConfig, err := hash.LoadPedersenConfig(rootPath + "/" + defaultPedersenConfigsPath + strings.ToLower(exchangeName) + ".json")
	if err != nil {
		return nil, err
	}
	starkCurve := &StarkConfig{
		CurveParams: &elliptic.CurveParams{
			P:       pedersenConfig.FieldPrime,
			N:       pedersenConfig.EcOrder,
			B:       pedersenConfig.BETA,
			Gx:      pedersenConfig.ConstantPoints[0][0],
			Gy:      pedersenConfig.ConstantPoints[0][1],
			BitSize: 252,
		},
		EcGenX:         pedersenConfig.ConstantPoints[1][0],
		EcGenY:         pedersenConfig.ConstantPoints[1][1],
		Alpha:          big.NewInt(int64(pedersenConfig.ALPHA)),
		ConstantPoints: pedersenConfig.ConstantPoints,
		PedersenHash:   pedersenConfig.PedersenHash,
	}
	starkCurve.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	starkCurve.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10)
	starkCurve.Max, _ = new(big.Int).SetString("3618502788666131106986593281521497120414687020801267626233049500247285301248", 10) // 2 ** 251
	return starkCurve, nil
}

// Sign generates a signature out using the users private key and signable order params.
func (sfg *StarkConfig) Sign(sgn Signable, starkPrivateKey string) (string, error) {
	pHash, err := sgn.GetPedersenHash(sfg.PedersenHash)
	if err != nil {
		return pHash, err
	}
	priKey, okay := big.NewInt(0).SetString(starkPrivateKey, 0)
	if !okay {
		return "", fmt.Errorf("%w, %v", ErrInvalidPrivateKey, starkPrivateKey)
	}
	msgHash, okay := new(big.Int).SetString(pHash, 10)
	if !okay {
		return "", ErrInvalidHashPayload
	}
	r, s, err := sfg.SignECDSA(msgHash, priKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%x", append(r.Bytes(), s.Bytes()...)), nil
}

// Add computes the sum of two points on the StarkCurve.
func (sc StarkConfig) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	yDelta := new(big.Int).Sub(y1, y2)
	xDelta := new(big.Int).Sub(x1, x2)

	m := math_utils.DivMod(yDelta, xDelta, sc.P)

	xm := new(big.Int).Mul(m, m)

	x = new(big.Int).Sub(xm, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, sc.P)

	y = new(big.Int).Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, sc.P)

	return x, y
}

// Double calculates the double of a point on a StarkCurve (equation y^2 = x^3 + alpha*x + beta mod p).
func (sc StarkConfig) Double(x1, y1 *big.Int) (x, y *big.Int) {
	xin := new(big.Int).Mul(big.NewInt(3), x1)
	xin = xin.Mul(xin, x1)
	xin = xin.Add(xin, sc.Alpha)

	yin := new(big.Int).Mul(y1, big.NewInt(2))

	m := math_utils.DivMod(xin, yin, sc.P)

	xout := new(big.Int).Mul(m, m)
	xmed := new(big.Int).Mul(big.NewInt(2), x1)
	xout = xout.Sub(xout, xmed)
	xout = xout.Mod(xout, sc.P)

	yout := new(big.Int).Sub(x1, xout)
	yout = yout.Mul(m, yout)
	yout = yout.Sub(yout, y1)
	yout = yout.Mod(yout, sc.P)

	return xout, yout
}

// ScalarMult performs scalar multiplication on a point (x1, y1) with a scalar value k.
func (sc StarkConfig) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	m := new(big.Int).SetBytes(k)
	x, y = sc.EcMult(m, x1, y1)
	return x, y
}

// ScalarBaseMult returns the result of multiplying the base point of the StarkCurve
// by the given scalar value.
func (sc StarkConfig) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return sc.ScalarMult(sc.Gx, sc.Gy, k)
}

// IsOnCurve checks if the given point (x, y) lies on the curve defined by the StarkCurve instance.
func (sc StarkConfig) IsOnCurve(x, y *big.Int) bool {
	left := new(big.Int).Mul(y, y)
	left = left.Mod(left, sc.P)

	right := new(big.Int).Mul(x, x)
	right = right.Mul(right, x)
	right = right.Mod(right, sc.P)

	ri := new(big.Int).Mul(big.NewInt(1), x)

	right = right.Add(right, ri)
	right = right.Add(right, sc.B)
	right = right.Mod(right, sc.P)

	if left.Cmp(right) == 0 {
		return true
	} else {
		return false
	}
}

// InvModCurveSize calculates the inverse modulus of a given big integer 'x' with respect to the StarkCurve 'sc'.
func (sc StarkConfig) InvModCurveSize(x *big.Int) *big.Int {
	return math_utils.DivMod(big.NewInt(1), x, sc.N)
}

// GetYCoordinate calculates the y-coordinate of a point on the StarkCurve.
func (sc StarkConfig) GetYCoordinate(starkX *big.Int) *big.Int {
	y := new(big.Int).Mul(starkX, starkX)
	y = y.Mul(y, starkX)
	yin := new(big.Int).Mul(sc.Alpha, starkX)

	y = y.Add(y, yin)
	y = y.Add(y, sc.B)
	y = y.Mod(y, sc.P)

	y = y.ModSqrt(y, sc.P)
	return y
}

// MimicEcMultAir performs a computation on the StarkCurve struct (m * point + shift_point)
func (sc StarkConfig) MimicEcMultAir(mout, x1, y1, x2, y2 *big.Int) (x *big.Int, y *big.Int, err error) {
	m := new(big.Int).Set(mout)
	if m.Cmp(big.NewInt(0)) != 1 || m.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("too many bits %v", m.BitLen())
	}

	psx := x2
	psy := y2
	for i := 0; i < 251; i++ {
		if psx == x1 {
			return x, y, fmt.Errorf("xs are the same")
		}
		if m.Bit(0) == 1 {
			psx, psy = sc.Add(psx, psy, x1, y1)
		}
		x1, y1 = sc.Double(x1, y1)
		m = m.Rsh(m, 1)
	}
	if m.Cmp(big.NewInt(0)) != 0 {
		return psx, psy, fmt.Errorf("m doesn't equal zero")
	}
	return psx, psy, nil
}

// EcMult multiplies a point (equation y^2 = x^3 + alpha*x + beta mod p) on the StarkCurve by a scalar value.
func (sc StarkConfig) EcMult(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMult func(m, x1, y1 *big.Int) (x, y *big.Int)

	_add := func(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
		yDelta := new(big.Int).Sub(y1, y2)
		xDelta := new(big.Int).Sub(x1, x2)

		m := math_utils.DivMod(yDelta, xDelta, sc.P)

		xm := new(big.Int).Mul(m, m)

		x = new(big.Int).Sub(xm, x1)
		x = x.Sub(x, x2)
		x = x.Mod(x, sc.P)

		y = new(big.Int).Sub(x1, x)
		y = y.Mul(m, y)
		y = y.Sub(y, y1)
		y = y.Mod(y, sc.P)

		return x, y
	}

	// alpha is our Y
	_ecMult = func(m, x1, y1 *big.Int) (x, y *big.Int) {
		if m.BitLen() == 1 {
			return x1, y1
		}
		mk := new(big.Int).Mod(m, big.NewInt(2))
		if mk.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int).Div(m, big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _ecMult(h, c, d)
		}
		n := new(big.Int).Sub(m, big.NewInt(1))
		e, f := _ecMult(n, x1, y1)
		return _add(e, f, x1, y1)
	}

	x, y = _ecMult(m, x1, y1)
	return x, y
}

// Verify verifies the validity of the signature for a given message hash using the StarkCurve.
func (sc StarkConfig) Verify(msgHash, r, s, pubX, pubY *big.Int) bool {
	w := sc.InvModCurveSize(s)

	if s.Cmp(big.NewInt(0)) != 1 || s.Cmp(sc.N) != -1 {
		return false
	}
	if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
		return false
	}
	if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
		return false
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
		return false
	}
	if !sc.IsOnCurve(pubX, pubY) {
		return false
	}

	zGx, zGy, err := sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if err != nil {
		return false
	}

	rQx, rQy, err := sc.MimicEcMultAir(r, pubX, pubY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}
	inX, inY := sc.Add(zGx, zGy, rQx, rQy)
	wBx, wBy, err := sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}

	outX, _ := sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if r.Cmp(outX) == 0 {
		return true
	} else {
		altY := new(big.Int).Neg(pubY)

		zGx, zGy, err = sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(r, pubX, new(big.Int).Set(altY), sc.Gx, sc.Gy)
		if err != nil {
			return false
		}
		inX, inY = sc.Add(zGx, zGy, rQx, rQy)
		wBx, wBy, err = sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}

		outX, _ = sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if r.Cmp(outX) == 0 {
			return true
		}
	}
	return false
}

// Sign calculates the signature of a message using the StarkCurve algorithm.
func (sc StarkConfig) SignECDSA(msgHash, privKey *big.Int, seed ...*big.Int) (x, y *big.Int, err error) {
	if msgHash == nil {
		return x, y, fmt.Errorf("nil msgHash")
	}
	if privKey == nil {
		return x, y, fmt.Errorf("nil privKey")
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("invalid bit length")
	}

	inSeed := big.NewInt(0)
	if len(seed) == 1 && inSeed != nil {
		inSeed = seed[0]
	}
	for {
		k := sc.GenerateSecret(big.NewInt(0).Set(msgHash), big.NewInt(0).Set(privKey), big.NewInt(0).Set(inSeed))
		// In case r is rejected k shall be generated with new seed
		inSeed = inSeed.Add(inSeed, big.NewInt(1))

		r, _ := sc.EcMult(k, sc.EcGenX, sc.EcGenY)

		// DIFF: in classic ECDSA, we take int(x) % n.
		if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		agg := new(big.Int).Mul(r, privKey)
		agg = agg.Add(agg, msgHash)

		if new(big.Int).Mod(agg, sc.N).Cmp(big.NewInt(0)) == 0 {
			// Bad value. This fails with negligible probability.
			continue
		}

		w := math_utils.DivMod(k, agg, sc.N)
		if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		s := sc.InvModCurveSize(w)
		return r, s, nil
	}
}

// GenerateSecret generates a secret using the StarkCurve struct.
func (sc StarkConfig) GenerateSecret(msgHash, privKey, seed *big.Int) (secret *big.Int) {
	alg := sha256.New
	holen := alg().Size()
	rolen := (sc.BitSize + 7) >> 3

	if msgHash.BitLen()%8 <= 4 && msgHash.BitLen() >= 248 {
		msgHash = msgHash.Mul(msgHash, big.NewInt(16))
	}

	by := append(math_utils.Int2Octets(privKey, rolen), math_utils.Bits2Octets(msgHash.Bytes(), sc.N, sc.BitSize, rolen)...)

	if seed.Cmp(big.NewInt(0)) == 1 {
		by = append(by, seed.Bytes()...)
	}

	v := bytes.Repeat([]byte{0x01}, holen)

	k := bytes.Repeat([]byte{0x00}, holen)

	k = math_utils.Mac(alg, k, append(append(v, 0x00), by...), k)

	v = math_utils.Mac(alg, k, v, v)

	k = math_utils.Mac(alg, k, append(append(v, 0x01), by...), k)

	v = math_utils.Mac(alg, k, v, v)

	for {
		var t []byte

		for len(t) < rolen {
			v = math_utils.Mac(alg, k, v, v)
			t = append(t, v...)
		}

		secret = math_utils.Bits2Int(t, sc.BitSize)
		if secret.Cmp(big.NewInt(0)) == 1 && secret.Cmp(sc.N) == -1 {
			return secret
		}
		k = math_utils.Mac(alg, k, append(v, 0x00), k)
		v = math_utils.Mac(alg, k, v, v)
	}
}
