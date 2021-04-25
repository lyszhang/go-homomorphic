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
	"math/big"
)

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
