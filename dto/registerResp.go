package dto

type RegisterResponse struct {
	UID   string `json:"uid" gorm:"not null"`
	Phone string `json:"phone" gorm:"not null"`
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
}
