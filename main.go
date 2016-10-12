package main

import (
	"fmt"
	"os"

	flag "github.com/juju/gnuflag"
	"github.com/ok-borg/borg/commands"
)

func main() {
	flag.Parse(true)
	if flag.NArg() == 0 {
		help()
		return
	}

	if c, ok := commands.Commands[flag.Arg(0)]; !ok {
		commands.Query(flag.Arg(0))
	} else {
		if err := c.F(flag.Args()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func help() {
	fmt.Print("\033[4mUsage:\033[0m\n\n")
	fmt.Print("\t$ \033[32mborg \"your question\"\033[0m\n")
	fmt.Print("\t$ \033[32mborg COMMAND\033[0m\n")
	fmt.Print("\n\t  BORG - A terminal based search engine for bash snippets\n\n")
	fmt.Print("\033[4mCommands:\033[0m\n\n")
	for k, v := range commands.Commands {
		fmt.Printf("\t\033[32m+ %-8s\t\033[0m%s\n", k, v.Summary)
	}
	// TODO: Display all the possible flags
	fmt.Print("\n\033[4mOptions:\033[0m\n\n")
	// TODO: Replace --help so that it displays this usage instead
	fmt.Printf("\t\033[34m%-8s\t\033[0m%s\n", "--help", "Show help")
}
