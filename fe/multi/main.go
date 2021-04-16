/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/25 3:24 PM
 */

package main

import (
	"fmt"
	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/simple"
	"github.com/fentec-project/gofe/sample"
	"math/big"
)

func main() {
	numClients := 2           // number of encryptors
	l := 3                    // length of input vectors
	bound := big.NewInt(1000) // upper bound for input vectors

	// Simulate collection of input data.
	// X and Y represent matrices of input vectors, where X are collected
	// from numClients encryptors (omitted), and Y is only known by a single decryptor.
	// Encryptor i only knows its own input vector X[i].
	sampler := sample.NewUniform(bound)
	X, _ := data.NewRandomMatrix(numClients, l, sampler)
	Y, _ := data.NewRandomMatrix(numClients, l, sampler)

	// Trusted entity instantiates scheme instance and generates
	// master keys for all the encryptors. It also derives the FE
	// key derivedKey for the decryptor.
	modulusLength := 2048
	multiDDH, _ := simple.NewDDHMultiPrecomp(numClients, l, modulusLength, bound)
	pubKey, secKey, _ := multiDDH.GenerateMasterKeys()
	derivedKey, _ := multiDDH.DeriveKey(secKey, Y)

	// Different encryptors may reside on different machines.
	// We simulate this with the for loop below, where numClients
	// encryptors are generated.
	encryptors := make([]*simple.DDHMultiClient, numClients)
	for i := 0; i < numClients; i++ {
		encryptors[i] = simple.NewDDHMultiClient(multiDDH.Params)
	}
	// Each encryptor encrypts its own input vector X[i] with the
	// keys given to it by the trusted entity.
	ciphers := make([]data.Vector, numClients)
	for i := 0; i < numClients; i++ {
		cipher, _ := encryptors[i].Encrypt(X[i], pubKey[i], secKey.OtpKey[i])
		ciphers[i] = cipher
	}

	// Ciphers are collected by decryptor, who then computes
	// inner product over vectors from all encryptors.
	decryptor := simple.NewDDHMultiFromParams(numClients, multiDDH.Params)
	prod, _ := decryptor.Decrypt(ciphers, derivedKey, Y)

	fmt.Printf("prod: %d\n", prod)
}
