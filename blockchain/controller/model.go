package main

type BlockInfo struct {
	Timestamp    string      `json:"timestamp"`
	PreviousHash string      `json:"previous_hash"`
	Trades       []TradeInfo `json:"trades"` // 假设Trade是一个已定义的结构体
	Hash         string      `json:"hash"`
	Pow          bool        `json:"pow"`
}

type TradeInfo struct {
	ID      string
	Inputs  []InputInfo
	Outputs []OutputInfo
}

type InputInfo struct {
	TradeID     string
	OutID       int
	FromAddress string
}

type OutputInfo struct {
	Num       int
	ToAddress string
}
