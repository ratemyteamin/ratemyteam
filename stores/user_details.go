package stores

import (
	"encoding/json"
	"github.com/ratemyteam/rmt/common"
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
	panic("implement me")
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
	res, err := rmt.Sql_Store.WithNewTransaction(f2)
	if res == nil {
		return true, err
	}
	return false, nil
}

func NewRmtUserStore(sqlStore *SqlStore) *RmtUserStore{
	userStore := &sync.Mutex{}
	return &RmtUserStore{Sql_Store:sqlStore, userStore:userStore}
}


