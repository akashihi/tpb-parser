package main

const (
	tpbUrl = "https://thepiratebay.cr"
)

func main() {
	InitLog()
	log.Info("Starting tpb parser...")

	log.Info("We will look for TPB at https://thepiratebay.cr")

	configuration := config()

	var outputModule OutputModule
	if configuration.csv {
		csv := newCsvOutputModule(configuration.output)
		outputModule = csv
		go csv.run()
	}

	d := newDownloader(outputModule, tpbUrl, getRecentId(tpbUrl))
	d.run()
}
