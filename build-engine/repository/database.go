package repository

import (
	"database/sql"
	"dezzles-apps/build-engine/model"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func (d *Database) getDB() *sql.DB {
	return d.db
}

func (d *Database) Connect(config *model.DatabaseConfig) error {
	cfg := mysql.NewConfig()

	cfg.User = config.Username
	cfg.Passwd = config.Password
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	cfg.DBName = config.Database
	// Get a database handle.
	var err error
	d.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	pingErr := d.db.Ping()
	if pingErr != nil {
		return pingErr
	}
	fmt.Println("Database connected!")
	return nil
}
