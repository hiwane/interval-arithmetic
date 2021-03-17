package ia

// I/F は big.Int に揃える
import (
	"fmt"
	"math/big"
)

type Interval struct {
	lv *big.Float // lower value
	uv *big.Float // upper value

	lb, ub bool // true if strict
}

func (z *Interval) String() string {
	var lb, ub rune
	if z.lb {
		lb = '('
	} else {
		lb = '['
	}
	if z.ub {
		ub = ')'
	} else {
		ub = ']'
	}
	return fmt.Sprintf("%c%s,%s%c", lb, z.lv.String(), z.uv.String(), ub)
}

func newInterval(lv, uv *big.Float, lb, ub bool) *Interval {
	// assume: lv <= uv
	lv.SetMode(big.ToNegativeInf)
	uv.SetMode(big.ToPositiveInf)
	f := new(Interval)
	f.lv = lv
	f.uv = uv
	f.lb = lb
	f.ub = ub
	return f
}

func NewInterval() *Interval {
	lv := new(big.Float)
	uv := new(big.Float)
	lv.SetMode(big.ToNegativeInf)
	uv.SetMode(big.ToPositiveInf)
	return newInterval(lv, uv, false, false)
}

func NewIntervalInt64(n int64, prec uint) *Interval {
	// lower bound
	lv := new(big.Float)
	lv.SetMode(big.ToNegativeInf)
	lv.SetPrec(prec)
	lv.SetInt64(n)

	// upper bound
	uv := new(big.Float)
	uv.SetMode(big.ToPositiveInf)
	uv.SetPrec(prec)
	uv.SetInt64(n)

	bound := (lv.Cmp(uv) != 0)
	return newInterval(lv, uv, bound, bound)
}

func NewIntervalFloat64(n float64, prec uint) *Interval {
	// lower bound
	lv := new(big.Float)
	lv.SetMode(big.ToNegativeInf)
	lv.SetPrec(prec)
	lv.SetFloat64(n)

	// upper bound
	uv := new(big.Float)
	uv.SetPrec(prec)
	uv.SetMode(big.ToPositiveInf)
	uv.SetFloat64(n)

	bound := (lv.Cmp(uv) != 0)
	return newInterval(lv, uv, bound, bound)
}

func NewIntervalStr(s string, base int, prec uint) (*Interval, error) {
	lv, _, err := big.ParseFloat(s, base, prec, big.ToNegativeInf)
	if err != nil {
		return nil, err
	}
	uv, _, err := big.ParseFloat(s, base, prec, big.ToPositiveInf)
	if err != nil {
		return nil, err
	}
	bound := false
	if lv.Cmp(uv) != 0 {
		bound = true
	}
	return newInterval(lv, uv, bound, bound), nil
}

func (z *Interval) SetPrec(prec uint) {
	z.lv.SetPrec(prec)
	z.uv.SetPrec(prec)
}

func (z *Interval) SetFloat64x(low, up float64) {
	z.lv.SetFloat64(low)
	z.uv.SetFloat64(up)
}

func (c *Interval) ContainsZero() bool {
	lsgn := c.lv.Sign()
	usgn := c.uv.Sign()

	return (lsgn < 0 || lsgn == 0 && !c.lb) &&
		(usgn > 0 || usgn == 0 && !c.ub)
}

func (c *Interval) ContainsFloat(x *big.Float) bool {
	lc := c.lv.Cmp(x)
	uc := x.Cmp(c.uv)
	return (lc < 0 || lc == 0 && !c.lb) && (uc > 0 || uc == 0 && !c.ub)
}

func (z *Interval) Neg(x *Interval) *Interval {
	z.uv.Neg(x.lv)
	z.lv.Neg(x.uv)

	z.lb = x.ub
	z.ub = x.lb
	return z
}

func MaxPrec(x, y *Interval) uint {
	if x.lv.Prec() <= y.lv.Prec() {
		return x.lv.Prec()
	} else {
		return y.lv.Prec()
	}
}

