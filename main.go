package main

const (
	tpbUrl = "https://thepiratebay.cr"
)

func main() {
	InitLog()
	log.Info("Starting tpb parser...")

	log.Info("We will look for TPB at https://thepiratebay.cr")

	d := newDownloader(tpbUrl, getRecentId(tpbUrl))
	d.run()
	/*p := newPaginator(tpbUrl)

	getTopCategories(tpbUrl, p)

	p.run()*/
}
