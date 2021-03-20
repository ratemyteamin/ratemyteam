package stores

import (
	"github.com/ratemyteam/rmt/common"
	"sync"
)

type RmtOrganisationStore struct {
	Sql_Store *SqlStore
	orgStore *sync.Mutex
}

func (orgStore *RmtOrganisationStore) GetUserOrganisation(user *common.User) (common.Organisation, error) {
	panic("implement me")
}

func NewOrganisationStore(sqlStore *SqlStore) *RmtOrganisationStore{
	orgStore := &sync.Mutex{}
	return &RmtOrganisationStore{Sql_Store: sqlStore, orgStore: orgStore}
}