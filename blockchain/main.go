package main

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

	//blockchain := blockchain.InitBlockChain()
	//time.Sleep(time.Second)
	//blockchain.AddBlock("here is the first block")
	//time.Sleep(time.Second)
	//blockchain.AddBlock("man!")
	//time.Sleep(time.Second)
	//blockchain.AddBlock("what can i say!")
	//time.Sleep(time.Second)
	//
	//for _, block := range blockchain.BlockList {
	//	fmt.Println("=====================Block Info=====================")
	//	fmt.Printf("Time: %d\n", block.Time.Format("2006-01-02 15:04:05"))
	//	fmt.Printf("hash: %x\n", block.Hash)
	//	fmt.Printf("Previous hash: %x\n", block.PrevHash)
	//	fmt.Printf("data: %s\n", block.Data)
	//	fmt.Println("====================================================\n\n")
	//
	//}

	//txPool := make([]*trade.Trade, 0)
	//var tempTx *trade.Trade
	//var ok bool
	//var property int
	//chain := blockchain.InitBlockChain()
	//property, _ = chain.FindUTXOs([]byte("Leo Cao"))
	//fmt.Println("Balance of Leo Cao: ", property)
	//
	//tempTx, ok = chain.CreateTrade([]byte("Leo Cao"), []byte("Krad"), 100)
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//chain.Mine(txPool)
	//txPool = make([]*trade.Trade, 0)
	//property, _ = chain.FindUTXOs([]byte("Leo Cao"))
	//fmt.Println("Balance of Leo Cao: ", property)
	//
	//tempTx, ok = chain.CreateTrade([]byte("Krad"), []byte("Exia"), 200) // this transaction is invalid
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//
	//tempTx, ok = chain.CreateTrade([]byte("Krad"), []byte("Exia"), 50)
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//
	//tempTx, ok = chain.CreateTrade([]byte("Leo Cao"), []byte("Exia"), 100)
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//chain.Mine(txPool)
	//txPool = make([]*trade.Trade, 0)
	//property, _ = chain.FindUTXOs([]byte("Leo Cao"))
	//fmt.Println("Balance of Leo Cao: ", property)
	//property, _ = chain.FindUTXOs([]byte("Krad"))
	//fmt.Println("Balance of Krad: ", property)
	//property, _ = chain.FindUTXOs([]byte("Exia"))
	//fmt.Println("Balance of Exia: ", property)
	//
	//for _, block := range chain.BlockList {
	//	fmt.Printf("Timestamp: %d\n", block.Time.Format("2006-01-02 15:04:05"))
	//	fmt.Printf("hash: %x\n", block.Hash)
	//	fmt.Printf("Previous hash: %x\n", block.PrevHash)
	//	fmt.Printf("nonce: %d\n", block.Nonce)
	//	fmt.Println("Proof of Work validation:", block.ValidatePoW())
	//}
	//
	////I want to show the bug at this version.
	//
	//tempTx, ok = chain.CreateTrade([]byte("Krad"), []byte("Exia"), 30)
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//
	//tempTx, ok = chain.CreateTrade([]byte("Krad"), []byte("Leo Cao"), 30)
	//if ok {
	//	txPool = append(txPool, tempTx)
	//}
	//
	//chain.Mine(txPool)
	//txPool = make([]*trade.Trade, 0)
	//
	//for _, block := range chain.BlockList {
	//	fmt.Printf("Timestamp: %d\n", block.Time.Format("2006-01-02 15:04:05"))
	//	fmt.Printf("hash: %x\n", block.Hash)
	//	fmt.Printf("Previous hash: %x\n", block.PrevHash)
	//	fmt.Printf("nonce: %d\n", block.Nonce)
	//	fmt.Println("Proof of Work validation:", block.ValidatePoW())
	//}
	//
	//property, _ = chain.FindUTXOs([]byte("Leo Cao"))
	//fmt.Println("Balance of Leo Cao: ", property)
	//property, _ = chain.FindUTXOs([]byte("Krad"))
	//fmt.Println("Balance of Krad: ", property)
	//property, _ = chain.FindUTXOs([]byte("Exia"))
	//fmt.Println("Balance of Exia: ", property)

	//service := Service{}
	//service.CreateBlockChain("LeoCao")
	//service.Balance("LeoCao")
	//service.GetBlockChainInfo()
	//service.Send("Leo", "Cao", 100)
	//service.Send("LeoCao", "Krad", 100)
	//service.Balance("LeoCao")
	//service.Mine()
	//service.Balance("LeoCao")
	//service.Balance("Krad")
	//service.GetBlockChainInfo()

}
