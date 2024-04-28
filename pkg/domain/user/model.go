package user

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true"`
	Username string    `json:"username" gorm:"unique;not null" validate:"required"`
	Nickname string    `json:"nickname" validate:"required"`
	Password string    `json:"password" validate:"required"`
	Avatar   string    `json:"avatar"`
	Email    string    `json:"email" validate:"required"`
	Phone    string    `json:"phone" validate:"required"`

	database.BaseModel
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

type Group struct {
	ID     uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title  string    `json:"title" gorm:"unique;not null" validate:"required"`
	Remark string    `json:"remark"`

	database.BaseModel
}

func (g *Group) TableName() string {
	return "user_group"
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	g.ID = uuid.New()
	return nil
}

type Relation struct {
	UserID  uuid.UUID `json:"user_id" gorm:"not null" validate:"required"`
	GroupID uuid.UUID `json:"group_id" gorm:"not null" validate:"required"`
}

func (r *Relation) TableName() string {
	return "user_group_relation"
}
