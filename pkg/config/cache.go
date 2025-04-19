package config

import "fmt"

type cache struct {
	driver string
	ip     string
	port   uint
	prefix string
}

func (c cache) Prefix() string {
	return c.prefix
}

func (c cache) Driver() string {
	return c.driver
}

func (c cache) ConnectionString() string {
	return fmt.Sprintf("%s:%d", c.ip, c.port)
}
