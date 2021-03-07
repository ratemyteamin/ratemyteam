package stores

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/ratemyteam/rmt/common"
	"github.com/sirupsen/logrus"
	"os"
)

type SqlStore struct {
	db           *sql.DB // DB handle
	migrationDir string
}

func (s *SqlStore) NewTx() (*sql.Tx, error) {
	return s.db.Begin()
}


func (s *SqlStore) WithNewTransaction(fn common.TxFn) (res interface{}, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback
			_ = tx.Rollback()
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fn(tx)
}

func (s *SqlStore) Close() {
	if s.db == nil {
		return
	}
	err := s.db.Close()
	if err != nil {
		logrus.Warnf("Failed to close sql db: %v", err)
	}
}

func NewSqlStore(db *sql.DB, migrationDir string) (*SqlStore, error){
	sql_store := &SqlStore{db:db, migrationDir:migrationDir}

	if migrationDir != ""{
		err := sql_store.migrateDbSchemaToLatest()
		if err != nil {
			return nil, err
		}
	}

	return sql_store, nil
}

func (s *SqlStore) migrateDbSchemaToLatest() error {
	driver, err := mysql.WithInstance(s.db, &mysql.Config{})

	// Read migrations and connect to a local mysql database.
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://migrations/%s", s.migrationDir), "mysql", driver)
	if err != nil {
		return fmt.Errorf("error setting up DB migration: %v", err)
	}

	// explanation for os.ErrNotExist check: https://github.com/golang-migrate/migrate/issues/79
	if err = m.Up(); err != nil && err != migrate.ErrNoChange && err != os.ErrNotExist {
		return fmt.Errorf("error upgrading schema: %v", err)
	}

	version, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("error getting schema version: %v", err)
	}
	logrus.Infof("schema is now at version: %d", version)

	return nil
}