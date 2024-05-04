package util

import "os"

const (
	Difficulty = 12
	InitNum    = 1000
	TradePool  = "./files/tradePool.data"  //this line is new
	BCPath     = "./files/blocks"          //this line is new
	BCFile     = "./files/blocks/MANIFEST" //this line is new
)

// 检查文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
