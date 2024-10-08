package authentication

import (
	"errors"
	"github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	once       sync.Once
	permission *Permission
)

type Permission struct {
	enforcer casbin.IEnforcer
}

func GetPermission() *Permission {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(database.GetDB().DB)
		if err != nil {
			logrus.Fatalf("Failed to initialize database: %v", err)
		}

		m, err := model.NewModelFromString(modelDefine)
		if err != nil {
			logrus.Fatalf("Failed to new model: %v", err)
		}

		e, err := casbin.NewSyncedCachedEnforcer(m, adapter)
		if err != nil {
			logrus.Fatalf("Failed to new enforcer: %v", err)
		}
		permission = &Permission{enforcer: e}
	})

	return permission
}

func (p *Permission) HasPermissionForRoleInDomain(domain, role, object, action string) (bool, error) {
	logrus.Debugf("HasPermissionForRoleInDomain: %s, %s, %s, %s", domain, role, object, action)
	return p.enforcer.Enforce(role, domain, object, action)
}

func (p *Permission) GetPolicyForRoleInDomain(domain, role string) ([][]string, error) {
	logrus.Debugf("GetPolicyForRoleInDomain: %s, %s", domain, role)
	return p.enforcer.GetPermissionsForUser(role, domain)
}

func (p *Permission) GetPolicyForObjectInDomain(domain, object string) [][]string {
	logrus.Debugf("GetPolicyForObjectInDomain: %s, %s", domain, object)
	res := make([][]string, 0)
	policies := p.enforcer.GetPolicy()
	for _, policy := range policies {
		if policy[1] == domain && policy[2] == object {
			res = append(res, policy)
		}
	}
	return res
}

func (p *Permission) RemovePoliciesForObjectInDomain(domain, object string) (bool, error) {
	logrus.Debugf("RemovePoliciesForObjectInDomain: %s, %s", domain, object)
	res := make([][]string, 0)
	policies := p.enforcer.GetPolicy()
	for _, policy := range policies {
		if policy[1] == domain && policy[2] == object {
			res = append(res, policy)
		}
	}
	return p.RemovePolicies(res)
}

func (p *Permission) AddPolicyForRoleInDomain(domain, role, object, action string) (bool, error) {
	logrus.Debugf("AddPolicyForRoleInDomain: %s, %s, %s, %s", domain, role, object, action)
	return p.AddPolicies([][]string{
		{role, domain, object, action},
	})
}

func (p *Permission) RemovePolicyForRoleInDomain(domain, role, object, action string) (bool, error) {
	logrus.Debugf("RemovePolicyForRoleInDomain: %s, %s, %s, %s", domain, role, object, action)
	return p.RemovePolicies([][]string{
		{role, domain, object, action},
	})
}

func (p *Permission) AddPolicies(rules [][]string) (bool, error) {
	logrus.Debugf("AddPolicies: %+v", rules)
	return p.enforcer.AddPolicies(rules)
}

func (p *Permission) RemovePolicies(rules [][]string) (bool, error) {
	logrus.Debugf("RemovePolicies: %+v", rules)
	return p.enforcer.RemovePolicies(rules)
}

func (p *Permission) FilterDataArray(data []any, actions []string, domain string, roles []string, fieldName string) ([]any, error) {
	res := make([]any, 0)
	for _, d := range data {
		for _, role := range roles {
			ok := false
			var err error
			for _, action := range actions {
				obj, objOk := structutil.GetMapFiledByName(d.(map[string]any), fieldName)
				if !objOk {
					continue
				}
				ok, err = p.HasPermissionForRoleInDomain(domain, role, obj.(string), action)
				if err != nil {
					logrus.Errorf("has permission for role in domain failed, error: %v", err)
					return make([]any, 0), errors.New("has permission for role in domain failed")
				}

				if ok {
					res = append(res, d)
				}
			}
		}
	}
	return res, nil
}
