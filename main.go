package main

const (
	tpbUrl = "https://thepiratebay.cr"
)

func main() {
	InitLog()
	log.Info("Starting tpb parser...")

	log.Info("We will look for TPB at https://thepiratebay.cr")

	p := newPaginator(tpbUrl)

	getTopCategories(tpbUrl, p)

	p.run()
}
