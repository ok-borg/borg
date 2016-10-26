package commands

import (
	"errors"

	"github.com/ok-borg/borg/conf"
)

func init() {
	var summary string = "Login summary"
	Commands["login"] = Command{
		F:       Login,
		Summary: summary,
	}
}

// Login saves a token acquired from the web page into the user config file
func Login(args []string) error {
	if len(args) != 2 {
		return errors.New("Please supply a github token to login with.")
	}
	token := args[1]
	if len(token) == 0 {
		return errors.New("Please supply a token. Don't have one? Go to https://ok-b.org and get it")
	}
	conf, err := conf.Get()
	if err != nil {
		return err
	}
	conf.Token = token
	return conf.Save()
}
