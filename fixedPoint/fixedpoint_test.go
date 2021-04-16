/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/2/25 2:36 PM
 */

package fixedPoint

import (
	"fmt"
	"math"
	"testing"
)

func Test_reciprocal(t *testing.T) {
	type args struct {
		c float64
		n int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test01",
			args{c: 1500, n: 20},
			0.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reciprocal(tt.args.c, tt.args.n); got != tt.want {
				t.Errorf("reciprocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

type polyMath struct {
	param []float64
}

func (p *polyMath) Calc(x float64) (r float64) {
	for i, v := range p.param {
		r = r + math.Pow(x, float64(i))*v
	}
	return
}

func (p *polyMath) DeepCopy() *polyMath {
	t := polyMath{}
	for _, v := range p.param {
		t.param = append(t.param, v)
	}
	return &t
}

func (p *polyMath) Mul(n float64) {
	for i, v := range p.param {
		p.param[i] = v * n
	}
	return
}

func (p *polyMath) Add(q *polyMath) {
	var tmp []float64
	if len(p.param) > len(q.param) {
		for i, _ := range p.param {
			if i < len(q.param) {
				tmp = append(tmp, p.param[i]+q.param[i])
			} else {
				tmp = append(tmp, p.param[i])
			}
		}
	} else {
		for i, _ := range q.param {
			if i < len(p.param) {
				tmp = append(tmp, p.param[i]+q.param[i])
			} else {
				tmp = append(tmp, q.param[i])
			}
		}
	}
	p.param = tmp
	return
}

//右移，左边补零
func (p *polyMath) RightShift(n int) {
	if n < 1 {
		return
	}
	p.param = append(p.param, make([]float64, n)...)
	for i := len(p.param) - 1; i >= 0; i-- {
		if i >= n {
			p.param[i] = p.param[i-n]
		} else {
			p.param[i] = 0 //左侧补零
		}
	}
	return
}

func (p *polyMath) AddConst(m float64) {
	if len(p.param) > 0 {
		p.param[0] = p.param[0] + m
	} else {
		p.param = append(p.param, m)
	}
}

func (p *polyMath) MulVector(q *polyMath) {
	var res polyMath
	for m, _ := range q.param {
		t := p.DeepCopy()
		t.RightShift(m)
		t.Mul(q.param[m])
		res.Add(t)
	}
	p.param = res.param
}

func TestCalculate(t *testing.T) {
	polys := polyMath{param: []float64{1.0, 1.0}}

	polys2 := polyMath{param: []float64{1.0, 1.0}}
	polys.MulVector(&polys2)
	fmt.Println(polys)
}

func TestCalculate2(t *testing.T) {
	n := 10
	p := getPolyInfo(n)
	fmt.Println(p)
	fmt.Println(p.Calc(3))
}

const xr = 2

func getPolyInfo(n int) *polyMath {
	if n == 1 {
		return &polyMath{param: []float64{2 * xr, -math.Pow(xr, 2)}}
	}
	x := getPolyInfo(n - 1)
	y := x.DeepCopy()
	x.MulVector(&polyMath{param: []float64{0, -1.0}})
	x.AddConst(2)
	x.MulVector(y)
	return x
}
