package oauth

import (
	"context"
	"fmt"
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/google/go-github/v61/github"
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
			RedirectURL:  config.Current().Server.BaseURL + config.Current().Server.Prefix + "/user/callback?oauth=github",
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

	client := github.NewClient(nil).WithAuthToken(token.AccessToken)

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logrus.Errorf("get user failed, error: %v", err)
		return nil, err
	}

	logrus.Debugf("GetInfo.user: %+v", structutil.Struct2String(user))
	userInfo := new(UserInfo)
	userInfo.ID = fmt.Sprintf("%d.github", structutil.ValueOfPtr(user.ID, 0))
	userInfo.Nickname = structutil.ValueOfPtr(user.Name, "unknown")
	userInfo.Avatar = structutil.ValueOfPtr(user.AvatarURL, "")
	userInfo.Email = structutil.ValueOfPtr(user.Email, "")
	userInfo.UserType = AuthTypeGithub
	userInfo.Username = structutil.ValueOfPtr(user.Login, "unknown")

	logrus.Debugf("GetInfo.userInfo: %+v", structutil.Struct2String(userInfo))

	return userInfo, nil
}
