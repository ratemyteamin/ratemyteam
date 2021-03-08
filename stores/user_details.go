package stores

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ratemyteam/rmt/common"
	log "github.com/sirupsen/logrus"
	"sync"
)

func userToSave(user *common.User) *common.User{
	res := &common.User{
		UserId:         user.UserId,
		Email:          user.Email,
		Password:       user.Password,
	}
	return res
}


type RmtUserStore struct {
	Sql_Store *SqlStore
	userStore *sync.Mutex
}

func (rmt *RmtUserStore) GetUserById(id string) (common.User, error) {
	panic("implement me")
}

func (rmt *RmtUserStore) GetUserByAccountId(id string) (common.User, error) {
	panic("implement me")
}

func (rmt *RmtUserStore) GetUserByEmail(id string) (common.User, error) {
	panic("implement me")
}

func (rmt *RmtUserStore) GetUserInTx(tx common.Transaction, id string,  queryColumn string) (common.User, error) {
	row := tx.QueryRow(fmt.Sprintf("select user_json from ltas_users where %s = $1", queryColumn), id)

	var rb []byte // Single-row scan doesn't support RawBytes
	err := row.Scan(&rb)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.User{}, err
		}
		return common.User{}, err
	}

	var user common.User
	err = json.Unmarshal(rb, &user)
	if err != nil {
		return common.User{}, err
	}
	return user, err
}

func (rmt *RmtUserStore) NewTx() (*sql.Tx, error) {
	return rmt.Sql_Store.NewTx()
}

func (rmt *RmtUserStore) WithNewTransaction(fn common.TxFn) (interface{}, error) {
	return rmt.Sql_Store.WithNewTransaction(fn)
}

func (rmt *RmtUserStore) Close() {
	log.Info("Closing the Connection With  LTAS DB")
	rmt.Sql_Store.Close()
}

func (rmt *RmtUserStore) StoreUser(user *common.User) (bool, error) {
	rmt.userStore.Lock()
	defer rmt.userStore.Unlock()
	result := userToSave(user)
	userrb, err := json.Marshal(&result)
	jsonStr := string(userrb)
	f2 := func(tx common.Transaction) (interface{}, error) {

		_, err := tx.Exec(jsonStr)
		if err != nil {
			return nil, err
		}

		return nil, err
	}
	res, err := rmt.WithNewTransaction(f2)
	if res == nil {
		return true, err
	}
	return false, nil

}
func NewRmtUserStore(dbInfo *common.DbInfo2) *RmtUserStore{
	db, _ := common.CreateDbConnection(dbInfo)
	sqlStore, _ := NewSqlStore(db, "")
	userStore := &sync.Mutex{}
	return &RmtUserStore{Sql_Store:sqlStore, userStore:userStore}
}


