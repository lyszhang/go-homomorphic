/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/13 3:18 PM
 */

package psi

import (
	"crypto/rand"
	"fmt"
	"github.com/renproject/secp256k1"
	"github.com/renproject/shamir/poly"
	"math/big"
	"testing"
)

const NUMBER_LIMIT = 100
const SET_LIMIT = 15

func isInSet(t []int64, n int64) bool {
	for _, value := range t {
		if value == n {
			return true
		}
	}
	return false
}

func randomSet() (t []int64) {
	for i := 0; i < SET_LIMIT; i++ {
		var n *big.Int
		for {
			n, _ = rand.Int(rand.Reader, big.NewInt(NUMBER_LIMIT))
			if !isInSet(t, n.Int64()) {
				break
			}
		}
		t = append(t, n.Int64())
	}
	return
}

func TestPSI(t *testing.T) {
	// 集合
	AliceSet := []int64{10, 2, 3, 6, 8, 34, 39, 99, 2349, 9832940}
	BobSet := []int64{1, 22, 3, 4, 5, 9, 8, 34, 49, 2349, 3249023, 9832940}
	//AliceSet := randomSet()
	//BobSet := randomSet()
	fmt.Println("AliceSet: ", AliceSet)
	fmt.Println("BobSet: ", BobSet)

	Process(AliceSet, BobSet)
}

func TestCheckPolyMatchSolution(t *testing.T) {
	p := poly.NewFromSlice([]secp256k1.Fn{secp256k1.NewFnFromU16(6),
		secp256k1.NewFnFromU16(11), secp256k1.NewFnFromU16(6), secp256k1.NewFnFromU16(1)})
	fmt.Println(CheckPolyMatchSolution(p, -4))
}

func TestCheckPolyMatchSolutionNew2(t *testing.T) {
	//AliceSet := []int64{1, 2, 3, 6, 8, 34, 39, 99, 2349, 23488734, 2359544389060, 2394932905909045, 23958804856}
	AliceSet := []int64{1, 2, 3, 6, 8, 34, 39, 99, 2349, 23488734, 2359544389060}
	AliceVector := NewPolyFromSet(AliceSet)
	a := NewFromVectorInt(&AliceVector)
	a.print()
	fmt.Println(CheckPolyMatchSolutionNew2(a, -2))
}

func TestProcessPoly(t *testing.T) {
	ProcessPoly()
}

func TestBigInt(t *testing.T) {
	a, _ := big.NewInt(0).SetString("12878155989941481344", 10)
	b, _ := big.NewInt(0).SetString("-18866153774597651776", 10)
	c := big.NewInt(0).Add(a, b)
	fmt.Println("c:", c)
}
