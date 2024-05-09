package wallet

import (
	"blockchain/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
)

// Wallet 钱包结构体
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
	Identity   util.Identity
}

// NewWallet 生成钱包
func NewWallet(identity util.Identity) *Wallet {
	privateKey, publicKey := util.GenNeKeyPair()
	wallet := Wallet{privateKey, publicKey, identity}
	return &wallet
}

// Address 获取钱包地址
func (w *Wallet) Address() []byte {
	publicHash := util.PublicKeyHash(w.PublicKey)
	return util.PublicHashToAddress(publicHash)
}

// SaveWallet 保存钱包
func (w *Wallet) SaveWallet() {
	filename := util.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer

	// 序列化私钥
	privKeyBytes := w.PrivateKey.D.Bytes() // 获取D值的字节表示
	lenPrivBytes := len(privKeyBytes)
	content.Write([]byte{byte(lenPrivBytes)}) // 首先写入D的长度
	content.Write(privKeyBytes)               // 写入D的值

	// 写入公钥
	pubKeyBytes := w.PublicKey
	lenPubBytes := len(pubKeyBytes)
	content.Write([]byte{byte(lenPubBytes)}) // 首先写入D的长度
	content.Write(pubKeyBytes)               // 写入D的值

	//content.Write(w.PublicKey)

	// 写入角色
	content.Write([]byte(w.Identity))

	err := ioutil.WriteFile(filename, content.Bytes(), 0644)
	util.Err(err)
}

// LoadWallet 加载钱包
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
	lenPubKeyBytes := int(fileContent[1+lenPrivBytes])
	pubKeyBytes := fileContent[2+lenPrivBytes : 2+lenPrivBytes+lenPubKeyBytes]
	privKey.PublicKey.X.SetBytes(pubKeyBytes[:len(pubKeyBytes)/2])
	privKey.PublicKey.Y.SetBytes(pubKeyBytes[len(pubKeyBytes)/2:])
	//pubKeyBytes := fileContent[1+lenPrivBytes:]
	//privKey.PublicKey.X.SetBytes(pubKeyBytes[:len(pubKeyBytes)/2])
	//privKey.PublicKey.Y.SetBytes(pubKeyBytes[len(pubKeyBytes)/2:])

	// 读取身份信息
	identity := string(fileContent[2+lenPrivBytes+lenPubKeyBytes:])
	util.Info(fmt.Sprintf("my identity is %s", identity))

	return &Wallet{PrivateKey: privKey, PublicKey: pubKeyBytes, Identity: util.Identity(identity)}
}
