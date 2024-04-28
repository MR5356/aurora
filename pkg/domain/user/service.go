package user

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	userDB  *database.BaseMapper[*User]
	groupDB *database.BaseMapper[*Group]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			userDB:  database.NewMapper(database.GetDB(), &User{}),
			groupDB: database.NewMapper(database.GetDB(), &Group{}),
		}
	})
	return service
}

// AddUser add user
func (s *Service) AddUser(user *User) error {
	if err := validate.Validate(user); err != nil {
		logrus.Errorf("validate user failed, error: %v", err)
		return err
	}

	if err := s.userDB.Insert(user); err != nil {
		logrus.Errorf("insert user failed, error: %v", err)
		return err
	}
	return nil
}

// DeleteUser delete user
func (s *Service) DeleteUser(user *User) error {
	if err := s.userDB.Delete(user); err != nil {
		logrus.Errorf("delete user failed, error: %v", err)
		return err
	}
	return nil
}

// UpdateUser update user
func (s *Service) UpdateUser(user *User) error {
	if err := validate.Validate(user); err != nil {
		logrus.Errorf("validate user failed, error: %v", err)
		return err
	}

	if err := s.userDB.Update(&User{ID: user.ID}, structutil.Struct2Map(user)); err != nil {
		logrus.Errorf("update user failed, error: %v", err)
		return err
	}
	return nil
}

// DetailUser detail user
func (s *Service) DetailUser(userID uuid.UUID) (*User, error) {
	if res, err := s.userDB.Detail(&User{ID: userID}); err != nil {
		logrus.Errorf("detail user failed, error: %v", err)
		return nil, err
	} else {
		return res, err
	}
}

// ListUser list user
func (s *Service) ListUser(user *User) ([]*User, error) {
	return s.userDB.List(user)
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&User{}, &Group{}, &Relation{}); err != nil {
		return err
	}
	return nil
}
