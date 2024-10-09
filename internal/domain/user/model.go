package user

import (
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	defaultAdminGroupID = "82599f6f-8ce0-47b4-9c6a-54278cc5e1ae"
	defaultAdminID      = "a0a9f0f9-0b3b-4e9a-9e5d-0a5b5d5a5b5b"
	defaultRelationID   = "b0b9f0f9-0b3b-4e9a-9e5d-0a5b5d5a5b5b"
)

const (
	StatusInactive = iota
	StatusActive
	StatusBan

	StatusInactiveReason = "inactive"
	StatusActiveReason   = "active"
	StatusBanReason      = "banned"
)

const (
	TypeLocal = iota
	TypeOAuth
)

type ListUserResponse struct {
	User
	Group string `json:"group"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	Username string `json:"username" validate:"required"`
	Old      string `json:"old"`
	New      string `json:"new" validate:"required,min=6"`
}

type User struct {
	ID           string `json:"id" gorm:"primary_key;" swaggerignore:"true"`
	Username     string `json:"username" gorm:"unique;not null" validate:"required"`
	Nickname     string `json:"nickname" validate:"required"`
	Password     string `json:"-"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email" validate:"required"`
	Phone        string `json:"phone"`
	Type         int    `json:"type"`
	Status       int    `json:"status"`
	StatusReason string `json:"statusReason"`

	database2.BaseModel
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		u.ID = uuid.New().String()
	}
	return nil
}

func (u *User) IsBanned() bool {
	if r, err := database2.NewMapper(database2.GetDB(), &User{}).Detail(&User{ID: u.ID}); err != nil {
		return true
	} else {
		return r.Status == StatusBan
	}
}

func (u *User) IsAdmin() bool {
	if r, err := database2.NewMapper(database2.GetDB(), &Relation{}).Detail(&Relation{UserID: u.ID}); err != nil {
		return false
	} else {
		return r.GroupID == uuid.MustParse(defaultAdminGroupID)
	}
}

type Uid struct {
	ID string `json:"id" binding:"required"`
}

type Group struct {
	ID     uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title  string    `json:"title" gorm:"unique;not null" validate:"required"`
	Remark string    `json:"remark"`

	database2.BaseModel
}

func (g *Group) TableName() string {
	return "user_group"
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type Relation struct {
	ID      uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	UserID  string    `json:"user_id" gorm:"not null" validate:"required"`
	GroupID uuid.UUID `json:"group_id" gorm:"not null;type:uuid;" validate:"required"`
}

func (r *Relation) TableName() string {
	return "user_group_relation"
}

func (r *Relation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
