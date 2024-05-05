package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 载入区块链服务
func loadBlockChain() *Service {
	s := Service{}
	//s.CreateBlockChain("shinku")
	return &s
}

func main() {
	//
	//s := loadBlockChain()
	//s.Send("shinku", "shizuka", 100)
	//s.Mine()

	s := loadBlockChain()
	//s.CreateBlockChain("122BZqDc5vPcgAG6DED8TcCHKQcosQujSu")
	//s.CreateWallet("")

	s.WalletsList()
	s.Balance("122BZqDc5vPcgAG6DED8TcCHKQcosQujSu")
	s.Balance("1MQ5dfua1bFmfAUqY3cZHuxeHhBPV2jGao")
	s.UpdateWallets()
	s.WalletsList()
	//s.Send("122BZqDc5vPcgAG6DED8TcCHKQcosQujSu", "1MQ5dfua1bFmfAUqY3cZHuxeHhBPV2jGao", 10)
	//s.Mine()

	//s.CreateBlockChain("1K2FNmJtHZZYKrJFSYfjAfXzGjmj5YdUPk")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"blocks": s.GetBlockChainInfo(),
		})
	})
	r.Run()
}
