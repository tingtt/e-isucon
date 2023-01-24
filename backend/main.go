package main

import (
	"fmt"
	"prc_hub_back/application/eisucon"
	"prc_hub_back/presentation/echo"

	"github.com/spf13/pflag"
)

// flags (コマンドライン引数)
var (
	logLTSV = pflag.Bool("log.ltsv", false, "Enable log with ltsv format")

	eisuconMigrationFile = pflag.String("migrate-sql-file", "./domain/model/eisucon/migrate.sql", "sql file for migrate with 'POST /reset'")
)

func main() {
	// コマンドライン引数の取得
	pflag.Parse()

	// Init application services
	eisucon.Init(*eisuconMigrationFile)

	// Migrate seed data
	err := eisucon.Migrate()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	echo.Start(*logLTSV)
}
