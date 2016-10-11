package endpoints

import (
	log "github.com/cihub/seelog"
	"github.com/jpillora/go-ogle-analytics"
	"github.com/ok-borg/borg/types"
	"gopkg.in/olivere/elastic.v3"
	"reflect"
)

// Query the borg
func (e *Endpoints) Query(q string, size int, private bool) ([]types.Problem, error) {
	if size > 50 {
		size = 50
	}
	ql := q
	if private {
		ql = "PRIVATE"
	}
	log.Infof("Querying %v with size '%v'", ql, size)
	if e.analytics != nil {
		err := e.analytics.Send(ga.NewEvent("search", "backend").Label(ql))
		if err != nil {
			log.Warnf("Failed to send analytics events: %v", err)
		}
	}
	res, err := e.client.Search().Index("borg").Type("problem").From(0).Size(size).Query(
		elastic.NewMultiMatchQuery(q).FieldWithBoost("Title", 5.0).Field("Solutions.Body")).Do()
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
