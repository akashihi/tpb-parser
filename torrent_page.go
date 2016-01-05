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
	"strings"
	"sync"
	"time"
)

const (
	streams = 128
)

type TorrentEntry struct {
	Id          int
	Title       string
	Size        int
	Files       int
	Category    string
	Subcategory string
	By          string
	Hash        string
	Uploaded    time.Time
	Magnet      string
	Info        string
}

type Downloader struct {
	topUrl    string
	initialId int
	pageId    chan int
	wg        sync.WaitGroup
	output    OutputModule
}

func newDownloader(o OutputModule, tU string, id int) *Downloader {
	return &Downloader{
		topUrl:    tU,
		initialId: id,
		pageId:    make(chan int, 1024),
		wg:        sync.WaitGroup{},
		output:    o,
	}
}

func (d *Downloader) run() {
	d.wg.Add(streams)
	for w := 0; w <= streams; w++ {
		go d.processPage()
	}
	for w := d.initialId; w >= 0; w-- {
		d.pageId <- w
	}
	close(d.pageId)
	log.Info("Processing complete, waiting for goroutines to finish")
	d.wg.Wait()
	d.output.Done()
}

func (d *Downloader) processPage() {
	for id := range d.pageId {
		var url bytes.Buffer
		url.WriteString(d.topUrl)
		url.WriteString("/torrent/")
		url.WriteString(strconv.Itoa(id))

		log.Info("Parsing torrent page at: %s", url.String())
		doc, err := goquery.NewDocument(url.String())

		if err != nil {
			log.Warning("Can't download torrent page %s from TPB: %v", url, err)
			continue
		}

		torrentData := doc.Find("#detailsframe")
		if torrentData.Length() < 1 {
			log.Warning("Erroneous torrent %d: \"%s\"", id, url.String())
			continue
		}

		torrent := TorrentEntry{Id: id}
		torrent.processTitle(torrentData)
		torrent.processFirstColumn(torrentData)
		torrent.processSecondColumn(torrentData)
		torrent.processHash(torrentData)
		torrent.processMagnet(torrentData)
		torrent.processInfo(torrentData)

		d.output.Put(&torrent)

		log.Info("Processed torrent %d: \"%s\"", id, torrent.Title)

	}
	d.wg.Done()
}

func (t *TorrentEntry) processTitle(torrentData *goquery.Selection) {
	t.Title = strings.TrimSpace(torrentData.Find("#title").Text())
}

func (t *TorrentEntry) processFirstColumn(torrentData *goquery.Selection) {
	data := torrentData.Find("#details .col1 dt")
	if data.Size() < 3 {
		log.Warning("Not enough data to parse in first coumnt")
	}

	//Categories
	categoryData, _ := data.First().Next().Find("a").Html()
	categoryData = strings.Replace(categoryData, "&gt;", ">", -1)
	categories := strings.Split(categoryData, ">")
	if len(categories) < 2 {
		log.Warning("Can't retrieve category and sub category of torrent")
		t.Category = "Unknown"
		t.Subcategory = "Unknown"
	} else {
		t.Category = strings.TrimSpace(categories[0])
		t.Subcategory = strings.TrimSpace(categories[1])
	}

	//Files
	filesData, _ := data.Eq(1).Next().Find("a").Html()
	files, err := strconv.Atoi(filesData)
	if err != nil {
		log.Warning("Can't retrieve number of files in torrent")
		t.Files = -1
	} else {
		t.Files = files
	}

	//Size
	sizeData, _ := data.Eq(2).Next().Html()
	sizeData = strings.TrimSpace(sizeData)
	rx, _ := regexp.Compile(`\((\d+)`)
	if rx.MatchString(sizeData) {
		size, err := strconv.Atoi(rx.FindStringSubmatch(sizeData)[1])
		if err != nil {
			log.Warning("Can't retrieve size of files in the torrent")
			t.Size = -1
		} else {
			t.Size = size
		}
	} else {
		log.Warning("Can't parse size of files in the torrent")
		t.Size = -1
	}
}

func (t *TorrentEntry) processSecondColumn(torrentData *goquery.Selection) {
	data := torrentData.Find("#details .col2 dt")
	if data.Size() < 2 {
		log.Warning("Not enough data to parse in second coumnt")
	}

	//Uploaded
	uploadedData := data.First().Next().Text()
	t.Uploaded, _ = time.Parse("2006-01-02 15:04:05 MST", uploadedData)

	//By
	t.By = strings.TrimSpace(data.Eq(1).Next().Find("a").Text())
}

func (t *TorrentEntry) processHash(torrentData *goquery.Selection) {
	t.Hash = strings.TrimSpace(torrentData.Find("#details .col2").Contents().Last().Text())
}

func (t *TorrentEntry) processMagnet(torrentData *goquery.Selection) {
	u, pU := torrentData.Find(".download a").First().Attr("href")
	if pU {
		t.Magnet = strings.TrimSpace(u)
	} else {
		t.Magnet = ""
	}
}

func (t *TorrentEntry) processInfo(torrentData *goquery.Selection) {
	t.Info, _ = torrentData.Find(".nfo pre").Html()
}
