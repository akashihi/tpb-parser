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
