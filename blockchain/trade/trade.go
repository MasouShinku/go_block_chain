package trade

import (
	"blockchain/util"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// 首先定义转入转出结构体
type TradeIn struct {
	TradeID     []byte // 订单标识
	OutID       int    // 订单的第几个Output
	FromAddress []byte //转出者地址
	//OutName     string
}

type TradeOut struct {
	Num       int    // 转出值
	ToAddress []byte // 接受者地址
	//InName    string
}

// 交易结构体
type Trade struct {
	ID      []byte // ID，哈希值表示
	Inputs  []TradeIn
	Outputs []TradeOut
}

// 计算交易哈希值
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

// 创建初始订单，将InitNum商品转入用户
func FirstTrade(toaddress []byte) *Trade {
	In := TradeIn{[]byte{}, -1, []byte{}}
	Out := TradeOut{util.InitNum, toaddress}
	tx := Trade{[]byte("The First Trade!"), []TradeIn{In}, []TradeOut{Out}}
	return &tx
}

// 判断是否为初始订单
func (t *Trade) IsFirstTrade() bool {
	return len(t.Inputs) == 1 && t.Inputs[0].OutID == -1
}

// 判断源地址是否正确
func (in *TradeIn) IsFromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

// 判断目标地址是否正确
func (out *TradeOut) IsToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
