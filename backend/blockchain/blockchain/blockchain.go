package blockchain

import (
	"blockchain/trade"
	"blockchain/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"github.com/dgraph-io/badger"
	"io/ioutil"
	"os"
	"runtime"
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
}

// BlockChain 区块链
type BlockChain struct {
	LastHash []byte     // 最后一个区块的Hash
	Database *badger.DB // 指向数据库
}

// BlockChainIterator 区块链迭代器
// 遍历区块链时使用
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// TradePool 交易池
type TradePool struct {
	TradeInfo []*trade.Trade // 收集到的交易信息
}

func (tp *TradePool) AddTrade(trade *trade.Trade) {
	tp.TradeInfo = append(tp.TradeInfo, trade)
}

// SaveFile 保存交易信息
func (tp *TradePool) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(tp)
	util.Err(err)
	err = ioutil.WriteFile(util.TradePool, content.Bytes(), 0644)
	util.Err(err)
}

// LoadFile 读取交易信息
func (tp *TradePool) LoadFile() error {
	if !util.FileExists(util.TradePool) {
		return nil
	}

	var tradePool TradePool

	fileContent, err := ioutil.ReadFile(util.TradePool)
	if err != nil {
		util.Err(err)
		return err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&tradePool)
	if err != nil {
		util.Err(err)
		return err
	}

	tp.TradeInfo = tradePool.TradeInfo
	return nil
}

// CreateTradePool 创建或加载交易池
func CreateTradePool() *TradePool {
	tradePool := TradePool{}
	err := tradePool.LoadFile()
	util.Err(err)

	return &tradePool
}

// RemoveTradePoolFile 移除交易池
func RemoveTradePoolFile() error {
	err := os.Remove(util.TradePool)
	return err
}

// InitIterator 初始化迭代器
func (blockChain *BlockChain) InitIterator() *BlockChainIterator {
	iterator := BlockChainIterator{blockChain.LastHash, blockChain.Database}
	return &iterator
}

// GetOGPrevHash 获取OGPrevHash
func (chain *BlockChain) GetOGPrevHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		util.Err(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		util.Err(err)
		return err
	})
	util.Err(err)

	return ogprevhash
}

// Next 迭代器的迭代函数
func (iterator *BlockChainIterator) Next() *Block {
	var tempBlock *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		util.Err(err)

		err = item.Value(func(val []byte) error {
			tempBlock = DeSerializeBlock(val)
			return nil
		})
		util.Err(err)
		return err
	})
	util.Err(err)

	iterator.CurrentHash = tempBlock.PrevHash

	return tempBlock
}

// SetHash 构建区块哈希值
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

// InitBlockChain 通过创世区块生成区块链
func InitBlockChain(address []byte) *BlockChain {

	var lastHash []byte
	if util.FileExists(util.BCFile) {
		util.Info("区块链已存在...")
		return ContinueBlockChain()
	}

	opts := badger.DefaultOptions(util.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	util.Err(err)

	err = db.Update(func(txn *badger.Txn) error {

		firstTrade := trade.FirstTrade(address)
		firstBlock := CreateBlock([]byte("无prevHash..."), []*trade.Trade{firstTrade})
		firstBlock.SetHash()

		util.Info("创世区块成功创建！")

		if err = txn.Set(firstBlock.Hash, firstBlock.Serialize()); err != nil {
			util.Err(err)
		}
		if err = txn.Set([]byte("lh"), firstBlock.Hash); err != nil {
			util.Err(err)
		}
		if err = txn.Set([]byte("ogprevhash"), firstBlock.PrevHash); err != nil {
			util.Err(err)
		}
		lastHash = firstBlock.Hash
		return err
	})

	util.Err(err)

	blockChain := BlockChain{lastHash, db}
	return &blockChain
}

// ContinueBlockChain 加载区块链
func ContinueBlockChain() *BlockChain {
	if util.FileExists(util.BCFile) == false {
		util.Err(errors.New("没有找到区块链..."))
		runtime.Goexit()
	}
	var lastHash []byte

	opts := badger.DefaultOptions(util.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	util.Err(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		util.Err(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		util.Err(err)
		return err
	})
	util.Err(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

// AddBlock 添加新区块
func (blockChain *BlockChain) AddBlock(newBlock *Block) {
	var lastHash []byte

	err := blockChain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		util.Err(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		util.Err(err)

		return err
	})
	util.Err(err)
	// 检查一下新区块的前向hash和区块链的末尾hash是否一致
	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		util.Err(errors.New("哈希检查不一致！"))
		runtime.Goexit()
	}

	err = blockChain.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		util.Err(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)
		blockChain.LastHash = newBlock.Hash
		return err
	})
	util.Err(err)
}

