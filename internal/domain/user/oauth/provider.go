package oauth

import (
	"errors"
	"github.com/MR5356/aurora/internal/config"
	"sync"
)

var (
	once    sync.Once
	manager *AuthManager
)

const (
	AuthTypeGithub = "github"
	AuthTypeGitlab = "gitlab"
)

var (
	ErrAuthTypeNotSupport = errors.New("auth type not support")
	ErrAuthFailed         = errors.New("auth failed, please try again")
)

type Provider interface {
	GetInfo(code string) (*UserInfo, error)
	GetAuthURL(redirectURL string) string
}

type AuthManager struct {
	config map[string]config.OAuthConfig
}

func NewOAuthManager(cfg *config.Config) *AuthManager {
	once.Do(func() {
		manager = &AuthManager{
			config: cfg.OAuthConfig,
		}
	})
	return manager
}

func GetOAuthManager() *AuthManager {
	return manager
}

type AvailableOAuth struct {
	OAuth string `json:"oauth" yaml:"oauth"`
	Type  string `json:"type" yaml:"type"`
}

func (m *AuthManager) GetAvailableOAuth() []AvailableOAuth {
	res := make([]AvailableOAuth, 0)
	for k, v := range m.config {
		res = append(res, AvailableOAuth{OAuth: k, Type: v.AuthType})
	}
	return res
}

func (m *AuthManager) GetAuthProvider(authName string) (Provider, error) {
	if conf, ok := m.config[authName]; !ok {
		return nil, ErrAuthTypeNotSupport
	} else {
		switch conf.AuthType {
		case AuthTypeGithub:
			return NewGithubProvider(conf), nil
		case AuthTypeGitlab:
			return NewGitlabAuth(authName, conf), nil
		default:
			return nil, ErrAuthTypeNotSupport
		}
	}
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	UserType string `json:"userType"`
}
