package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crufter/borg/conf"
	"github.com/crufter/borg/types"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func extractPost(s string) (string, string, error) {
	ss := strings.Split(s, "\n")
	if len(ss) < 2 {
		return "", "", errors.New("Content too short")
	}
	title := strings.TrimSpace(ss[0])
	body := strings.TrimSpace(ss[1])
	return title, body, nil
}

func New() error {
	cmd := exec.Command("vim", conf.HomeDir+"/.borg/edit")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
	bs, err := ioutil.ReadFile(conf.HomeDir + "/.borg/edit")
	if err != nil {
		return err
	}
	t, b, err := extractPost(string(bs))
	if err != nil {
		return err
	}
	p := types.Problem{
		Title: t,
		Solutions: []types.Solution{
			{
				Body: []string{b},
			},
		},
	}
	bs, err = json.Marshal(p)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/v1/p", host()), bytes.NewReader(bs))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create request: %v", err))
	}
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	rsp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while making request: %v", err))
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code: %v: %s", rsp.StatusCode, body))
	}
	return nil
}
