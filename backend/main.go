package main

import (
	"fmt"
	"os"
	"prc_hub_back/application/eisucon"
	"prc_hub_back/application/event"
	"prc_hub_back/application/user"
	"prc_hub_back/domain/model/logger"
	"prc_hub_back/domain/model/logrus"
	"prc_hub_back/presentation/echo"
	"time"

	logruss "github.com/sirupsen/logrus"

	"github.com/spf13/pflag"
)

// flags (コマンドライン引数)
var (
	port = pflag.Uint("port", 1323, "publish port")

	logLevel     = pflag.String("log.level", "info", "Only log messages with the given severity or above. One of: [panic, fatal, error, warn, info, debug, trace]")
	logTimestamp = pflag.Bool("log.timestamp", false, "Enable log timestamp.")
	issuer       = pflag.String("jwt.issuer", "", "jwt issuer")
	secret       = pflag.String("jwt.secret", "", "jwt secret")

	mysqlHost     = pflag.String("mysql.host", "localhost", "MySQL host")
	mysqlPort     = pflag.Uint("mysql.port", 3306, "MySQL port")
	mysqlDB       = pflag.String("mysql.db", "prc_hub", "MySQL db")
	mysqlUser     = pflag.String("mysql.user", "prc_hub", "MySQL username")
	mysqlPassword = pflag.String("mysql.password", "", "MySQL password")

	eisuconMigrationFile = pflag.String("migrate-sql-file", "./domain/model/eisucon/migrate.sql", "sql file for migrate with 'POST /reset'")
)

func main() {
	logger.Init(logrus.New(logrus.Param{
		RepeatCaller: func() *bool { var b = true; return &b }(),
		Formatter: &logruss.TextFormatter{
			FullTimestamp:   *logTimestamp,
			TimestampFormat: time.RFC3339Nano,
		},
	}))

	// コマンドライン引数の取得
	pflag.Parse()

	// `--log.level`
	ok, lv := convertLogLevel(*logLevel)
	if !ok {
		logger.Logger().Fatalf("`--log.level` must be specified as \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\" or \"trace\"")
	}
	logger.Logger().SetLevel(lv)

	// `--jwt.issuer`
	if *issuer == "" {
		fmt.Println("`--jwt.issuer` option is required")
		os.Exit(1)
	}

	// `--jwt.secret`
	if *secret == "" {
		fmt.Println("`--jwt.secret` option is required")
		os.Exit(1)
	}

	// Init application services
	user.Init(*mysqlUser, *mysqlPassword, *mysqlHost, *mysqlPort, *mysqlDB)
	event.Init(*mysqlUser, *mysqlPassword, *mysqlHost, *mysqlPort, *mysqlDB)
	eisucon.Init(*mysqlUser, *mysqlPassword, *mysqlHost, *mysqlPort, *mysqlDB, *eisuconMigrationFile)

	// Migrate seed data
	err := eisucon.Migrate()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	echo.Start(*port, *issuer, *secret, []string{"*"})
}
