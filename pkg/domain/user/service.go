package user

import (
	"fmt"
	"github.com/MR5356/aurora/pkg/domain/authentication"
	"github.com/MR5356/aurora/pkg/domain/notify"
	"github.com/MR5356/aurora/pkg/domain/user/oauth"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
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
func (s *Service) DeleteUser(userID string) error {
	if err := s.userDB.Delete(&User{ID: userID}); err != nil {
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
func (s *Service) DetailUser(userID string) (*User, error) {
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

// GetOAuthURL get oauth url
func (s *Service) GetOAuthURL(authType string, redirectURL string) (string, error) {
	if provider, err := oauth.GetOAuthManager().GetAuthProvider(authType); err != nil {
		return "", err
	} else {
		return provider.GetAuthURL(redirectURL), nil
	}
}

// GetUserInfo get user info
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

	_, err = s.userDB.Detail(&User{ID: user.ID})
	if err != nil {
		err = s.userDB.Insert(user)
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
			Title:  "admin",
			Remark: "admin group",
		},
	}

	for _, adminGroup := range adminGroups {
		if err := s.groupDB.Insert(adminGroup); err == nil {
			_, _ = authentication.GetPermission().AddPolicyForRoleInDomain("*", adminGroup.ID.String(), "*", "*")
		}
	}

	return nil
}
