package domain

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // "-" agar password tidak muncul di JSON
	Roles    string    `json:"roles"` // admin, staff
	CreatedAt time.Time `json:"created_at"`
}
