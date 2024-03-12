package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	FirstName string `json:"first_name" gorm:"type:char(50);not null"`
	LastName  string `json:"last_name" gorm:"type:char(50);not null"`
	Email     string `json:"email" gorm:"type:char(50);not null;unique"`
	Phone     string `json:"phone" gorm:"type:char(30);not null"`
	//Course    *Course        `gorm:"-"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	Deleted   gorm.DeletedAt `json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return
}
