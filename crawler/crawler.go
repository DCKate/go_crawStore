package crawler

import (
	"bytes"
	"go_test/inerfun"
	"log"
	"os"

	"golang.org/x/net/html"
)

//BaseCrawler use for crawlering different website
type BaseCrawler interface {
	StartCrawlering(map[string]string)
	GetCrawlingData(map[string]string) []ParseData
}

//ParseData the date after crawlering
type ParseData struct {
	Name    string
	Price   string
	Image   string
	Profile string
	Amount  string
	Note    string
}

type FuncParseHtml func(*html.Tokenizer) (interface{}, interface{}) //(map[int]rturl, []rtData)

func parseHtmlDoc(fname string, rf FuncParseHtml) (interface{}, interface{}) { //(map[int]rturl, []rtData) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	domDoc := html.NewTokenizer(f)
	return rf(domDoc)
}

func parseRespHtml(apiurl string, para map[string]string, rf FuncParseHtml) (int, interface{}, interface{}) {
	aa, bb := inerfun.MakeGet(apiurl, nil, para)
	if aa == 200 {
		domDoc := html.NewTokenizer(bytes.NewReader(bb))
		rm, rd := rf(domDoc)
		return 0, rm, rd
	}
	return -1, nil, nil

}
