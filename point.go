package dstu4145go

import (
	"fmt"
	"math/big"
)

type Point struct {
	Curve *Curve
	X     *big.Int
	Y     *big.Int
}

func ZERO() *Point {
	return &Point{
		Curve: nil,
		X:     big.NewInt(0),
		Y:     big.NewInt(0),
	}
}

func NewPoint(curve *Curve, x *big.Int, y *big.Int) *Point {
	return &Point{
		Curve: curve,
		X:     new(big.Int).Set(x),
		Y:     new(big.Int).Set(y),
	}
}

func (p *Point) IsInfinity() bool {
	return (p.X.Cmp(ZERO().X) == 0) && (p.Y.Cmp(ZERO().Y) == 0)
}

func (p *Point) Negate() *Point {
	if p.IsInfinity() {
		return p
	}
	return NewPoint(p.Curve, p.X, p.Curve.Field.Add(p.X, p.Y))
}

func (p *Point) Double() *Point {
	if p.IsInfinity() {
		return p
	}
	if p.X.Sign() == 0 {
		if p.Curve.Field.Sqr(p.Y).Cmp(p.Curve.Params.B) == 0 {
			return ZERO()
		}
		panic("Invalid point: x=0 but not a point of order 2")
	}

	// λ = x + y/x
	lambda := p.Curve.Field.Add(p.X, p.Curve.Field.Div(p.Y, p.X))

	// x3 = λ² + λ + a
	x3 := p.Curve.Field.Add(p.Curve.Field.Add(p.Curve.Field.Sqr(lambda), lambda), p.Curve.Params.A)

	// y3 = x² + (λ + 1) * x3
	y3 := p.Curve.Field.Add(p.Curve.Field.Sqr(p.X), p.Curve.Field.Mul(p.Curve.Field.Add(lambda, big.NewInt(1)), x3))

	return NewPoint(p.Curve, x3, y3)
}

func (p *Point) Add(other *Point) *Point {
	if p == nil || other == nil {
		panic("Cannot add nil points")
	}
	if p.IsInfinity() {
		return other
	}
	if other.IsInfinity() {
		return p
	}

	if p.X.Cmp(other.X) == 0 {
		if p.Y.Cmp(other.Y) == 0 {
			return p.Double()
		}
		if other.Y.Cmp(p.Curve.Field.Add(p.X, p.Y)) == 0 {
			return ZERO()
		}
		panic("Invalid points: same x but not inverses")
	}

	// λ = (y2 + y1) / (x2 + x1)
	xDiff := p.Curve.Field.Add(other.X, p.X)
	yDiff := p.Curve.Field.Add(other.Y, p.Y)
	lambda := p.Curve.Field.Div(yDiff, xDiff)

	// x3 = λ² + λ + x1 + x2 + a
	x3 := p.Curve.Field.Add(
		p.Curve.Field.Add(p.Curve.Field.Add(p.Curve.Field.Sqr(lambda), lambda), p.Curve.Field.Add(p.X, other.X)),
		p.Curve.Params.A,
	)

	// y3 = λ(x1 + x3) + x3 + y1
	y3 := p.Curve.Field.Add(
		p.Curve.Field.Add(p.Curve.Field.Mul(lambda, p.Curve.Field.Add(p.X, x3)), x3),
		p.Y,
	)

	return NewPoint(p.Curve, x3, y3)
}

func (p *Point) MultiplySlow(scalar *big.Int) *Point {
	if p == nil {
		panic("point is nil")
	}
	if p.Curve == nil {
		panic("point is not associated with a curve")
	}
	if scalar.Sign() < 0 {
		panic("scalar must be non-negative")
	}
	if p.IsInfinity() {
		return ZERO()
	}
	exp := new(big.Int).Set(scalar)
	result := ZERO()
	base := p
	zero := big.NewInt(0)
	one := big.NewInt(1)
	for exp.Cmp(zero) > 0 {
		if new(big.Int).And(exp, one).Cmp(one) == 0 {
			result = result.Add(base)
		}
		base = base.Double()
		exp.Rsh(exp, 1)
	}

	return result
}

func (p *Point) Multiply(scalar *big.Int) *Point {
	num := new(big.Int).Set(scalar)
	acc_x, acc_y := new(big.Int), new(big.Int)
	doubler_x, doubler_y := new(big.Int).Set(p.X), new(big.Int).Set(p.Y)
	tmp_x, tmp_y := new(big.Int), new(big.Int)
	for num.Sign() > 0 {
		if num.Bit(0) != 0 {
			tmp_x.Set(acc_x)
			tmp_y.Set(acc_y)
			add(acc_x, acc_y, tmp_x, tmp_y, doubler_x, doubler_y, p)
		}
		num.Rsh(num, 1)
		tmp_x.Set(doubler_x)
		tmp_y.Set(doubler_y)
		double(doubler_x, doubler_y, tmp_x, tmp_y, p)
	}

	return NewPoint(p.Curve, acc_x, acc_y)
}

func (p *Point) String() string {
	if p.IsInfinity() {
		return "O"
	}
	return fmt.Sprintf("Point(%x, %x)", p.X, p.Y)
}

func (p *Point) Compress() []byte {
	if p.IsInfinity() {
		return []byte{0x00}
	}
	parity := byte(0)
	if p.Y.Bit(0) == 1 {
		parity = 1
	}
	prefix := 0x02 + parity
	byteLen := (p.Curve.Field.M+1)/8 + 1
	xBytes := make([]byte, byteLen)
	p.X.FillBytes(xBytes)
	return append([]byte{prefix}, xBytes...)
}
