/*
   tpb-parser
   Copyright (C) 2016 Denis V Chapligin <akashihi@gmail.com>
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
