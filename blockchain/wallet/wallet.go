package wallet

import (
	"blockchain/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"math/big"
)

// 钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// 生成密钥对
func GenNeKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	util.Err(err)
	// 公钥为点坐标，需要拼接起来保存
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

// 生成钱包
func NewWallet() *Wallet {
	privateKey, publicKey := GenNeKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

// hash公钥
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hash := ripemd160.New()
	_, err := hash.Write(hashedPublicKey[:])
	util.Err(err)
	publicRipeMd := hash.Sum(nil)
	return publicRipeMd
}

// 检查位生成函数
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:util.ChecksumLength]
}

// base256转base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

// base58转base256
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	util.Err(err)
	return decode
}

// 从公钥哈希生成钱包地址
func PublicHashToAddress(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{util.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

// 从钱包地址转公钥哈希
func AddressToPublicHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-util.ChecksumLength]
	return pubKeyHash
}

// 获取钱包地址
func (w *Wallet) Address() []byte {
	publicHash := PublicKeyHash(w.PublicKey)
	return PublicHashToAddress(publicHash)
}

// 保存钱包
func (w *Wallet) SaveWallet() {
	filename := util.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer

	// 序列化私钥
	privKeyBytes := w.PrivateKey.D.Bytes() // 获取D值的字节表示
	lenPrivBytes := len(privKeyBytes)
	content.Write([]byte{byte(lenPrivBytes)}) // 首先写入D的长度
	content.Write(privKeyBytes)               // 写入D的值

	// 写入公钥
	content.Write(w.PublicKey)

	err := ioutil.WriteFile(filename, content.Bytes(), 0644)
	util.Err(err)
	//filename := util.Wallets + string(w.Address()) + ".wlt"
	//var content bytes.Buffer
	//gob.Register(elliptic.P256())
	//encoder := gob.NewEncoder(&content)
	//err := encoder.Encode(w)
	//util.Err(err)
	//err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	//util.Err(err)
}

// 加载钱包
func LoadWallet(address string) *Wallet {
	filename := util.Wallets + address + ".wlt"
	if !util.FileExists(filename) {
		util.Err(errors.New("该地址无法载入钱包！"))
	}

	fileContent, err := ioutil.ReadFile(filename)
	util.Err(err)

	// 读取私钥
	lenPrivBytes := int(fileContent[0])
	d := new(big.Int).SetBytes(fileContent[1 : 1+lenPrivBytes])
	privKey := ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int),
			Y:     new(big.Int),
		},
	}

	// 读取公钥
	pubKeyBytes := fileContent[1+lenPrivBytes:]
	privKey.PublicKey.X.SetBytes(pubKeyBytes[:len(pubKeyBytes)/2])
	privKey.PublicKey.Y.SetBytes(pubKeyBytes[len(pubKeyBytes)/2:])

	return &Wallet{PrivateKey: privKey, PublicKey: pubKeyBytes}

	//filename := util.Wallets + address + ".wlt"
	//if !util.FileExists(filename) {
	//	util.Err(errors.New("该地址无法载入钱包！"))
	//}
	//var w Wallet
	//gob.Register(elliptic.P256())
	//fileContent, err := ioutil.ReadFile(filename)
	//util.Err(err)
	//decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	//err = decoder.Decode(&w)
	//util.Err(err)
	//return &w
}
