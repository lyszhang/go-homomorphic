/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/13 3:18 PM
 */

package psi

import (
	"fmt"
	"github.com/renproject/secp256k1"
	"github.com/renproject/shamir/poly"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// 集合
	AliceSet := []int64{1, 2, 3, 6, 8, 34, 39, 99}
	BobSet := []int64{1, 2, 3, 4, 5, 9, 8, 34, 49}

	Process(AliceSet, BobSet)
}

func TestCheckPolyMatchSolution(t *testing.T) {
	p := poly.NewFromSlice([]secp256k1.Fn{secp256k1.NewFnFromU16(6),
		secp256k1.NewFnFromU16(11), secp256k1.NewFnFromU16(6), secp256k1.NewFnFromU16(1)})
	fmt.Println(CheckPolyMatchSolution(p, -4))
}

func TestProcessPoly(t *testing.T) {
	ProcessPoly()
}
