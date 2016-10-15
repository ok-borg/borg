package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/crufter/slugify"
)

//TODO: Should I leave the structs here or move them to types package ?
type SlackMessage struct {
	Text        string             `json:"text"`
	Mrkdwn      bool               `json:"mrkdwn"`
	Attachments []SlackAttachments `json:"attachments"`
}

type SlackAttachments struct {
	Title     string   `json:"title"`
	TitleLink string   `json:"title_link"`
	Text      string   `json:"text"`
	Mrkdwn_in []string `json:"mrkdwn_in"`
	Color     string   `json:"color"`
}

// TODO: Return a specific message if there is no results
func (e Endpoints) Slack(text string) (string, error) {
	problems, err := e.Query(text, 3, false)
	if err != nil {
		log.Errorf("[endpoint.Slack] error processing slack command: %s ", err.Error())
		return "", err
	}

	var m SlackMessage
	m.Text = fmt.Sprint("_", text, "_")
	m.Mrkdwn = true
	attachments := []SlackAttachments{}
	for _, prob := range problems {
		var buffer bytes.Buffer
		for x, sol := range prob.Solutions {
			buffer.WriteString(fmt.Sprintf("[%v] ```%s```\n", x, strings.Join(sol.Body, "\n")))
		}
		attachments = append(attachments, SlackAttachments{
			Title: prob.Title,
			// TODO: The URL should be defined somewhere else for the whole project
			TitleLink: "https://ok-b.org/t/" + fmt.Sprintf("%v/%v", prob.Id, slugify.S(prob.Title)),
			Text:      buffer.String(),
			Mrkdwn_in: []string{"text"},
			Color:     "#69DBF8",
		})
	}
	m.Attachments = attachments

	json, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
