package sitemap

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/crufter/slugify"
	"github.com/joeguo/sitemap"
	"github.com/ok-borg/borg/types"
	"gopkg.in/olivere/elastic.v3"
	"reflect"
	"time"
)

// GenerateSitemap grabs all entries (now only ones added with borg) and saves a sitemap.xml.gz file in `sitemapPath`
func GenerateSitemap(sitemapPath string, client *elastic.Client) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Sitemap generation failed: %v", r)
		}
	}()
	// this query is because we only want to show user submitted content for now - not ones scraped from somewhere else - to not piss of google
	// @TODO include ones which were changed substantially
	// @TODO this is going to get dog slow
	res, err := client.Search().Query(elastic.NewRegexpQuery("CreatedBy", ".{3,}")).Size(500).Do()
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
	items := []*sitemap.Item{}
	for _, v := range all {
		item := &sitemap.Item{
			Loc:        "https://ok-b.org/t/" + fmt.Sprintf("%v/%v", v.Id, slugify.S(v.Title)),
			LastMod:    time.Now(),
			Priority:   0.5,
			Changefreq: "daily",
		}
		items = append(items, item)
	}
	err = sitemap.SiteMap(sitemapPath+"/sitemap.xml.gz", items)
	if err != nil {
		panic(err)
	}
	log.Info("Generated sitemap successfully")
}
