package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
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

// hash公钥
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hash := ripemd160.New()
	_, err := hash.Write(hashedPublicKey[:])
	Err(err)
	publicRipeMd := hash.Sum(nil)
	return publicRipeMd
}

// 生成密钥对
func GenNeKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	Err(err)
	// 公钥为点坐标，需要拼接起来保存
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

// 检查位生成函数
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:ChecksumLength]
}

// base256转base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

// base58转base256
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Err(err)
	return decode
}

// 从公钥哈希生成钱包地址
func PublicHashToAddress(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

// 从钱包地址转公钥哈希
func AddressToPublicHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLength]
	return pubKeyHash
}
