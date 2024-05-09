package main

type BlockInfo struct {
	Timestamp    string
	PreviousHash string
	Trades       []TradeInfo
	Hash         string
	Pow          bool
}

type TradeInfo struct {
	ID          string
	Inputs      []InputInfo
	Outputs     []OutputInfo
	Description string
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
	Identity      string
}

type UpdateWalletsResult struct {
	Message string
}

type WalletsListResult struct {
	Wallets []WalletInfoResult
}

type getAllBalanceResult struct {
	A_dealer_balance   string
	B_dealer_balance   string
	C_dealer_balance   string
	Raw_balance        string
	A_producer_balance string
	User_balance       string
}

type TraceTrade struct {
	Time        string
	Description string
}

type BuyResult struct {
	Success     bool
	Message     string
	TraceTrades []TraceTrade
}

type ProducerBuyResult struct {
	Success bool
	Message string
}

type DealerBuyResult struct {
	Success bool
	Message string
}
