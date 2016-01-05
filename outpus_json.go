package main

import (
	"encoding/json"
	"io"
	"os"
)

type JsonOutputModule struct {
	filename string
	file     io.WriteCloser
	writer   *json.Encoder
	entries  chan *TorrentEntry
}

func newJsonOutputModule(filename string) *JsonOutputModule {
	file, err := os.Create(filename)
	if err != nil {
		log.Critical("Unable to create output file: %v", err)
		os.Exit(1)
	}
	return &JsonOutputModule{
		filename: filename,
		file:     file,
		writer:   json.NewEncoder(file),
		entries:  make(chan *TorrentEntry, 1024),
	}
}

func (o *JsonOutputModule) Done() {
	o.file.Close()
	close(o.entries)
}

func (o *JsonOutputModule) Put(t *TorrentEntry) {
	o.entries <- t
}

func (o *JsonOutputModule) Run() {
	for entry := range o.entries {
		o.writer.Encode(entry)
	}
}
