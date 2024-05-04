package blockchain

import "fmt"

// 挖矿函数
func (blockchain *BlockChain) Mine() {
	tradePool := CreateTradePool()
	candidateBlock := CreateBlock(blockchain.LastHash, tradePool.TradeInfo)
	if candidateBlock.ValidatePoW() {
		blockchain.AddBlock(candidateBlock)
		if err := RemoveTradePoolFile(); err != nil {
			fmt.Println(err)
		}
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
