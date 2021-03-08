package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)


func InitConfig(path string){
	if path != ""{
		viper.SetConfigFile(path)
	} else{
		path = "/tmp/rmt.yaml"
		log.Printf("Path is not provides, we will use default %s", path)

	}

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Can't read config: %s: err: %s\n", path, err)
		os.Exit(1)
	}
	viper.SetDefault("log-file", "ltas.log")
	logFileName := viper.GetString("log-file")
	logfile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)
	var logLevel log.Level
	logLevelStr := viper.GetString("log-level")
	if strings.EqualFold(logLevelStr, "error") {
		logLevel = log.ErrorLevel
	} else if strings.EqualFold(logLevelStr, "warn") {
		logLevel = log.WarnLevel
	} else if strings.EqualFold(logLevelStr, "info") {
		logLevel = log.InfoLevel
	} else {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02T15:04:05.999Z07:00"})
}
