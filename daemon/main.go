package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/crufter/borg/types"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
	"net/http"
	"reflect"
	"strconv"
)

var (
	client *elastic.Client
)

func init() {
	cl, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}
	client = cl
}

func main() {
	r := httpr.New()
	r.GET("/v1/query", query)
	log.Info("Starting http server")
	log.Critical(http.ListenAndServe(fmt.Sprintf(":%v", 9992), r))
}

func query(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	size := 5
	s, err := strconv.ParseInt(r.FormValue("l"), 10, 32)
	if err == nil && s > 0 {
		size = int(s)
	}
	if size > 50 {
		size = 50
	}
	q := r.FormValue("q")
	if r.FormValue("p") == "true" {
		q = "PRIVATE"
	}
	log.Infof("Querying '%v' with size '%v'", q, size)
	res, err := client.Search().Index("borg").Type("problem").From(0).Size(size).Query(
		elastic.NewQueryStringQuery(q)).Do()
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
	fmt.Fprint(w, string(bs))
}
