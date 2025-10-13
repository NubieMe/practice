package models

type User struct {
	Name     string `gorm:"column:name;size:255" json:"name"`
	Email    string `gorm:"column:email;unique;size:255" json:"email"`
	Password string `gorm:"column:password;size:255" json:"-"`

	// Todos []Todo `gorm:"foreignKey:UserID" json:"todos"`

	Base
}

type UserRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
