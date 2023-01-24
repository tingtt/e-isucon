package eisucon

import (
	"errors"
	"prc_hub_back/domain/model/eisucon"
)

// Singleton field
var migrateSqlFile string

func Init(migrateSqlFilePath string) {
	migrateSqlFile = migrateSqlFilePath
}

func Migrate() error {
	if migrateSqlFile == "" {
		return errors.New("migrate sql file does not set")
	}
	return eisucon.Migrate(migrateSqlFile)
}
