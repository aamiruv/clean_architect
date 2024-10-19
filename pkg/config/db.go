package config

import (
	"fmt"
	"path/filepath"
)

type db struct {
	driver   string
	ip       string
	port     uint
	userName string
	password string
	name     string
	path     string
}

func (db db) Driver() string {
	return db.driver
}

func (db db) ConnectionString() string {
	switch db.driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			db.userName,
			db.password,
			db.ip,
			db.port,
			db.name,
		)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			db.userName,
			db.password,
			db.ip,
			db.port,
			db.name,
		)
	case "sqlite":
		return fmt.Sprintf("%s.sqlite", filepath.Join(db.path, db.name))
	default:
		return ""
	}

}
