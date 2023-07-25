package dto

import "github.com/shopspring/decimal"

type NewWallet struct {
	UserId       string `json:"user-id"`
	UserType     string `json:"merch-type"`
	LastChangeBy string `json:"last_change_by"`
}

type EditWallet struct {
	UserID       uint64          `json:"user_id"`
	PhoneNumber  string          `json:"phone_number"`
	EWalletType  string          `json:"e_wallet"`
	Operand      bool            `json:"operand"`
	Amount       decimal.Decimal `json:"amount"`
	LastChangeBy string          `json:"last_change_by"`
}

type TopUpWallet struct {
	WalletID    uint64  `json:"wallet_id"`
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	Discount    float64 `json:"discount"`
}
