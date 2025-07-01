package module

import (
	"time"
)

type Module struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"unique;not null;length:64" json:"name"`
	Owner          string `gorm:"not null;length:64" json:"owner"`
	OwnerID        int64  `gorm:"not null" json:"ownerID"`
	SCMType        string `gorm:"not null;length:32" json:"scmType"`
	Description    string `gorm:"length:1024" json:"description"`
	Language       string `gorm:"length:128" json:"language"`
	Private        bool   `json:"private"`
	HtmlURL        string `gorm:"length:256" json:"htmlURL"`
	CloneURL       string `gorm:"length:256" json:"cloneURL"`
	SSHURL         string `gorm:"length:256" json:"sshURL"`
	SVNURL         string `gorm:"length:256" json:"svnURL"`
	InstallationID int64  `gorm:"not null" json:"installationID"`

	CreatedAt time.Time `json:"createdAt" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updatedAt" swaggerignore:"true"`
}

type InstallationIDRelation struct {
	InstallationID int64  `gorm:"primaryKey;not null" json:"installationID"`
	Owner          string `gorm:"not null;length:64" json:"owner"`
}
