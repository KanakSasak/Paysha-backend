package dto

import "github.com/shopspring/decimal"

type WalletCliResponse struct {
	WalletID uint64          `json:"wallet_id"`
	Amount   decimal.Decimal `json:"amount" gorm:"type:numeric"`
	UserType string          `json:"merch_type"`
}
