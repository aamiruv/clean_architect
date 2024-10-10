package config

import "fmt"

type db struct {
	ip       string
	port     uint
	userName string
	password string
	dbName   string
}

func (db db) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.userName,
		db.password,
		db.ip,
		db.port,
		db.dbName,
	)
}
