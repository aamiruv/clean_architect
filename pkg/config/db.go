package config

import "fmt"

type db struct {
	IP       string `json:"ip"`
	Port     uint   `json:"port"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
}

func (db db) GetConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.UserName,
		db.Password,
		db.IP,
		db.Port,
		db.DBName,
	)
}
