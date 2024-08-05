package hash

// Note: currently some methods implementations here are directly copied from the github.com/yaune/starkex repository, and will be removed/update and tested

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
)

// LoadPedersenConfig loads a pedersen configuration from a json file.
func LoadPedersenConfig(path string) (*PedersenCfg, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var resp *PedersenCfg
	return resp, json.Unmarshal([]byte(file), &resp)
}

func (cfg *PedersenCfg) PedersenHash(str ...string) string {
	NElementBitsHash := cfg.FieldPrime.BitLen()
	point := cfg.ConstantPoints[0]
	for i, s := range str {
		x, _ := big.NewInt(0).SetString(s, 10)
		pointList := cfg.ConstantPoints[2+i*NElementBitsHash : 2+(i+1)*NElementBitsHash]
		n := big.NewInt(0)
		for _, pt := range pointList {
			n.And(x, big.NewInt(1))
			if n.Cmp(big.NewInt(0)) > 0 {
				point = eccAdd(point, pt, cfg.FieldPrime)
			}
			x = x.Rsh(x, 1)
		}
	}
	return point[0].String()
}

// eccAdd Gets two points on an elliptic curve mod p and returns their sum.
// Assumes the points are given in affine form (x, y) and have different x coordinates.
func eccAdd(point1 [2]*big.Int, point2 [2]*big.Int, p *big.Int) [2]*big.Int {
	// m = div_mod(point1[1] - point2[1], point1[0] - point2[0], p)
	d1 := big.NewInt(0).Sub(point1[1], point2[1])
	d2 := big.NewInt(0).Sub(point1[0], point2[0])
	m := divMod(d1, d2, p)

	// x = (m * m - point1[0] - point2[0]) % p
	x := big.NewInt(0)
	x.Sub(big.NewInt(0).Mul(m, m), point1[0])
	x.Sub(x, point2[0])
	x.Mod(x, p)

	// y := (m*(point1[0]-x) - point1[1]) % p
	y := big.NewInt(0)
	y.Mul(m, big.NewInt(0).Sub(point1[0], x))
	y.Sub(y, point1[1])
	y.Mod(y, p)

	return [2]*big.Int{x, y}
}

// divMod Finds a nonnegative integer 0 <= x < p such that (m * x) % p == n
func divMod(n, m, p *big.Int) *big.Int {
	a, _, _ := igcdex(m, p)
	// (n * a) % p
	tmp := big.NewInt(0).Mul(n, a)
	return tmp.Mod(tmp, p)
}

/*
igcdex
Returns x, y, g such that g = x*a + y*b = gcd(a, b).

	>>> from sympy.core.numbers import igcdex
	>>> igcdex(2, 3)
	(-1, 1, 1)
	>>> igcdex(10, 12)
	(-1, 1, 2)

	>>> x, y, g = igcdex(100, 2004)
	>>> x, y, g
	(-20, 1, 4)
	>>> x*100 + y*2004
	4
*/
/**
from sympy.core.numbers import igcdex
source code:
   if (not a) and (not b):
	   return (0, 1, 0)

   if not a:
	   return (0, b//abs(b), abs(b))
   if not b:
	   return (a//abs(a), 0, abs(a))

   if a < 0:
	   a, x_sign = -a, -1
   else:
	   x_sign = 1

   if b < 0:
	   b, y_sign = -b, -1
   else:
	   y_sign = 1

   x, y, r, s = 1, 0, 0, 1

   while b:
	   (c, q) = (a % b, a // b)
	   (a, b, r, s, x, y) = (b, c, x - q*r, y - q*s, r, s)

   return (x*x_sign, y*y_sign, a)
*/
func igcdex(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	var zero = big.NewInt(0)
	if a.Cmp(zero) == 0 && b.Cmp(zero) == 0 {
		return big.NewInt(0), big.NewInt(1), big.NewInt(0)
	}
	if a.Cmp(zero) == 0 {
		return big.NewInt(0), big.NewInt(0).Quo(b, big.NewInt(0).Abs(b)), big.NewInt(0).Abs(b)
	}
	if b.Cmp(zero) == 0 {
		return big.NewInt(0).Quo(a, big.NewInt(0).Abs(a)), big.NewInt(0), big.NewInt(0).Abs(a)
	}
	xSign := big.NewInt(1)
	ySign := big.NewInt(1)
	if a.Cmp(zero) == -1 {
		a, xSign = a.Neg(a), big.NewInt(-1)
	}
	if b.Cmp(zero) == -1 {
		b, ySign = b.Neg(b), big.NewInt(-1)
	}
	x, y, r, s := big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)
	for b.Cmp(zero) > 0 {
		c, q := big.NewInt(0).Mod(a, b), big.NewInt(0).Quo(a, b)
		a, b, r, s, x, y = b, c, big.NewInt(0).Sub(x, big.NewInt(0).Mul(q, r)), big.NewInt(0).Sub(y, big.NewInt(0).Mul(big.NewInt(0).Neg(q), s)), r, s
	}
	return x.Mul(x, xSign), y.Mul(y, ySign), a
}
