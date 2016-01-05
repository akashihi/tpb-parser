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
