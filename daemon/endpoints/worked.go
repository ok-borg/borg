package endpoints

import (
	"gopkg.in/olivere/elastic.v3"
)

// Worked tells the borg server that a result works for a given query
func (e Endpoints) Worked(id, query string) error {
	_, err := e.client.Update().
		Index("borg").
		Type("problem").
		Id(id).Script(elastic.NewScriptInline("ctx._source.worked += query").Param("query", query)). // possibly injection?
		Do()
	return err
}
