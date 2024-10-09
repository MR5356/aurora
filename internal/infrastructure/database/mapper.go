package database

import "gorm.io/gorm"

type Mapper[T any] interface {
	Insert(T, ...*gorm.DB) error
	Delete(T, ...*gorm.DB) error
	Update(T, map[string]interface{}, ...*gorm.DB) error
	Detail(T) (T, error)
	List(T, ...string) ([]T, error)
	Count(T) (int64, error)
	Page(T, int64, int64) (*Pager[T], error)

	GetDB() *gorm.DB
}
