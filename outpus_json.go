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
