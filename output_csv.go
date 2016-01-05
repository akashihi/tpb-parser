package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type CsvOutputModule struct {
	filename string
	file     io.WriteCloser
	writer   *csv.Writer
	entries  chan *TorrentEntry
}

func newCsvOutputModule(filename string) *CsvOutputModule {
	file, err := os.Create(filename)
	if err != nil {
		log.Critical("Unable to create output file: %v", err)
		os.Exit(1)
	}
	return &CsvOutputModule{
		filename: filename,
		file:     file,
		writer:   csv.NewWriter(file),
		entries:  make(chan *TorrentEntry, 1024),
	}
}

func (o *CsvOutputModule) Done() {
	o.writer.Flush()
	o.file.Close()
	close(o.entries)
}

func (o *CsvOutputModule) Put(t *TorrentEntry) {
	o.entries <- t
}

func (o *CsvOutputModule) run() {
	for entry := range o.entries {
		record := []string{
			entry.title,
			strconv.Itoa(entry.size),
			strconv.Itoa(entry.files),
			entry.category,
			entry.subcategory,
			entry.by,
			entry.hash,
			entry.uploaded.Format(time.RFC822),
			entry.magnet,
		}
		o.writer.Write(record)
	}
}
