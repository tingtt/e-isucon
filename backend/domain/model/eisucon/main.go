package eisucon

import (
	"prc_hub_back/domain/model/mysql"

	"github.com/jmoiron/sqlx"
)

func Migrate(sqlFile string) error {
	db, err := mysql.Open()
	if err != nil {
		return err
	}

	_, err = sqlx.LoadFile(db, sqlFile)
	if err != nil {
		return err
	}
	return nil
}
