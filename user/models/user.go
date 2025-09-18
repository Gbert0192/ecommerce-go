package models

type RegisterParameter struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginParameter struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(50);default:'user'" json:"role"`
}
