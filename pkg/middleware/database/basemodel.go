package database

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeleteAt  gorm.DeletedAt `json:"-"`
}
