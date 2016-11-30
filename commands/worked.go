package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ok-borg/borg/conf"
)

func init() {
	var summary string = "Worked summary"
	Commands["worked"] = Command{
		F:       Worked,
		Summary: summary,
	}
}

// Worked lets you mark a result as relevant one for a query
func Worked(args []string) error {
	if len(args) != 2 {
		return errors.New("Please supply a query index.")
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
	query, err := getLastQuery()
	if err != nil {
		return err
	}
	return saveWorked(id, query)
}

func getLastQuery() (string, error) {
	bs, err := ioutil.ReadFile(conf.QueryFile)
	if err != nil {
		return "", err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(bs, &m)
	if err != nil {
		return "", err
	}
	s, ok := m["query"].(string)
	if !ok {
		return "", errors.New("Can't find last query")
	}
	return s, nil
}

func saveWorked(id, query string) error {
	bs, err := json.Marshal(map[string]string{
		"id":    id,
		"query": query,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/v1/worked", host()), bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	c, err := conf.Get()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", c.Token)
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	rsp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while making request: %v", err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("Status code: %v: %s", rsp.StatusCode, body)
	}
	return nil
}
