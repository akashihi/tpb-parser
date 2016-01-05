package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	streams = 1
)

type TorrentEntry struct {
	title       string
	size        int
	files       int
	category    string
	subcategory string
	by          string
	hash        string
	uploaded    time.Time
}

type Downloader struct {
	topUrl    string
	initialId int
	pageId    chan int
}

func newDownloader(tU string, id int) *Downloader {
	return &Downloader{
		topUrl:    tU,
		initialId: id,
		pageId:    make(chan int),
	}
}

func (d *Downloader) run() {
	for w := 0; w <= streams; w++ {
		go d.processPage()
	}
	for w := d.initialId; w >= 0; w-- {
		d.pageId <- w
	}
	close(d.pageId)
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
			return
		}

		torrentData := doc.Find("#detailsframe")

		torrent := TorrentEntry{}
		torrent.processTitle(torrentData)
		torrent.processFirstColumn(torrentData)
		torrent.processSecondColumn(torrentData)
		torrent.processHash(torrentData)

		log.Info("Torrent title is: %s", torrent.title)
	}

	/*if doc.Is("#main-content p.info") {
		log.Info("Empty page, finishing paginator processing")
		return
	}

	u, pU := doc.Find("img[alt=\"Next\"]").Parent().Attr("href")
	if pU {
		log.Info("Found next page at %s", u)
		p.queue(u)
	}*/

}

func (t *TorrentEntry) processTitle(torrentData *goquery.Selection) {
	t.title, _ = torrentData.Find("#title").Html()
	t.title = strings.TrimSpace(t.title)
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
		t.category = "Unknown"
		t.subcategory = "Unknown"
	} else {
		t.category = strings.TrimSpace(categories[0])
		t.subcategory = strings.TrimSpace(categories[1])
	}

	//Files
	filesData, _ := data.Eq(1).Next().Find("a").Html()
	files, err := strconv.Atoi(filesData)
	if err != nil {
		log.Warning("Can't retrieve number of files in torrent")
		t.files = -1
	} else {
		t.files = files
	}

	//Size
	sizeData, _ := data.Eq(2).Next().Html()
	sizeData = strings.TrimSpace(sizeData)
	rx, _ := regexp.Compile(`\((\d+)`)
	if rx.MatchString(sizeData) {
		size, err := strconv.Atoi(rx.FindStringSubmatch(sizeData)[1])
		if err != nil {
			log.Warning("Can't retrieve size of files in the torrent")
			t.size = -1
		} else {
			t.size = size
		}
	} else {
		log.Warning("Can't parse size of files in the torrent")
		t.size = -1
	}
}

func (t *TorrentEntry) processSecondColumn(torrentData *goquery.Selection) {
	data := torrentData.Find("#details .col2 dt")
	if data.Size() < 2 {
		log.Warning("Not enough data to parse in second coumnt")
	}

	//Uploaded
	uploadedData := data.First().Next().Text()
	t.uploaded, _ = time.Parse("2006-01-02 15:04:05 MST", uploadedData)

	//By
	t.by = strings.TrimSpace(data.Eq(1).Next().Find("a").Text())
}

func (t *TorrentEntry) processHash(torrentData *goquery.Selection) {
	data := torrentData.Find("#details .col2 dt")
	if data.Size() < 2 {
		log.Warning("Not enough data to parse in second coumnt")
	}

	//Uploaded
	uploadedData := data.First().Next().Text()
	t.uploaded, _ = time.Parse("2006-01-02 15:04:05 MST", uploadedData)

	//By
	t.by = strings.TrimSpace(data.Eq(1).Next().Find("a").Text())
}
