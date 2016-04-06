# tpb-parser

## What is this?

thepiratebay.cr database downloader. Parses The Pirate Bay site and stores it's data into csv or json files.

Written as a exercise in Go, so do not expect too much.

## Building it

1. Install [go](http://golang.org/doc/install)

2. Install "goquery" go get -u github.com/PuerkitoBio/goquery

4. Compile tpb-parser

        git clone git://github.com/akashihi/tpb-parser.git
        cd tpb-parser
        go build .

## Running it

Generally:

    tpb-parser -output tpb.csv 
or 
    tpb-parser -output tpb.json -json

All parameters could be omited. Run with --help to get parameters description

## License 

See LICENSE file.

Copyright 2016 Denis V Chapligin <akashihi@gmail.com>
