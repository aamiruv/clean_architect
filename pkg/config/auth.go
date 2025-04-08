package config

type auth struct {
	secret string
}

func (a auth) Secret() string {
	return a.secret
}
