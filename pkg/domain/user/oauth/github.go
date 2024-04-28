package oauth

import (
	"github.com/MR5356/aurora/pkg/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	config *oauth2.Config
}

func NewGithubProvider(conf config.OAuthConfig) *GithubProvider {
	return &GithubProvider{
		config: &oauth2.Config{
			Scopes: []string{"user:email", "read:user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			ClientID:     conf.ClientId,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  config.Current().Server.BaseURL + config.Current().Server.Prefix + "/user/callback?authType=github",
		},
	}
}

func (p *GithubProvider) GetAuthURL(redirectURL string) string {
	return p.config.AuthCodeURL(redirectURL, oauth2.AccessTypeOffline)
}

func (p *GithubProvider) GetInfo(code string) (*UserInfo, error) {
	logrus.Debugf("GetInfo.code: %s", code)

	token, err := p.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		logrus.Errorf("get token failed, error: %v", err)
		return nil, err
	}
	
	logrus.Debugf("GetInfo.token: %v", token)
	return &UserInfo{}, nil
}
