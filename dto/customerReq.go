package dto

type NewCustomer struct {
	FireUid   string `json:"fire-uid"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	CreatedBy string `json:"created-by"`
}

type EditCustomer struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	LastChangeBy string `json:"last_change_by"`
	WalletId     string `json:"wallet-id"`
	Role         string `json:"role"`
}
