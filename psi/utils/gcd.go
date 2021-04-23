/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/23 9:50 AM
 */

package utils

import (
	"math/big"
)

func Gcd(a, b *big.Int) *big.Int {
	a = big.NewInt(0).Abs(a)
	b = big.NewInt(0).Abs(b)

	if a.Cmp(b) < 0 {
		a, b = b, a
	}
	if b.Cmp(big.NewInt(0)) == 0 {
		return a
	} else {
		return Gcd(b, big.NewInt(0).Mod(a, b))
	}
}

// 最小公倍数， 及各自相差倍数
func Lcm(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	g := Gcd(a, b)
	div := big.NewInt(0).Div(a, g)
	return big.NewInt(0).Mul(div, b), big.NewInt(0).Div(b, g), big.NewInt(0).Div(a, g)
}
