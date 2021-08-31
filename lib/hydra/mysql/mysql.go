package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Check(Host, Username, Password string, Port int) (bool, error) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v", Username, Password, Host, Port, 5*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return false, err
	}
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(0)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}
