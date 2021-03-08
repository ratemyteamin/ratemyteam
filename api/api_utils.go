package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ratemyteam/rmt/common"
	"net/http"
	"time"
)

func respondAccepted(response interface{}, c *gin.Context) {
	respond(http.StatusAccepted, response, c)
}

func respondError(response interface{}, c *gin.Context) {
	respond(http.StatusInternalServerError, response, c)
}

func respondBadRequest(response interface{}, c *gin.Context) {
	respond(http.StatusBadRequest, response, c)
}

func respond(status int, response interface{}, c *gin.Context) {
	log := common.GetLog(c)
	resp, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshling Response")
	}
	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.Writer.WriteHeader(status)
	log.Printf("response --> "+string(resp))
	c.Writer.WriteString(string(resp))
}


func (l *RmtContext) GetSessionContext() *common.Session{

	session := new(common.Session)
	session.Store = l.ServerContext.RmtStore
	return session
}



func assertParamPresent(c *gin.Context, params ...string) string {
	for _, param := range params {
		ps := c.Param(param)
		if ps != "" {
			return ps
		}
	}

	panic(fmt.Sprintf("required parameter missing: %v", params))

}

func GenerateIdChecked(prefix string, check func(string) bool) string {
	safety := 5
	for {
		id := GenerateId(prefix)
		if check(id) {
			return id
		}
		safety = safety - 1
		if safety == 0 {
			panic(common.NewRuntimeError(common.ErrTODO).WithMessagef("Cannot create a non-colliding ID"))
		}
	}
}

func GenerateId(prefix string) string {
	// Note: YARN branch was using UnixNano() for this
	return fmt.Sprintf("%s-%d", prefix, time.Now().Unix())
}