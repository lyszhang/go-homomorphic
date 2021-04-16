/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/2/25 2:53 PM
 */

package main

import (
	"fmt"
	"math"

	"github.com/ldsec/lattigo/v2/ckks"
	"github.com/ldsec/lattigo/v2/utils"
)

type ckksParams struct {
	params      *ckks.Parameters
	encoder     ckks.Encoder
	kgen        ckks.KeyGenerator
	sk          *ckks.SecretKey
	pk          *ckks.PublicKey
	rlk         *ckks.EvaluationKey
	encryptorPk ckks.Encryptor
	encryptorSk ckks.Encryptor
	decryptor   ckks.Decryptor
	evaluator   ckks.Evaluator
}

func genTestParams(defaultParam *ckks.Parameters, hw uint64) (testContext *ckksParams, err error) {

	testContext = new(ckksParams)

	testContext.params = defaultParam.Copy()

	testContext.kgen = ckks.NewKeyGenerator(testContext.params)

	if hw == 0 {
		testContext.sk, testContext.pk = testContext.kgen.GenKeyPair()
	} else {
		testContext.sk, testContext.pk = testContext.kgen.GenKeyPairSparse(hw)
	}

	testContext.rlk = testContext.kgen.GenRelinKey(testContext.sk)

	testContext.encoder = ckks.NewEncoder(testContext.params)

	testContext.encryptorPk = ckks.NewEncryptorFromPk(testContext.params, testContext.pk)
	testContext.encryptorSk = ckks.NewEncryptorFromSk(testContext.params, testContext.sk)
	testContext.decryptor = ckks.NewDecryptor(testContext.params, testContext.sk)

	testContext.evaluator = ckks.NewEvaluator(testContext.params)

	return testContext, nil

}

func newTestVectors(testContext *ckksParams, encryptor ckks.Encryptor, num float64) (values []complex128,
	plaintext *ckks.Plaintext, ciphertext *ckks.Ciphertext) {

	logSlots := testContext.params.LogSlots()

	values = make([]complex128, 1<<logSlots)

	//a := complex(-1, -1)
	//b := complex(1, 1)
	a := complex(0.1, 0)
	b := complex(1, 0)

	for i := uint64(0); i < 1<<logSlots; i++ {
		values[i] = complex(utils.RandFloat64(real(a), real(b)), utils.RandFloat64(imag(a), imag(b)))
	}

	values[0] = complex(num, 0)

	plaintext = testContext.encoder.EncodeNTTAtLvlNew(testContext.params.MaxLevel(), values, logSlots)

	if encryptor != nil {
		ciphertext = encryptor.EncryptNew(plaintext)
	}

	return values, plaintext, ciphertext
}

func decode(contextParams *ckksParams, encoder ckks.Encoder, decryptor ckks.Decryptor, element interface{}) float64 {
	var valuesTest []complex128

	logSlots := contextParams.params.LogSlots()
	//slots := uint64(1 << logSlots)

	switch element := element.(type) {
	case *ckks.Ciphertext:
		valuesTest = encoder.Decode(decryptor.DecryptNew(element), logSlots)
	case *ckks.Plaintext:
		valuesTest = encoder.Decode(element, logSlots)
	case []complex128:
		valuesTest = element
	}
	ch := real(valuesTest[0])
	fmt.Printf("decrypt result: %f\n", ch)
	return ch
}

func add() {
	paramStd := ckks.DefaultParams[ckks.PN14QP438]
	params, _ := genTestParams(paramStd, 0)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 4.0)
	_, _, ciphertext2 := newTestVectors(params, params.encryptorSk, 2.0)

	params.evaluator.Add(ciphertext1, ciphertext2, ciphertext1)
	decode(params, params.encoder, params.decryptor, ciphertext1)
}

func mul() {
	paramStd := ckks.DefaultParams[ckks.PN14QP438]
	params, _ := genTestParams(paramStd, 0)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 4.0)
	_, _, ciphertext2 := newTestVectors(params, params.encryptorSk, 2.0)

	params.evaluator.MulRelin(ciphertext1, ciphertext2, params.rlk, ciphertext1)
	decode(params, params.encoder, params.decryptor, ciphertext1)
}

func inverse() {
	paramStd := ckks.DefaultParams[ckks.PN16QP1761]
	params, _ := genTestParams(paramStd, 0)
	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 10000)
	//params.evaluator.MultByConst(ciphertext1, 0.0001, ciphertext1)
	//decode(params, params.encoder, params.decryptor, ciphertext1)

	n := uint64(7)
	ciphertext1 = params.evaluator.InverseNew(ciphertext1, n, params.rlk)
	fmt.Println(ciphertext1)
	//decode(params, params.encoder, params.decryptor, ciphertext1)
	//params.evaluator.MultByConst(ciphertext1, 10000, ciphertext1)
	decode(params, params.encoder, params.decryptor, ciphertext1)
}

