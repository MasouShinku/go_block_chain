package blockchain

import (
	"blockchain/util"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"math/big"
)

// GetTarget 获取区块的目标值
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-util.Difficulty))
	return target.Bytes()
}

// GetBase4Nonce 根据给定nonce生成基础数据用于计算哈希
func (b *Block) GetBase4Nonce(nonce int64) []byte {
	// time和nonce应当先转换为字节切片,否则无法join
	var timeBuf bytes.Buffer
	var nonceBuf bytes.Buffer

	err := binary.Write(&timeBuf, binary.BigEndian, b.Time.Unix())
	if err != nil {
		util.Err(errors.New("时间写入失败"))
		return nil
	}

	err = binary.Write(&nonceBuf, binary.BigEndian, nonce)
	if err != nil {
		util.Err(errors.New("nonce写入失败"))
		return nil
	}

	buf := []byte{}
	tradeIDs := make([][]byte, 0)
	for _, t := range b.TradeList {
		tradeIDs = append(tradeIDs, t.ID)
	}
	for _, tradeID := range tradeIDs {
		buf = append(buf, tradeID...)
	}

	data := bytes.Join([][]byte{
		timeBuf.Bytes(),
		b.PrevHash,
		nonceBuf.Bytes(),
		b.Target,
		buf,
	},
		[]byte{},
	)
	return data
}

// FindNonce 寻找有效nonce值
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0
	intTarget.SetBytes(b.Target)

	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

// ValidatePoW 验证PoW
func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
