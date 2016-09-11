package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/crufter/borg/types"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
	"io/ioutil"
	"net/http"
	"reflect"
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
	//r.PUT("/v1/problem/:id", update)
	//r.POST("/v1/problem", save)
	//r.POST()
	log.Info("Starting http server")
	log.Critical(http.ListenAndServe(fmt.Sprintf(":%v", 9992), r))
}

func query(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	res, err := client.Search().Index("borg").Type("problem").From(0).Size(5).Query(
		elastic.NewQueryStringQuery(r.FormValue("q"))).Do()
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

func save(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	b64encoded := []string{}
	err = json.Unmarshal(body, &b64encoded)
	if err != nil {
		panic(err)
	}
	ss := []types.Solution{}
	log.Debugf("Putting services %v", ss)
}
