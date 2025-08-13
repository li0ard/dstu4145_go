package dstu4145go

import "math/big"

var one = big.NewInt(1)

type bfi struct {
	v *big.Int
}

func newBFI() *bfi {
	return &bfi{new(big.Int)}
}
func newBFI64(x int64) *bfi {
	return &bfi{big.NewInt(x)}
}
func copyBFI(v *big.Int) *bfi {
	return &bfi{new(big.Int).Set(v)}
}
func wrapBFI(v *big.Int) *bfi {
	return &bfi{v}
}

func (z *bfi) BitLen() int {
	return z.v.BitLen()
}

func (z *bfi) Set(o *bfi) *bfi {
	z.v.Set(o.v)
	return z
}

func (z *bfi) SetBigInt(o *big.Int) *bfi {
	z.v.Set(o)
	return z
}

func (z *bfi) Clone() *bfi {
	return copyBFI(z.v)
}
func (z *bfi) CloneBigInt() *big.Int {
	return new(big.Int).Set(z.v)
}

func (z *bfi) DivMod(num, den, p *bfi) *bfi {
	inv, _ := _extended_gcd(den, p)
	z.Mul(inv, num)
	return z
}

func (z *bfi) Add(x, y *bfi) *bfi {
	z.v.Xor(x.v, y.v)
	return z
}

func (z *bfi) Mul(self, y *bfi) *bfi {
	acc := new(big.Int)
	shift := uint(0)
	o := big.NewInt(0).Set(y.v)
	tmp := big.NewInt(0)
	for o.Sign() > 0 {
		if o.Bit(0) != 0 {
			acc.Xor(acc, tmp.Lsh(self.v, shift))
		}
		shift++
		o.Rsh(o, 1)
	}
	z.SetBigInt(acc)
	return z
}

func (z *bfi) Mod(self, base *bfi) *bfi {
	_, r := _bf_div(self, base)
	z.SetBigInt(r.v)
	return z
}

func (z *bfi) Div(self, o *bfi) *bfi {
	q, _ := _bf_div(self, o)
	z.Set(q)
	return z
}

func _extended_gcd(a_, b_ *bfi) (*bfi, *bfi) {
	a := a_.Clone()
	b := b_.Clone()
	x := newBFI64(0)
	last_x := newBFI64(1)
	y := newBFI64(1)
	last_y := newBFI64(0)
	tmp := newBFI()
	quot := newBFI()
	for b.v.Sign() > 0 {
		quot.Div(a, b)
		tmp.Mod(a, b)
		a.Set(b)
		b.Set(tmp)
		tmp.Add(last_x, tmp.Mul(x, quot))
		last_x.Set(x)
		x.Set(tmp)
		tmp.Add(last_y, tmp.Mul(y, quot))
		last_y.Set(y)
		y.Set(tmp)
	}
	return last_x, last_y
}

func _bf_div(a, b *bfi) (*bfi, *bfi) {
	r := a.CloneBigInt()
	q := new(big.Int)
	rlen := a.BitLen()
	blen := b.BitLen()
	sweeper := new(big.Int).Lsh(one, uint(rlen-1))
	tmp := new(big.Int)
	for rlen >= blen {
		shift := uint(rlen - blen)
		q.Or(q, tmp.Lsh(one, shift))
		r.Xor(r, tmp.Lsh(b.v, shift))
		if r.Sign() == 0 {
			break
		}
		for r.Sign() != 0 && tmp.And(sweeper, r).Sign() == 0 {
			sweeper.Rsh(sweeper, 1)
			rlen--
		}
	}
	return wrapBFI(q), wrapBFI(r)
}

func add(x_, y_, px_, py_, qx_, qy_ *big.Int, c *Point) (*big.Int, *big.Int) {
	if px_.Sign() == 0 && py_.Sign() == 0 {
		x_.Set(qx_)
		y_.Set(qy_)
		return x_, y_
	}
	if qx_.Sign() == 0 && qy_.Sign() == 0 {
		x_.Set(px_)
		y_.Set(py_)
		return x_, y_
	}
	if px_.Cmp(qx_) == 0 && py_.Cmp(qy_) == 0 {
		return double(x_, y_, px_, py_, c)
	}
	x := wrapBFI(x_)
	y := wrapBFI(y_)
	px, py := wrapBFI(px_), wrapBFI(py_)
	qx, qy := wrapBFI(qx_), wrapBFI(qy_)
	s := newBFI()
	f := wrapBFI(c.Curve.Field.Poly)
	a := wrapBFI(c.Curve.Params.A)
	tmp := newBFI()
	s.Add(py, qy)
	tmp.Add(px, qx)
	s.DivMod(s, tmp, f)
	x.Mul(s, s)
	x.Add(x, s)
	x.Add(x, px)
	x.Add(x, qx)
	x.Add(x, a)
	x.Mod(x, f)
	tmp.Add(px, x)
	y.Mul(s, tmp)
	y.Add(y, x)
	y.Add(y, py)
	y.Mod(y, f)

	return x_, y_
}
func double(x_, y_, px_, py_ *big.Int, c *Point) (*big.Int, *big.Int) {
	x, y := wrapBFI(x_), wrapBFI(y_)
	px, py := wrapBFI(px_), wrapBFI(py_)
	s := newBFI()
	f := wrapBFI(c.Curve.Field.Poly)
	a := wrapBFI(c.Curve.Params.A)
	s.DivMod(py, px, f)
	s.Add(px, s)
	x.Mul(s, s)
	x.Add(x, s)
	x.Add(x, a)
	x.Mod(x, f)
	y.Mul(px, px)
	s.Add(s, wrapBFI(one))
	s.Mul(s, x)
	y.Add(y, s)
	y.Mod(y, f)

	return x_, y_
}
