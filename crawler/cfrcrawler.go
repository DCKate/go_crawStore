package crawler

import (
	"encoding/json"
	"fmt"
	"go_test/inerfun"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	crfCrawlerDomain = "https://online.carrefour.com.tw"
	//RtSearch used to call the action of search on rt-mart
	CrfSearch OperationCmd = "CarrefourECProduct/GetSearchJson"
	CrfHtml   OperationCmd = "parse_carrefour_html"
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

type crfProduct struct {
	Store       string `json:"store"`
	ID          string `json:"id"`
	SeName      string `json:"sename"`
	PictureURL  string `json:"pictureUrl"`
	Price       string `json:"price"`
	Name        string `json:"name"`
	ItemPerPack string `json:"profile"`
	Category    string `json:"note"`
}

//https://online.carrefour.com.tw/search\?key\=%E8%9D%A6\&categoryId\=1
type CfrCrawler struct {
	QueryProducts []crfProduct
}

func (cr *CfrCrawler) StartCrawlering(data CrawleCmd) {
	switch data.Cmd {
	case CrfSearch:
		domain := fmt.Sprintf("%s/%s", crfCrawlerDomain, CrfSearch)
		rtd := parseCfrRespJson(domain, data.Parameter)
		if len(rtd.Products) > 0 {
			cr.makeQueryProduct(rtd.Products)
		}
	case CrfHtml:
		rtd, _ := ParseHtmlDoc(data.Parameter["file"], parseCfrHtml)
		cr.makeQueryProduct(rtd.(crfSearchContent).Products)
	}
}

func (cr *CfrCrawler) GetCrawlingData(cond map[string]string) []ParseData {
	tda := make([]ParseData, len(cr.QueryProducts))
	jdata, err := json.Marshal(cr.QueryProducts)
	if err == nil {
		err = json.Unmarshal(jdata, &tda)
	}
	if err != nil {
		for ii, vv := range cr.QueryProducts {
			tmp := ParseData{
				Store:   vv.Store,
				ID:      vv.ID,
				Name:    vv.Name,
				Price:   vv.Price,
				Image:   vv.PictureURL,
				Profile: vv.ItemPerPack,
				SeName:  vv.SeName,
				Note:    vv.Category,
			}
			tda[ii] = tmp
		}
	}
	return tda
}

func (cr *CfrCrawler) GetProductDetal(pro ParseData) interface{} {
	return fmt.Sprintf("%s/%s", crfCrawlerDomain, pro.SeName)
}

func (crw *CfrCrawler) makeQueryProduct(pitems []crfData) int {
	crw.QueryProducts = make([]crfProduct, len(pitems))
	for ii, vv := range pitems {
		tmp := crfProduct{
			Store:       "CARREFOUR",
			ID:          strconv.Itoa(vv.ID),
			SeName:      vv.SeName,
			PictureURL:  vv.PictureURL,
			Price:       vv.Price,
			Name:        vv.Name,
			ItemPerPack: vv.ItemQtyPerPackFormat,
			Category:    vv.Specification,
		}
		crw.QueryProducts[ii] = tmp
	}
	return len(crw.QueryProducts)
}

//https://online.carrefour.com.tw/search\?key\=%E8%9D%A6\&categoryId\=1

func parseCfrHtml(domDoc *html.Tokenizer) (interface{}, interface{}) {
	var rda crfSearchContent
	previousStartToken := domDoc.Token()
	for {
		tt := domDoc.Next()
		switch {
		case tt == html.ErrorToken:
			return rda, nil
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
					if err := json.Unmarshal([]byte(jstr), &rda); err != nil {
						log.Println(err)
					}
				}

			}

		}
	}
}

func parseCfrRespJson(apiurl string, postform map[string]string) crfSearchContent {
	var rda crfSearchContent
	// apiurl := "https://online.carrefour.com.tw/CarrefourECProduct/GetSearchJson"
	// postform := map[string]string{"key": "澳洲梅花牛排", "orderBy": "0", "pageSize": "2", "pageIndex": "1", "minPrice": "0", "maxPrice": "1000"}
	aa, bb := inerfun.MakePostForm(apiurl, nil, postform)
	if aa == 200 {
		if err := json.Unmarshal(bb, &rda); err != nil {
			log.Println(err)
		}
		// log.Printf("%v", crfSearchContent)
	}
	return rda
}

// func parseCfrRespHtml(apiurl string, para map[string]string) {
// 	// apiurl := "https://online.carrefour.com.tw/search"
// 	// para := map[string]string{"key": "澳洲梅花牛排"}
// 	aa, bb := inerfun.MakeGet(apiurl, nil, para)
// 	if aa == 200 {
// 		domDoc := html.NewTokenizer(bytes.NewReader(bb))
// 		parseCfrHtml(domDoc)
// 	}
// }
