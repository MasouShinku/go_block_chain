package trade

import (
	"blockchain/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"math/big"
)

// TradeIn 首先定义转入转出结构体
type TradeIn struct {
	TradeID   []byte // 订单标识
	OutID     int    // 订单的第几个Output
	PublicKey []byte // 公钥
	Sign      []byte // 签名
}

type TradeOut struct {
	Num           int    // 转出值
	HashPublicKey []byte // 公钥哈希
}

// Trade 交易结构体
type Trade struct {
	ID          []byte // ID，哈希值表示
	Inputs      []TradeIn
	Outputs     []TradeOut
	Description string
}

// GetTradeHash 计算交易哈希值
func (t *Trade) GetTradeHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte
	// 用gob序列化一下结构体
	encoder := gob.NewEncoder(&encoded)
	if err := encoder.Encode(t); err != nil {
		return nil
	}
	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (t *Trade) SetID() {
	t.ID = t.GetTradeHash()
}

// FirstTrade 创建初始订单，将InitNum商品转入用户
func FirstTrade(toaddress []byte) *Trade {
	In := TradeIn{[]byte{}, -1, []byte{}, nil}
	Out := TradeOut{util.InitNum, toaddress}
	tx := Trade{[]byte("The First Trade!"), []TradeIn{In}, []TradeOut{Out}, "first trade"}
	return &tx
}

// IsFirstTrade 判断是否为初始订单
func (t *Trade) IsFirstTrade() bool {
	return len(t.Inputs) == 1 && t.Inputs[0].OutID == -1
}

// IsFromAddressRight 判断源地址是否正确
func (in *TradeIn) IsFromAddressRight(address []byte) bool {
	return bytes.Equal(in.PublicKey, address)
}

// IsToAddressRight 判断目标地址是否正确
func (out *TradeOut) IsToAddressRight(address []byte) bool {
	return bytes.Equal(out.HashPublicKey, util.PublicKeyHash(address))
}

// Sign 构造签名
func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	util.Err(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

// Verify 验证签名是否合法
func Verify(msg []byte, pubkey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}

// PlainCopy 描述交易信息
func (t *Trade) PlainCopy() Trade {
	var inputs []TradeIn
	var outputs []TradeOut

	for _, tin := range t.Inputs {
		inputs = append(inputs, TradeIn{tin.TradeID, tin.OutID, nil, nil})
	}

	for _, tout := range t.Outputs {
		outputs = append(outputs, TradeOut{tout.Num, tout.HashPublicKey})
	}

	tradeCopy := Trade{t.ID, inputs, outputs, t.Description}

	return tradeCopy
}

// Sign 对交易信息进行签名
func (t *Trade) Sign(privKey ecdsa.PrivateKey) {
	if t.IsFirstTrade() {
		return
	}
	for i, input := range t.Inputs {

		tradeCopy := t.PlainCopy()
		tradeCopy.Inputs[i].PublicKey = input.PublicKey
		plainhash := tradeCopy.GetTradeHash()
		signature := Sign(plainhash, privKey)
		t.Inputs[i].Sign = signature
	}
}

// Verify 验证整个交易是否合法
func (t *Trade) Verify() bool {
	// 使用ECDSA算法的公钥验证签名
	for i, input := range t.Inputs {

		tradeCopy := t.PlainCopy()
		tradeCopy.Inputs[i].PublicKey = input.PublicKey
		plainhash := tradeCopy.GetTradeHash()

		if !Verify(plainhash, input.PublicKey, input.Sign) {
			return false
		}
	}
	return true
}
