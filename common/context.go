package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ratemyteam/rmt/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var serverContext *ServerContext

func GetServerContext() *ServerContext {
	if serverContext == nil {
		panic("trying to access uninitialized RMTContext")
	}
	return serverContext
}

func CreateServerContext() *ServerContext{
	config.InitConfig("rmt.yaml")

	rmtInfo :=	&DbInfo{
		DbDriver:             viper.GetString("rmt-db-driver"),
		DbHost:               viper.GetString("rmt-db-host"),
		DbPort:               viper.GetInt("rmt-db-port"),
		DbName:               viper.GetString("rmt-db-name"),
		DbUser:               viper.GetString("rmt-db-user"),
		DbPassword:           viper.GetString("rmt-db-password"),
		DbSslMode:            viper.GetString("rmt-db-ssl-mode"),
		DbSslRootCert:        viper.GetString("rmt-db-ssl-cert"),
		DBSnapshotIdentifier: viper.GetString("rmt-db-snapshot-identifier"),
	}
	serverContext = &ServerContext{
		RMTInfo: rmtInfo,
	}
	return serverContext
}

var once sync.Once
func SetServerContext(serverContext *ServerContext) {
	set := false
	once.Do(func() {
		set = true
		serverContext = serverContext
	})
	if !set {
		panic("trying to change already saved serverContext")
	}
}

func GetContext(c *gin.Context) context.Context {
	goCtx, found := c.Get(string(GoContextKey))
	if !found {
		panic("No Context Error")
	}
	return goCtx.(context.Context)
}

func GetLogFromCtx(ctx context.Context) *logrus.Entry {
	entry := ctx.Value(LKEY)
	if entry == nil {
		return logrus.WithField(string(RequestIdKey), "unknown")
	}
	return entry.(*logrus.Entry)
}
func GetLog(c *gin.Context) *logrus.Entry {
	ctx := GetContext(c)
	return GetLogFromCtx(ctx)
}

func GetContextAndLog(c *gin.Context) (context.Context, *logrus.Entry) {
	ctx := GetContext(c)
	log := GetLog(c)

	return ctx, log
}

