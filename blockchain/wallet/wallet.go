package wallet

import (
	"blockchain/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"io/ioutil"
	"math/big"
)

// 钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// 生成钱包
func NewWallet() *Wallet {
	privateKey, publicKey := util.GenNeKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

// 获取钱包地址
func (w *Wallet) Address() []byte {
	publicHash := util.PublicKeyHash(w.PublicKey)
	return util.PublicHashToAddress(publicHash)
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
