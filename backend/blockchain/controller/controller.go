package main

import (
	"blockchain/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 载入区块链服务
func loadBlockChain() *Service {
	util.CleanData()
	s := Service{}
	return &s
}

func main() {

	s := loadBlockChain()
	s.InitBlockChain()
	r := gin.Default()
	r.Use(cors.Default())

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

	r.GET("/send/:from/:to/:amount/:des", func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		amount, _ := strconv.Atoi(c.Param("amount"))
		des := c.Param("des")
		result := s.Send(from, to, amount, des)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/send_by_ref/:from/:to/:amount/:des", func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		amount, _ := strconv.Atoi(c.Param("amount"))
		des := c.Param("des")
		result := s.SendRefName(from, to, amount, des)
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

	r.GET("/create_wallet/:refname/:identity", func(c *gin.Context) {
		refname := c.Param("refname")
		identity := util.Identity(c.Param("identity"))
		result := s.CreateWallet(refname, identity)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/blockchain_info", func(c *gin.Context) {
		result := s.GetBlockChainInfo()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/trace_currency", func(c *gin.Context) {
		result := s.traceCurrency()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/all_balance", func(c *gin.Context) {
		result := s.getAllBalance()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/buy", func(c *gin.Context) {
		result := s.buy()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/producer_buy", func(c *gin.Context) {
		result := s.producerBuy()
		c.JSON(http.StatusOK, result)
	})

	r.GET("/dealer_buy", func(c *gin.Context) {
		result := s.dealerBuy()
		c.JSON(http.StatusOK, result)
	})

	r.Run(":8081")
}
