package main

type OutputModule interface {
	Put(t *TorrentEntry)
	Done()
	Run()
}
