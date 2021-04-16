// https://asecuritysite.com/encryption/go_lattice_cc4
package main

import (
	"fmt"
	"github.com/ldsec/lattigo/ckks"
	"log"
	"math/rand"
)

//const (
//	logN   = 10
//	logQ   = 30
//	levels = 8
//	scale  = logQ
//	sigma  = 3.19
//)

const (
	logN   = 9
	logQ   = 49
	levels = 12
	scale  = 49
	sigma  = 3.19
)

func newCkksContext() *ckks.CkksContext {
	var ckkscontext *ckks.CkksContext
	ckkscontext, _ = ckks.NewCkksContext(logN, logQ, scale, levels, sigma)
	return ckkscontext
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func encryptNumber(context *ckks.CkksContext, encryptor *ckks.Encryptor, f float64) (*ckks.Ciphertext, error) {
	slots := 1 << (logN - 1)
	values := make([]complex128, slots)
	for i := 0; i < slots; i++ {
		values[i] = complex(randomFloat(0.1, 1), 0)
	}
	values[0] = complex(f, 0)
	plaintext1 := context.NewPlaintext(levels-1, scale)
	plaintext1.EncodeComplex(values)

	// Encryption process
	return encryptor.EncryptNew(plaintext1)
}

func main() {
	a := 1.2
	b := 1.0

	// Context
	ckkscontext := newCkksContext()
	kgen := ckkscontext.NewKeyGenerator()

	// Keys
	var sk *ckks.SecretKey
	var pk *ckks.PublicKey
	sk, pk, _ = kgen.NewKeyPair()

	//rlk, _ := kgen.NewRelinKey(sk, 40)

	// Encryptor
	var encryptor *ckks.Encryptor
	encryptor, _ = ckkscontext.NewEncryptor(pk)

	// Decryptor
	var decryptor *ckks.Decryptor
	decryptor, _ = ckkscontext.NewDecryptor(sk)

	// Encryption process
	var ciphertext1, ciphertext2 *ckks.Ciphertext
	ciphertext1, _ = encryptNumber(ckkscontext, encryptor, a)
	ciphertext2, _ = encryptNumber(ckkscontext, encryptor, b)

	// Polynomial values can be viewed with: ciphertext1.Value()[0]
	fmt.Printf("\nCipher (a) pointer: %v\n", ciphertext1)
	fmt.Printf("\nCipher (b) pointer: %v\n", ciphertext2)

	evaluator := ckkscontext.NewEvaluator()

	log.Printf("Generating relinearization keys")

	evaluator.MulRelin(ciphertext1, ciphertext2, nil, ciphertext1)

	plaintext1, _ := decryptor.DecryptNew(ciphertext1)

	valuesTest := plaintext1.DecodeComplex()

	fmt.Printf("\nInput: %f-%f", a, b)

	fmt.Printf("\nCipher (a+b) pointer: %v", *ciphertext1)

	ch := real(valuesTest[0])
	fmt.Printf("\n\nDecrypted: %.2f", ch)

	//ciphertextNew := reciprocal(evaluator, decryptor, ciphertext1, ciphertext2, 5)
	////ciphertextNew, err := evaluator.InverseNew(ciphertext1, 10, rlk)
	////if err != nil {
	////	fmt.Println("err: ", err.Error())
	////}
	////evaluator.MulRelin(ciphertextNew, ciphertextNew, rlk, ciphertextNew)
	//plaintextNew, err := decryptor.DecryptNew(ciphertextNew)
	//if err != nil {
	//	fmt.Println("err: ", err.Error())
	//}
	//
	//valuesTestNew := plaintextNew.DecodeComplex()
	//chNew := real(valuesTestNew[0])
	//fmt.Printf("\n\nDecrypted: %f", chNew)
}

//func reciprocal(evaluator *ckks.Evaluator, decryptor *ckks.Decryptor, c *ckks.Ciphertext, d *ckks.Ciphertext,
//	n int) *ckks.Ciphertext {
//	x := 0.1
//
//	evaluator.MultConst(d, complex(x, 0), d)
//	evaluator.Neg(d, d)
//	evaluator.AddConst(d, complex(2, 0), d)
//	evaluator.MultConst(d, complex(x, 0), d)
//
//	for i := 1; i < n; i++ {
//		evaluator.MulRelin(d, c, rlk, d)
//		evaluator.Neg(d, d)
//		evaluator.AddConst(d, complex(2, 0), d)
//		evaluator.MulRelin(d, complex(x, 0), d)
//
//
//		plaintextNew, err := decryptor.DecryptNew(c)
//		if err != nil {
//			fmt.Println("err: ", err.Error())
//		}
//
//		valuesTestNew := plaintextNew.DecodeComplex()
//		chNew := real(valuesTestNew[0])
//		fmt.Printf("\n\nDecrypted: %f", chNew)
//	}
//
//	return c
//}

func reciprocalPlain(c float64, n int) float64 {
	x := 0.1
	for i := 0; i < n; i++ {
		x = x * (2 - c*x)
		fmt.Println(x)
	}
	return x
}
