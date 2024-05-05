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
						TradeID:     hex.EncodeToString(input.TradeID),
						OutID:       input.OutID,
						FromAddress: string(input.FromAddress),
					}
				}
				for i, output := range trade.Outputs {
					tInfo.Outputs[i] = OutputInfo{
						Num:       output.Num,
						ToAddress: string(output.ToAddress),
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

	//var buffer bytes.Buffer
	//
	//chain := blockchain.ContinueBlockChain()
	//defer chain.Database.Close()
	//iterator := chain.InitIterator()
	//ogprevhash := chain.GetOGPrevHash()
	//
	//for {
	//	block := iterator.Next()
	//
	//	buffer.WriteString("--------------------------------------------------------------------------------------------------------------\n")
	//	buffer.WriteString(fmt.Sprintf("Timestamp:%s\n", block.Time.Format("2006-01-02 15:04:05")))
	//	buffer.WriteString(fmt.Sprintf("Previous hash:%x\n", block.PrevHash))
	//	buffer.WriteString(fmt.Sprintf("Trades:%v\n", block.TradeList))
	//	buffer.WriteString(fmt.Sprintf("Hash:%x\n", block.Hash))
	//	buffer.WriteString(fmt.Sprintf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW())))
	//	buffer.WriteString("--------------------------------------------------------------------------------------------------------------\n")
	//	buffer.WriteString("\n")
	//
	//	if bytes.Equal(block.PrevHash, ogprevhash) {
	//		break
	//	}
	//}
	//
	//return buffer.String()

	//chain := blockchain.ContinueBlockChain()
	//defer chain.Database.Close()
	//iterator := chain.InitIterator()
	//ogprevhash := chain.GetOGPrevHash()
	//for {
	//	block := iterator.Next()
	//	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	//	fmt.Printf("Timestamp:%d\n", block.Time.Format("2006-01-02 15:04:05"))
	//	fmt.Printf("Previous hash:%x\n", block.PrevHash)
	//	fmt.Printf("Trades:%v\n", block.TradeList)
	//	fmt.Printf("hash:%x\n", block.Hash)
	//	fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
	//	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	//	fmt.Println()
	//	if bytes.Equal(block.PrevHash, ogprevhash) {
	//		break
	//	}
	//}
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
