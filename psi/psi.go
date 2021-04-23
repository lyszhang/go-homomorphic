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
	"github.com/lyszhang/go-homomorphic/psi/utils"
	"github.com/renproject/secp256k1"
	"github.com/renproject/shamir/eea"
	"github.com/renproject/shamir/poly"
	"math/big"
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
func CheckPolyMatchSolutionNew(p *VectorBigInt, s int64) bool {
	spow := big.NewInt(1)
	sum := big.NewInt(0)
	for i := 0; i < len(p.Data); i++ {
		if i > 0 {
			spow = big.NewInt(0).Mul(spow, big.NewInt(s))
		}
		mul := big.NewInt(0).Mul(p.Data[i], spow)
		sum = big.NewInt(0).Add(sum, mul)
	}
	if sum.Cmp(big.NewInt(0)) == 0 {
		return true
	}
	return false
}

const INT64LIMIT = "9223372036854775808"

func CheckPolyMatchSolutionNew2(p *VectorBigInt, s int64) bool {
	spow := big.NewInt(1)
	sum := big.NewInt(0)
	sbig := big.NewInt(s)

	for i := 0; i < len(p.Data); i++ {
		if i > 0 {
			spow = big.NewInt(0).Mul(spow, sbig)
		}

		pdbig := p.Data[i]
		mulbig := big.NewInt(0).Mul(pdbig, spow)
		sum = big.NewInt(0).Add(sum, mulbig)
	}
	//fmt.Println("sum: ", sum)
	int64Limit, _ := big.NewInt(0).SetString(INT64LIMIT, 10)
	rem := big.NewInt(0).Mod(sum, int64Limit)
	if big.NewInt(0).Cmp(rem) == 0 {
		return true
	}

	return false
}

func SearchMatchSolutionSetNew(p *VectorBigInt, v []int64) []int64 {
	// 如果是常数式，则肯定没有解
	if len(p.Data) <= 1 {
		return nil
	}
	var inter []int64
	for _, value := range v {
		neg := int64(0 - value)
		if CheckPolyMatchSolutionNew2(p, neg) == true {
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
	//AliceVector.Print()
	//BobVector.Print()

	// encrypt
	encALice := AliceVector.Encrypt(&privKey.PublicKey)
	encBob := BobVector.Encrypt(&privKey.PublicKey)

	// E(f*r+g)
	///TODO: 验证有效性
	encALice.Mul(&RandVector, &privKey.PublicKey)
	encALice.Add(encBob, &privKey.PublicKey)

	// decrypt
	finalVector := encALice.Decrypt(privKey)
	finalVector.Print()

	// 公因式
	a := NewFromVectorInt(&AliceVector)
	b := NewFromVectorInt(finalVector)
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

	//fmt.Println(a.String())
	//fmt.Println(b.String())

	eea.Trial(a, b)

}

// 5个元素为1组，一一求并集
func ProcessByGroup(aliceSet, bobSet []int64) {
	// Generate a 128-bit private key.
	privKey, err := paillier.GenerateKey(rand.Reader, 128)
	if err != nil {
		fmt.Println(err)
		return
	}

	AliceVectors := NewPolysFromSet(aliceSet)
	BobVectors := NewPolysFromSet(bobSet)
	fmt.Println("AliceVectors: ", AliceVectors)
	fmt.Println("BobVectors: ", BobVectors)

	// encrypt
	encALice := EncryptVectors(AliceVectors, &privKey.PublicKey)
	encBob := EncryptVectors(BobVectors, &privKey.PublicKey)

	// E(f*r+g)
	///TODO: 验证有效性
	var encVectors []*EncVector
	for _, vectorAlice := range encALice {
		for _, vectorBoob := range encBob {
			vectorTmp := *vectorAlice

			vectorTmp.Mul(&RandVector, &privKey.PublicKey)
			vectorTmp.Add(vectorBoob, &privKey.PublicKey)
			encVectors = append(encVectors, &vectorTmp)
		}
	}

	// decrypt
	//finalVectors := DecryptVectors(encVectors, privKey)

	// 公因式
	ss := utils.SplitArray(aliceSet, 5)
	var inter []int64
	for i, vectA := range AliceVectors {
		for _, vectF := range BobVectors {
			a := NewFromVectorInt(&vectA)
			b := NewFromVectorInt(&vectF)
			poly := HCF(a, b)

			inter = append(inter, SearchMatchSolutionSetNew(poly, ss[i])...)

		}
	}
	fmt.Println("intersection:", inter)
}
