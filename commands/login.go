package commands

import (
	"errors"
	"github.com/crufter/borg/conf"
)

func Login(token string) error {
	if len(token) == 0 {
		return errors.New("Please supply a token. Don't have one? Go to http://ok-b.org and get it")
	}
	conf, err := conf.Get()
	if err != nil {
		return err
	}
	conf.Token = token
	return conf.Save()
}
