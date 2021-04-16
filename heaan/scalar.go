/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/2/25 2:53 PM
 */

package main

import (
	"fmt"
	"github.com/ldsec/lattigo/ckks"
)

func m() {

	var logN, logQ, levels, scale uint64

	// Scheme params
	logN = 1
	logQ = 30
	levels = 8
	scale = logQ
	sigma := 3.19

	a := 6.0
	b := 3.0

	// Context
	var ckkscontext *ckks.CkksContext
	ckkscontext, _ = ckks.NewCkksContext(logN, logQ, scale, levels, sigma)

	kgen := ckkscontext.NewKeyGenerator()

	// Keys
	var sk *ckks.SecretKey
	var pk *ckks.PublicKey
	sk, pk, _ = kgen.NewKeyPair()

	// Encryptor
	var encryptor *ckks.Encryptor
	encryptor, _ = ckkscontext.NewEncryptor(pk)

	// Decryptor
	var decryptor *ckks.Decryptor
	decryptor, _ = ckkscontext.NewDecryptor(sk)

	// Values to encrypt
	values := make([]complex128, 1<<(logN-1))

	values[0] = complex(a, 0)

	fmt.Printf("HEAAN parameters : logN = %d, logQ = %d, levels = %d (%d bits), logPrecision = %d, logScale = %d, sigma = %f \n", logN, logQ, levels, 60+(levels-1)*logQ, ckkscontext.Precision(), scale, sigma)

	// Plaintext creation and encoding process
	plaintext := ckkscontext.NewPlaintext(levels-1, scale)

	plaintext.EncodeComplex(values)

	// Encryption process
	var ciphertext *ckks.Ciphertext
	ciphertext, _ = encryptor.EncryptNew(plaintext)

	evaluator := ckkscontext.NewEvaluator()
	evaluator.MultConst(ciphertext, complex(b, 0), ciphertext)

	// Decryption process
	plaintext, _ = decryptor.DecryptNew(ciphertext)

	// Decoding process
	valuesTest := plaintext.DecodeComplex()

	fmt.Printf("\na: %f, b (multiplier): %f\n", a, b)
	fmt.Printf("\nCipher: %v\n\nDecrypted: ", ciphertext)

	ch := real(valuesTest[0])
	fmt.Printf("a times b is %.2f", ch)

}
