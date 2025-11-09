package Models

import "time"

// User models struct
type User struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	UUID        string `gorm:"uniqueIndex"`
	Fullname    string
	Email       string `gorm:"uniqueIndex"`
	Password    string
	PhoneNumber string
	Role        string
	IsActive    bool `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
