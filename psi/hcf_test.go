/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/15 5:26 PM
 */

package psi

import (
	"fmt"
	"testing"
)

func TestHCF(t *testing.T) {
	//m := VectorFloat64{[]float64{-1, -2, 3}}
	//n := VectorFloat64{[]float64{-1, 1}}
	//
	//HCF(&m, &n)

	m1 := &VectorFloat64{[]float64{1, -2, 2, -2, 1}}
	n1 := &VectorFloat64{[]float64{-1, -1, 1, 1}}

	HCF(m1, n1).print()
}

func TestVectorFloat64_sub(t *testing.T) {
	m := &VectorFloat64{[]float64{-1, -2, 3}}
	n := &VectorFloat64{[]float64{-1, 1}}

	m.sub(n).print()
}

func TestVectorFloat64_mulConst(t *testing.T) {
	m := &VectorFloat64{[]float64{-1, 1}}

	m.mulConst(3, 1).print()
}

func TestVectorFloat64_divide(t *testing.T) {
	//m1 := &VectorFloat64{[]float64{-1, 1, 0}}
	//n1 := &VectorFloat64{[]float64{-1, 1}}
	//
	//fmt.Println("+++++++++++")
	//m1.divide(n1).print()

	m2 := &VectorFloat64{[]float64{1, -1, 3, -3}}
	n2 := &VectorFloat64{[]float64{-1, -1, 1, 1}}

	fmt.Println("+++++++++++")
	m2.divide(n2)
}

func TestVectorFloat64_reduce(t *testing.T) {
	m1 := VectorFloat64{[]float64{1, -1, 3, -3, 0}}
	m2 := VectorFloat64{[]float64{-1, 1, 0}}

	m1.reduce()
	m1.print()

	m2.reduce()
	m2.print()
}
