package commands

import (
	"errors"
	"github.com/ok-borg/borg/conf"
)

// Login saves a token acquired from the web page into the user config file
func Editor(editor string) error {
	if len(editor) == 0 {
		return errors.New("Please supply an editor. The default is vim, so if you are happy with that, do nothing.")
	}
	conf, err := conf.Get()
	if err != nil {
		return err
	}
	conf.Editor = editor
	return conf.Save()
}
