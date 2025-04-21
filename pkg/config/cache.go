package config

import "fmt"

type cache struct {
	driver   string
	ip       string
	port     uint
	prefix   string
	user     string
	password string
}

func (c cache) Prefix() string {
	return c.prefix
}

func (c cache) Driver() string {
	return c.driver
}

func (c cache) ConnectionString() string {
	url := fmt.Sprintf("%s:%d", c.ip, c.port)
	if c.user != "" && c.password != "" {
		url = fmt.Sprintf("%s:%s@%s", c.user, c.password, url)
	}
	switch c.driver {
	case "redis":
		return fmt.Sprintf("redis://%s", url)
	case "memcached":
		return url
	default:
		return ""
	}
}
