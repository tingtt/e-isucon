package mysql

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func Open() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", "root:secret@unix(/var/lib/mysql/mysql.sock)/prc_hub?parseTime=true&multiStatements=true")
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(16)
	return db, err
}
