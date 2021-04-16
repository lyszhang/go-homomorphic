/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/2/25 2:53 PM
 */

package main

import (
	"fmt"
	"github.com/ldsec/lattigo/ckks"
	"math/rand"
	"time"
)

type ckksParams struct {
	params      *ckks.Parameters
	ckkscontext *ckks.Context
	encoder     ckks.Encoder
	kgen        ckks.KeyGenerator
	sk          *ckks.SecretKey
	pk          *ckks.PublicKey
	encryptorPk ckks.Encryptor
	encryptorSk ckks.Encryptor
	decryptor   ckks.Decryptor
	evaluator   ckks.Evaluator
}

type ckksParameters struct {
	verbose    bool
	medianprec float64
	slots      uint64

	ckksParameters []*ckks.Parameters
}

var testParams = new(ckksParameters)

func init() {
	rand.Seed(time.Now().UnixNano())

	testParams.medianprec = 15
	testParams.verbose = false

	testParams.ckksParameters = ckks.DefaultParams
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomComplex(min, max float64) complex128 {
	return complex(randomFloat(min, max), randomFloat(min, max))
}

func genCkksParams(contextParameters *ckks.Parameters) (params *ckksParams) {

	params = new(ckksParams)

	params.params = contextParameters.Copy()

	params.kgen = ckks.NewKeyGenerator(contextParameters)

	params.sk, params.pk = params.kgen.GenKeyPairSparse(128)

	params.encoder = ckks.NewEncoder(contextParameters)

	params.encryptorPk = ckks.NewEncryptorFromPk(contextParameters, params.pk)
	params.encryptorSk = ckks.NewEncryptorFromSk(contextParameters, params.sk)
	params.decryptor = ckks.NewDecryptor(contextParameters, params.sk)

	params.evaluator = ckks.NewEvaluator(contextParameters)

	return
}

func newTestVectors(contextParams *ckksParams, encryptor ckks.Encryptor, a float64) (values []complex128,
	plaintext *ckks.Plaintext, ciphertext *ckks.Ciphertext) {

	slots := uint64(1 << contextParams.params.LogSlots)

	values = make([]complex128, slots)

	for i := uint64(0); i < slots; i++ {
		values[i] = randomComplex(0.1, 1)
	}

	values[0] = complex(a, 0)

	plaintext = ckks.NewPlaintext(contextParams.params, contextParams.params.MaxLevel(), contextParams.params.Scale)

	contextParams.encoder.Encode(plaintext, values, slots)

	if encryptor != nil {
		ciphertext = encryptor.EncryptNew(plaintext)
	}

	return values, plaintext, ciphertext
}

func decode(contextParams *ckksParams, decryptor ckks.Decryptor, element interface{}) float64 {
	var plaintextTest *ckks.Plaintext
	var valuesTest []complex128

	switch element.(type) {
	case *ckks.Ciphertext:
		plaintextTest = decryptor.DecryptNew(element.(*ckks.Ciphertext))
	case *ckks.Plaintext:
		plaintextTest = element.(*ckks.Plaintext)
	}

	slots := uint64(1 << contextParams.params.LogSlots)

	valuesTest = contextParams.encoder.Decode(plaintextTest, slots)

	ch := real(valuesTest[0])
	fmt.Printf("decrypt result: %f\n", ch)
	return ch
}

func add() {
	paramStd := testParams.ckksParameters[ckks.PN14QP438]
	params := genCkksParams(paramStd)

	values1, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 1)
	values2, _, ciphertext2 := newTestVectors(params, params.encryptorSk, 1)

	for i := range values1 {
		values1[i] += values2[i]
	}

	params.evaluator.Add(ciphertext1, ciphertext2, ciphertext1)

	decode(params, params.decryptor, ciphertext1)
}

func mul() {
	paramStd := testParams.ckksParameters[ckks.PN14QP438]
	params := genCkksParams(paramStd)

	values1, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 1)
	values2, _, ciphertext2 := newTestVectors(params, params.encryptorSk, 1)

	for i := range values1 {
		values1[i] += values2[i]
	}

	rlk := params.kgen.GenRelinKey(params.sk)
	params.evaluator.MulRelin(ciphertext1, ciphertext2, rlk, ciphertext1)

	decode(params, params.decryptor, ciphertext1)
}

func neg() {
	paramStd := testParams.ckksParameters[ckks.PN14QP438]
	params := genCkksParams(paramStd)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 0.002348)
	params.evaluator.Neg(ciphertext1, ciphertext1)

	decode(params, params.decryptor, ciphertext1)
}

func inverse() {
	paramStd := testParams.ckksParameters[ckks.PN14QP438]
	params := genCkksParams(paramStd)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 0.5)

	rlk := params.kgen.GenRelinKey(params.sk)
	ciphertextNew := params.evaluator.InverseNew(ciphertext1, 7, rlk)

	decode(params, params.decryptor, ciphertextNew)
}

func inverseNew() {
	paramStd := testParams.ckksParameters[ckks.PN14QP438]
	params := genCkksParams(paramStd)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 2)

	rlk := params.kgen.GenRelinKey(params.sk)
	ciphertextNew := reciprocal(params, params.evaluator, rlk, ciphertext1, 2)

	decode(params, params.decryptor, ciphertextNew)
}

func reciprocal(contextParameters *ckksParams, evaluator ckks.Evaluator, rlk *ckks.EvaluationKey,
	c *ckks.Ciphertext, n int) *ckks.Ciphertext {
	x := 0.1

	var m *ckks.Ciphertext
	for i := 0; i < n; i++ {
		if i == 0 {
			m = evaluator.MultByConstNew(c, complex(x, 0))
			evaluator.Neg(m, m)
			evaluator.AddConst(m, complex(2, 0), m)
			evaluator.MultByConst(m, complex(x, 0), m)
			fmt.Println("0")

			decode(contextParameters, contextParameters.decryptor, m)
		} else {
			decode(contextParameters, contextParameters.decryptor, c)
			decode(contextParameters, contextParameters.decryptor, m)
			fmt.Println("====")
			n := evaluator.MulRelinNew(c, m, rlk)
			decode(contextParameters, contextParameters.decryptor, n)
			evaluator.Sub(n, n, n)
			decode(contextParameters, contextParameters.decryptor, n)
			evaluator.AddConst(n, complex(2, 0), n)
			decode(contextParameters, contextParameters.decryptor, n)
			evaluator.MulRelin(n, m, rlk, m)
			decode(contextParameters, contextParameters.decryptor, m)
		}
	}

	return m
}

func main() {
	inverse()
}
