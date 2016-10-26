package commands

import (
	"errors"
	"fmt"
	"strconv"
)

func init() {
	var summary string = "Link summary"
	Commands["link"] = Command{
		F:       Link,
		Summary: summary,
	}
}

// Link prints the url to a query result
func Link(args []string) error {
	if len(args) != 2 {
		return errors.New("Please supply a query index to generate the link.")
	}
	queryIndex := args[1]
	i, err := strconv.ParseInt(queryIndex, 10, 32)
	if err != nil {
		return err
	}
	id, err := findIdFromQueryIndex(int(i - 1))
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("https://ok-b.org/t/%v/x", id))
	return nil
}
