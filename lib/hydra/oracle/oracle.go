package oracle

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"time"
)

var ServiceName = []string{
	"orcl",
	"xe",
	"oracle",
}

func Check(Host, Username, Password string, Port int) (bool, error) {
	var db *sql.DB
	var err error
	for _, service := range ServiceName {
		dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", Username, Password, Host, Port, service)
		db, err = sql.Open("oracle", dataSourceName)
		if err == nil {
			break
		}
	}
	if db == nil {
		return false, err
	}
	defer db.Close()
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(0)
	err = db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}
