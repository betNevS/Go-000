package dao

import (
	"database/sql"
	"sync"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type dbsMap struct {
	dbs map[string]*sql.DB
	sync.RWMutex
}

var dbsM = dbsMap{dbs: make(map[string]*sql.DB)}

func NewDB(dbname string) (*sql.DB, error) {
	dbsM.RLock()
	if db, ok := dbsM.dbs[dbname]; ok {
		dbsM.RUnlock()
		return db, nil
	}
	dbsM.RUnlock()
	dbsM.Lock()
	defer dbsM.Unlock()
	db, err := sql.Open("mysql", "root:666666@tcp(127.0.0.1)/"+dbname)
	if err != nil {
		return nil, errors.Wrapf(err, "connect mysql[%s] dbname[%s] failed!", "127.0.0.1", dbname)
	}
	dbsM.dbs[dbname] = db
	return db, err
}
