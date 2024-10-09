package database

import (
	"fmt"
	"github.com/MR5356/aurora/internal/infrastructure/cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
)

type CachedMapper[T any] struct {
	BaseMapper[T]
	cache cache.Cache
}

func NewCachedMapper[T any](db *Database, model T, cache cache.Cache) *CachedMapper[T] {
	return &CachedMapper[T]{
		BaseMapper: *NewMapper(db, model),
		cache:      cache,
	}
}

func (m *CachedMapper[T]) Insert(entity T, db ...*gorm.DB) error {
	err := m.BaseMapper.Insert(entity, db...)
	if err != nil {
		return err
	}
	key := generateKey(entity)
	logrus.Infof("key: %+v", key)
	_ = m.cache.Del(key)
	return nil
}

func (m *CachedMapper[T]) Delete(entity T, db ...*gorm.DB) error {
	err := m.BaseMapper.Delete(entity, db...)
	if err != nil {
		return err
	}
	key := generateKey(entity)
	_ = m.cache.Del(key)
	return nil
}

func (m *CachedMapper[T]) Update(entity T, fields map[string]interface{}, db ...*gorm.DB) error {
	err := m.BaseMapper.Update(entity, fields, db...)
	if err != nil {
		return err
	}
	key := generateKey(entity)
	_ = m.cache.Del(key)
	return nil
}

func (m *CachedMapper[T]) Detail(entity T) (res T, err error) {
	key := generateKey(entity)
	logrus.Infof("key: %+v", key)
	if cachedValue, err := m.cache.Get(key); err == nil && cachedValue != nil {
		logrus.Infof("detail key from cache")
		res = cachedValue.(T)
		return res, err
	}
	res, err = m.BaseMapper.Detail(entity)
	if err != nil {
		return res, err
	}
	_ = m.cache.Set(key, res)
	return
}

// 生成唯一的缓存键
func generateKey[T any](entity T) string {
	v := reflect.ValueOf(entity)

	// 如果是指针，获取指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID") // 假设实体有 ID 字段
	if idField.IsValid() {
		return fmt.Sprintf("entity:%v", idField.Interface())
	}
	return ""
}
