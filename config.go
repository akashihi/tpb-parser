package main

import (
	"flag"
)

type Configuration struct {
	output string
	csv    bool
	json   bool
}

func config() Configuration {
	outputPtr := flag.String("outfile", "tpb.txt", "Output file name")
	csvPtr := flag.Bool("csv", true, "Write output as csv. Disables other output formats")
	jsonPtr := flag.Bool("json", false, "Write output as json. Disables other output formats")

	flag.Parse()

	return Configuration{output: *outputPtr, csv: *csvPtr, json: *jsonPtr}
}
