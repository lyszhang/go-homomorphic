/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/14 2:45 PM
 */

package psi

import (
	"fmt"
	paillier "github.com/lyszhang/go-go-gadget-paillier"
	"math/big"
)

type Vector struct {
	Data []int64
}

var RandVector = Vector{
	Data: []int64{1, 1},
}

// 中间迭代
// (g[0]x^n + g[1]x^n-1 + ... + g[n]x^0)(x+g[n+1])
func IterationN(v Vector, item int64) Vector {
	tmp := []int64{0}
	for _, value := range v.Data {
		tmp = append(tmp, value)
	}
	for i, _ := range v.Data {
		tmp[i] = tmp[i] + v.Data[i]*item
	}
	return Vector{Data: tmp}
}

// 构造vector
// (x+s1)(x+s2)...(x+sn)展开式的系数，构成vector
func NewPolyFromSet(s set) Vector {
	length := len(s)
	if length == 1 {
		return Vector{
			Data: []int64{s[0], 1},
		}
	}
	return IterationN(NewPolyFromSet(s[:length-1]), s[length-1])
}

// string
func (v *Vector) Print() {
	fmt.Println(v)
}

// index
func (v *Vector) Value(i int) int64 {
	if i < len(v.Data) {
		return v.Data[i]
	}
	return 0
}

// 加密
func (v *Vector) Encrypt(privKey *paillier.PrivateKey) *EncVector {
	tmpEnc := make([][]byte, 0)
	for _, value := range v.Data {
		mValue := new(big.Int).SetInt64(value)

		///TODO: big int负数问题
		eValue, err := paillier.Encrypt(&privKey.PublicKey, mValue.Bytes())
		if err != nil {
			fmt.Println(err)
			panic(nil)
		}
		tmpEnc = append(tmpEnc, eValue)
	}
	return &EncVector{
		PubKey:    privKey.PublicKey,
		Encrypted: tmpEnc,
	}
}

type EncVector struct {
	PubKey    paillier.PublicKey
	Encrypted [][]byte
}

// 解密
func (v *EncVector) Decrypt(privKey *paillier.PrivateKey) *Vector {
	// Decrypt.
	tmpValue := make([]int64, 0)
	for _, value := range v.Encrypted {
		d, err := paillier.Decrypt(privKey, value)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		plainText := new(big.Int).SetBytes(d)
		tmpValue = append(tmpValue, plainText.Int64())
	}
	return &Vector{Data: tmpValue}
}

// index
func (v *EncVector) Value(i int) []byte {
	if i < len(v.Encrypted) {
		return v.Encrypted[i]
	}
	tmp, _ := paillier.Encrypt(&v.PubKey, big.NewInt(0).Bytes())
	return tmp
}

// 加法
// E(h) = E(f+g)
// f,g皆为多项式密文
func (g *EncVector) Add(f *EncVector, privKey *paillier.PrivateKey) {
	lenG := len(g.Encrypted)
	lenF := len(f.Encrypted)

	if lenF <= lenG {
		for i := 0; i < lenG; i++ {
			if i < lenF {
				g.Encrypted[i] = paillier.AddCipher(&privKey.PublicKey, g.Encrypted[i], f.Encrypted[i])
			}
		}
		return
	} else {
		for i := 0; i < lenF; i++ {
			if i < lenG {
				g.Encrypted[i] = paillier.AddCipher(&privKey.PublicKey, g.Encrypted[i], f.Encrypted[i])
			} else {
				g.Encrypted = append(g.Encrypted, f.Encrypted[i])
			}
		}
		return
	}
}

// 乘法
// E(h) = E(f*g)
// g为多项式密文, f为多项式明文
func (g *EncVector) Mul(f *Vector, privKey *paillier.PrivateKey) {
	limit := len(g.Encrypted) + len(f.Data) - 1
	// TODO: limit checker
	var encs [][]byte
	for i := 0; i < limit; i++ {
		cipherTmp, _ := paillier.Encrypt(&privKey.PublicKey, big.NewInt(0).Bytes())
		for j := 0; j <= i; j++ {
			fBigBytes := big.NewInt(f.Value(j)).Bytes()

			mulCipher := paillier.Mul(&privKey.PublicKey, g.Value(i-j), fBigBytes)
			cipherTmp = paillier.AddCipher(&privKey.PublicKey, mulCipher, cipherTmp)

			//plain, _ := paillier.Decrypt(privKey, cipherTmp)
			//plainText := new(big.Int).SetBytes(plain)
			//
			//plaing, _ := paillier.Decrypt(privKey, g.Value(i-j))
			//plainTextg := new(big.Int).SetBytes(plaing)
			//fmt.Println("f: ", f.Value(j))
			//fmt.Println("g ", plainTextg.Int64())
			//
			//fmt.Printf("i: %d, j: %d, cipher: %d\n", i, j, plainText.Int64())
		}
		encs = append(encs, cipherTmp)
	}
	g.Encrypted = encs
	return
}
