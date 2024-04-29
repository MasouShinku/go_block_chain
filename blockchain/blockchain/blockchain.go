package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

// Block 区块结构体
type Block struct {
	Time     time.Time // 时间戳
	Hash     []byte    // 哈希值
	PrevHash []byte    // 上一个区块的哈希值
	Target   []byte    // 难度值
	Nonce    int64     // 是否进行获取
	Data     []byte    // 数据属性
}

// BlockChain 区块链
type BlockChain struct {
	BlockList []*Block
}

// 构建区块哈希值
func (b *Block) SetHash() {
	// 创建一个buf用于存放要生成哈希的数据
	var buf bytes.Buffer
	timeErr := binary.Write(&buf, binary.BigEndian, b.Time.Unix())
	if timeErr != nil {
		return
	}
	buf.Write(b.PrevHash)
	buf.Write(b.Target)
	nonceErr := binary.Write(&buf, binary.BigEndian, b.Nonce)
	if nonceErr != nil {
		return
	}
	buf.Write(b.Data)
	// 计算哈希值
	hash := sha256.Sum256(buf.Bytes())
	b.Hash = hash[:]
}

// CreateBlock 创建区块
func CreateBlock(prevHash, data []byte) *Block {
	block := Block{time.Now(), []byte{}, prevHash, []byte{}, 0, data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 通过创世区块生成区块链
func InitBlockChain() *BlockChain {
	blockChain := BlockChain{}
	data := "BlockChain,Start!"
	block := CreateBlock([]byte{}, []byte(data))
	blockChain.BlockList = append(blockChain.BlockList, block)
	return &blockChain
}

func (blockChain *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(blockChain.BlockList[len(blockChain.BlockList)-1].Hash, []byte(data))
	blockChain.BlockList = append(blockChain.BlockList, newBlock)
}
