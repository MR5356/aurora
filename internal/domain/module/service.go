package module

import (
	"context"
	"errors"
	"github.com/MR5356/aurora/internal/config"
	"github.com/MR5356/aurora/internal/domain/user"
	"github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v72/github"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	moduleDB                 *database.BaseMapper[*Module]
	installationIDRelationDB *database.BaseMapper[*InstallationIDRelation]
	userDB                   *database.BaseMapper[*user.User]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			moduleDB:                 database.NewMapper(database.GetDB(), &Module{}),
			installationIDRelationDB: database.NewMapper(database.GetDB(), &InstallationIDRelation{}),
			userDB:                   database.NewMapper(database.GetDB(), &user.User{}),
		}
	})
	return service
}

func (s *Service) PageModule(ctx context.Context, page, size int, owner string) (*database.Pager[*Module], error) {
	iids, err := s.installationIDRelationDB.List(&InstallationIDRelation{Owner: owner})
	if err != nil {
		return nil, err
	}
	var installationIDs []int64
	for _, iid := range iids {
		installationIDs = append(installationIDs, iid.InstallationID)
	}

	res := new(database.Pager[*Module])
	res.CurrentPage = int64(page)
	res.PageSize = int64(size)
	database.GetDB().Model(&Module{}).Where("installation_id in ?", installationIDs).Count(&res.Total)
	if res.Total == 0 {
		res.Data = make([]*Module, 0)
		return res, nil
	}

	if err := database.GetDB().Where("installation_id in ?", installationIDs).Order("created_at desc").Scopes(database.Pagination(res)).Find(&res.Data).Error; err != nil {
		logrus.Errorf("PageModule: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *Service) UpdateGithubModule(ctx context.Context, action string, installationID int64) error {
	if action == "deleted" {
		if err := database.GetDB().Model(&Module{}).Where("installation_id = ?", installationID).Delete(&Module{}).Error; err != nil {
			return err
		}
		if err := database.GetDB().Model(&InstallationIDRelation{}).Where("installation_id = ?", installationID).Delete(&InstallationIDRelation{}).Error; err != nil {
			return err
		}
		return nil
	}
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Module{}).Where("installation_id = ?", installationID).Delete(&Module{}).Error; err != nil {
			return err
		}

		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.Current().GithubApp.AppID, installationID, config.Current().GithubApp.PrivateKey)
		if err != nil {
			logrus.Errorf("ghinstallation.NewKeyFromFile: %v", err)
			return err
		}

		client := github.NewClient(&http.Client{Transport: itr})

		all, err := listAllRepos(ctx, client)
		if err != nil {
			logrus.Errorf("listAllRepos: %v", err)
			return err
		}

		for _, repo := range all {
			logrus.Debugf("repo: %+v", repo)
			if err = s.moduleDB.Insert(&Module{
				Name:           repo.GetName(),
				SCMType:        "GitHub",
				Owner:          repo.GetOwner().GetLogin(),
				OwnerID:        repo.GetOwner().GetID(),
				Description:    repo.GetDescription(),
				Language:       repo.GetLanguage(),
				Private:        repo.GetPrivate(),
				HtmlURL:        repo.GetHTMLURL(),
				CloneURL:       repo.GetCloneURL(),
				SSHURL:         repo.GetSSHURL(),
				SVNURL:         repo.GetSVNURL(),
				InstallationID: installationID,
			}, tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) RegisterInstallationID(ctx context.Context, installationID int64, owner string) error {
	r, err := s.installationIDRelationDB.Detail(&InstallationIDRelation{InstallationID: installationID})

	if err == nil {
		r.Owner = owner
		return s.installationIDRelationDB.Update(r, structutil.Struct2Map(r))
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		installation := &InstallationIDRelation{
			InstallationID: installationID,
			Owner:          owner,
		}
		return s.installationIDRelationDB.Insert(installation)
	} else {
		return err
	}
}

func listAllRepos(ctx context.Context, client *github.Client) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	opts := &github.ListOptions{PerPage: 100} // 每页最多100个

	for {
		repos, resp, err := client.Apps.ListRepos(ctx, opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos.Repositories...)

		if resp.NextPage == 0 {
			break // 没有更多页
		}
		opts.Page = resp.NextPage
	}
	return allRepos, nil
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Module{}, &InstallationIDRelation{}); err != nil {
		return err
	}
	return nil
}
