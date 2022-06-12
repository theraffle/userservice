package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
	"time"
)

// Connect opens postgresql DB
func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PWD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Collation:            "utf8mb4_general_ci",
		Loc:                  time.UTC,
		MaxAllowedPacket:     4 << 20.,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		DBName:               os.Getenv("DB_NAME"),
	}
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(connector)
	return db, nil
}
