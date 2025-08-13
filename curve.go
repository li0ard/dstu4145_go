package dstu4145go

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Curve struct {
	Params *DstuParameters
	Field  *Field
	Base   *Point
}

type SignOptions struct {
	Rand *big.Int
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func NewCurve(params *DstuParameters) *Curve {
	d := &Curve{
		Params: params,
		Field:  NewField(params.M, params.KS),
	}
	d.Base = d.Point(d.Params.Gx, d.Params.Gy)
	return d
}

func (d *Curve) Point(x *big.Int, y *big.Int) *Point {
	return NewPoint(d, x, y)
}

// Generate public key from private
func (d *Curve) GetPublicKey(privateKey []byte) *Point {
	// Q = -dP
	return d.Base.Multiply(new(big.Int).SetBytes(privateKey)).Negate()
}

func (d *Curve) ComputePreSign() *big.Int {
	byteLen := d.Params.M / 8
	randomBytes := make([]byte, byteLen)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}
	randNum := new(big.Int).SetBytes(randomBytes)
	return new(big.Int).Mod(randNum, d.Params.N)
}

// Compute signature
func (c *Curve) Sign(privateKey []byte, hash []byte, opts *SignOptions) *Signature {
	var e *big.Int
	if opts != nil && opts.Rand != nil {
		e = opts.Rand
	} else {
		e = c.ComputePreSign()
	}

	d := new(big.Int).SetBytes(privateKey)
	// h = hash_to_field(H(T))
	h := c.Field.HashToField(hash)
	// Fe = eP.x
	Fe := c.Base.Multiply(e).X
	// r = h * Fe
	r := c.Field.Mul(h, Fe)
	// s = e + dr (mod n)
	s := new(big.Int)
	s.Add(e, new(big.Int).Mul(d, r)).Mod(s, c.Params.N)

	return &Signature{r, s}
}

// Verify signature
func (c *Curve) Verify(publicKey *Point, hash []byte, signature *Signature) bool {
	Q := publicKey
	// h = hash_to_field(H(T))
	h := c.Field.HashToField(hash)
	// R = sP +rQ
	R := c.Base.Multiply(signature.S).Add(Q.Multiply(signature.R))
	// y = h * R.x
	r := c.Field.Mul(h, R.X)

	return r.Cmp(signature.R) == 0
}

func (curve *Curve) Decompress(compressed []byte) (*Point, error) {
	if len(compressed) == 0 {
		return nil, fmt.Errorf("compressed data is empty")
	}
	if compressed[0] == 0x00 {
		return ZERO(), nil
	}
	if compressed[0] != 0x02 && compressed[0] != 0x03 {
		return nil, fmt.Errorf("invalid compression prefix: 0x%02x", compressed[0])
	}
	parity := compressed[0] & 0x01
	xBytes := compressed[1:]
	expectedLen := (curve.Field.M + 7) / 8
	if len(xBytes) != expectedLen {
		return nil, fmt.Errorf("invalid x length: expected %d bytes, got %d", expectedLen, len(xBytes))
	}
	x := curve.Field.HashToField(xBytes)
	x2 := curve.Field.Sqr(x)
	x3 := curve.Field.Mul(x2, x)
	aX2 := curve.Field.Mul(curve.Params.A, x2)
	c := curve.Field.Add(curve.Field.Add(x3, aX2), curve.Params.B)
	var y *big.Int
	if x.Sign() == 0 {
		y = curve.Field.Sqrt(c)
		if y.Sign() == 0 {
			return nil, fmt.Errorf("no valid square root for c when x=0")
		}
		if y.Bit(0) != uint(parity) {
			return nil, fmt.Errorf("parity mismatch for x=0: expected %d, got %d", parity, y.Bit(0))
		}
	} else {
		v := curve.Field.Div(c, x2)
		if curve.Field.Trace(v).Sign() != 0 {
			return nil, fmt.Errorf("no solution exists: trace(v) = %d", curve.Field.Trace(v))
		}
		z := curve.Field.HTrace(v)
		y = curve.Field.Mul(z, x)
		yAlt := curve.Field.Add(y, x)
		if y.Bit(0) == uint(parity) {
		} else if yAlt.Bit(0) == uint(parity) {
			y = yAlt
		} else {
			return nil, fmt.Errorf("parity mismatch: neither solution matches (y=%d, yAlt=%d)", y.Bit(0), yAlt.Bit(0))
		}
	}
	return NewPoint(curve, x, y), nil
}
