package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ok-borg/borg/conf"
	"github.com/ok-borg/borg/types"
)

// Query the borg server
func Query(q string) error {
	c, err := conf.Get()
	if err != nil {
		return err
	}
	if len(c.PipeTo) > 0 && *conf.DontPipe == false {
		c1 := exec.Command("borg", "--dontpipe", `"`+q+`"`)
		c2 := exec.Command(c.PipeTo)
		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout
		if err = c2.Start(); err != nil {
			return err
		}
		if err = c1.Run(); err != nil {
			return err
		}
		return c2.Wait()
	}
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/query?l=%v&p=%v&q=%v", host(), *conf.L, *conf.P, url.QueryEscape(q)), nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %s", err.Error())
	}
	rsp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while making request: %s", err.Error())
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if *conf.D {
		fmt.Println(fmt.Sprintf("json response: %v", string(body)))
	}
	problems := []types.Problem{}
	err = json.Unmarshal(body, &problems)
	if err != nil {
		return errors.New("Malformed response from server")
	}
	err = writeToFile(q, problems)
	if err != nil {
		fmt.Println(err)
	}
	renderQuery(problems)
	return nil
}

func renderQuery(problems []types.Problem) {
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight)
	for i, prob := range problems {
		if i > 0 {
			fmt.Fprintln(w, "")
		}
		fmt.Fprintln(w, fmt.Sprintf("(%v)", i+1), prob.Title)
		line := 0
	Loop:
		for x, sol := range prob.Solutions {
			fmt.Fprintf(w, "\t[%v]", toChar(x))
			for i, bodyPart := range sol.Body {
				if i > 0 {
					fmt.Fprintln(w, "\t\t", "")
				}
				bodyPartLines := strings.Split(bodyPart, "\n")
				for j, bodyPartLine := range bodyPartLines {
					t := "\t\t"
					if i == 0 && j == 0 {
						t = "\t"
					}
					if len(strings.TrimSpace(bodyPartLine)) == 0 {
						continue
					}
					fmt.Fprintln(w, t, strings.Trim(bodyPartLine, "\n"))
					line++
					if line == 10 && *conf.F == false {
						fmt.Fprintln(w, "\t", "...", "\t")
						break Loop
					}
				}
			}
		}
	}
	w.Flush()
}

func writeToFile(query string, ps []types.Problem) error {
	m := map[string]interface{}{
		"query": query,
	}
	ids := []string{}
	for _, v := range ps {
		ids = append(ids, v.Id)
	}
	m["ids"] = ids
	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(conf.QueryFile, bs, 0755)
}

func host() string {
	return fmt.Sprintf("http://%v:9992", *conf.S)
}

func toChar(i int) string {
	return string('a' + i)
}
