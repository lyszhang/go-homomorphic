/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/23 9:50 AM
 */

package utils

import "fmt"

func Gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	if a < b {
		a, b = b, a
	}
	if b == 0 {
		return a
	} else {
		return Gcd(b, a%b)
	}
}

// 最小公倍数， 及各自相差倍数
func Lcm(a, b int64) (int64, int64, int64) {
	g := Gcd(a, b)
	fmt.Println(g)
	return (a / g) * b, b / g, a / g
}
