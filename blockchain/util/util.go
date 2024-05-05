package util

import (
	"log"
	"os"
	"sync"
)

// blockchain相关
const (
	Difficulty     = 12
	InitNum        = 1000
	TradePool      = "./files/tradePool.data"
	BCPath         = "./files/blocks"
	BCFile         = "./files/blocks/MANIFEST"
	ChecksumLength = 4
	NetworkVersion = byte(0x00)
	Wallets        = "./files/wallets/"
	WalletsRefList = "./files/ref_list/"
)

// 日志相关
const (
	preError   = "\033[31m[Error]\033[0m"
	preInfo    = "\033[34m[Info]\033[0m"
	InfoLevel  = 0
	ErrorLevel = 1
)

var (
	errorLog = log.New(os.Stdout, preError, log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, preInfo, log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// 暴露log方法
var (
	errorln = errorLog.Println
	infoln  = infoLog.Println
)

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

func Err(err error) {
	if err != nil {
		errorln(err)
	}
}

func Info(s string) {
	infoln(s)
}
