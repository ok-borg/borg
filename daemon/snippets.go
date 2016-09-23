package main

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/crufter/borg/types"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func getSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	res, err := client.Get(). // GetService
					Index("borg").
					Type("problem").
					Id(id).
					Do()

	if err != nil || res.Found == false {
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid id")
		return
	}

	jsonSnipp, _ := res.Source.MarshalJSON()

	writeResponse(w, http.StatusOK, string(jsonSnipp))
}

func createSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}

	// first validate body
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[createSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}

	// no empty fields allowed
	if snipp.Title == "" || len(snipp.Solutions) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: snippet title and solutins cannot be empty")
		return
	}

	// create an id for this snippet
	snipp.Id = uuid.NewV4().String()

	// insert it in elastic
	_, err = client.Index().
		Index("borg").
		Type("problem").
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

func updateSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}

	// first validate body
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[updateSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}

	// no empty fields allowed
	if snipp.Id == "" {
		writeResponse(w, http.StatusBadRequest, "borg-api: snippet id must not be nil")
		return
	}
	// insert it in elastic
	uRes, err := client.Index().
		Index("borg").
		Type("problem").
		Id(snipp.Id).
		BodyJson(snipp).
		Refresh(true).
		Do()

	if err != nil {
		log.Errorf("[updateSnippet] error updating snippet id: %s: %v", snipp.Id, err)
		writeResponse(w, http.StatusInternalServerError, "borg-api: error")
		return
	}
	writeResponse(w, http.StatusOK, "{}")
}

func deleteSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {

}
