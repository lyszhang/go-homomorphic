/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/11 2:18 PM
 */

package main

import (
	"fmt"
	snark "github.com/arnaucube/go-snark"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"github.com/arnaucube/go-snark/fields"
	"github.com/arnaucube/go-snark/r1csqap"
	"math/big"
	"regexp"
	"strings"
)

func TrimSpaceNewlineInString(s string) string {
	re := regexp.MustCompile(`LF`)
	return re.ReplaceAllString(s, "\n")
}

func main() {
	str := "func exp3(private a):LF\tb = a * aLF\tc = a * bLF\treturn cLFLFfunc main(private s0, public s1):LF\ts3 = exp3(s0)LF\ts4 = s3 + s0LF\ts5 = s4 + 5LF\tequals(s1, s5)LF\tout = 1 * 1"
	x := 3
	y := 35

	fmt.Printf("x (Private input): %d, y (Public input): %d\n", x, y)

	str = TrimSpaceNewlineInString(str)
	fmt.Printf("Flat: %s\n", str)

	parser := circuitcompiler.NewParser(strings.NewReader(str))
	circuit, _ := parser.Parse()

	fmt.Println("\nThe circuit is defined as:")
	fmt.Println(circuit)

	val1 := big.NewInt(int64(x))
	privateVal := []*big.Int{val1}
	val2 := big.NewInt(int64(y))
	publicVal := []*big.Int{val2}

	// witness
	w, _ := circuit.CalculateWitness(privateVal, publicVal)

	fmt.Println("\nWitness: ", w)

	fmt.Println("\nR1CS from flat code")
	a, b, c := circuit.GenerateR1CS()
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)
	fmt.Printf("c: %d\n\n", c)

	r, _ := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

	f := fields.NewFq(r)
	// new Polynomial Field
	pf := r1csqap.NewPolynomialField(f)
	alphas, betas, gammas, _ := pf.R1CSToQAP(a, b, c)

	_, _, _, px := pf.CombinePolynomials(w, alphas, betas, gammas)

	setup, _ := snark.GenerateTrustedSetup(len(w), *circuit, alphas, betas, gammas)

	proof, _ := snark.GenerateProofs(*circuit, setup.Pk, w, px)

	yVerif := big.NewInt(int64(y))
	publicSignalsVerif := []*big.Int{yVerif}

	rtn := snark.VerifyProof(setup.Vk, proof, publicSignalsVerif, true)
	if rtn == true {
		fmt.Printf("Valid proofs!!!")
	}
}
