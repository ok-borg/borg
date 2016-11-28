package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ok-borg/borg/conf"
	"github.com/ok-borg/borg/types"
)

func init() {
	var summary string = "Edit summary"
	Commands["edit"] = Command{
		F:       Edit,
		Summary: summary,
	}
}

func findIdFromQueryIndex(queryIndex int) (string, error) {
	bs, err := ioutil.ReadFile(conf.QueryFile)
	if err != nil {
		return "", err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(bs, &m)
	if err != nil {
		return "", err
	}
	li, ok := m["ids"].([]interface{})
	l := []string{}
	for _, v := range li {
		l = append(l, v.(string))
	}
	if !ok {
		return "", errors.New("Can't find ids in query history")
	}
	if len(l) <= queryIndex {
		return "", errors.New("Can't find index in ids")
	}
	return l[queryIndex], nil
}

// Edit a snippet based on index from last query results
func Edit(args []string) error {
	if len(args) != 2 {
		return errors.New("Please supply a query index to edit.")
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
	s, err := getSnippet(id)
	if err != nil {
		return err
	}
	str := problemToText(s)
	c, err := conf.Get()
	if err != nil {
		return err
	}
	ioutil.WriteFile(conf.EditFile, []byte(str), 0755)
	cmd := exec.Command(c.Editor, conf.EditFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	bs, err := ioutil.ReadFile(conf.EditFile)
	if err != nil {
		return err
	}
	p, err := textToProblem(string(bs))
	if err != nil {
		return err
	}
	p.Id = s.Id
	if *conf.D {
		fmt.Println(s)
		fmt.Println(p)
		return nil
	}
	return saveSnippet(p)
}

func getSnippet(id string) (*types.Problem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/p/%v", host(), id), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	c, err := conf.Get()
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", c.Token)
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while making request: %v", err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Status code: %v: %s", rsp.StatusCode, body)
	}
	ret := types.Problem{}
	return &ret, json.Unmarshal(body, &ret)
}

var tpl = `{{.Title}}
{{range $index, $element := .Solutions}}
[{{toChar $index}}]
{{join $element.Body}}{{end}}
`

func problemToText(p *types.Problem) string {
	b := []byte{}
	buffer := bytes.NewBuffer(b)
	funcs := map[string]interface{}{
		"join": func(s []string) string {
			return strings.Join(s, "\n")
		},
		"toChar": toChar,
	}
	template.Must(template.New("edit").Funcs(funcs).Parse(tpl)).Execute(buffer, p)
	return string(buffer.Bytes())
}

// very primitive right now...
func textToProblem(s string) (types.Problem, error) {
	ret := types.Problem{}
	lines := strings.Split(s, "\n")
	// wish i had a parser combinator library to do this...
	title := ""
	solutions := []types.Solution{}
	if len(lines) < 3 {
		return ret, errors.New("Edit is too short")
	}
	buf := []string{}
	for i, v := range lines {
		if i == 0 {
			title = v
			continue
		}
		if i == 1 && len(strings.TrimSpace(v)) == 0 { // skip the newline after the title
			continue
		}
		if i == len(lines)-1 || (len(strings.TrimSpace(v)) == 3 && string(v[0]) == "[" && regexp.MustCompile("^[a-z]$").Match([]byte{v[1]}) && string(v[2]) == "]") {
			if len(buf) > 0 {
				solutions = append(solutions, types.Solution{
					Body: []string{strings.Join(buf, "\n")},
				})
				buf = []string{}
			}
			continue
		}
		buf = append(buf, v)
	}
	ret.Title = title
	ret.Solutions = solutions
	return ret, nil
}
