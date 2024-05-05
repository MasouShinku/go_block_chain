package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 载入区块链服务
func loadBlockChain() *Service {
	s := Service{}
	s.CreateBlockChain("shinku")
	return &s
}

func main() {

	s := loadBlockChain()
	s.Send("shinku", "shizuka", 100)
	s.Mine()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"blocks": s.GetBlockChainInfo(),
		})
	})
	r.Run()
}
