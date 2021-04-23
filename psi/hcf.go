/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/15 4:46 PM
 */

package psi

import (
	"fmt"
	"github.com/lyszhang/go-homomorphic/psi/utils"
	"math/big"
)

// 存储多项式的阶数
// 系数排列从低到高
// 例如 x^2+3x+5
// Data = 5，3，1
type VectorBigInt struct {
	Data []*big.Int
}

func NewFromVectorInt(vi *Vector) *VectorBigInt {
	var t VectorBigInt
	for _, value := range vi.Data {
		t.Data = append(t.Data, big.NewInt(0).SetInt64(value))
	}
	return &t
}

// 多项式打印
func (v *VectorBigInt) print() {
	fmt.Println(v.string())
}

// 多项式打印
func (v *VectorBigInt) string() string {
	var str string
	for i, value := range v.Data {
		if i == 0 {
			str = fmt.Sprintf("%d ", value)
			continue
		}
		str += fmt.Sprintf("+ %d x^%d", value, i)
	}
	return str
}

// 多项式阶数，移除高阶系数为0
func (v *VectorBigInt) reduce() {
	if len(v.Data) == 0 {
		return
	}

	cut := 0
	for i := len(v.Data) - 1; i >= 0; i-- {
		if v.Data[i].Cmp(big.NewInt(0)) != 0 {
			cut = i
			break
		}
	}
	v.Data = v.Data[:cut+1]
}

// 多项式阶数
func (v *VectorBigInt) degree() uint {
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
func (v *VectorBigInt) sub(d *VectorBigInt) *VectorBigInt {
	if v.degree() < d.degree() {
		panic("v degree lower than d")
	}

	var tmp []*big.Int
	for i, _ := range v.Data {
		if i < len(d.Data) {
			fmt.Println("v", v.Data[i])
			fmt.Println("d", d.Data[i])
			tmp = append(tmp, big.NewInt(0).Sub(v.Data[i], d.Data[i]))
		} else {
			tmp = append(tmp, v.Data[i])
		}
	}

	rem := &VectorBigInt{Data: tmp}
	rem.reduce()
	return rem
}

// v*d
// 多项式系数乘上常数
func (v *VectorBigInt) mulConst(d *big.Int, rshift uint) *VectorBigInt {
	fmt.Println("rshift: ", rshift)
	tmp := make([]*big.Int, rshift)
	for i, _ := range tmp {
		tmp[i] = big.NewInt(0)
	}

	for i := 0; i < len(v.Data); i++ {
		tmp = append(tmp, big.NewInt(0).Mul(v.Data[i], d))
	}
	fmt.Println(tmp)
	mul := &VectorBigInt{Data: tmp}
	mul.reduce()
	return mul
}

// 求最高阶系数
func (v *VectorBigInt) largestParameter() *big.Int {
	if len(v.Data) == 0 {
		panic(1)
	}
	for i := len(v.Data) - 1; i >= 0; i-- {
		if v.Data[i].Cmp(big.NewInt(0)) != 0 {
			return v.Data[i]
		}
	}
	return big.NewInt(0)
}

// v/divsor
// v阶数大于等于divsor的阶数
// 多项式相除，返回余数多项式
func (v *VectorBigInt) divide(d *VectorBigInt) (rem, divsor *VectorBigInt) {
	// f阶数大于等于divsor的阶数
	x := v.largestParameter()
	y := d.largestParameter()
	fmt.Println("+++++++++")
	v.print()
	d.print()
	fmt.Println("x: ", x)
	fmt.Println("y: ", y)

	_, xa, ya := utils.Lcm(x, y)
	///// 为什么出现负数的情况
	//fmt.Println("xa", xa)
	//fmt.Println("ya", ya)

	// you know why, avoid float missing precision
	t1 := v.mulConst(xa, 0)
	t2 := d.mulConst(ya, 0)

	/// 为什么出现负数的情况
	fmt.Println(t1.degree())
	fmt.Println(t2.degree())
	t3 := t2.mulConst(big.NewInt(1), t1.degree()-t2.degree())
	rem = t1.sub(t3)
	rem.reduce()

	divsor = d
	return
}

func (v *VectorBigInt) isZero() bool {
	v.reduce()
	if (len(v.Data) == 1 && v.Data[0].Cmp(big.NewInt(0)) == 0) || (len(v.Data) == 0) {
		return true
	}
	return false
}

func HCF(m, n *VectorBigInt) *VectorBigInt {
	var maxV, minV *VectorBigInt
	var t1, t2 *VectorBigInt
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
		//fmt.Println("++++++++++")
		maxV.print()
		minV.print()

		rem, disvor := maxV.divide(minV)
		//rem.print()
		//disvor.print()
		if rem.isZero() {
			minV.print()
			return minV
		}

		t1 = rem
		t2 = disvor
	}
}
