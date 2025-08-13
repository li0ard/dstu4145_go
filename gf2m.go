package dstu4145go

import "math/big"

type Field struct {
	M       int
	Poly    *big.Int
	Modulus *big.Int
}

func NewField(m int, ks []int) *Field {
	f := new(Field)
	f.M = m
	f.Modulus = new(big.Int).Lsh(big.NewInt(1), uint(m))

	lower := big.NewInt(1)
	for _, exp := range ks {
		temp := new(big.Int).Lsh(big.NewInt(1), uint(exp))
		lower.Or(lower, temp)
	}

	f.Poly = new(big.Int).Or(f.Modulus, lower)

	return f
}

func (f *Field) Add(a *big.Int, b *big.Int) *big.Int {
	return new(big.Int).Xor(a, b)
}

func (f *Field) Mul(a *big.Int, b *big.Int) *big.Int {
	result := big.NewInt(0)
	a = new(big.Int).Set(a)
	b = new(big.Int).Set(b)

	zero := big.NewInt(0)
	one := big.NewInt(1)

	for b.Cmp(zero) > 0 {
		if new(big.Int).And(b, one).Cmp(one) == 0 {
			result.Xor(result, a)
		}

		a.Lsh(a, 1)

		if new(big.Int).And(a, f.Modulus).Cmp(zero) != 0 {
			a.Xor(a, f.Poly)
		}

		b.Rsh(b, 1)
	}

	return result
}

func (f *Field) Sqr(a *big.Int) *big.Int {
	return f.Mul(a, a)
}

func (f *Field) Pow(a, exponent *big.Int) *big.Int {
	if a.Sign() == 0 {
		return big.NewInt(0)
	}

	result := big.NewInt(1)
	base := new(big.Int).Set(a)
	exp := new(big.Int).Set(exponent)
	zero := big.NewInt(0)
	one := big.NewInt(1)

	for exp.Cmp(zero) > 0 {
		if new(big.Int).And(exp, one).Cmp(one) == 0 {
			result = f.Mul(result, base)
		}

		base = f.Sqr(base)
		exp.Rsh(exp, 1)
	}

	return result
}

func (f *Field) Inv(a *big.Int) *big.Int {
	if a.Sign() == 0 {
		panic("Cannot invert zero")
	}

	exponent := new(big.Int).Lsh(big.NewInt(1), uint(f.M))
	exponent.Sub(exponent, big.NewInt(2))
	return f.Pow(a, exponent)
}

func (f *Field) Div(a, b *big.Int) *big.Int {
	return f.Mul(a, f.Inv(b))
}

func (f *Field) Sqrt(a *big.Int) *big.Int {
	if a.Sign() == 0 {
		return big.NewInt(0)
	}

	exponent := new(big.Int).Lsh(big.NewInt(1), uint(f.M-1))
	return f.Pow(a, exponent)
}

func (f *Field) Trace(a *big.Int) *big.Int {
	if a.Sign() == 0 {
		return big.NewInt(0)
	}

	result := big.NewInt(0)
	current := new(big.Int).Set(a)
	one := big.NewInt(1)

	for i := 0; i < f.M; i++ {
		result.Xor(result, current)
		current = f.Sqr(current)
	}

	return new(big.Int).And(result, one)
}

func (f *Field) HTrace(a *big.Int) *big.Int {
	if a.Sign() == 0 {
		return big.NewInt(0)
	}

	result := big.NewInt(0)
	current := new(big.Int).Set(a)
	t := (f.M - 1) / 2

	for i := 0; i <= t; i++ {
		result.Xor(result, current)
		current = f.Sqr(f.Sqr(current))
	}

	return result
}

func (f *Field) HashToField(hash []byte) *big.Int {
	num := new(big.Int).SetBytes(hash)

	mask := new(big.Int).Sub(f.Modulus, big.NewInt(1))

	result := new(big.Int)
	result.And(num, mask)

	return result
}

func (f *Field) Reduce(a *big.Int) *big.Int {
	if a.Cmp(f.Modulus) < 0 {
		return new(big.Int).Set(a)
	}

	result := new(big.Int).Set(a)

	for {
		deg := result.BitLen() - 1
		if deg < f.M {
			break
		}

		shift := deg - f.M

		polyShifted := new(big.Int).Lsh(f.Poly, uint(shift))

		result.Xor(result, polyShifted)
	}

	return result
}

func (f *Field) ToBytes(a *big.Int) []byte {
	byteLen := (f.M + 7) / 8
	result := make([]byte, byteLen)
	a.FillBytes(result)

	return result
}
