package main

import (
	"bytes"
	"encoding/csv"
	//"fmt"
	"io/ioutil"
	//"regexp"
	"strconv"
	//"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/crufter/borg/types"
	"github.com/olivere/elastic"
	"sort"
)

func main() {
	b1, err := ioutil.ReadFile("./TopBashOver4.csv")
	if err != nil {
		panic(err)
	}
	b2, err := ioutil.ReadFile("./AnswersToTopBashOver4.csv")
	if err != nil {
		panic(err)
	}
	qrs, err := csv.NewReader(bytes.NewReader(b1)).ReadAll()
	if err != nil {
		panic(err)
	}
	ars, err := csv.NewReader(bytes.NewReader(b2)).ReadAll()
	if err != nil {
		panic(err)
	}
	qs := []types.Problem{}
	for _, v := range qrs {
		qs = append(qs, types.Problem{
			Id:    v[0],
			Title: v[1],
			ImportMeta: types.ImportMeta{
				Source: 0,
				Id:     v[0],
			},
		})
	}
	as := map[string][]types.Solution{}
	for i, v := range ars {
		if i == 0 {
			continue
		}
		score, err := strconv.ParseInt(v[1], 10, 64)
		if err != nil {
			panic(err)
		}
		bodies := []string{}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(v[0])))
		doc.Find("pre code").Each((func(i int, s *goquery.Selection) {
			bodies = append(bodies, s.Text())
		}))
		if len(bodies) == 0 {
			doc.Find("code").Each((func(i int, s *goquery.Selection) {
				bodies = append(bodies, s.Text())
			}))
		}
		if len(bodies) == 0 {
			continue
		}
		as[v[2]] = append(as[v[2]], types.Solution{
			Score: int(score),
			Body:  bodies,
		})

	}
	for i, problem := range qs {
		problem.Solutions = as[problem.Id]
		sort.Sort(types.Solutions(problem.Solutions))
		qs[i] = problem
	}
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://192.168.99.100:9200"))
	if err != nil {
		panic(err)
	}
	client.CreateIndex("borg").Do()
	// Add a document to the index
	for _, p := range qs {
		_, err = client.Index().
			Index("borg").
			Type("problem").
			Id(p.Id).
			BodyJson(p).
			Refresh(true).
			Do()
		if err != nil {
			// Handle error
			panic(err)
		}
	}
}
