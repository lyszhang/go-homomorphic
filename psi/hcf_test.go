/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/15 5:26 PM
 */

package psi

import (
	"fmt"
	"math/big"
	"testing"
)

func TestHCF(t *testing.T) {
	//m := VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(-2), big.NewInt(3)}}
	//n := VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(1)}}
	//
	//HCF(&m, &n)

	m1 := &VectorBigInt{[]*big.Int{big.NewInt(-6), big.NewInt(11), big.NewInt(-6), big.NewInt(1)}}
	n1 := &VectorBigInt{[]*big.Int{big.NewInt(-8), big.NewInt(14), big.NewInt(-7), big.NewInt(1)}}

	HCF(m1, n1).print()
}

func TestVectorFloat64_sub(t *testing.T) {
	m := &VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(-2), big.NewInt(3)}}
	n := &VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(1)}}

	m.sub(n).print()
}

func TestVectorFloat64_mulConst(t *testing.T) {
	m := &VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(1)}}

	m.mulConst(big.NewInt(3), 1).print()
}

func TestVectorFloat64_divide(t *testing.T) {
	//m1 := &VectorBigInt{[]float64{-1, 1, 0}}
	//n1 := &VectorBigInt{[]float64{-1, 1}}
	//
	//fmt.Println("+++++++++++")
	//m1.divide(n1).print()

	m2 := &VectorBigInt{[]*big.Int{big.NewInt(1), big.NewInt(-1), big.NewInt(3), big.NewInt(-3)}}
	n2 := &VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(-1), big.NewInt(1), big.NewInt(1)}}

	fmt.Println("+++++++++++")
	m2.divide(n2)
}

func TestVectorFloat64_reduce(t *testing.T) {
	m1 := VectorBigInt{[]*big.Int{big.NewInt(1), big.NewInt(-1), big.NewInt(3), big.NewInt(-3), big.NewInt(0)}}
	m2 := VectorBigInt{[]*big.Int{big.NewInt(-1), big.NewInt(1), big.NewInt(0)}}

	m1.reduce()
	m1.print()

	m2.reduce()
	m2.print()
}
