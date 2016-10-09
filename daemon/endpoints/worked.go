package endpoints

import (
	"gopkg.in/olivere/elastic.v3"
)

var workedScript = `
if (ctx._source.containsKey("worked")) {
	if (!ctx._source.worked.contains(query)) {
		ctx._source.worked += query;
	}
} else {
	ctx._source.worked = [query]
}
`

// Worked tells the borg server that a result works for a given query
func (e Endpoints) Worked(id, query string) error {
	_, err := e.client.Update().
		Index("borg").
		Type("problem").
		Id(id).Script(elastic.NewScriptInline(workedScript).Param("query", query)). // possibly injection?
		Do()
	return err
}
