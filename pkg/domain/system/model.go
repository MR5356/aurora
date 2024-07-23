package system

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Record struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	UserID    string    `json:"userId" gorm:"not null;default:'unknown'"`
	Path      string    `json:"path" gorm:"not null;"`
	Method    string    `json:"method" gorm:"not null;"`
	Code      string    `json:"code" gorm:"not null;"`
	ClientIP  string    `json:"clientIp" gorm:"not null;"`
	UserAgent string    `json:"userAgent" gorm:"not null;"`
	Cost      int64     `json:"cost" gorm:"not null;"`
	IsApi     bool      `json:"-"`

	database.BaseModel
}

func (r *Record) TableName() string {
	return "system_record"
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

type Statistic struct {
	Name  string `json:"name"`
	Count string `json:"count"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
}

type Version struct {
	Version       string `json:"version"`
	LatestVersion string `json:"latestVersion"`
	LatestInfo    string `json:"latestInfo"`
	LatestUrl     string `json:"latestUrl"`
}
