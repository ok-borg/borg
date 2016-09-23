package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

type Snippet struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func getSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	res, err := client.Get(). // GetService
					Index("borg").
					Type("snippet").
					Id(id).
					Do()

	if err != nil || res.Found == false {
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid id")
		return
	}

	jsonSnipp, _ := res.Source.MarshalJSON()

	writeResponse(w, http.StatusOK, string(jsonSnipp))
}

func createSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}

	// first validate body
	var snipp Snippet
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[createSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}

	// no empty fields allowed
	if snipp.Title == "" || snipp.Text == "" {
		writeResponse(w, http.StatusBadRequest, "borg-api: snippet title and body cannot be empty")
		return
	}

	// create an id for this snippet
	snipp.Id = uuid.NewV4().String()

	// insert it in elastic
	_, err = client.Index().
		Index("borg").
		Type("snippet").
		Id(snipp.Id).
		BodyJson(snipp).
		Refresh(true).
		Do()

	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to save snippet")
	} else {
		writeJsonResponse(w, http.StatusOK, snipp)
	}
}

func updateSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {

}

func deleteSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {

}

func searchSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {

}
