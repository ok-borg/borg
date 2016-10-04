package commands

import (
	"github.com/crufter/borg/conf"
)

func Login(token string) error {
	conf, err := conf.Get()
	if err != nil {
		return err
	}
	conf.Token = token
	return conf.Save()
}
