package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

type Paginator struct {
	topUrl        string
	categoryPages chan string
}

func newPaginator(tU string) *Paginator {
	return &Paginator{
		topUrl:        tU,
		categoryPages: make(chan string, 256),
	}
}

func (p *Paginator) queue(page string) {
	p.categoryPages <- page
}

func (p *Paginator) run() {
	for categoryUrl := range p.categoryPages {
		var url bytes.Buffer
		url.WriteString(p.topUrl)
		url.WriteString(categoryUrl)

		log.Info("Parsing paginator at: %s", url.String())
		doc, err := goquery.NewDocument(url.String())

		if err != nil {
			log.Critical("Can't download category page from TPB: %v", err)
			continue
		}

		doc.Find("td[colspan=\"9\"]").Each(func(i int, s *goquery.Selection) {
			log.Info("Found pginator")
		})
	}
}
