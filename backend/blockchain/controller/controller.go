package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 载入区块链服务
func loadBlockChain() *Service {
	s := Service{}
	//s.CreateBlockChain("shinku")
	return &s
}

func main() {

	s := loadBlockChain()
	r := gin.Default()

	r.GET("/create_blockchain/:address", func(c *gin.Context) {
		address := c.Param("address")
		result := s.CreateBlockChain(address)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/wallet_info/:address", func(c *gin.Context) {
		address := c.Param("address")
		result := s.WalletInfo(address)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/wallet_info_ref/:refname", func(c *gin.Context) {
		refname := c.Param("refname")
		result := s.WalletInfoRefName(refname)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/wallets_list", func(c *gin.Context) {
		result := s.WalletsList()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/balance/:address", func(c *gin.Context) {
		address := c.Param("address")
		result := s.Balance(address)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/balance_ref/:refname", func(c *gin.Context) {
		address := c.Param("refname")
		result := s.BalanceRefName(address)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/send/:from/:to/:amount", func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		amount, _ := strconv.Atoi(c.Param("amount"))
		result := s.Send(from, to, amount)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/send_by_ref/:from/:to/:amount", func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		amount, _ := strconv.Atoi(c.Param("amount"))
		result := s.SendRefName(from, to, amount)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/mine", func(c *gin.Context) {
		result := s.Mine()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/update_wallets", func(c *gin.Context) {
		result := s.UpdateWallets()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/create_wallet/:refname", func(c *gin.Context) {
		refname := c.Param("refname")
		result := s.CreateWallet(refname)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/blockchain_info", func(c *gin.Context) {
		result := s.GetBlockChainInfo()
		c.JSON(http.StatusOK, result)
	})

	r.Run()
}
