package pluginhub

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
)

type Plugin struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" example:"00000000-0000-0000-0000-000000000000"`
	Label       string    `json:"label" validate:"required"`
	Abstract    string    `json:"abstract"`
	Author      string    `json:"author"`
	DownloadUrl string    `json:"downloadUrl"`
	ProjectUrl  string    `json:"projectUrl"`
	Icon        string    `json:"icon"`
	Version     string    `json:"version"`
	Usage       string    `json:"usage"`

	database.BaseModel
}