func (z *Interval) Add(x *Interval, y *Interval) *Interval {
	z.lv.Add(x.lv, y.lv)
	z.uv.Add(x.uv, y.uv)

	z.lb = x.lb || y.lb
	z.ub = x.ub || y.ub
	return z
}

func (z *Interval) Sub(x *Interval, y *Interval) *Interval {
	z.lv.Sub(x.lv, y.lv)
	z.uv.Sub(x.uv, y.uv)

	z.lb = x.lb || y.lb
	z.ub = x.ub || y.ub
	return z
}

func (z *Interval) Mul(x *Interval, y *Interval) *Interval {
	if x.lv.Sign() >= 0 {
		if y.lv.Sign() >= 0 {
			z.lv.Mul(x.lv, y.lv)
			z.uv.Mul(x.uv, y.uv)
			z.lb = x.lb || y.lb
			z.ub = x.ub || y.ub
		} else if y.uv.Sign() <= 0 {
			// x >= 0, y <= 0
			z.lv.Mul(x.uv, y.lv)
			z.uv.Mul(x.lv, y.uv)
			z.lb = x.ub || y.lb
			z.ub = x.lb || y.ub
		} else {
			z.lv.Mul(x.uv, y.lv)
			z.uv.Mul(x.uv, y.uv)
			z.lb = x.ub || y.lb
			z.ub = x.ub || y.ub
		}
	} else if x.uv.Sign() <= 0 {
		if y.lv.Sign() >= 0 {
			z.lv.Mul(x.lv, y.uv)
			z.uv.Mul(x.uv, y.lv)
			z.lb = x.lb || y.ub
			z.ub = x.ub || y.lb
		} else if y.uv.Sign() <= 0 {
			// [-xl, -xu] * [-yl, -yu] => [xu*yu, xl*yl]
			z.lv.Mul(x.uv, y.uv)
			z.uv.Mul(x.lv, y.lv)
			z.lb = x.ub || y.ub
			z.ub = x.lb || y.lb
		} else {
			// [-xl, -xu] * [-yl, +yu] => [xu*yu, xl*yl]
			z.lv.Mul(x.lv, y.uv)
			z.uv.Mul(x.lv, y.lv)
			z.lb = x.lb || y.ub
			z.ub = x.lb || y.lb
		}
	} else {
		if y.lv.Sign() >= 0 {
			// [-xl, xu] * [yl, yu]
			z.lv.Mul(x.lv, y.uv)
			z.uv.Mul(x.uv, y.uv)
			z.lb = x.lb || y.ub
			z.ub = x.ub || y.ub
		} else if y.uv.Sign() <= 0 {
			// [-xl, xu] * [-yl, -yu]
			// [-xl, xu] * (-yl, -yu]
			z.lv.Mul(x.uv, y.lv)
			z.uv.Mul(x.lv, y.lv)
			z.lb = x.ub || y.lb
			z.ub = x.lb || y.lb
		} else {
			// [-xl, +xu] * [-yl, +yu] => [min(-xl*yu,-xu,yl), max(xl*yl,xu*yu)]
			u := new(big.Float)
			u.SetPrec(z.lv.Prec())
			u.SetMode(big.ToNegativeInf)
			u.Mul(x.lv, y.uv)
			z.lv.Mul(x.uv, y.lv)
			cmp := u.Cmp(z.lv)
			if cmp < 0 {
				z.lv.Set(u)
				z.lb = x.lb || y.ub
			} else if cmp > 0 {
				z.lb = x.ub || y.lb
			} else {
				z.lb = (x.ub || y.lb) && (x.lb || y.ub)
			}
			u.SetMode(big.ToPositiveInf)
			u.Mul(x.lv, y.lv)
			z.uv.Mul(x.uv, y.uv)
			cmp = u.Cmp(z.uv)
			if cmp > 0 {
				z.uv.Set(u)
				z.ub = x.lb || y.lb
			} else if cmp < 0 {
				z.ub = x.ub || y.ub
			} else {
				z.ub = (x.lb || y.lb) && (x.ub || y.ub)
			}
		}
	}

	return z
}
