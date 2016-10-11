package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	//"regexp"
	"strconv"
	//"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/ok-borg/borg/types"
	"gopkg.in/olivere/elastic.v3"
	"sort"
	"strings"
)

func problems() []types.Problem {
	b1, err := ioutil.ReadFile("./QuestionsGitOver4.csv")
	if err != nil {
		panic(err)
	}
	qrs, err := csv.NewReader(bytes.NewReader(b1)).ReadAll()
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
	return qs
}

func solutions() map[string][]types.Solution {
	b2, err := ioutil.ReadFile("./AnswersToGitOver4.csv")
	if err != nil {
		panic(err)
	}
	ars, err := csv.NewReader(bytes.NewReader(b2)).ReadAll()
	if err != nil {
		panic(err)
	}
	as := map[string][]types.Solution{}
	filtered := 0
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
		if err != nil {
			panic(err)
		}
		doc.Find("pre code").Each((func(i int, s *goquery.Selection) {
			bodies = append(bodies, s.Text())
		}))
		if len(bodies) == 0 {
			doc.Find("code").Each((func(i int, s *goquery.Selection) {
				bodies = append(bodies, s.Text())
			}))
		}
		bodies = func(bs []string) []string {
			ret := []string{}
			for _, v := range bs {
				if len(strings.TrimSpace(v)) > 0 {
					ret = append(ret, v)
				} else {
					filtered++
				}
			}
			return ret
		}(bodies)
		if len(bodies) == 0 {
			continue
		}
		as[v[2]] = append(as[v[2]], types.Solution{
			Score: int(score),
			Body:  bodies,
		})
	}
	fmt.Println(fmt.Sprintf("Filtered out %v bodies", filtered))
	return as
}

func main() {
	qs := problems()
	as := solutions()
	for i, problem := range qs {
		problem.Solutions = as[problem.Id]
		sort.Sort(types.Solutions(problem.Solutions))
		qs[i] = problem
	}
	filtered := 0
	origiSize := len(qs)
	qs = func(fqs []types.Problem) []types.Problem {
		ret := []types.Problem{}
		for _, v := range fqs {
			if len(v.Solutions) > 0 {
				ret = append(ret, v)
				fmt.Println(v)
			} else {
				filtered++
			}
		}
		return ret
	}(qs)
	fmt.Println(fmt.Sprintf("Filtered %v ps out of %v", filtered, origiSize))
	saveProblems(qs)
}

func saveProblems(qs []types.Problem) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://borg.crufter.com:9200"))
	if err != nil {
		panic(err)
	}
	//client.CreateIndex("borg").Do()
	// Add a document to the index
	for i, p := range qs {
		if i%100 == 0 {
			fmt.Println(fmt.Sprintf("Done %v out of %v", i, len(qs)))
		}
		_, err := client.Index().
			Index("borg").
			Type("problem").
			Id(p.Id).
			BodyJson(p).
			Refresh(true).
			Do()
		//_, err := client.Delete().
		//	Index("borg").
		//	Type("problem").
		//	Id(p.Id).Do()
		if err != nil {
			// Handle error
			panic(err)
		}
	}
}
