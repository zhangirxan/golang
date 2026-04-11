package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `json:"DeletedAt" gorm:"index"`
	Username  string         `json:"Username" gorm:"uniqueIndex"`
	Email     string         `json:"Email" gorm:"uniqueIndex"`
	Password  string         `json:"Password"`
	Role      string         `json:"Role"`
	Verified  bool           `json:"Verified"`
}
