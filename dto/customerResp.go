package dto

type CustAdminResponse struct {
	ID             uint64 `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FireUid        string `json:"fire-uid"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Address        string `json:"address"`
	Role           string `json:"role"`
	WalletId       string `json:"wallet-id"`
	CreateDate     string `json:"create_date"`
	CreateTime     string `json:"create_time"`
	CreatedBy      string `json:"created_by"`
	LastChangeBy   string `json:"last_change_by"`
	LastUpdateDate string `json:"last_update_date"`
	LastUpdateTime string `json:"last_update_time"`
	Active         bool   `json:"active"`
}

type CustClientResp struct {
	ID       uint64      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FireUid  string      `json:"fire-uid"`
	Name     string      `json:"name"`
	Phone    string      `json:"phone"`
	Email    string      `json:"email"`
	Address  string      `json:"address"`
	Role     string      `json:"role"`
	WalletId string      `json:"wallet-id"`
	Active   bool        `json:"active"`
	Wallet   interface{} `json:"wallet"`
}
