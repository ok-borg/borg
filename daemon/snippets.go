package main

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	log "github.com/cihub/seelog"
	"github.com/crufter/borg/types"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/ventu-io/go-shortid"
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

func getLatestSnippets(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	//byUserId := r.FormValue("by")
	res, err := client.Search().
		Index("borg").
		Type("problem").
		From(0).
		Size(200).
		Sort("Created", false).
		Do()
	if err != nil {
		panic(err)
	}
	all := []types.Problem{}
	var ttyp types.Problem
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(types.Problem); ok {
			all = append(all, t)
		}
	}
	bs, err := json.Marshal(all)
	if err != nil {
		panic(err)
	}
	writeResponse(w, http.StatusOK, string(bs))
}

func createSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[createSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}
	if snipp.Title == "" || len(snipp.Solutions) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: snippet title and solutins cannot be empty")
		return
	}
	snipp.Id = shortid.MustGenerate()
	snipp.CreatedBy = ctx.Value("userId").(string)
	snipp.Created = time.Now()
	log.Infof("Snippet with id %v is created by %v", snipp.Id, snipp.CreatedBy)
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
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[updateSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}
	if snipp.Id == "" {
		writeResponse(w, http.StatusBadRequest, "borg-api: snippet id must not be nil")
		return
	}
	snipp.LastUpdatedBy = ctx.Value("userId").(string)
	snipp.LastUpdated = time.Now()
	log.Infof("Snippet %v is being updated by %v", snipp.Id, snipp.LastUpdatedBy)
	_, err = client.Index().
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
