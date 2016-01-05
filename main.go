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
	if configuration.json {
		outputModule = newJsonOutputModule(configuration.output)
	} else if configuration.csv {
		outputModule = newCsvOutputModule(configuration.output)
	}
	go outputModule.Run()

	d := newDownloader(outputModule, tpbUrl, getRecentId(tpbUrl))
	d.run()
}
