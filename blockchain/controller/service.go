package main

import (
	"blockchain/blockchain"
	"blockchain/util"
	"blockchain/wallet"
	"bytes"
	"encoding/hex"
	"fmt"
)

type Service struct{}

func (s *Service) CreateBlockChain(address string) {
	newChain := blockchain.InitBlockChain(util.AddressToPublicHash([]byte(address)))
	newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}

func (s *Service) Balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	wallet := wallet.LoadWallet(address)

	balance, _ := chain.FindUTXOs(wallet.PublicKey)
	fmt.Printf("Address:%s, Balance:%d \n", address, balance)
}

func (s *Service) GetBlockChainInfo() []BlockInfo {

	var blocks []BlockInfo

	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iterator := chain.InitIterator()
	ogprevhash := chain.GetOGPrevHash()

	for {
		block := iterator.Next()
		if block == nil {
			break
		}

		blockInfo := BlockInfo{
			Timestamp:    block.Time.Format("2006-01-02 15:04:05"),
			PreviousHash: fmt.Sprintf("%x", block.PrevHash),
			Trades:       nil,
			Hash:         fmt.Sprintf("%x", block.Hash),
			Pow:          block.ValidatePoW(),
		}

		// 接下来处理tradelist
		var tradesInfo []TradeInfo
		for _, trade := range block.TradeList {
			if trade != nil { // 确保指针非空
				tInfo := TradeInfo{
					ID:      hex.EncodeToString(trade.ID),
					Inputs:  make([]InputInfo, len(trade.Inputs)),
					Outputs: make([]OutputInfo, len(trade.Outputs)),
				}
				for i, input := range trade.Inputs {
					tInfo.Inputs[i] = InputInfo{
						TradeID: hex.EncodeToString(input.TradeID),
						OutID:   input.OutID,
						PubKey:  string(input.PublicKey),
					}
				}
				for i, output := range trade.Outputs {
					tInfo.Outputs[i] = OutputInfo{
						Num:        output.Num,
						HashPubKey: string(output.HashPublicKey),
					}
				}
				tradesInfo = append(tradesInfo, tInfo)
			}
		}

		// 现在 tradesInfo 包含实际的数据而非指针
		blockInfo.Trades = tradesInfo

		blocks = append(blocks, blockInfo)

		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}

	return blocks
}

func (s *Service) Send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	fromWallet := wallet.LoadWallet(from)
	trade, ok := chain.CreateTrade(fromWallet.PublicKey, util.AddressToPublicHash([]byte(to)), amount, fromWallet.PrivateKey)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}
	tp := blockchain.CreateTradePool()
	tp.AddTrade(trade)
	tp.SaveFile()
	fmt.Println("Success!")
}

func (s *Service) Mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.Mine()
	fmt.Println("Finish Mining")
}

func (s *Service) CreateWallet(refname string) {
	newWallet := wallet.NewWallet()

	newWallet.SaveWallet()
	refList := wallet.LoadRefList()
	refList.SetRef(string(newWallet.Address()), refname)
	refList.Save()
	fmt.Println("Succeed in creating wallet.")
}

func (s *Service) WalletInfoRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	util.Err(err)
	s.WalletInfo(address)
}

func (s *Service) WalletInfo(address string) {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	fmt.Printf("Wallet address:%x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
	fmt.Printf("Reference Name:%s\n", (*refList)[address])
}

func (s *Service) UpdateWallets() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("Succeed in updating wallets.")
}

func (s *Service) WalletsList() {
	refList := wallet.LoadRefList()
	for address, _ := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s\n", (*refList)[address])
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}

func (s *Service) SendRefName(fromRefname, toRefname string, amount int) {
	refList := wallet.LoadRefList()
	fromAddress, err := refList.FindRef(fromRefname)
	util.Err(err)
	toAddress, err := refList.FindRef(toRefname)
	util.Err(err)
	s.Send(fromAddress, toAddress, amount)
}

func (s *Service) CreateBlockChainRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	util.Err(err)
	s.CreateBlockChain(address)
}

func (s *Service) BalanceRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	util.Err(err)
	s.Balance(address)
}
