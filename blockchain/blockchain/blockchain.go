package blockchain

import (
	"blockchain/trade"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// Block 区块结构体
type Block struct {
	Time      time.Time // 时间戳
	Hash      []byte    // 哈希值
	PrevHash  []byte    // 上一个区块的哈希值
	Target    []byte    // 难度值
	Nonce     int64     // 是否进行获取
	TradeList []*trade.Trade
	//Data      []byte // 数据属性
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
	//buf.Write(b.Data)
	// 还要加入交易
	tradeIDs := make([][]byte, 0)
	for _, t := range b.TradeList {
		tradeIDs = append(tradeIDs, t.ID)
	}
	for _, tradeID := range tradeIDs {
		buf.Write(tradeID)
	}
	// 计算哈希值
	hash := sha256.Sum256(buf.Bytes())
	b.Hash = hash[:]
}

// CreateBlock 创建区块
func CreateBlock(prevHash []byte, trades []*trade.Trade) *Block {
	block := Block{time.Now(), []byte{}, prevHash, []byte{}, 0, trades}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 通过创世区块生成区块链
func InitBlockChain() *BlockChain {
	blockChain := BlockChain{}
	//data := "BlockChain,Start!"
	firstTrade := trade.FirstTrade([]byte("Leo Cao"))
	firstBlock := CreateBlock([]byte{}, []*trade.Trade{firstTrade})
	blockChain.BlockList = append(blockChain.BlockList, firstBlock)

	//fmt.Println(len(blockChain.BlockList))

	return &blockChain
}

// 添加新区块
func (blockChain *BlockChain) AddBlock(trades []*trade.Trade) {
	newBlock := CreateBlock(blockChain.BlockList[len(blockChain.BlockList)-1].Hash, trades)
	blockChain.BlockList = append(blockChain.BlockList, newBlock)
}

//// 寻找可用交易信息
//func (blockChain *BlockChain) FindValidTrades(address []byte) {
//	// 存放可用交易信息
//	var unSpentTrades []trade.Trade
//	// 存放已使用交易信息
//	spentTrades := make(map[string][]int)
//	// 开始遍历交易
//	for i := len(blockChain.BlockList) - 1; i > 0; i-- {
//		block := blockChain.BlockList[i]
//		for _, trade := range block.TradeList {
//			tradeID := hex.EncodeToString(trade.ID)
//
//		IterOutputs:
//		}
//
//	}
//}

// 寻找可用交易信息
func (blockChain *BlockChain) FindUnspentTrades(address []byte) []trade.Trade {
	// 存放可用交易信息
	var unSpentTrades []trade.Trade
	// 存放已使用交易信息
	spentTrades := make(map[string][]int)
	// 开始遍历交易

	for i := len(blockChain.BlockList) - 1; i >= 0; i-- {
		block := blockChain.BlockList[i]
		for _, trade := range block.TradeList {
			tradeID := hex.EncodeToString(trade.ID)

		IterOutputs:

			for outId, out := range trade.Outputs {
				// 若tradeID已经spent则跳过
				if spentTrades[tradeID] != nil {
					for _, spentOut := range spentTrades[tradeID] {
						if spentOut == outId {
							continue IterOutputs
						}
					}
				}
				// 否则查看地址，当正确则为查找的信息
				if out.IsToAddressRight(address) {
					unSpentTrades = append(unSpentTrades, *trade)
				}
			}

			// 判断是不是初始交易
			// 不是的话判断in是否包含目标地址，有的话将out信息加入spentTrades
			if !trade.IsFirstTrade() {
				for _, in := range trade.Inputs {
					if in.IsFromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TradeID)
						spentTrades[inTxID] = append(spentTrades[inTxID], in.OutID)
					}
				}
			}
		}

	}
	return unSpentTrades
}

// 找到一个地址的全部UTXO
func (blockChain *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTrades := blockChain.FindUnspentTrades(address)
	accumulated := 0

Work:
	for _, trade := range unspentTrades {
		txID := hex.EncodeToString(trade.ID)
		for outIdx, out := range trade.Outputs {
			if out.IsToAddressRight(address) {
				accumulated += out.Num
				unspentOuts[txID] = outIdx
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

// 找到可用的UTXO
// 即资产量大于转账额
func (blockChain *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := blockChain.FindUnspentTrades(address)
	accumulated := 0

Work:
	for _, trade := range unspentTxs {
		tradeID := hex.EncodeToString(trade.ID)
		for outId, out := range trade.Outputs {
			if out.IsToAddressRight(address) && accumulated < amount {
				accumulated += out.Num
				unspentOuts[tradeID] = outId
				if accumulated >= amount {
					break Work
				}
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

// 创建交易
func (blockChain *BlockChain) CreateTransaction(from, to []byte, amount int) (*trade.Trade, bool) {
	var inputs []trade.TradeIn
	var outputs []trade.TradeOut

	acc, validOutputs := blockChain.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("余额不足!")
		return &trade.Trade{}, false
	}
	for tradeID, outID := range validOutputs {
		txID, err := hex.DecodeString(tradeID)
		if err != nil {
			return nil, false
		}
		input := trade.TradeIn{txID, outID, from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, trade.TradeOut{amount, to})
	if acc > amount {
		outputs = append(outputs, trade.TradeOut{acc - amount, from})
	}
	tx := trade.Trade{nil, inputs, outputs}
	tx.SetID()

	return &tx, true
}

// 暂时预留一个挖矿接口
func (blockChain *BlockChain) Mine(txs []*trade.Trade) {
	blockChain.AddBlock(txs)
}
