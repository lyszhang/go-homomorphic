/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/14 3:12 PM
 */

package psi

import (
	"fmt"
	"testing"
)

func TestIterationN(t *testing.T) {
	t1 := Vector{
		Data: []int64{1, 2, 3, 1},
	}
	fmt.Println(IterationN(t1, 4))
}

func TestNewPolyFromSet(t *testing.T) {
	t1 := set{1, 2, 3, 4}
	fmt.Println(NewPolyFromSet(t1))
}
