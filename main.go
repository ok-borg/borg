package main

import (
	"fmt"
	"github.com/ok-borg/borg/commands"
	flag "github.com/juju/gnuflag"
	"os"
)

func main() {
	flag.Parse(true)
	if flag.NArg() == 0 {
		help()
		return
	}
	var err error
	switch flag.Arg(0) {
	case "new":
		err = commands.New()
	case "login":
		err = commands.Login(flag.Arg(1))
	case "edit":
		err = commands.Edit(flag.Arg(1))
	case "worked":
		err = commands.Worked(flag.Arg(1))
	case "link":
		err = commands.Link(flag.Arg(1))
	default:
		err = commands.Query(flag.Arg(0))
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func help() {
	fmt.Println("Usage: borg \"your question\"")
}
