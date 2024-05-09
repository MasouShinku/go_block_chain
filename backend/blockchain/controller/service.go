package main

import (
	"blockchain/blockchain"
	"blockchain/util"
	"blockchain/wallet"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Service struct{}

func (s *Service) InitBlockChain() {
	s.CreateWallet("原料厂", util.Identity("Raw"))
	s.CreateWallet("贵州生产商", util.Identity("Producer"))
	//s.CreateWallet("重庆生产商", util.Identity("Producer"))
	s.CreateWallet("北京经销商", util.Identity("Dealer"))
	s.CreateWallet("上海经销商", util.Identity("Dealer"))
	s.CreateWallet("天津经销商", util.Identity("Dealer"))
	s.CreateWallet("用户", util.Identity("User"))
	address := s.WalletInfoRefName("原料厂").Address
	s.CreateBlockChain(address)
	s.SendRefName("原料厂", "贵州生产商", 900, "贵州生产商进货")
	s.Mine()
	time.Sleep(1 * time.Second)
	s.SendRefName("贵州生产商", "北京经销商", 189, "北京经销商进货")
	s.Mine()
	s.SendRefName("贵州生产商", "上海经销商", 243, "上海经销商进货")
	s.Mine()
	s.SendRefName("贵州生产商", "天津经销商", 199, "天津经销商进货")
	s.Mine()
	time.Sleep(1 * time.Second)
	s.SendRefName("天津经销商", "用户", 1, "用户购买")
	s.Mine()
	s.traceCurrency()

}

func (s *Service) CreateBlockChain(address string) BlockchainCreationResult {
	newChain := blockchain.InitBlockChain(util.AddressToPublicHash([]byte(address)))
	err := newChain.Database.Close()
	if err != nil {
		util.Err(err)
		return BlockchainCreationResult{
			Success: false,
			Message: err.Error(),
		}
	}
	return BlockchainCreationResult{
		Success: true,
		Message: "Finished creating blockchain,a and the owner is: " + address,
	}
}

func (s *Service) Balance(address string) BalanceResult {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	wallet := wallet.LoadWallet(address)

	balance, _ := chain.FindUTXOs(wallet.PublicKey)
	util.Info(fmt.Sprintf("Address:%s, Balance:%d \n", address, balance))
	return BalanceResult{
		Address: address,
		Balance: balance,
	}
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
					ID:          hex.EncodeToString(trade.ID),
					Inputs:      make([]InputInfo, len(trade.Inputs)),
					Outputs:     make([]OutputInfo, len(trade.Outputs)),
					Description: trade.Description,
				}
				for i, input := range trade.Inputs {
					tInfo.Inputs[i] = InputInfo{
						TradeID: hex.EncodeToString(input.TradeID),
						OutID:   input.OutID,
						PubKey:  fmt.Sprintf("%x", input.PublicKey),
					}
					//util.Info("Key : " + fmt.Sprintf("%x", input.PublicKey))
				}
				for i, output := range trade.Outputs {
					tInfo.Outputs[i] = OutputInfo{
						Num:        output.Num,
						HashPubKey: fmt.Sprintf("%x", output.HashPublicKey),
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

func (s *Service) Send(from, to string, amount int, des string) TradeResult {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	fromWallet := wallet.LoadWallet(from)

	toPubHash := util.AddressToPublicHash([]byte(to))
	trade, ok := chain.CreateTrade(fromWallet.PublicKey, toPubHash, amount, fromWallet.PrivateKey, des)
	if !ok {
		return TradeResult{Success: false, Message: "Failed to create trade"}
	}

	tp := blockchain.CreateTradePool()
	tp.AddTrade(trade)
	tp.SaveFile()

	return TradeResult{Success: true, Message: "Trade successful"}
}

func (s *Service) Mine() MiningResult {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.Mine()

	return MiningResult{
		Message: "Mining successful",
	}
}

func (s *Service) CreateWallet(refname string, identity util.Identity) CreateWalletResult {
	newWallet := wallet.NewWallet(identity)

	newWallet.SaveWallet()
	refList := wallet.LoadRefList()
	refList.SetRef(string(newWallet.Address()), refname)
	refList.Save()
	return CreateWalletResult{
		Message: "Succeed in creating wallet.",
	}
}

func (s *Service) WalletInfo(address string) WalletInfoResult {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()

	return WalletInfoResult{
		Address:       address,
		PublicKey:     fmt.Sprintf("%x", wlt.PublicKey),
		ReferenceName: (*refList)[address],
		Identity:      string(wlt.Identity),
	}
}

func (s *Service) WalletInfoRefName(refname string) WalletInfoResult {
	refList := wallet.LoadRefList()
	address, _ := refList.FindRef(refname) // 这里忽略错误处理，根据您的指示

	return s.WalletInfo(address)
}

func (s *Service) UpdateWallets() UpdateWalletsResult {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	return UpdateWalletsResult{
		Message: ("Succeed in updating wallets."),
	}
}

func (s *Service) WalletsList() WalletsListResult {
	refList := wallet.LoadRefList()
	var wallets []WalletInfoResult

	for address := range *refList {
		wlt := wallet.LoadWallet(address)
		walletInfo := WalletInfoResult{
			Address:       address,
			PublicKey:     fmt.Sprintf("%x", wlt.PublicKey),
			ReferenceName: (*refList)[address],
			Identity:      string(wlt.Identity),
		}
		wallets = append(wallets, walletInfo)
	}
	return WalletsListResult{
		Wallets: wallets,
	}
}

func (s *Service) SendRefName(fromRefname, toRefname string, amount int, des string) TradeResult {
	refList := wallet.LoadRefList()

	fromAddress, err := refList.FindRef(fromRefname)
	if err != nil {
		return TradeResult{Success: false, Message: "Failed to find sender address: " + err.Error()}
	}

	toAddress, err := refList.FindRef(toRefname)
	if err != nil {
		return TradeResult{Success: false, Message: "Failed to find receiver address: " + err.Error()}
	}

	return s.Send(fromAddress, toAddress, amount, des)
}

func (s *Service) CreateBlockChainRefName(refname string) BlockchainCreationResult {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	util.Err(err)
	return s.CreateBlockChain(address)
}

func (s *Service) BalanceRefName(refname string) BalanceResult {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	util.Err(err)
	return s.Balance(address)
}

func (s *Service) getAllBalance() getAllBalanceResult {
	return getAllBalanceResult{
		Raw_balance:        strconv.Itoa(s.BalanceRefName("原料厂").Balance),
		A_producer_balance: strconv.Itoa(s.BalanceRefName("贵州生产商").Balance),
		A_dealer_balance:   strconv.Itoa(s.BalanceRefName("北京经销商").Balance),
		B_dealer_balance:   strconv.Itoa(s.BalanceRefName("上海经销商").Balance),
		C_dealer_balance:   strconv.Itoa(s.BalanceRefName("天津经销商").Balance),
		User_balance:       strconv.Itoa(s.BalanceRefName("用户").Balance),
	}
}

func (s *Service) traceCurrency() []TraceTrade {

	hasDealer := false
	hasProducer := false

	blocks := s.GetBlockChainInfo()
	// 检查区块链是否为空
	if len(blocks) == 0 {
		util.Info("区块链暂无区块！")
		return nil
	}
	// 从最新的区块的第一个交易开始追溯
	currentTrade := blocks[0].Trades[0] // 最新区块的第一个交易
	//fmt.Printf("Starting trace from trade %s\n", currentTrade.ID)
	tradeTrade := []TraceTrade{}
	if strings.Contains(blocks[0].Trades[0].Description, "用户") {
		tradeTrade = append(tradeTrade, TraceTrade{
			Time:        blocks[0].Timestamp,
			Description: currentTrade.Description,
		},
		)
		util.Info(currentTrade.Description)
	}

	for {
		// 打印当前交易信息
		for _, input := range currentTrade.Inputs {
			//fmt.Printf("Input from trade %s, output index %d, pubKey %s\n", input.TradeID, input.OutID, input.PubKey)
			// 寻找输入来源交易
			found := false
			for _, block := range blocks {
				for _, trade := range block.Trades {
					if trade.ID == input.TradeID {
						currentTrade = trade
						if strings.Contains(currentTrade.Description, "经销商") && hasDealer == false {

							tradeTrade = append(tradeTrade, TraceTrade{
								Time:        block.Timestamp,
								Description: currentTrade.Description,
							},
							)
							util.Info(currentTrade.Description)
							hasDealer = true
						}
						if strings.Contains(currentTrade.Description, "生产商") && hasProducer == false {

							tradeTrade = append(tradeTrade, TraceTrade{
								Time:        block.Timestamp,
								Description: currentTrade.Description,
							},
							)
							util.Info(currentTrade.Description)
							hasProducer = true
						}

						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				util.Err(errors.New("源交易未找到..."))
				return nil
			}
			if currentTrade.Description == "first trade" {
				util.Info("到达首次交易！")
				return tradeTrade
			}
		}
	}
}

func (s *Service) buy() BuyResult {
	if s.BalanceRefName("北京经销商").Balance+s.BalanceRefName("上海经销商").Balance+s.BalanceRefName("天津经销商").Balance <= 0 {

		return BuyResult{
			Success: false,
			Message: "经销商余量不足",
		}
	}
	// 创建一个随机数生成器
	rand.Seed(time.Now().UnixNano())
	selected := false
	for selected == false {
		randomIndex := rand.Intn(3)
		switch randomIndex {
		case 0:
			if s.BalanceRefName("北京经销商").Balance > 0 {
				s.SendRefName("北京经销商", "用户", 1, "用户购买")
				selected = true
			}
		case 1:
			if s.BalanceRefName("上海经销商").Balance > 0 {
				s.SendRefName("上海经销商", "用户", 1, "用户购买")
				selected = true
			}
		case 2:
			if s.BalanceRefName("天津经销商").Balance > 0 {
				s.SendRefName("天津经销商", "用户", 1, "用户购买")
				selected = true
			}
		}
	}
	s.Mine()
	return BuyResult{
		Success:     true,
		Message:     "购买成功",
		TraceTrades: s.traceCurrency(),
	}

}

func (s *Service) producerBuy() ProducerBuyResult {
	if s.BalanceRefName("原料厂").Balance <= 0 {
		return ProducerBuyResult{
			Success: false,
			Message: "原料厂余量不足",
		}
	}

	s.SendRefName("原料厂", "贵州生产商", 1, "贵州生产商进货")

	s.Mine()
	return ProducerBuyResult{
		Success: true,
		Message: "进货成功",
	}

}

func (s *Service) dealerBuy() DealerBuyResult {
	if s.BalanceRefName("贵州生产商").Balance <= 0 {

		return DealerBuyResult{
			Success: false,
			Message: "生产商余量不足",
		}
	}
	// 创建一个随机数生成器
	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(3)
	switch randomIndex {
	case 0:
		s.SendRefName("贵州生产商", "北京经销商", 1, "北京经销商进货")

	case 1:
		s.SendRefName("贵州生产商", "上海经销商", 1, "上海经销商进货")

	case 2:
		s.SendRefName("贵州生产商", "天津经销商", 1, "天津经销商进货")

	}
	s.Mine()
	return DealerBuyResult{
		Success: true,
		Message: "购买成功",
	}

}
