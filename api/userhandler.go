package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ratemyteam/rmt/common"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func unmarshalBody(c *gin.Context, v interface{}) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic("Error in Reading Request Body")
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		panic("Error in Parsing Request Body")
	}
}

func (rmt *RmtContext) StoreUser(c *gin.Context){
	ctx , _ := common.GetContextAndLog(c)
	var usercreatereq common.CreateUserRequest
	unmarshalBody(c, &usercreatereq)
	s := rmt.GetSessionContext()
	res, err := s.Store.WithNewTransaction(func(tx common.Transaction) (interface{} ,error) {
		id := GenerateIdChecked("user", func(id string) bool {
			_, err := s.Store.GetUserById(id)
			logrus.Error(err)
			return common.IsNotFound(err)
		})
		_, err := s.Store.GetUserByEmail(string(usercreatereq.Email))

		exists := common.IsNotFound(err)

		if !exists{
			logrus.Error(common.ErrUserExists)
			return nil, common.NewRuntimeError(common.ErrUserExists).WithHttpResponseCode(http.StatusBadRequest)
		}
		user := common.User{
			Email:usercreatereq.Email,
			UserId: common.UserDet(id),
			Password: usercreatereq.Password,
			CompanyName: usercreatereq.CompanyName,
			CreateTimestamp: &time.Time{},
		}
		err = s.CreateNewUser(ctx, tx, &user)
		return user, err

	})

	if err == nil{
		userId := res.(common.User).UserId
		userResp := common.CreateUserResponse{Status: common.Unverified, UserId:userId}
		respondAccepted(userResp, c)
	}else {
		userResp := common.CreateUserResponse{Status: common.Failed}
		respondError(userResp, c)
	}

}
