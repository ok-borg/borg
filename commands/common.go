package commands

var Commands map[string]Command = map[string]Command{}

type Command struct {
	F       func([]string) error
	Summary string
}
