package user

import (
	"errors"
	"fmt"
	"github.com/MR5356/aurora/pkg/domain/authentication"
	"github.com/MR5356/aurora/pkg/domain/notify"
	"github.com/MR5356/aurora/pkg/domain/user/oauth"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	userDB     *database.BaseMapper[*User]
	groupDB    *database.BaseMapper[*Group]
	relationDB *database.BaseMapper[*Relation]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			userDB:     database.NewMapper(database.GetDB(), &User{}),
			groupDB:    database.NewMapper(database.GetDB(), &Group{}),
			relationDB: database.NewMapper(database.GetDB(), &Relation{}),
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
	if user.Type == TypeLocal && len(user.Password) == 0 {
		return errors.New("password is required")
	}

	tx := s.userDB.DB.Begin()
	defer tx.Rollback()

	if err := s.userDB.Insert(user, tx); err != nil {
		logrus.Errorf("insert user failed, error: %v", err)
		return err
	}

	if count, err := s.userDB.Count(&User{}); err != nil {
		return err
	} else {
		// insert default admin group
		if count == 0 {
			relation := &Relation{
				UserID:  user.ID,
				GroupID: uuid.MustParse(defaultAdminGroupID),
			}

			if err := s.relationDB.Insert(relation, tx); err != nil {
				logrus.Errorf("insert user relation failed, error: %v", err)
				return err
			}
		}
	}
	tx.Commit()

	return nil
}

// DeleteUser delete user
func (s *Service) DeleteUser(userID string) error {
	if err := s.userDB.Delete(&User{ID: userID}); err != nil {
		logrus.Errorf("delete user failed, error: %v", err)
		return err
	}
	return nil
}

func (s *Service) SetUserStatus(user *User, status int) error {
	return s.userDB.DB.Model(&User{}).Where(user).Update("status", status).Error
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
func (s *Service) DetailUser(userID string) (*User, error) {
	if res, err := s.userDB.Detail(&User{ID: userID}); err != nil {
		logrus.Errorf("detail user failed, error: %v", err)
		return nil, err
	} else {
		return res, err
	}
}

// ListUser list user
func (s *Service) ListUser(user *User) ([]*ListUserResponse, error) {
	var users []*ListUserResponse
	err := s.userDB.DB.Table("user").
		Select("user.*, user_group.title as `group`").
		Joins("left join user_group_relation on user.id = user_group_relation.user_id").
		Joins("left join user_group on user_group_relation.group_id = user_group.id").
		Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) ResetPassword(user *ResetPasswordRequest, IsAdmin bool) error {
	if err := validate.Validate(user); err != nil {
		return err
	}
	if u, err := s.userDB.Detail(&User{ID: user.Username}); err != nil {
		return errors.New("user not exist")
	} else {
		if !IsAdmin && u.Password != user.Old {
			return errors.New("old password not correct")
		} else if u.Type != TypeLocal {
			return errors.New("only local user can reset password")
		} else {
			return s.userDB.DB.Model(&User{}).Where("id = ?", u.ID).Update("password", user.New).Error
		}
	}
}

func (s *Service) Login(user *LoginRequest) (string, error) {
	if err := validate.Validate(user); err != nil {
		return "", err
	}
	if u, err := s.userDB.Detail(&User{Username: user.Username}); err != nil {
		return "", errors.New("username not exist")
	} else {
		if u.Type != TypeLocal {
			return "", errors.New("only local user can login")
		} else if u.Password != user.Password {
			return "", errors.New("password not correct")
		} else if u.Status == StatusInactive {
			return "", errors.New("user is inactive")
		} else if u.Status == StatusBan {
			return "", errors.New("user has been banned")
		} else {
			return GetJWTService().CreateToken(u)
		}
	}
}

// GetOAuthURL get oauth url
func (s *Service) GetOAuthURL(authType string, redirectURL string) (string, error) {
	if provider, err := oauth.GetOAuthManager().GetAuthProvider(authType); err != nil {
		return "", err
	} else {
		return provider.GetAuthURL(redirectURL), nil
	}
}

// GetOAuthUserInfo get user info
func (s *Service) GetOAuthUserInfo(authType string, code string) (*oauth.UserInfo, error) {
	if provider, err := oauth.GetOAuthManager().GetAuthProvider(authType); err != nil {
		return nil, err
	} else {
		return provider.GetInfo(code)
	}
}

func (s *Service) GetOAuthToken(authType string, code string) (token string, err error) {
	userinfo, err := s.GetOAuthUserInfo(authType, code)
	if err != nil {
		return "", err
	}

	user := new(User)
	user.ID = userinfo.ID
	user.Username = userinfo.Username
	user.Nickname = userinfo.Nickname
	user.Email = userinfo.Email
	user.Phone = userinfo.Phone
	user.Avatar = userinfo.Avatar
	user.Status = StatusActive
	user.Type = TypeOAuth

	_, err = s.userDB.Detail(&User{ID: user.ID})
	if err != nil {
		err = s.AddUser(user)
		if err != nil {
			logrus.Errorf("insert user failed, error: %v", err)
			return "", err
		}
	} else {
		err = s.userDB.Update(&User{ID: user.ID}, structutil.Struct2Map(user))
		if err != nil {
			logrus.Errorf("update user failed, error: %v", err)
			return "", err
		}
	}

	u, err := s.userDB.Detail(&User{ID: user.ID})
	if err != nil {
		return "", err
	}

	if err = eventbus.GetEventBus().Publish(notify.TopicSendMessage, &notify.MessageTemplate{
		Event:   notify.EventLogin,
		Subject: "登录通知",
		Body:    fmt.Sprintf("您好，%s，您在%s登录成功，欢迎您。", u.Nickname, time.Now().Format(time.DateTime)),
		Level:   "info",

		Receivers: notify.MessageReceiver{
			Receivers: []string{u.Email},
			Type:      "email",
		},
	}); err != nil {
		logrus.Errorf("publish message failed, error: %v", err)
	}

	return GetJWTService().CreateToken(u)
}

// GetAvailableOAuth get available oauth
func (s *Service) GetAvailableOAuth() []oauth.AvailableOAuth {
	return oauth.GetOAuthManager().GetAvailableOAuth()
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&User{}, &Group{}, &Relation{}); err != nil {
		return err
	}

	// init admin group
	adminGroups := []*Group{
		{
			ID:     uuid.MustParse(defaultAdminGroupID),
			Title:  "admin",
			Remark: "admin group",
		},
	}

	for _, adminGroup := range adminGroups {
		if err := s.groupDB.DB.FirstOrCreate(adminGroup).Error; err == nil {
			_, _ = authentication.GetPermission().AddPolicyForRoleInDomain("*", adminGroup.ID.String(), "*", "*")
		}
	}

	// init admin user
	adminUser := &User{
		ID:       defaultAdminID,
		Username: "admin",
		Nickname: "admin",
		Password: "admin",
		Avatar:   "/logo.svg",
		Status:   StatusActive,
		Type:     TypeLocal,
	}

	if err := s.userDB.DB.FirstOrCreate(adminUser).Error; err != nil {
		return err
	}

	relation := &Relation{
		ID:      uuid.MustParse(defaultRelationID),
		GroupID: uuid.MustParse(defaultAdminGroupID),
		UserID:  adminUser.ID,
	}
	if err := s.relationDB.DB.FirstOrCreate(relation).Error; err != nil {
		return err
	}

	return nil
}
