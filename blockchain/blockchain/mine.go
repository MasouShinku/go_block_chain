package blockchain

import (
	"blockchain/trade"
	"blockchain/util"
	"bytes"
	"encoding/hex"
	"fmt"
)

func isInputRight(trades []trade.Trade, in trade.TradeIn) (bool, int) {
	for _, tx := range trades {
		if bytes.Equal(tx.ID, in.TradeID) {
			return true, tx.Outputs[in.OutID].Num
		}
	}
	return false, 0
}

// 验证交易信息有效性
func (blockChain *BlockChain) VerifyTrades(trades []*trade.Trade) bool {
	if len(trades) == 0 {
		return true
	}
	spentOutputs := make(map[string]int)
	for _, tx := range trades {
		pubKey := tx.Inputs[0].PublicKey
		unspentOutputs := blockChain.FindUnspentTrades(pubKey)
		inputAmount := 0
		OutputAmount := 0

		for _, input := range tx.Inputs {
			if outid, ok := spentOutputs[hex.EncodeToString(input.TradeID)]; ok && outid == input.OutID {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TradeID)] = input.OutID
		}

		for _, output := range tx.Outputs {
			OutputAmount += output.Num
		}
		if inputAmount != OutputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

// 挖矿函数
func (blockchain *BlockChain) Mine() {
	tradePool := CreateTradePool()
	if !blockchain.VerifyTrades(tradePool.TradeInfo) {
		util.Info("交易验证失败...")
		err := RemoveTradePoolFile()
		util.Err(err)
		return
	}

	candidateBlock := CreateBlock(blockchain.LastHash, tradePool.TradeInfo)
	if candidateBlock.ValidatePoW() {
		blockchain.AddBlock(candidateBlock)
		err := RemoveTradePoolFile()
		util.Err(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
