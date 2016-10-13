package commands

import (
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
	queryIndex := args[1]
	i, err := strconv.ParseInt(queryIndex, 10, 32)
	if err != nil {
		return err
	}
	id, err := findIdFromQueryIndex(int(i - 1))
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("http://ok-b.org/t/%v/x", id))
	return nil
}
