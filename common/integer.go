package common

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const Precision = 8

var Zero Integer

func init() {
	Zero = NewInteger(0)
}

type Integer struct {
	i big.Int
}

func NewIntegerFromString(x string) (v Integer, err error) {
	d, err := decimal.NewFromString(x)
	if err != nil {
		return
	}
	if d.Sign() <= 0 {
		err = fmt.Errorf("Invalid x: %s", x)
		return
	}
	s := d.Mul(decimal.New(1, Precision)).StringFixed(0)
	v.i.SetString(s, 10)
	return
}

func NewInteger(x uint64) (v Integer) {
	p := new(big.Int).SetUint64(x)
	d := big.NewInt(int64(math.Pow(10, Precision)))
	v.i.Mul(p, d)
	return
}

func (x Integer) Add(y Integer) (v Integer) {
	v.i.Add(&x.i, &y.i)
	return
}

func (x Integer) Sub(y Integer) (v Integer) {
	v.i.Sub(&x.i, &y.i)
	return
}

func (x Integer) Mul(y int) (v Integer) {
	v.i.Mul(&x.i, big.NewInt(int64(y)))
	return
}

func (x Integer) Div(y int) (v Integer) {
	v.i.Div(&x.i, big.NewInt(int64(y)))
	return
}

func (x Integer) Cmp(y Integer) int {
	return x.i.Cmp(&y.i)
}

func (x Integer) Sign() int {
	return x.i.Sign()
}

func (x Integer) String() string {
	s := x.i.String()
	p := len(s) - Precision
	if p > 0 {
		return s[:p] + "." + s[p:]
	}
	return "0." + strings.Repeat("0", -p) + s
}

func (x Integer) MarshalMsgpack() ([]byte, error) {
	return x.i.Bytes(), nil
}

func (x *Integer) UnmarshalMsgpack(data []byte) error {
	x.i.SetBytes(data)
	return nil
}

func (x Integer) MarshalJSON() ([]byte, error) {
	s := x.String()
	return []byte(strconv.Quote(s)), nil
}

func (x *Integer) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	i, err := NewIntegerFromString(unquoted)
	if err != nil {
		return err
	}
	x.i.SetBytes(i.i.Bytes())
	return nil
}
