package wallet

import (
	"blockchain/util"
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 记录钱包信息
// 从地址到别名的映射
type RefList map[string]string

// 为钱包设定别名
func (r *RefList) SetRef(address, refname string) {
	(*r)[address] = refname
}

// 通过别名获取钱包
func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("通过别名未找到钱包...")
		return temp, err
	}
	return temp, nil
}

// 保存RefList
func (r *RefList) Save() {
	filename := util.WalletsRefList + "ref_list"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	util.Err(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	util.Err(err)
}

// 扫描所有保存的钱包文件
func (r *RefList) Update() {
	err := filepath.Walk(util.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	util.Err(err)
}

// 加载RefList信息
func LoadRefList() *RefList {
	filename := util.WalletsRefList + "ref_list"
	var reflist RefList
	if util.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		util.Err(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		util.Err(err)
	} else {
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}
