package main

import (
	"fmt"
	"github.com/MR5356/aurora/pkg/domain/runner/proto"
	"github.com/MR5356/aurora/pkg/domain/runner/shared"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type CheckoutTask struct {
	repository string
	branch     string
	submodules bool

	shared.UnimplementedITask
}

func (t *CheckoutTask) GetInfo() *proto.TaskInfo {
	return &proto.TaskInfo{
		Label:       "Checkout",
		Abstract:    "Checkout a Git repository at a particular version",
		Author:      "Rui Ma",
		DownloadUrl: "",
		ProjectUrl:  "",
		Icon:        "",
		Version:     "v1.0.0",
		Usage:       "",
	}
}

func (t *CheckoutTask) GetParams() *proto.TaskParams {
	return &proto.TaskParams{
		Params: []*proto.TaskParam{
			{
				Title:       "Repository",
				Placeholder: "repository address",
				Order:       1,
				Type:        "string",
				Required:    true,
				Key:         "repository",
				Value:       "",
			},
			{
				Title:       "Branch",
				Placeholder: "repository branch",
				Order:       2,
				Type:        "string",
				Required:    true,
				Key:         "branch",
				Value:       "",
			},
			{
				Title:       "Submodules",
				Placeholder: "whether to download submodules",
				Order:       3,
				Type:        "switch",
				Required:    false,
				Key:         "submodules",
				Value:       "false",
			},
			{
				Title:       "Token",
				Placeholder: "token",
				Order:       4,
				Type:        "string",
				Required:    false,
				Key:         "token",
				Value:       "",
			},
		},
	}
}

func (t *CheckoutTask) SetParams(params *proto.TaskParams) {}

func (t *CheckoutTask) Start() error {
	srcAuth, repoUrl, err := getAuth(t.repository, "", "")
	if err != nil {
		return err
	}

	dirName := fmt.Sprintf("/tmp/%s", strings.ReplaceAll(filepath.Base(repoUrl), ".git", ""))
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		_ = os.RemoveAll(dirName)
	}

	shared.Logger.Debug("clone %s to %s", repoUrl, dirName)

	co := &git.CloneOptions{
		URL:               repoUrl,
		Auth:              srcAuth,
		ReferenceName:     plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", t.branch)),
		RecurseSubmodules: git.NoRecurseSubmodules,
		InsecureSkipTLS:   true,
	}

	repo, err := git.PlainClone(dirName, false, co)
	if err != nil {
		return err
	}

	if t.submodules {
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}
		submodules, err := worktree.Submodules()
		if err != nil {
			return err
		}

		for _, sm := range submodules {
			logrus.Debugf("clone submodule: %s", sm.Config().URL)
			smAuth, _, _ := getAuth(sm.Config().URL, "", "")
			err = sm.Update(&git.SubmoduleUpdateOptions{
				Init: true,
				Auth: smAuth,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *CheckoutTask) Stop() error {
	return nil
}

func (t *CheckoutTask) Pause() error {
	return nil
}

func (t *CheckoutTask) Resume() error {
	return nil
}

func getAuth(repo, privateKeyFile, privateKeyPassword string) (auth transport.AuthMethod, repoUrl string, err error) {
	/**
	支持以下形式：
		1. https://github.com/MR5356/syncer.git
		2. git@github.com:MR5356/syncer.git
		3. https://username:password@github.com/MR5356/syncer.git
		4. https://<token>@github.com/MR5356/syncer.git
		5. https://oauth2:access_token@github.com/MR5356/syncer.git
	*/

	repoUrl = repo
	switch getUrlType(repo) {
	case gitUrlType:
		if privateKeyFile == "" {
			u, err := user.Current()
			if err == nil {
				privateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
			}
		}
		_, err = os.Stat(privateKeyFile)
		if err != nil {
			return auth, repoUrl, err
		}
		shared.Logger.Debug("privateKeyFile: %s, privateKeyPassword: %s", privateKeyFile, privateKeyPassword)
		auth, err = ssh.NewPublicKeysFromFile("git", privateKeyFile, privateKeyPassword)
		return auth, repoUrl, err
	case httpUrlType:
		auth = nil
		return
	case tokenizedHttpUrlType:
		token := strings.ReplaceAll(strings.ReplaceAll(strings.Split(repo, "@")[0], "https://", ""), "http://", "")
		shared.Logger.Info(token)
		auth = &http.TokenAuth{
			Token: token,
		}
		repoUrl = strings.ReplaceAll(repo, token+"@", "")
		return
	case basicHttpUrlType:
		basicInfo := strings.ReplaceAll(strings.ReplaceAll(strings.Split(repo, "@")[0], "https://", ""), "http://", "")
		fields := strings.Split(basicInfo, ":")
		auth = &http.BasicAuth{
			Username: fields[0],
			Password: fields[1],
		}
		repoUrl = strings.ReplaceAll(repo, basicInfo+"@", "")
		return
	default:
		return nil, "", fmt.Errorf("unsupported repo url: %s", repo)
	}
}

type urlType int

const (
	unknownUrlType urlType = iota
	gitUrlType
	httpUrlType
	tokenizedHttpUrlType
	basicHttpUrlType
)

var (
	isGitUrl           = regexp.MustCompile(`^git@[-\w.:]+:[-\/\w.]+\.git$`)
	isHttpUrl          = regexp.MustCompile(`^(https|http)://[-\w.:]+/[-\/\w.]+\.git$`)
	isTokenizedHttpUrl = regexp.MustCompile(`^(https|http)://[a-zA-Z0-9_]+@[-\w.:]+/[-\/\w.]+\.git$`)
	isBasicHttpUrl     = regexp.MustCompile(`^(https|http)://[a-zA-Z0-9]+:[\w]+@[-\w.:]+/[-\/\w.]+\.git$`)
)

func getUrlType(url string) (t urlType) {
	if isGitUrl.MatchString(url) {
		t = gitUrlType
	} else if isHttpUrl.MatchString(url) {
		t = httpUrlType
	} else if isTokenizedHttpUrl.MatchString(url) {
		t = tokenizedHttpUrlType
	} else if isBasicHttpUrl.MatchString(url) {
		t = basicHttpUrlType
	} else {
		t = unknownUrlType
	}
	shared.Logger.Debug("getUrlType: %v", t)
	return t
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	//shared.Serve(&CheckoutTask{})
	task := &CheckoutTask{
		branch:     "master",
		repository: "git@github.com:MR5356/aurora.git",
		submodules: true,
	}

	err := task.Start()
	if err != nil {
		panic(err)
	}
}
