package common

import (
	"database/sql"
	"time"
)
// API Reference Details

type Status string

const (
	Active Status = "Active"
	Unverified  Status = "Unverified"
	Disabled Status = "Disabled"
)

type UserDet string
type User struct{
	UserId UserDet `json:"userid,omitempty"`
	Email UserDet  `json:"email,omitempty"`
	Password UserDet  `json:"password,omitempty"`
	CompanyName UserDet  `json:"companyName,omitempty"`
	CreateTimestamp  *time.Time    `json:"createTimestamp,omitempty"`
}
type CreateUserResponse struct {
	Status  Status
	UserId UserDet `json:"userid,omitempty"`
}

// Server Context
type ServerContext struct {
	RMTInfo *DbInfo
	RmtStore RMTStore
}


type str string

const (
	GoContextKey     str = "goContext"
	LKEY           str = "logEntry"
	RequestIdKey     str = "requestId"
)


// DB Interfaces and Structs

type Session struct {
	Store 	RMTStore
}

type DbInfo struct {
	DbDriver             string `json:"dbDriver,omitempty"`
	DbHost               string `json:"dbHost,omitempty"`
	DbPort               int    `json:"dbPort,omitempty"`
	DbName               string `json:"dbName,omitempty"`
	DbUser               string `json:"dbUser,omitempty"`
	DbPassword           string `json:"dbPassword,omitempty"` // TODO: redact or change to jceks file pointer
	DbSslMode            string `json:"dbSslMode,omitempty"`
	DbSslRootCert        string `json:"dbSslRootCert,omitempty"`
	DBSnapshotIdentifier string `json:"dbSnapshotIdentifier,omitempty"`
}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type TxFn func(Transaction) (interface{}, error)

type SqlStore interface {
	NewTx() (*sql.Tx, error)
	WithNewTransaction(fn TxFn) (interface{}, error)
	Close()
}

type UserDetailStore interface {
	SqlStore
	StoreUser(user *User) (bool, error)
	GetUserInTx(tx Transaction, id string, queryColumn string) (User, error)
	GetUserById( id string) (User, error)
	GetUserByAccountId(id string) (User, error)
	GetUserByEmail(id string) (User, error)
}

type RMTStore interface {
	UserDetailStore
}