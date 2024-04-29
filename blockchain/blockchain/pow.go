package blockchain

import (
	"blockchain/util"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-util.Difficulty))
	return target.Bytes()
}

func (b *Block) GetBase4Nonce(nonce int64) []byte {
	// time和nonce应当先转换为字节切片,否则无法join
	var timeBuf bytes.Buffer
	var nonceBuf bytes.Buffer

	err := binary.Write(&timeBuf, binary.BigEndian, b.Time.Unix())
	if err != nil {
		fmt.Println("Error writing time:", err)
		return nil
	}

	err = binary.Write(&nonceBuf, binary.BigEndian, nonce)
	if err != nil {
		fmt.Println("Error writing nonce:", err)
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
