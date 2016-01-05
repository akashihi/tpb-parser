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

func (o *CsvOutputModule) Run() {
	for entry := range o.entries {
		record := []string{
			strconv.Itoa(entry.Id),
			entry.Title,
			strconv.Itoa(entry.Size),
			strconv.Itoa(entry.Files),
			entry.Category,
			entry.Subcategory,
			entry.By,
			entry.Hash,
			entry.Uploaded.Format(time.RFC822),
			entry.Magnet,
		}
		o.writer.Write(record)
	}
}
