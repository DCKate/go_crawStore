package crawler

import (
	"bytes"
	"encoding/json"
	"go_test/inerfun"
	"log"
	"os"

	"golang.org/x/net/html"
)

//ParseData the date after crawlering
type ParseData struct {
	Store   string `json:"store"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Image   string `json:"pictureUrl"`
	Profile string `json:"profile"`
	SeName  string `json:"sename"`
	Note    string `json:"note"`
}

type OperationCmd string

type CrawleCmd struct {
	Cmd       OperationCmd
	Parameter map[string]string
}

//BaseCrawler use for crawlering different website
type BaseCrawler interface {
	StartCrawlering(CrawleCmd)
	GetCrawlingData(map[string]string) []ParseData
	GetProductDetal(ParseData) interface{}
}

type FuncParseHtml func(*html.Tokenizer) (interface{}, interface{}) //(map[int]rturl, []rtData)

func ParseHtmlDoc(fname string, rf FuncParseHtml) (interface{}, interface{}) { //(map[int]rturl, []rtData) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	domDoc := html.NewTokenizer(f)
	return rf(domDoc)
}

func parseRespHtml(method string, apiurl string, para map[string]string, rf FuncParseHtml) (int, interface{}, interface{}) {
	var (
		status int
		data   []byte
	)

	if method == "POST" {
		status, data = inerfun.MakePostForm(apiurl, nil, para)
	} else {
		status, data = inerfun.MakeGet(apiurl, nil, para)
	}

	if status == 200 {
		domDoc := html.NewTokenizer(bytes.NewReader(data))
		rm, rd := rf(domDoc)
		return 0, rm, rd
	}
	return -1, nil, nil
}

func parsePostRespJson(apiurl string, postform map[string]string) {
	// apiurl := "https://online.carrefour.com.tw/CarrefourECProduct/GetSearchJson"
	// postform := map[string]string{"key": "澳洲梅花牛排", "orderBy": "0", "pageSize": "2", "pageIndex": "1", "minPrice": "0", "maxPrice": "1000"}
	aa, bb := inerfun.MakePostForm(apiurl, nil, postform)
	log.Printf("%v\n", string(bb))
	if aa == 200 {
		var dd crfSearchResp
		if err := json.Unmarshal(bb, &dd); err != nil {
			log.Println(err)
		}
		log.Printf("%v", dd)
	}
}
