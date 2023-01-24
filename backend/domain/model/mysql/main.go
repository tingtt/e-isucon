package mysql

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func Open() (*sqlx.DB, error) {
	return sqlx.Open("mysql", "root:secret@unix(/var/lib/mysql/mysql.sock)/prc_hub?parseTime=true&multiStatements=true")
}
