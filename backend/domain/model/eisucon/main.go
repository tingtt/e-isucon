package eisucon

import (
	"github.com/jmoiron/sqlx"
)

func Migrate(sqlFile string) error {
	db, err := OpenMysql()
	if err != nil {
		return err
	}

	_, err = sqlx.LoadFile(db, sqlFile)
	if err != nil {
		return err
	}
	return nil
}
