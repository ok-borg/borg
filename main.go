package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/fatih/color"
	flag "github.com/juju/gnuflag"
	"github.com/ok-borg/borg/commands"
	"github.com/ok-borg/borg/conf"
)

var versionNumber, operatingSystem, architecture string

func main() {
	flag.Parse(true)
	if *conf.Version || *conf.V {
		printVersion()
		return
	}
	if flag.NArg() == 0 {
		help()
		return
	}

	var err error
	if c, ok := commands.Commands[flag.Arg(0)]; !ok {
		err = commands.Query(strings.Join(flag.Args(), " "))
	} else {
		err = c.F(flag.Args())
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func help() {
	underline := color.New(color.Underline)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)

	underline.Println("Usage:")
	fmt.Print("\t$ ")
	green.Println("borg \"your question\"\n")
	fmt.Print("\t$ ")
	green.Println("borg COMMAND\n")
	fmt.Print("\n\t  BORG - Terminal based search for bash snippets\n\n")
	underline.Println("Commands:\n\n")
	for k, v := range commands.Commands {
		green.Printf("\t+ %-8s\t", k)
		fmt.Println(v.Summary)
	}
	// TODO: Display all the possible flags
	underline.Println("\nOptions:\n\n")
	// TODO: Replace --help so that it displays this usage instead
	blue.Printf("\t%-8s\t", "--help")
	fmt.Println("Show help")
}
func printVersion() {
	fmt.Printf("\tBorg version: %s (%s/%s)\n", versionNumber, operatingSystem, architecture)
}
