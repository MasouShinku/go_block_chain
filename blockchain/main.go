package main

import (
	"blockchain/blockchain"
	"fmt"
	"time"
)

//// Block 区块结构体
//type Block struct{
//	Time int		// 时间戳
//	Hash []byte		// 哈希值
//	PrevHash []byte	    // 上一个区块的哈希值
//	Data []byte // 数据属性
//}
//
//// BlockChain 区块链
//type BlockChain struct {
//	BlockList []*Block
//}

func main() {

	blockchain := blockchain.InitBlockChain()
	time.Sleep(time.Second)
	blockchain.AddBlock("here is the first block")
	time.Sleep(time.Second)
	blockchain.AddBlock("man!")
	time.Sleep(time.Second)
	blockchain.AddBlock("what can i say!")
	time.Sleep(time.Second)

	for _, block := range blockchain.BlockList {
		fmt.Println("=====================Block Info=====================")
		fmt.Printf("Time: %d\n", block.Time.Format("2006-01-02 15:04:05"))
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Println("====================================================\n\n")

	}
}
