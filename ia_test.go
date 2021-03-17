package ia

import (
	"math/big"
	"testing"
)

func TestNewIntervalInt64(t *testing.T) {
	for _, v := range []struct {
		n    int64
		prec uint
	}{
		{3, 10}, {5, 10},
	} {
		f := NewIntervalInt64(v.n, v.prec)
		if f.lv.Cmp(f.uv) != 0 || f.lb || f.ub {
			t.Errorf("invalid1 input=%v cmp=%d: %s", v, f.lv.Cmp(f.uv), f.String())
		}
	}

	for _, v := range []struct {
		n    int64
		prec uint
	}{
		{65535, 10},
	} {
		f := NewIntervalInt64(v.n, v.prec)
		if f.lv.Cmp(f.uv) >= 0 || !f.lb || !f.ub {
			t.Errorf("invalid2 input=%v cmp=%d: %s", v, f.lv.Cmp(f.uv), f.String())
		}
	}
}

func ContaintsFloat(t *testing.T) {
	x := NewInterval()
	x.SetFloat64x(2, 5)
	f := big.NewFloat(3)
	if !x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	f = big.NewFloat(2)
	if !x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	x.ub = true
	if !x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	x.lb = true
	if x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	x.lb = false
	x.ub = false
	f = big.NewFloat(5)
	if !x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	x.lb = true
	if !x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
	x.lb = false
	if x.ContainsFloat(f) {
		t.Errorf("invalid2 x=%v, f=%v\n", x, f)
	}
}

func test_equals(z *Interval, low, up float64, lb, ub bool) bool {
	return z.lv.Cmp(big.NewFloat(low)) == 0 && z.uv.Cmp(big.NewFloat(up)) == 0 &&
		z.lb == lb && z.ub == ub
}

func TestNeg(t *testing.T) {
	x := NewInterval()
	x.SetFloat64x(-3, 5)

	z := NewInterval()
	z.Neg(x)

	if !test_equals(z, -5, 3, x.ub, x.lb) {
		t.Errorf("x=%v, z=%v\n", x, z)
	}

	x.lb = true
	z.Neg(x)
	if !test_equals(z, -5, 3, x.ub, x.lb) {
		t.Errorf("x=%v, z=%v\n", x, z)
	}

	x.ub = true
	z.Neg(x)
	if !test_equals(z, -5, 3, x.ub, x.lb) {
		t.Errorf("x=%v, z=%v\n", x, z)
	}

	x.lb = false
	z.Neg(x)
	if !test_equals(z, -5, 3, x.ub, x.lb) {
		t.Errorf("x=%v, z=%v\n", x, z)
	}
}

func testMul(t *testing.T, x, y *Interval, flow, fup float64, lb, ub bool) {
	z := NewInterval()
	z.Mul(x, y)
	if !test_equals(z, flow, fup, lb, ub) {
		t.Errorf("x=%v, y=%v, z=%v %v.%v", x, y, z, lb, ub)
	}

	z = NewInterval()
	z.Mul(y, x)
	if !test_equals(z, flow, fup, lb, ub) {
		t.Errorf("y=%v, x=%v, z=%v %v.%v", y, x, z, lb, ub)
	}

	xm := NewInterval()
	xm.Neg(x)

	ym := NewInterval()
	ym.Neg(y)

	z = NewInterval()
	z.Mul(xm, ym)
	if !test_equals(z, flow, fup, lb, ub) {
		t.Errorf("-x=%v, -y=%v, z=%v %v.%v", x, y, z, lb, ub)
	}

	z = NewInterval()
	z.Mul(ym, xm)
	if !test_equals(z, flow, fup, lb, ub) {
		t.Errorf("-y=%v, -x=%v, z=%v %v.%v", y, x, z, lb, ub)
	}
}

func TestMul(t *testing.T) {
	x := NewInterval()
	y := NewInterval()

	// ++
	x.SetFloat64x(2, 3) // [2, 3]
	y.SetFloat64x(4, 5) // [4, 5]

	testMul(t, x, y, 8, 15, false, false)

	x.lb = true
	testMul(t, x, y, 8, 15, true, false)

	x.lb = false
	y.lb = true
	testMul(t, x, y, 8, 15, true, false)

	x.lb = false
	y.lb = false
	x.ub = true
	y.ub = false
	testMul(t, x, y, 8, 15, false, true)

	// reset
	x = NewInterval()
	y = NewInterval()
	x.SetFloat64x(-2, 3)
	y.SetFloat64x(4, 5)

	testMul(t, x, y, -10, 15, false, false)

	y.lb = true
	testMul(t, x, y, -10, 15, false, false)

	x.lb = true
	testMul(t, x, y, -10, 15, true, false)

	x.lb = false
	x.ub = true
	testMul(t, x, y, -10, 15, false, true)

	x.lb = true
	x.ub = true
	testMul(t, x, y, -10, 15, true, true)

	x.lb = false
	x.ub = false
	y.ub = true
	testMul(t, x, y, -10, 15, true, true)

	x = NewInterval()
	y = NewInterval()
	x.SetFloat64x(-2, 3)
	y.SetFloat64x(-4, 5)
	testMul(t, x, y, -12, 15, false, false)

	y.lb = true
	testMul(t, x, y, -12, 15, true, false)

	y.lb = false
	x.ub = true
	testMul(t, x, y, -12, 15, true, true)

	x = NewInterval()
	y = NewInterval()
	x.SetFloat64x(-2, 2) // (-2, 2]
	y.SetFloat64x(-3, 3) // [-3, 3)
	x.lb = true
	y.ub = true
	testMul(t, x, y, -6, 6, false, true)

	x = NewInterval()
	y = NewInterval()
	x.SetFloat64x(-2, 2) // (-2, 2]
	y.SetFloat64x(4, 5)  // [4, 5]
	x.lb = true
	testMul(t, x, y, -10, 10, true, false)

	// [-2, 2] * (-4, 5]
	y.lb = true
	testMul(t, x, y, -10, 10, true, false)
}
