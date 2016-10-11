package endpoints

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	log "github.com/cihub/seelog"
	"github.com/ok-borg/borg/types"
	"github.com/ventu-io/go-shortid"
)

// GetSnippet by id
func (e Endpoints) GetSnippet(id string) (*types.Problem, error) {
	res, err := e.client.Get().
		Index("borg").
		Type("problem").
		Id(id).
		Do()
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, nil
	}
	jsonSnipp, _ := res.Source.MarshalJSON() // must be a better way to do this
	ret := types.Problem{}
	return &ret, json.Unmarshal(jsonSnipp, &ret)
}

// GetLatestSnippets in reverse chronological order
func (e *Endpoints) GetLatestSnippets() ([]types.Problem, error) {
	res, err := e.client.Search().
		Index("borg").
		Type("problem").
		From(0).
		Size(50).
		Sort("Created", false).
		Do()
	if err != nil {
		return nil, err
	}
	all := []types.Problem{}
	var ttyp types.Problem
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(types.Problem); ok {
			all = append(all, t)
		}
	}
	return all, nil
}

// CreateSnippet saves a snippet, generates id
func (e Endpoints) CreateSnippet(snipp *types.Problem, userId string) error {
	if snipp.Title == "" || len(snipp.Solutions) == 0 {
		return errors.New("Title or solutions missing")
	}
	snipp.Id = shortid.MustGenerate()
	snipp.CreatedBy = userId
	snipp.Created = time.Now()
	log.Infof("Snippet with id %v is created by %v", snipp.Id, snipp.CreatedBy)
	_, err := e.client.Index().
		Index("borg").
		Type("problem").
		Id(snipp.Id).
		BodyJson(snipp).
		Refresh(true).
		Do()
	return err
}

// UpdateSnippet saves a snippet
func (e Endpoints) UpdateSnippet(snipp *types.Problem, userId string) error {
	if snipp.Id == "" {
		return errors.New("No id found")
	}
	if snipp.Title == "" || len(snipp.Solutions) == 0 {
		return errors.New("Title or solutions missing")
	}
	snipp.LastUpdatedBy = userId
	snipp.LastUpdated = time.Now()
	log.Infof("Snippet %v is being updated by %v", snipp.Id, snipp.LastUpdatedBy)
	_, err := e.client.Index().
		Index("borg").
		Type("problem").
		Id(snipp.Id).
		BodyJson(snipp).
		Refresh(true).
		Do()
	if err != nil {
		log.Errorf("[updateSnippet] error updating snippet id: %s: %v", snipp.Id, err)
		return err
	}
	return nil
}

func deleteSnippet() {

}
