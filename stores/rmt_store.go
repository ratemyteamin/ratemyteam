package stores

import (
	"database/sql"
	"github.com/ratemyteam/rmt/common"
	log "github.com/sirupsen/logrus"
)

type RMTStore struct {
	Sql_Store *SqlStore
	RmtUserStore *RmtUserStore
	OrganisationStore *RmtOrganisationStore
}

func (rmt *RMTStore) NewTx() (*sql.Tx, error) {
	return rmt.Sql_Store.NewTx()
}

func (rmt *RMTStore) WithNewTransaction(fn common.TxFn) (interface{}, error) {
	return rmt.Sql_Store.WithNewTransaction(fn)
}

func (rmt *RMTStore) StoreNewPostMetadata(postMetaData *common.PostMetadata) (bool, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetUserOrganisation(user *common.User) (common.Organisation, error) {
	return rmt.OrganisationStore.GetUserOrganisation(user)
}

func (rmt *RMTStore) GetUserTeam(user *common.User) (common.Team, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetUserTeamAndOrganisation(user *common.User) (common.Team, common.Organisation, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetUserPostsMetadata(user *common.User) ([]common.PostMetadata, error) {
	panic("implement me")
}

func (rmt *RMTStore) StoreOrganisation(organistation *common.Organisation) (bool, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetTeamsInOrgaisation(organistation *common.Organisation) ([]common.Team, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetUsersInOrganisation(organistation *common.Organisation) ([]common.User, error) {
	panic("implement me")
}

func (rmt *RMTStore) GetOrganisationPostsMetadata(organisation *common.Organisation) ([]common.PostMetadata, error) {
	panic("implement me")
}

func (rmt *RMTStore) StoreNewPost(postMetaData *common.PostMetadata) (bool, error) {
	panic("implement me")
}

func (rmt *RMTStore) StoreUser(user *common.User) (bool, error) {
	return rmt.RmtUserStore.StoreUser(user)
}

func (rmt *RMTStore) GetUserInTx(tx common.Transaction, id string, queryColumn string) (common.User, error) {
	return rmt.RmtUserStore.GetUserInTx(tx, id, queryColumn)
}

func (rmt *RMTStore) GetUserById(id string) (common.User, error) {
	return rmt.RmtUserStore.GetUserById(id)
}

func (rmt *RMTStore) GetUserByAccountId(id string) (common.User, error) {
	return rmt.RmtUserStore.GetUserByAccountId(id)
}

func (rmt *RMTStore) GetUserByEmail(id string) (common.User, error) {
	return rmt.RmtUserStore.GetUserByEmail(id)
}


func (rmt *RMTStore) Close() {
	log.Info("Closing the Connection With RMT DB")
	rmt.Sql_Store.Close()
}

func NewRmtStore(dbInfo *common.DbInfo2) *RMTStore{
	db, _ := common.CreateDbConnection(dbInfo)
	sqlStore, _ := NewSqlStore(db, "")
	rmtUserStore := NewRmtUserStore(sqlStore)
	rmtOrgStore := NewOrganisationStore(rmtUserStore.Sql_Store)
	return &RMTStore{RmtUserStore: rmtUserStore, OrganisationStore: rmtOrgStore}
}