/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 6:17 PM
 */

package main

import (
	"fmt"
	"github.com/dedis/lago/bigint"
	"github.com/dedis/lago/crypto"
	"github.com/dedis/lago/encoding"
	"os"
	"strconv"
)

func main() {

	a := 1
	b := 3
	argCount := len(os.Args[1:])

	if argCount > 0 {
		a, _ = strconv.Atoi(os.Args[1])
	}

	if argCount > 1 {
		b, _ = strconv.Atoi(os.Args[2])
	}

	msg1 := bigint.NewInt(int64(a))
	msg2 := bigint.NewInt(int64(b))

	N := uint32(32)                                        // polynomial degree
	T := bigint.NewInt(10)                                 // plaintext moduli
	Q := bigint.NewInt(8380417)                            // ciphertext moduli
	BigQ := bigint.NewIntFromString("4611686018326724609") // big ciphertext moduli, used in homomorphic multiplication and should be greater than q^2

	// create FV context and generate keys

	fv := crypto.NewFVContext(N, *T, *Q, *BigQ)
	key := crypto.GenerateKey(fv)

	// encode messages

	encoder := encoding.NewEncoder(fv)
	plaintext1 := crypto.NewPlaintext(N, *Q, fv.NttParams)
	plaintext2 := crypto.NewPlaintext(N, *Q, fv.NttParams)

	encoder.Encode(msg1, plaintext1)
	encoder.Encode(msg2, plaintext2)

	// encrypt plainetexts
	encryptor := crypto.NewEncryptor(fv, &key.PubKey)
	ciphertext1 := encryptor.Encrypt(plaintext1)
	ciphertext2 := encryptor.Encrypt(plaintext2)

	// evaluate ciphertexts
	evaluator := crypto.NewEvaluator(fv, &key.EvaKey, key.EvaSize)

	add_cipher := evaluator.Add(ciphertext1, ciphertext2)
	mul_cipher := evaluator.Multiply(ciphertext1, ciphertext2)
	sub_cipher := evaluator.Sub(ciphertext1, ciphertext2)

	// decrypt ciphertexts
	decryptor := crypto.NewDecryptor(fv, &key.SecKey)
	new_plaintext1 := decryptor.Decrypt(ciphertext1)
	new_plaintext2 := decryptor.Decrypt(ciphertext2)

	add_plaintext := decryptor.Decrypt(add_cipher)
	mul_plaintext := decryptor.Decrypt(mul_cipher)
	sub_plaintext := decryptor.Decrypt(sub_cipher)

	// decode messages
	new_msg1 := new(bigint.Int)
	new_msg2 := new(bigint.Int)
	add_msg := new(bigint.Int)
	mul_msg := new(bigint.Int)
	sub_msg := new(bigint.Int)

	encoder.Decode(new_msg1, new_plaintext1)
	encoder.Decode(new_msg2, new_plaintext2)
	encoder.Decode(add_msg, add_plaintext)
	encoder.Decode(mul_msg, mul_plaintext)
	encoder.Decode(sub_msg, sub_plaintext)

	fmt.Printf("N=%d, Q=%s\n", fv.N, fv.Q)

	fmt.Printf("a=%d, b=%d\n", a, b)
	fmt.Printf("Cipher (a)=%d, Cipher (b)=%d\n", ciphertext1, ciphertext2)
	//  fmt.Printf("Public key=%s\n Private key=%s\n\n", key.PubKey,key.SecKey)

	fmt.Printf("%v + %v = %v\n", msg1.Int64(), msg2.Int64(), add_msg.Int64())
	fmt.Printf("%v * %v = %v\n", msg1.Int64(), msg2.Int64(), mul_msg.Int64())
	fmt.Printf("%v - %v = %v\n", msg1.Int64(), msg2.Int64(), sub_msg.Int64())

}
