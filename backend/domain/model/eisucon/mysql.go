package eisucon

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Singleton field
var dsn string

func Init(user string, password string, db string) {
	dsn = fmt.Sprintf("%s:%s@unix(/var/lib/mysql/mysql.sock)/%s?parseTime=true&multiStatements=true", user, password, db)
}

// MySQLサーバーに接続
func OpenMysql() (*sqlx.DB, error) {
	if dsn == "" {
		return nil, errors.New("dsn does not set")
	}
	return sqlx.Open("mysql", dsn)
}
