package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crufter/borg/types"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	servAddr = "http://borg.crufter.com:9992"
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	default:
		query(flag.Arg(0))
	}
}

func query(q string) {
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("GET", servAddr+"/v1/query?q="+url.QueryEscape(q), nil)
	if err != nil {
		fmt.Println("Failed to create request: " + err.Error())
		os.Exit(1)
	}

	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while making request: " + err.Error())
		os.Exit(1)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	problems := []types.Problem{}
	err = json.Unmarshal(body, &problems)
	if err != nil {
		fmt.Println("Malformed response from server")
		os.Exit(1)
	}
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
			fmt.Fprintf(w, "\t[%v%v]", i+1, x+1)
			for i, bodyPart := range sol.Body {
				if i > 0 {
					fmt.Fprintln(w, "\t\t", "-")
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
					if line == 10 {
						fmt.Fprintln(w, "\t", "...", "\t")
						break Loop
					}
				}
			}
		}
	}
	w.Flush()
}
