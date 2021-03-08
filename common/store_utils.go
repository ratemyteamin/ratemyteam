package common

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-openapi/errors"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)
const DbConnectionMaxLifetime = "db-connectionmaxlifetime"
const DbMaxOpenConnections = "db-maxopenconnections"
const DbMaxIdleConnections = "db-maxidleconnections"


var db *sql.DB
var dbInfo DbInfo2

type DbInfo2 struct {
	DriverType          string `json:"driverType,omitempty"`
	HostName            string `json:"hostName,omitempty"`
	Port                int    `json:"port,omitempty"`
	DatabaseName        string `json:"databaseName,omitempty"`
	RestoreFromDatabase string `json:"restoreFromDatabase,omitempty"`
	MaintDbName         string `json:"maintDbName,omitempty"` // Optional, additional maintenance database name for PG connection. May be unset.
	User                string `json:"user,omitempty"`
	Password            string `json:"password,omitempty"`
	Options             string `json:"options,omitempty"`
}

var gosqlToJdbcDriverMappings = map[string]string{
	"postgres": "postgresql",
}

func (dbInfo *DbInfo2) GetConnectionString() string {
	connectionString := fmt.Sprintf("%s:rmtuser@tcp(%s:%d)/%s",
		dbInfo.User, dbInfo.HostName, dbInfo.Port, dbInfo.DatabaseName)
	return connectionString
}

func (dbInfo *DbInfo2) GetPrintableConnectionString() string {

	connectionString := fmt.Sprintf("%s:rmtuser@tcp(%s:%d)/%s",
		 dbInfo.User, dbInfo.HostName, dbInfo.Port, dbInfo.DatabaseName)
	return connectionString
}

func (dbInfo *DbInfo2) GetJdbcConnectionString() string {
	jdbcDriver, found := gosqlToJdbcDriverMappings[dbInfo.DriverType]
	if !found {
		panic(fmt.Sprintf("Unable to find JDBC Driver for go sql driver type %s", dbInfo.DriverType))
	}

	connectionString := fmt.Sprintf("jdbc:%s://%s:%d/%s",
		jdbcDriver, dbInfo.HostName, dbInfo.Port, dbInfo.DatabaseName)
	if dbInfo.Options != "" {
		connectionString = fmt.Sprintf("%s?%s", connectionString, dbInfo.Options)
	}
	return connectionString
}

// Check if everything is connectable or not should happen before start of API
func DatabaseSetup(info *DbInfo) error {
	dbInfo = ConvertDbInfo(info)
	var err error
	db, err = CreateDbConnection(&dbInfo)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	log.Infof("Ping successful! dbHost: %s", dbInfo.HostName)
	return nil
}

func ConvertDbInfo(info *DbInfo) DbInfo2 {
	info2 := DbInfo2{
		DriverType:   info.DbDriver,
		HostName:     info.DbHost,
		Port:         info.DbPort,
		DatabaseName: info.DbName,
		User:         info.DbUser,
		Password:     info.DbPassword,
	}
	if len(info.DbSslMode) > 0 {
		if len(info.DbSslRootCert) > 0 {
			info2.Options = "sslmode=" + info.DbSslMode + " sslrootcert=" + info.DbSslRootCert
		} else {
			info2.Options = "sslmode=" + info.DbSslMode
		}
	}
	return info2
}

func CreateDbConnection(info *DbInfo2) (*sql.DB, error) {
	connMaxLifetime, err := time.ParseDuration(viper.GetString(DbConnectionMaxLifetime))
	if err != nil {
		return nil, fmt.Errorf("error parsing connection lifetime: %v", err)
	}
	maxOpenConns := viper.GetInt(DbMaxOpenConnections)
	maxIdleConns := viper.GetInt(DbMaxIdleConnections)

	dbConnectString := info.GetConnectionString()

	printableConnectString := info.GetPrintableConnectionString()
	log.Infof("Connceting to Database %s", printableConnectString)
	db, err := sql.Open(dbInfo.DriverType, dbConnectString)
	if err != nil {
		errMsg := fmt.Sprintf("Error connecting to database with %s: %v", printableConnectString, err)
		log.Errorf(errMsg)
		return nil, errors.New(-1, errMsg)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	log.Infof("Connected to db: %s", printableConnectString)
	return db, nil
}