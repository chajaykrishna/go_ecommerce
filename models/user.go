package models

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Name      string `json:"name" gorm:"not null"`
	Address   string `json:"address" gorm:"not null"`
	Phone     string `json:"phone" gorm:"unique"`
	Password  string `json:"password" gorm:"not null"`
	CreatedAt string `json:"createdAt"`
}
