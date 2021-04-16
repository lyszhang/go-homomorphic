/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/15 4:46 PM
 */

package psi

import "fmt"

// 存储多项式的阶数
// 系数排列从低到高
// 例如 x^2+3x+5
// Data = 5，3，1
type VectorFloat64 struct {
	Data []float64
}

func NewFromVectorInt(vi *Vector) *VectorFloat64 {
	var t VectorFloat64
	for _, value := range vi.Data {
		t.Data = append(t.Data, float64(value))
	}
	return &t
}

// 多项式打印
func (v *VectorFloat64) print() {
	fmt.Println(v)
}

// 多项式阶数，移除高阶系数为0
func (v *VectorFloat64) reduce() {
	if len(v.Data) == 0 {
		return
	}

	cut := 0
	for i := len(v.Data) - 1; i >= 0; i-- {
		if v.Data[i] != 0 {
			cut = i
			break
		}
	}
	v.Data = v.Data[:cut+1]
}

// 多项式阶数
func (v *VectorFloat64) degree() uint {
	// 先做降阶处理
	v.reduce()
	if len(v.Data) == 0 {
		return 0
	}
	return uint(len(v.Data) - 1)
}

// v - d
// f阶数大于等于divsor的阶数
// 多项式相减
func (v *VectorFloat64) sub(d *VectorFloat64) *VectorFloat64 {
	if v.degree() < d.degree() {
		panic("v degree lower than d")
	}

	var tmp []float64
	for i, _ := range v.Data {
		if i < len(d.Data) {
			tmp = append(tmp, v.Data[i]-d.Data[i])
		} else {
			tmp = append(tmp, v.Data[i])
		}
	}

	rem := &VectorFloat64{Data: tmp}
	rem.reduce()
	return rem
}

// v*d
// 多项式系数乘上常数
func (v *VectorFloat64) mulConst(d float64, rshift uint) *VectorFloat64 {
	tmp := make([]float64, rshift)
	for i := 0; i < len(v.Data); i++ {
		tmp = append(tmp, v.Data[i]*d)
	}
	mul := &VectorFloat64{Data: tmp}
	mul.reduce()
	return mul
}

// 求最高阶系数
func (v *VectorFloat64) largestParameter() float64 {
	if len(v.Data) == 0 {
		panic(1)
	}
	for i := len(v.Data) - 1; i >= 0; i-- {
		if v.Data[i] != 0 {
			return v.Data[i]
		}
	}
	return 0
}

// v/divsor
// v阶数大于等于divsor的阶数
// 多项式相除，返回余数多项式
func (v *VectorFloat64) divide(d *VectorFloat64) (rem, divsor *VectorFloat64) {
	// f阶数大于等于divsor的阶数
	x := v.largestParameter()
	y := d.largestParameter()

	// you know why, avoid float missing precision
	t1 := v.mulConst(y, 0)
	t2 := d.mulConst(x, 0)

	fmt.Println("t1: ", t1)
	fmt.Println("t2: ", t2)
	t3 := t2.mulConst(1, t1.degree()-t2.degree())
	rem = t1.sub(t3)
	rem.reduce()

	divsor = d
	return
}

func (v *VectorFloat64) isZero() bool {
	v.reduce()
	if (len(v.Data) == 1 && v.Data[0] == 0) || (len(v.Data) == 0) {
		return true
	}
	return false
}

func HCF(m, n *VectorFloat64) *VectorFloat64 {
	var maxV, minV *VectorFloat64
	var t1, t2 *VectorFloat64
	t1 = m
	t2 = n

	for {
		if t1.degree() >= t2.degree() {
			maxV = t1
			minV = t2
		} else {
			maxV = t2
			minV = t1
		}
		fmt.Println("++++++++++")
		maxV.print()
		minV.print()

		rem, disvor := maxV.divide(minV)
		rem.print()
		disvor.print()
		if rem.isZero() {
			minV.print()
			return minV
		}

		t1 = rem
		t2 = disvor
	}
}
