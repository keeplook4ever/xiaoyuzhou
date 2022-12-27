package auth_service

import (
	"xiaoyuzhou/models/manager"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, error) {
	return manager.CheckAuth(a.Username, a.Password)
}
