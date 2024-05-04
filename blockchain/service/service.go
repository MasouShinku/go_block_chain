package service

import (
	"blockchain/blockchain"
	"bytes"
	"fmt"
	"strconv"
)

type Service struct{}

func (s *Service) CreateBlockChain(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}

func (s *Service) Balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	balance, _ := chain.FindUTXOs([]byte(address))
	fmt.Printf("Address:%s, Balance:%d \n", address, balance)
}

func (s *Service) GetBlockChainInfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iterator := chain.InitIterator()
	ogprevhash := chain.GetOGPrevHash()
	for {
		block := iterator.Next()
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Time.Format("2006-01-02 15:04:05"))
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Trades:%v\n", block.TradeList)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}
}

func (s *Service) Send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	tx, ok := chain.CreateTrade([]byte(from), []byte(to), amount)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}
	tp := blockchain.CreateTradePool()
	tp.AddTrade(tx)
	tp.SaveFile()
	fmt.Println("Success!")
}

func (s *Service) Mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.Mine()
	fmt.Println("Finish Mining")
}
