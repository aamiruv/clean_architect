package config

import (
	"fmt"
)

type event struct {
	driver   string
	ip       string
	port     uint
	user     string
	password string
}

func (e event) Driver() string {
	return e.driver
}

func (e event) ConnectionString() string {
	url := fmt.Sprintf("%s:%d", e.ip, e.port)
	if e.user != "" && e.password != "" {
		url = fmt.Sprintf("%s:%s@%s", e.user, e.password, url)
	}
	switch e.driver {
	case "rabbitmq":
		return fmt.Sprintf("amqp://%s", url)
	case "redis":
		return fmt.Sprintf("redis://%s", url)
	case "nats":
		return fmt.Sprintf("nats://%s", url)
	default:
		return ""
	}
}
