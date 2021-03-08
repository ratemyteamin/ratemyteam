package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ratemyteam/rmt/common"
	"github.com/ratemyteam/rmt/stores"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)


type RmtContext struct {
	ServerContext *common.ServerContext
}

func NewContext() *RmtContext{
	return &RmtContext{}
}


func initContext() context.Context {
	entry := logrus.WithField("process", "init")
	ctx := context.WithValue(context.Background(), common.LKEY, entry)
	return ctx
}

func Run(serverContext *common.ServerContext){
	ctx := initContext()
	log := common.GetLogFromCtx(ctx)
	rmtContext := NewContext()
	rmtContext.ServerContext = serverContext
	common.SetServerContext(serverContext)
	log.Info("Starting WebApp")
	if err := common.DatabaseSetup(rmtContext.ServerContext.RMTInfo); err != nil {
		panic(fmt.Errorf("unable to init db err: %v", err))
	}
	db2info := common.ConvertDbInfo(rmtContext.ServerContext.RMTInfo)
	rmtStore := stores.NewRmtStore(&db2info)
	rmtContext.ServerContext.RmtStore = rmtStore
	router := gin.New()
	gin.DisableConsoleColor()
	router.Use(gin.Logger())

	user := router.Group("/api/v1/user")
	{
		user.POST("/create", rmtContext.StoreUser)

	}

	portStr := strconv.Itoa(8989)
	transport := "http"
	log.Infof("Starting to serve on %s://localhost:%s", transport, portStr)

	srv := http.Server{
		Addr:    ":" + portStr,
		Handler: router,
		// TODO: Enable TLS verify when we have a proper internal CA to verify the certs.
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Start listening in background
	go func() {

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}

	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down...")

	bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(bgCtx); err != nil {
		log.Fatal("error during shutdown of http server:", err)
	}
	log.Println("exiting")

}