func inverseNew() {
	paramStd := ckks.DefaultParams[ckks.PN16QP1761]
	params, _ := genTestParams(paramStd, 0)

	_, _, ciphertext1 := newTestVectors(params, params.encryptorSk, 12)
	ciphertextNew := reciprocal(params, params.evaluator, params.rlk, ciphertext1, 8)

	decode(params, params.encoder, params.decryptor, ciphertextNew)
}

var (
	rawList = []float64{1, 100, 10000, 1000000, 0.01, 0.0001, 0.000001}
)

func reciprocal(contextParameters *ckksParams, evaluator ckks.Evaluator, rlk *ckks.EvaluationKey,
	c *ckks.Ciphertext, n int) *ckks.Ciphertext {
	x := 0.01

	var m *ckks.Ciphertext
	for i := 0; i < n; i++ {
		if i == 0 {
			m = evaluator.MultByConstNew(c, complex(x, 0))
			evaluator.Neg(m, m)
			evaluator.AddConst(m, complex(2, 0), m)
			evaluator.MultByConst(m, complex(x, 0), m)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, m)
		} else {
			fmt.Println("====")
			x := evaluator.MultByConstNew(m, complex(2, 0))
			fmt.Println("x")
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, x)
			y := evaluator.PowerNew(m, 2, rlk)
			fmt.Println("y")
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, y)
			y = evaluator.MulRelinNew(c, y, rlk)
			fmt.Println("y")
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, y)
			m = evaluator.SubNew(x, y)
			fmt.Println("m")

			//m = evaluator.DropLevelNew(m, 1)

			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, m)
		}
	}
	return m
}

func inverSqrt(contextParameters *ckksParams, evaluator ckks.Evaluator, rlk *ckks.EvaluationKey,
	c *ckks.Ciphertext, n int) *ckks.Ciphertext {

	var chalf, m *ckks.Ciphertext
	chalf = evaluator.MultByConstNew(c, 0.5)
	m = evaluator.MultByConstNew(c, -0.5)
	m = evaluator.AddConstNew(m, 0x5f3759df)

	fmt.Println("====")
	x := evaluator.MultByConstNew(m, 1.5)
	fmt.Println("1,5fx")
	decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, x)
	y := evaluator.PowerNew(m, 3, rlk)
	fmt.Println("x^3")
	decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, y)
	y = evaluator.MulRelinNew(chalf, y, rlk)
	fmt.Println("xhalf*x^3")
	decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, y)

	m = evaluator.SubNew(x, y)
	fmt.Println("m")
	decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, m)

	return m
}

func reciprocalNew(contextParameters *ckksParams, evaluator ckks.Evaluator, rlk *ckks.EvaluationKey,
	c *ckks.Ciphertext, n int) *ckks.Ciphertext {
	x := 0.1

	var m *ckks.Ciphertext
	for i := 0; i < n; i++ {
		if i == 0 {
			m = evaluator.MultByConstNew(c, complex(x, 0))
			evaluator.Neg(m, m)
			evaluator.AddConst(m, complex(2, 0), m)
			evaluator.MultByConst(m, complex(x, 0), m)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, m)
		} else {
			fmt.Println("====")
			n := evaluator.MulRelinNew(c, m, rlk)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, n)
			evaluator.Neg(n, n)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, n)
			evaluator.AddConst(n, complex(2, 0), n)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, n)
			m = evaluator.MulRelinNew(n, m, rlk)
			decode(contextParameters, contextParameters.encoder, contextParameters.decryptor, m)
		}
	}
	return m
}

func inverse2() {
	paramStd := ckks.DefaultParams[ckks.PN13QP218]
	params, _ := genTestParams(paramStd, 0)
	_, _, ciphertextRaw := newTestVectors(params, params.encryptorSk, 3)

	// 4x-6cx^2+4c^2x^3-c^3x^4
	x := 0.001
	//6cx^2
	ciphertextTmp := params.evaluator.MultByConstNew(ciphertextRaw, 6*math.Pow(x, 2))
	// -6cx^2
	params.evaluator.Neg(ciphertextTmp, ciphertextTmp)
	// 4x-6cx^2
	params.evaluator.AddConst(ciphertextTmp, 4*x, ciphertextTmp)
	// 4x-6cx^2+4c^2x^3
	// 4c^2x^3
	ciphertext2 := params.evaluator.PowerNew(ciphertextRaw, 2, params.rlk)
	params.evaluator.MultByConst(ciphertext2, 4*math.Pow(x, 3), ciphertext2)
	params.evaluator.Add(ciphertextTmp, ciphertext2, ciphertextTmp)
	// 4x-6cx^2+4c^2x^3-c^3x^4
	// -c^3x^4
	ciphertext3 := params.evaluator.PowerNew(ciphertextRaw, 3, params.rlk)
	params.evaluator.MultByConst(ciphertext3, math.Pow(x, 4), ciphertext3)
	params.evaluator.Neg(ciphertext3, ciphertext3)
	params.evaluator.Add(ciphertextTmp, ciphertext3, ciphertextTmp)

	decode(params, params.encoder, params.decryptor, ciphertextTmp)

}

func main() {
	inverseNew()
}
