/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/13 2:27 PM
 */

package psi

import (
	"crypto/rand"
	"fmt"
	paillier "github.com/lyszhang/go-go-gadget-paillier"
	"github.com/renproject/secp256k1"
	"github.com/renproject/shamir/eea"
	"github.com/renproject/shamir/poly"
)

// 验证是否为其中的解
func CheckPolyMatchSolution(p poly.Poly, s int64) bool {
	var spow, sum int64
	for i := 0; i < len(p); i++ {
		if i == 0 {
			spow = 1
		} else {
			spow = spow * s
		}
		sum += p[i].Int().Int64() * spow
	}
	if sum == 0 {
		return true
	}
	return false
}

func SearchMatchSolutionSet(p poly.Poly, v []int64) []int64 {
	var inter []int64
	for _, value := range v {
		neg := 0 - value
		fmt.Println("neg: ", neg)
		if CheckPolyMatchSolution(p, neg) == true {
			inter = append(inter, value)
		}
	}
	return inter
}

// 验证是否为其中的解
func CheckPolyMatchSolutionNew(p *VectorFloat64, s float64) bool {
	var spow, sum float64
	for i := 0; i < len(p.Data); i++ {
		if i == 0 {
			spow = 1
		} else {
			spow = spow * s
		}
		sum += p.Data[i] * spow
	}
	if sum == 0 {
		return true
	}
	return false
}

func SearchMatchSolutionSetNew(p *VectorFloat64, v []int64) []int64 {
	var inter []int64
	for _, value := range v {
		neg := float64(0 - value)
		fmt.Println("neg: ", neg)
		if CheckPolyMatchSolutionNew(p, neg) == true {
			inter = append(inter, value)
		}
	}
	return inter
}

// 生成poly形式的多项式
func newPolyFromVector(v *Vector) poly.Poly {
	var fns []secp256k1.Fn
	for _, value := range v.Data {
		fns = append(fns, secp256k1.NewFnFromU32(uint32(value)))
	}
	return poly.NewFromSlice(fns)
}

func Process(aliceSet, bobSet []int64) {
	// Generate a 128-bit private key.
	privKey, err := paillier.GenerateKey(rand.Reader, 128)
	if err != nil {
		fmt.Println(err)
		return
	}

	AliceVector := NewPolyFromSet(aliceSet)
	BobVector := NewPolyFromSet(bobSet)
	AliceVector.Print()
	BobVector.Print()

	// encrypt
	encALice := AliceVector.Encrypt(privKey)
	encBob := BobVector.Encrypt(privKey)

	// E(f*r+g)
	///TODO: 验证有效性
	encALice.Mul(&RandVector, privKey)
	encALice.Add(encBob, privKey)

	// decrypt
	finalVector := encALice.Decrypt(privKey)
	finalVector.Print()

	// 公因式
	// 有个限制, a的degree需要小于b的degree, 否则会报错
	// (x+1)(x+2)(x+3)
	// (x+1)(x+2)(x+4)
	//a := poly.NewFromSlice([]secp256k1.Fn{secp256k1.NewFnFromU32(2),
	//	secp256k1.NewFnFromU32(3), secp256k1.NewFnFromU32(1)})
	//b := poly.NewFromSlice([]secp256k1.Fn{secp256k1.NewFnFromU32(42),
	//	secp256k1.NewFnFromU32(83),
	//	secp256k1.NewFnFromU32(53), secp256k1.NewFnFromU32(13), secp256k1.NewFnFromU32(1)})
	//a := newPolyFromVector(&AliceVector)
	//b := newPolyFromVector(&BobVector)
	//poly, _ := eea.Trial(a, b)
	//
	//inter := SearchMatchSolutionSet(poly, AliceSet)
	//fmt.Println("intersection:", inter)
	a := NewFromVectorInt(&AliceVector)
	b := NewFromVectorInt(&BobVector)
	poly := HCF(a, b)

	inter := SearchMatchSolutionSetNew(poly, aliceSet)
	fmt.Println("intersection:", inter)
}

func newPoly(v Vector) poly.Poly {
	var fn []secp256k1.Fn
	for _, value := range v.Data {
		fn = append(fn, secp256k1.NewFnFromU16(uint16(value)))
	}
	return poly.NewFromSlice(fn)
}

func ProcessPoly() {
	// 集合
	AliceVector := NewPolyFromSet([]int64{1, 2, 3, 7})
	BobVector := NewPolyFromSet([]int64{1, 2, 4, 5, 8})

	// 公因式
	// 有个限制, a的degree需要小于b的degree, 否则会报错
	// (x+1)(x+2)(x+3)
	// (x+1)(x+2)(x+4)
	a := newPoly(AliceVector)
	b := newPoly(BobVector)

	fmt.Println(a.String())
	fmt.Println(b.String())

	eea.Trial(a, b)

}

// 两边的阶数必须是不一样的，否则结果有问题
// 目前只能支持统一的正数
