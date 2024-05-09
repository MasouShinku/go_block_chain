package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// blockchain相关
const (
	Difficulty     = 12                        // 区块链挖掘的难度
	InitNum        = 1000                      // 初始币数量
	TradePool      = "./files/tradePool.data"  // 交易池数据存储文件的路径
	BCPath         = "./files/blocks"          // 存放区块链数据的目录路径
	BCFile         = "./files/blocks/MANIFEST" // 区块链的清单文件路径
	ChecksumLength = 4                         // 用于验证数据完整性的校验和长度
	NetworkVersion = byte(0x00)                // 网络版本号，用于版本控制
	Wallets        = "./files/wallets/"        // 钱包文件存储目录的路径
	WalletsRefList = "./files/ref_list/"       // 钱包引用列表文件存储的目录路径
)

// Identity 枚举身份角色
type Identity string

const (
	Raw      Identity = "Raw"
	Producer Identity = "Producer"
	Dealer   Identity = "Dealer"
	User     Identity = "User"
)

// 日志相关
// 使用ANSI颜色代码来区分不同级别的日志
const (
	preError = "\033[31m[Error]\033[0m"
	preInfo  = "\033[34m[Info]\033[0m"
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

// FileExists 检查指定路径的文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

// Err 处理err,若存在则打印日志
func Err(err error) {
	if err != nil {
		errorln(err)
	}
}

// Info 处理info,打印日志
func Info(s string) {
	infoln(s)
}

// PublicKeyHash hash公钥
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hash := ripemd160.New()
	_, err := hash.Write(hashedPublicKey[:])
	Err(err)
	publicRipeMd := hash.Sum(nil)
	return publicRipeMd
}

// GenNeKeyPair 生成密钥对
func GenNeKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	Err(err)
	// 公钥为点坐标，需要拼接起来保存
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

// CheckSum 检查位生成函数
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:ChecksumLength]
}

// Base58Encode base256转base58
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

// Base58Decode base58转base256
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Err(err)
	return decode
}

// PublicHashToAddress 从公钥哈希生成钱包地址
func PublicHashToAddress(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

// AddressToPublicHash 从钱包地址转公钥哈希
func AddressToPublicHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLength]
	return pubKeyHash
}

// CleanData 清空数据
func CleanData() {
	emptyDir(BCPath + "/")
	emptyDir(Wallets)
	emptyDir(WalletsRefList)
}

func emptyDir(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := os.RemoveAll(filepath.Join(dirPath, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
