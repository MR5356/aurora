package oauth

import (
	"context"
	"fmt"
	"github.com/MR5356/aurora/internal/config"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"net/url"
)

type GitlabAuth struct {
	config    *oauth2.Config
	oAuthName string
}

func NewGitlabAuth(oAuthName string, conf config.OAuthConfig) *GitlabAuth {
	return &GitlabAuth{
		config: &oauth2.Config{
			Scopes: []string{"read_user", "email", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  conf.AuthURL,
				TokenURL: conf.TokenURL,
			},
			ClientID:     conf.ClientId,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  config.Current().Server.BaseURL + "/api/v1/user/callback?oauth=" + oAuthName,
		},
		oAuthName: oAuthName,
	}
}

func (p *GitlabAuth) GetAuthURL(redirectURL string) string {
	return p.config.AuthCodeURL(redirectURL, oauth2.AccessTypeOffline)
}

func (p *GitlabAuth) GetInfo(code string) (*UserInfo, error) {
	logrus.Debugf("GetInfo.code: %s", code)
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		logrus.Errorf("GetInfo.config exchange error: %+v", err)
		return nil, ErrAuthFailed
	}
	logrus.Debugf("token: %+v", token)

	parsedURL, err := url.Parse(p.config.Endpoint.AuthURL)
	if err != nil {
		return nil, ErrAuthFailed
	}
	extractedURL := url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	}

	git, err := gitlab.NewOAuthClient(token.AccessToken, gitlab.WithBaseURL(extractedURL.String()))
	if err != nil {
		logrus.Errorf("GetInfo.Create gitlab client error: %+v", err)
		return nil, ErrAuthFailed
	}

	user, _, err := git.Users.CurrentUser()
	if err != nil {
		logrus.Errorf("GetInfo.Git.Users.CurrentUser error: %+v", err)
		return nil, ErrAuthFailed
	}
	logrus.Debugf("user: %+v", structutil.Struct2String(user))

	userInfo := new(UserInfo)
	userInfo.ID = fmt.Sprintf("%d.%s", user.ID, p.oAuthName)
	userInfo.Username = user.Username
	userInfo.Nickname = user.Name
	userInfo.Avatar = user.AvatarURL
	userInfo.Email = user.Email
	userInfo.UserType = AuthTypeGitlab

	logrus.Debugf("userInfo: %+v", structutil.Struct2String(userInfo))
	return userInfo, nil
}
