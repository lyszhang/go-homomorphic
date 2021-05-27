/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/13 3:18 PM
 */

package psi

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

const NUMBER_LIMIT = 10000
const SET_LIMIT = 100

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
	t = append(t, 8878)
	return
}

func TestPSI(t *testing.T) {
	// 集合
	//AliceSet := []int64{10, 2, 3, 6, 8, 34, 6643, 39, 99, 2349, 9340, 8128, 5845, 1234, 8999, 8888, 8788, 1888,
	//	2999, 2383, 6372, 5468, 9892, 9999, 4589, 4444, 9043, 4893, 5343, 4434}
	//BobSet := []int64{1, 22, 3, 4, 5, 9, 8, 34, 49, 2349, 3243, 3240, 8128, 4545, 1234, 233, 8999, 8888, 8788, 1888,
	//	2999, 6372, 5468, 9892, 9589, 4444, 9043, 3045, 3434, 4343}

	AliceSet := randomSet()
	BobSet := randomSet()
	//AliceSet := []int64{260, 186, 146, 278, 286, 166, 39, 134, 40, 250, 225, 212, 36, 156, 257, 258, 263, 83, 176, 20}
	//BobSet := []int64{209, 58, 128, 82, 278, 42, 93, 125, 71, 109, 238, 149, 180, 199, 235, 25, 47, 12, 61, 285}

	fmt.Println("AliceSet: ", AliceSet)
	fmt.Println("BobSet: ", BobSet)

	var inter []int64
	for _, a := range AliceSet {
		for _, b := range BobSet {
			if a == b {
				inter = append(inter, a)
			}
		}
	}

	//Process(AliceSet, BobSet)

	ProcessByGroup(AliceSet, BobSet)
	fmt.Println("real inter: ", inter)

}

func TestPSINew(t *testing.T) {
	// 集合
	//AliceSet := []int64{10, 2, 3, 6, 8, 34, 6643, 39, 99, 2349, 9340, 8128, 5845, 1234, 8999, 8888, 8788, 1888,
	//	2999, 2383, 6372, 5468, 9892, 9999, 4589, 4444, 9043, 4893, 5343, 4434}
	//BobSet := []int64{1, 22, 3, 4, 5, 9, 8, 34, 49, 2349, 3243, 3240, 8128, 4545, 1234, 233, 8999, 8888, 8788, 1888,
	//	2999, 6372, 5468, 9892, 9589, 4444, 9043, 3045, 3434, 4343}

	//AliceSet := randomSet()
	//BobSet := randomSet()
	AliceSet := []int64{20, 53, 41, 34, 216, 286}
	BobSet := []int64{268, 77, 34, 138, 191, 82}

	fmt.Println("AliceSet: ", AliceSet)
	fmt.Println("BobSet: ", BobSet)

	var inter []int64
	for _, a := range AliceSet {
		for _, b := range BobSet {
			if a == b {
				inter = append(inter, a)
			}
		}
	}

	Process(AliceSet, BobSet)
	fmt.Println("real inter: ", inter)

}

func TestCheckPolyMatchSolutionNew2(t *testing.T) {
	//AliceSet := []int64{1, 2, 3, 6, 8, 34, 39, 99, 2349, 23488734, 2359544389060, 2394932905909045, 23958804856}
	AliceSet := []int64{1, 2, 3, 6, 8, 34, 39, 99, 2349, 23488734, 2359544389060}
	AliceVector := NewPolyFromSet(AliceSet)
	a := NewFromVectorInt(&AliceVector)
	a.print()
	fmt.Println(CheckPolyMatchSolutionNew2(a, -2))
}

func TestBigInt(t *testing.T) {
	a, _ := big.NewInt(0).SetString("12878155989941481344", 10)
	b, _ := big.NewInt(0).SetString("-18866153774597651776", 10)
	c := big.NewInt(0).Add(a, b)
	fmt.Println("c:", c)
}
