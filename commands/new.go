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
	"strings"
	"time"

	"github.com/ok-borg/borg/conf"
	"github.com/ok-borg/borg/types"
)

func init() {
	var summary string = "Save summary"
	Commands["new"] = Command{
		F:       New,
		Summary: summary,
	}
}

func extractPost(s string) (string, string, error) {
	ss := strings.Split(s, "\n")
	if len(ss) < 3 {
		return "", "", errors.New("Content too short must be at least 3 lines")
	}
	title := strings.TrimSpace(ss[0])
	body := strings.TrimSpace(strings.Join(ss[1:], "\n"))
	return title, body, nil
}

// New saves a new snippet into the borg mind
func New([]string) error {
	c, err := conf.Get()
	if err != nil {
		return err
	}
	cmd := exec.Command(c.Editor, conf.EditFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	bs, err := ioutil.ReadFile(conf.EditFile)
	if err != nil {
		return err
	}
	t, b, err := extractPost(string(bs))
	if err != nil {
		return err
	}
	if len(t) == 0 || len(b) == 0 {
		return errors.New("Title or body is empty")
	}
	p := types.Problem{
		Title: t,
		Solutions: []types.Solution{
			{
				Body: []string{b},
			},
		},
	}
	return saveSnippet(p)
}

// POSTS or PUTS based on id existence
func saveSnippet(p types.Problem) error {
	method := "POST"
	if len(p.Id) > 0 {
		method = "PUT"
	}
	bs, err := json.Marshal(p)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%v/v1/p", host()), bytes.NewReader(bs))
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
