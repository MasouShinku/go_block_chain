package main

type BlockInfo struct {
	Timestamp    string
	PreviousHash string
	Trades       []TradeInfo
	Hash         string
	Pow          bool
}

type TradeInfo struct {
	ID      string
	Inputs  []InputInfo
	Outputs []OutputInfo
}

type InputInfo struct {
	TradeID string
	OutID   int
	PubKey  string
}

type OutputInfo struct {
	Num        int
	HashPubKey string
}

type BlockchainCreationResult struct {
	Success bool
	Message string
}

type BalanceResult struct {
	Address string
	Balance int
}

type TradeResult struct {
	Success bool
	Message string
}

type MiningResult struct {
	Message string
}

type CreateWalletResult struct {
	Message string
}

type WalletInfoResult struct {
	Address       string
	PublicKey     string
	ReferenceName string
}

type UpdateWalletsResult struct {
	Message string
}

type WalletsListResult struct {
	Wallets []WalletInfoResult
}
