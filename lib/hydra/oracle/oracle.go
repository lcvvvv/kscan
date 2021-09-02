package oracle

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"time"
)

func Check(Host, Username, Password string, Port int) (bool, error) {
	dataSourceName := fmt.Sprintf(`user="%v" password="%v" connectString="%v:%d/orcl?connect_timeout=%v"`, Username, Password, Host, Port, time.Second*5)
	db, err := sql.Open("godror", dataSourceName)
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
