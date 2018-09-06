package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := GetKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}
func GetKeyPair() ([]byte, []byte) {
	cure := elliptic.P256()
	private, err := ecdsa.GenerateKey(cure, rand.Reader)
	if err != nil {
		fmt.Println(err.Error())
	}
	d := private.D.Bytes()
	b := make([]byte, 0, 32)
	prikey := paddedbytes(32, b, d)
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return prikey, pubkey
}

func paddedbytes(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func (w Wallet) GetAddress() (address string) {
	pubkey := w.PublicKey
	pubhash := sha256.New()
	pubhash.Reset()
	pubhash.Write(pubkey)
	hash1 := pubhash.Sum(nil)
	fmt.Printf("hash1:%x\n=========================\n", hash1)
	riphash := ripemd160.New()
	riphash.Reset()
	riphash.Write(hash1)
	hash2 := riphash.Sum(nil)
	fmt.Printf("hash2 is:%x\n===============================\n", hash2)
	address = mybase58decode(0x00, hash2)
	return address
}

func mybase58decode(ver uint8, hash []byte) string {
	hash = append([]byte{ver}, hash...)
	fmt.Printf("verhash is:%x\n====================================\n", hash)
	onehash := sha256.New()
	onehash.Reset()
	onehash.Write(hash)
	hash1 := onehash.Sum(nil)
	fmt.Printf("hash1 is:%x\n=================================\n", hash1)
	onehash.Reset()
	onehash.Write(hash1)
	hash2 := onehash.Sum(nil)
	fmt.Printf("hash2 is:%x\n===============================\n", hash2)
	boot := hash2[0:4]
	hash = append(hash, boot...)
	fmt.Printf("last hash is:%x", hash)
	address := Base58Encode(hash)
	//return "1111"
	fmt.Printf("address is:%s", string(address))
	return string(address)
}

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	//ReverseBytes(result)

	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)

		} else {
			break
		}
	}

	return result

}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for _, b := range input {
		if b != b58Alphabet[0] {
			break
		}

		zeroBytes++
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(int64(len(b58Alphabet))))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}
