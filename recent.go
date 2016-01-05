package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

func getRecentId(topUrl string) int {
	var url bytes.Buffer
	url.WriteString(topUrl)
	url.WriteString("/recent")

	log.Info("Processing recent torrents page at: %s", url.String())
	doc, err := goquery.NewDocument(url.String())
	if err != nil {
		log.Critical("Can't dowbload recent torrents page from TPB: %v", err)
		return 0
	}

	topTorrent := doc.Find("#searchResult .detName a").First()
	t, pT := topTorrent.Attr("title")
	u, pU := topTorrent.Attr("href")
	if pT && pU {
		rx, _ := regexp.Compile(`\/torrent\/(\d+)\/.*`)
		if rx.MatchString(u) {
			id, err := strconv.Atoi(rx.FindStringSubmatch(u)[1])
			if err != nil {
				log.Critical("Can't retrieve latest torrent id")
				return 0
			}
			log.Info("The most recent torrent is %s and it's id is %d", t, id)
			return id
		}
	}
	return 0
}
