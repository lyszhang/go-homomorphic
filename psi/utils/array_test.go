/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/22 11:10 AM
 */

package utils

import (
	"fmt"
	"testing"
)

func Test_splitArray(t *testing.T) {
	test := []int64{1, 2, 3, 4, 5, 5, 6, 6, 7, 8, 2, 4, 5, 4, 6, 5}
	res := SplitArray(test, 5)
	fmt.Println(res)
}
