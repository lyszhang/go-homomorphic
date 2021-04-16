/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/2/25 2:31 PM
 */

package fixedPoint

import "fmt"

func reciprocal(c float64, n int) float64 {
	x := 0.0001
	fmt.Println("x: ", x)
	for i := 0; i < n; i++ {
		x = x * (2 - c*x)
		fmt.Println(x)
	}
	return x
}
