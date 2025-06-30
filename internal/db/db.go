package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"final_assessment/internal/utils"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB Database

func Connect() error {
	connStr := os.Getenv(utils.EnvDatabaseURL)
	if connStr == "" {
		return fmt.Errorf("%s not set", utils.EnvDatabaseURL)
	}
	baseDB, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).Error(utils.LogDBConnectFail)
		return err
	}
	baseDB.SetMaxOpenConns(10)
	baseDB.SetMaxIdleConns(5)
	baseDB.SetConnMaxLifetime(30 * time.Minute)
	if err := baseDB.Ping(); err != nil {
		logrus.WithError(err).Error("Failed to ping DB")
		return err
	}
	DB = &PostgresDB{baseDB}
	logrus.Info(utils.LogDBConnected)
	return nil
}

func Close() {
	if DB != nil {
		if pg, ok := DB.(*PostgresDB); ok {
			pg.DB.Close()
			logrus.Info(utils.LogDBClosed)
		}
	}
}