// FindUnspentTrades 寻找可用交易信息
func (blockChain *BlockChain) FindUnspentTrades(address []byte) []trade.Trade {
	// 存放可用交易信息
	var unSpentTrades []trade.Trade
	// 存放已使用交易信息
	spentTrades := make(map[string][]int)
	iterator := blockChain.InitIterator()

all:
	for {
		block := iterator.Next()

		for _, trade := range block.TradeList {
			txID := hex.EncodeToString(trade.ID)

		IterOutputs:
			for outIdx, out := range trade.Outputs {
				if spentTrades[txID] != nil {
					for _, spentOut := range spentTrades[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}

				if out.IsToAddressRight(address) {
					unSpentTrades = append(unSpentTrades, *trade)
				}
			}
			if !trade.IsFirstTrade() {
				for _, in := range trade.Inputs {
					if in.IsFromAddressRight(address) {
						inTradeID := hex.EncodeToString(in.TradeID)
						spentTrades[inTradeID] = append(spentTrades[inTradeID], in.OutID)
					}
				}
			}
		}
		if bytes.Equal(block.PrevHash, blockChain.GetOGPrevHash()) {
			break all
		}
	}
	return unSpentTrades

}

// FindUTXOs 找到一个地址的全部UTXO
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

// FindSpendableOutputs 找到可用的UTXO
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

// CreateTrade 创建交易
func (blockChain *BlockChain) CreateTrade(fromPublicKey, toHashPublicKey []byte, amount int, privateKey ecdsa.PrivateKey, des string) (*trade.Trade, bool) {

	var inputs []trade.TradeIn
	var outputs []trade.TradeOut

	acc, validOutputs := blockChain.FindSpendableOutputs(fromPublicKey, amount)
	if acc < amount {
		util.Info("余额不足")
		return &trade.Trade{}, false
	}
	for tradeID, outID := range validOutputs {
		tID, err := hex.DecodeString(tradeID)
		util.Err(err)
		input := trade.TradeIn{tID, outID, fromPublicKey, nil}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, trade.TradeOut{amount, toHashPublicKey})
	if acc > amount {
		outputs = append(outputs, trade.TradeOut{acc - amount, util.PublicKeyHash(fromPublicKey)})
	}
	t := trade.Trade{nil, inputs, outputs, des}
	t.SetID()
	t.Sign(privateKey)
	return &t, true
}

// Serialize 序列化区块
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(b); err != nil {
		util.Err(err)
	}
	return res.Bytes()
}

// DeSerializeBlock 反序列化区块
func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		util.Err(err)
	}
	return &block
}

type TradeInfo struct {
	ID          string
	Inputs      []InputInfo
	Outputs     []OutputInfo
	Description string
}

type InputInfo struct {
	TradeID string
	OutID   int
	PubKey  string
}

type OutputInfo struct {
	Num        int
	HashPubKey string
}

// TraceCurrency 追踪货币的流向
func (chain *BlockChain) TraceCurrency(address []byte) ([]TradeInfo, error) {
	var trades []TradeInfo
	visited := make(map[string]bool) // 防止重复处理同一个交易

	iterator := chain.InitIterator()
	for {
		block := iterator.Next()
		if block == nil || len(block.PrevHash) == 0 {
			break
		}

		for _, t := range block.TradeList {
			txID := hex.EncodeToString(t.ID)
			// 检查是否已处理过这个交易
			if visited[txID] {
				continue
			}

			for _, in := range t.Inputs {
				if bytes.Equal(in.PublicKey, address) {
					trades = append(trades, createTradeInfo(t))
					visited[txID] = true
					break
				}
			}

			for _, out := range t.Outputs {
				if bytes.Equal(util.PublicKeyHash(address), out.HashPublicKey) {
					trades = append(trades, createTradeInfo(t))
					visited[txID] = true
					break
				}
			}
		}
	}

	return trades, nil
}

func createTradeInfo(t *trade.Trade) TradeInfo {
	inputs := make([]InputInfo, len(t.Inputs))
	for i, in := range t.Inputs {
		inputs[i] = InputInfo{
			TradeID: hex.EncodeToString(in.TradeID),
			OutID:   in.OutID,
			PubKey:  hex.EncodeToString(in.PublicKey),
		}
	}
	outputs := make([]OutputInfo, len(t.Outputs))
	for i, out := range t.Outputs {
		outputs[i] = OutputInfo{
			Num:        out.Num,
			HashPubKey: hex.EncodeToString(out.HashPublicKey),
		}
	}
	return TradeInfo{
		ID:          hex.EncodeToString(t.ID),
		Inputs:      inputs,
		Outputs:     outputs,
		Description: t.Description,
	}
}
