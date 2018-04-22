package crawler

import (
	"bytes"
	"go_test/inerfun"
	"log"
	"os"
	"strconv"

	"golang.org/x/net/html"
)

type QueryData struct {
	Code  int
	Store string
	Data  []ParseData
}

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

//ParseDataGroup : use to implement the interface used for sort.Sort
type ParseDataGroup []ParseData

func (p ParseDataGroup) Len() int {
	return len(p)
}
func (p ParseDataGroup) Less(i, j int) bool {
	pi, err := strconv.Atoi(p[i].Price)
	pj, err := strconv.Atoi(p[j].Price)
	if err != nil {
		return false
	}
	return pi < pj
}
func (p ParseDataGroup) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type OperationCmd string

type CrawleCmd struct {
	Cmd       OperationCmd
	Parameter map[string]string
}

//BaseCrawler use for crawlering different website
type BaseCrawler interface {
	GetStoreName() string
	StartCrawlering(CrawleCmd) interface{}
	GetCrawlingData(interface{}) []ParseData
	GetProductDetail(ParseData) interface{}
	MakeCrawCmd(map[string]interface{}) CrawleCmd
}

// FuncParseHtml : this kind of function is used to parse html document
type FuncParseHtml func(*html.Tokenizer) (interface{}, interface{})

// ParseHtmlDoc : parse the html file
func ParseHtmlDoc(fname string, rf FuncParseHtml) (interface{}, interface{}) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	domDoc := html.NewTokenizer(f)
	return rf(domDoc)
}

// parseRespHtml : send http request and parse the response
// 		method: support GET, POST
// 		apiurl: the url
//		para: request parameter
//		rf:  function used to parse response
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
