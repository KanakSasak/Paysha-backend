package dto

type RegisterRequest struct {
	Phone    string `json:"phone" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}
