module github.com/lyszhang/go-homomorphic

go 1.12

require (
	github.com/LoCCS/bliss v0.0.0-20180223025823-07585ac9b817 // indirect
	github.com/arnaucube/go-snark v0.0.4
	github.com/dedis/lago v0.0.0-20181016124759-789763ed5fb5
	github.com/fentec-project/bn256 v0.0.0-20190726093940-0d0fc8bfeed0 // indirect
	github.com/fentec-project/gofe v0.0.0-20201027091113-8b2d8e42d985
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lyszhang/go-go-gadget-paillier v0.0.0-20201116065834-16752f36117c
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/renproject/secp256k1 => ../../../github.com/renproject/secp256k1
	github.com/renproject/shamir => ../../../github.com/renproject/shamir
)
