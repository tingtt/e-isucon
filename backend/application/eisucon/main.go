package eisucon

import (
	"errors"
	"prc_hub_back/domain/model/eisucon"
)

// Singleton field
var migrateSqlFile string

func Init(user string, password string, host string, port uint, db string, migrateSqlFilePath string) {
	eisucon.Init(user, password, host, port, db)
	migrateSqlFile = migrateSqlFilePath
}

func Migrate() error {
	if migrateSqlFile == "" {
		return errors.New("migrate sql file does not set")
	}
	return eisucon.Migrate(migrateSqlFile)
}
