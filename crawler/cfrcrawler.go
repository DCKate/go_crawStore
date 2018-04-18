package crawler

import (
	"bytes"
	"encoding/json"
	"go_test/inerfun"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type crfData struct {
	ID                   int    `json:"Id"`
	PictureURL           string `json:"PictureUrl"`
	Price                string `json:"Price"`
	Name                 string `json:"Name"`
	ItemQtyPerPack       int    `json:"ItemQtyPerPack"`
	ItemQtyPerPackFormat string `json:"ItemQtyPerPackFormat"`
	Specification        string `json:"Specification"`
	SeName               string `json:"SeName"`
	// Data                 json.RawMessage
}

type crfSearchContent struct {
	CategoryID int       `json:"CategoryId"`
	CurrentID  int       `json:"CurrentCategoryId"`
	SearchKey  string    `json:"Key"`
	Products   []crfData `json:"ProductListModel"`
	// Data       json.RawMessage
}

type crfSearchResp struct {
	Status   int              `json:"success"`
	Contents crfSearchContent `json:"content"`
	// Data     json.RawMessage
}

//https://online.carrefour.com.tw/search\?key\=%E8%9D%A6\&categoryId\=1
type CfrCrawler struct {
}

func (*CfrCrawler) StartCrawlering(key map[string]string) {

}

// func (*CfrCrawler) GetCrawlingData(cond map[string]string) []ParseData {

// }

//https://online.carrefour.com.tw/search\?key\=%E8%9D%A6\&categoryId\=1

func parseCfrHtml(domDoc *html.Tokenizer) {
	previousStartToken := domDoc.Token()
	for {
		tt := domDoc.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			previousStartToken = domDoc.Token()
		case tt == html.TextToken:
			if previousStartToken.Data != "script" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDoc.Text())))
			if len(TxtContent) > 0 && strings.Contains(TxtContent, "searchProductListModel=") {
				re := regexp.MustCompile("searchProductListModel=\"(.*?)\";")
				match := re.FindStringSubmatch(TxtContent)
				if len(match) > 1 {
					jstr := match[1]
					jstr = strings.Replace(jstr, "\\", "", -1)
					var dd crfSearchContent
					if err := json.Unmarshal([]byte(jstr), &dd); err != nil {
						log.Println(err)
					}
				}

			}

		}
	}
}

func parseCfrRespHtml() { //url string, urlbody []byte) {
	apiurl := "https://online.carrefour.com.tw/search"
	para := map[string]string{"key": "澳洲梅花牛排"}

	aa, bb := inerfun.MakeGet(apiurl, nil, para)
	if aa == 200 {
		domDoc := html.NewTokenizer(bytes.NewReader(bb))
		parseCfrHtml(domDoc)

	}

}

func parseCfrRespJson(apiurl string, postform map[string]string) {
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
