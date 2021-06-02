/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2021/4/14 2:45 PM
 */

package psi

import (
	"fmt"
	paillier "github.com/lyszhang/go-go-gadget-paillier"
	"github.com/lyszhang/go-homomorphic/psi/utils"
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

const LengthLimit = 3

func NewPolysFromSet(raw set) []Vector {
	ss := utils.SplitArray(raw, LengthLimit)
	var vectors []Vector
	for _, s := range ss {
		vectors = append(vectors, NewPolyFromSet(s))
	}
	return vectors
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
func (v *Vector) Encrypt(pubKey *paillier.PublicKey) *EncVector {
	tmpEnc := make([][]byte, 0)
	for _, value := range v.Data {
		mValue := new(big.Int).SetInt64(value)

		///TODO: big int负数问题
		eValue, err := paillier.Encrypt(pubKey, mValue.Bytes())
		if err != nil {
			fmt.Println(err)
			panic(nil)
		}
		tmpEnc = append(tmpEnc, eValue)
	}
	return &EncVector{
		PubKey:    *pubKey,
		Encrypted: tmpEnc,
	}
}

func EncryptVectors(ss []Vector, pubKey *paillier.PublicKey) []*EncVector {
	var encVectors []*EncVector
	for _, value := range ss {
		encVectors = append(encVectors, value.Encrypt(pubKey))
	}
	return encVectors
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

func DecryptVectors(ss []*EncVector, privKey *paillier.PrivateKey) []Vector {
	var vectors []Vector
	for _, value := range ss {
		vectors = append(vectors, *value.Decrypt(privKey))
	}
	return vectors
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
func (g *EncVector) Add(f *EncVector, pubKey *paillier.PublicKey) {
	lenG := len(g.Encrypted)
	lenF := len(f.Encrypted)

	if lenF <= lenG {
		for i := 0; i < lenG; i++ {
			if i < lenF {
				g.Encrypted[i] = paillier.AddCipher(pubKey, g.Encrypted[i], f.Encrypted[i])
			}
		}
		return
	} else {
		for i := 0; i < lenF; i++ {
			if i < lenG {
				g.Encrypted[i] = paillier.AddCipher(pubKey, g.Encrypted[i], f.Encrypted[i])
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
func (g *EncVector) Mul(f *Vector, pubkey *paillier.PublicKey) {
	limit := len(g.Encrypted) + len(f.Data) - 1
	// TODO: limit checker
	var encs [][]byte
	for i := 0; i < limit; i++ {
		cipherTmp, _ := paillier.Encrypt(pubkey, big.NewInt(0).Bytes())
		for j := 0; j <= i; j++ {
			fBigBytes := big.NewInt(f.Value(j)).Bytes()

			mulCipher := paillier.Mul(pubkey, g.Value(i-j), fBigBytes)
			cipherTmp = paillier.AddCipher(pubkey, mulCipher, cipherTmp)
		}
		encs = append(encs, cipherTmp)
	}
	g.Encrypted = encs
	return
}
