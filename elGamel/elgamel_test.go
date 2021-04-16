/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 3:30 PM
 */

package elGamel

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestElgamel(t *testing.T) {
	val1 := "3"
	val2 := "4"

	fmt.Printf("Prime number size: %s\n\n", plen)

	priv := CreatePrivateKey()

	e1, _ := Encrypt(rand.Reader, &priv.PublicKey, val1)
	e2, _ := Encrypt(rand.Reader, &priv.PublicKey, val2)

	m := Mul(e1, e2)

	message2, _ := Decrypt(priv, m.X, m.Y)

	fmt.Printf("====Values (val1=%s and val2=%s)\n", val1, val2)
	fmt.Printf("====Private key (x):\nX=%d", priv.X)
	fmt.Printf("\n\n====Public key (Y,G,P):\nY=%d\nG=%d\nP=%d", priv.Y, priv.PublicKey.G, priv.PublicKey.P)

	fmt.Printf("\n\n====Cipher (a1=%s)\n\n(b1=%s): ", e1.X, e1.Y)
	fmt.Printf("\n\n====Decrypted: %d", valint(message2))
}

func BenchmarkMul(b *testing.B) {
	val1 := "3"
	val2 := "4"
	priv := CreatePrivateKey()

	e1, _ := Encrypt(rand.Reader, &priv.PublicKey, val1)
	e2, _ := Encrypt(rand.Reader, &priv.PublicKey, val2)

	for i := 0; i < b.N; i++ {
		Mul(e1, e2)
	}
}
