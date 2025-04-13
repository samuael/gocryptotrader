package helper

import (
	"reflect"

	"github.com/thrasher-corp/gocryptotrader/internal/utils/zklink/bn256/fr"
)

func NewElement() *fr.Element {
	typ := reflect.TypeOf((*fr.Element)(nil)).Elem()
	val := reflect.New(typ.Elem())
	return val.Interface().(*fr.Element)
}

func zero() *fr.Element {
	return NewElement().SetZero()
}

func one() *fr.Element {
	return NewElement().SetOne()
}

// func Modulus() *big.Int {
// 	e := NewElement().SetZero()
// 	e.Sub(e, NewElement().SetOne())
// 	b := e.BigInt(new(big.Int))
// 	return b.Add(b, big.NewInt(1))
// }

// func Bits() int {
// 	return Modulus().BitLen()
// }

// func Bytes() int {
// 	return (Bits() + 7) / 8
// }

// Exp is a copy of gnark-crypto's implementation, but takes a pointer argument
// func Exp(z, x *fr.Element, k *big.Int) {
// 	if k.IsUint64() && k.Uint64() == 0 {
// 		z.SetOne()
// 	}

// 	e := k
// 	if k.Sign() == -1 {
// 		// negative k, we invert
// 		// if k < 0: xᵏ (mod q) == (x⁻¹)ᵏ (mod q)
// 		x.Inverse(x)

// 		// we negate k in a temp big.Int since
// 		// Int.Bit(_) of k and -k is different
// 		e = pool.BigInt.Get()
// 		defer pool.BigInt.Put(e)
// 		e.Neg(k)
// 	}

// 	z.Set(x)

// 	for i := e.BitLen() - 2; i >= 0; i-- {
// 		z.Square(z)
// 		if e.Bit(i) == 1 {
// 			z.Mul(z, x)
// 		}
// 	}
// }
