package stores

import (
	"database/sql"
	"github.com/ratemyteam/rmt/common"
)

type RMTStore struct {
	RmtUserStore *RmtUserStore
}

func (R *RMTStore) NewTx() (*sql.Tx, error) {
	panic("implement me")
}

func (R *RMTStore) WithNewTransaction(fn common.TxFn) (interface{}, error) {
	panic("implement me")
}

func (R *RMTStore) Close() {
	panic("implement me")
}

func (R *RMTStore) StoreUser(user *common.User) (bool, error) {
	return R.RmtUserStore.StoreUser(user)
}

func (R *RMTStore) GetUserInTx(tx common.Transaction, id string, queryColumn string) (common.User, error) {
	return R.RmtUserStore.GetUserInTx(tx, id, queryColumn)
}

func (R *RMTStore) GetUserById(id string) (common.User, error) {
	return R.RmtUserStore.GetUserById(id)
}

func (R *RMTStore) GetUserByAccountId(id string) (common.User, error) {
	return R.RmtUserStore.GetUserByAccountId(id)
}

func (R *RMTStore) GetUserByEmail(id string) (common.User, error) {
	return R.RmtUserStore.GetUserByEmail(id)
}

func NewRmtStore(dbInfo *common.DbInfo2) *RMTStore{
	rmtUserStore := NewRmtUserStore(dbInfo)
	return &RMTStore{RmtUserStore: rmtUserStore}
}