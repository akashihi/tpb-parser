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
	w := <-p.categoryPages
	p.categoryPages <- w
	for {
		select {
		case page, ok := <-p.categoryPages:
			if !ok {
				log.Warning("Paginator channel closed")
				return
			}
			p.processPaginator(page)
		default:
			log.Info("All category pages pagination processed")
			close(p.categoryPages)
		}
	}
}

func (p *Paginator) processPaginator(page string) {
	var url bytes.Buffer
	url.WriteString(p.topUrl)
	url.WriteString(page)

	log.Info("Parsing paginator at: %s", url.String())
	doc, err := goquery.NewDocument(url.String())

	if err != nil {
		log.Critical("Can't download category page from TPB: %v", err)
		return
	}

	if doc.Is("#main-content p.info") {
		log.Info("Empty page, finishing paginator processing")
		return
	}

	u, pU := doc.Find("img[alt=\"Next\"]").Parent().Attr("href")
	if pU {
		log.Info("Found next page at %s", u)
		p.queue(u)
	}
}
