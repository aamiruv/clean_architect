package config

import "time"

type auth struct {
	secret   string
	lifeTime int
}

func (a auth) Secret() string {
	return a.secret
}

func (a auth) LifeTime() time.Duration {
	return time.Duration(a.lifeTime) * time.Minute
}
