package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

func getTopCategories(topUrl string, p *Paginator) {
	var url bytes.Buffer
	url.WriteString(topUrl)
	url.WriteString("/browse")

	log.Info("Processing top categories at: %s", url.String())
	topDoc, err := goquery.NewDocument(url.String())
	if err != nil {
		log.Critical("Can't dowbload top categories page from TPB: %v", err)
		return
	}

	topDoc.Find(".categoriesContainer dt a").Each(func(i int, s *goquery.Selection) {
		t, pT := s.Attr("title")
		u, pU := s.Attr("href")
		if pU && pT {
			log.Info("Top category %s href: %s", t, u)
			p.queue(u)
		}
	})
}
