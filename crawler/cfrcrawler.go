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
	CrfStoreName     = "carrefour"
	crfCrawlerDomain = "https://online.carrefour.com.tw"
	//RtSearch used to call the action of search on rt-mart
	CrfSearch OperationCmd = "CarrefourECProduct/GetSearchJson"
	CrfHtml   OperationCmd = "parse_carrefour_html"
)

var priceMap = map[int][]string{
	1: []string{"0", "500"},
	2: []string{"501", "1000"},
	3: []string{"1001", "3000"},
	4: []string{"3001", "10000"},
	5: []string{"10001", "30000"},
	6: []string{"30001", "0"},
}

// the following stucture is json format return by carrefour
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

// crfProduct : the product info used in cfrcrawler, and the key of json format is consistent with crawler.ParseData
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

type CfrCrawler struct {
}

func (cr CfrCrawler) GetStoreName() string {
	return CrfStoreName
}

func (cr CfrCrawler) MakeCrawCmd(para map[string]interface{}) CrawleCmd {

	cmd := CrawleCmd{
		Cmd:       "",
		Parameter: make(map[string]string),
	}
	if vv, ok := para["cmd"]; ok {
		switch vv {
		case "search":
			cmd.Cmd = CrfSearch
			cmd.Parameter["key"] = para["key"].(string)
		}
	}
	if vv, ok := para["price_range"]; ok {
		ranges := priceMap[vv.(int)]
		cmd.Parameter["minPrice"] = ranges[0]
		if ranges[1] != "0" {
			cmd.Parameter["maxPrice"] = ranges[1]
		}
	}
	// log.Println(cmd)
	return cmd
}

func (cr CfrCrawler) StartCrawlering(data CrawleCmd) interface{} {
	switch data.Cmd {
	case CrfSearch:
		domain := fmt.Sprintf("%s/%s", crfCrawlerDomain, CrfSearch)
		rtd := parseCfrRespJson(domain, data.Parameter)
		if len(rtd.Products) > 0 {
			return cr.makeQueryProduct(rtd.Products)
		}
	case CrfHtml:
		rtd, _ := ParseHtmlDoc(data.Parameter["file"], parseCfrHtml)
		return cr.makeQueryProduct(rtd.(crfSearchContent).Products)
	}
	return nil
}

func (cr CfrCrawler) GetCrawlingData(data interface{}) []ParseData {
	qpros := data.([]crfProduct)
	tda := make([]ParseData, len(qpros))
	jdata, err := json.Marshal(qpros)
	if err == nil {
		err = json.Unmarshal(jdata, &tda)
	}
	if err != nil {
		for ii, vv := range qpros {
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

func (cr CfrCrawler) GetProductDetail(pro ParseData) interface{} {
	return fmt.Sprintf("%s/%s", crfCrawlerDomain, pro.SeName)
}

func (crw CfrCrawler) makeQueryProduct(pitems []crfData) []crfProduct {
	qpros := make([]crfProduct, len(pitems))
	for ii, vv := range pitems {
		tmp := crfProduct{
			Store:       CrfStoreName,
			ID:          strconv.Itoa(vv.ID),
			SeName:      vv.SeName,
			PictureURL:  vv.PictureURL,
			Price:       vv.Price,
			Name:        vv.Name,
			ItemPerPack: vv.ItemQtyPerPackFormat,
			Category:    vv.Specification,
		}
		qpros[ii] = tmp
	}
	return qpros
}

// parseCfrHtml : use to parse the html response return by carrefour
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

// parseCfrRespJson : carrefour is support the api for returning json data
func parseCfrRespJson(apiurl string, postform map[string]string) crfSearchContent {
	var rda crfSearchResp
	aa, bb := inerfun.MakePostForm(apiurl, nil, postform)
	if aa == 200 {
		if err := json.Unmarshal(bb, &rda); err != nil {
			log.Println(err)
		}
	}
	return rda.Contents
}

// note: how to search in carrefour
// json url : https://online.carrefour.com.tw/CarrefourECProduct/GetSearchJson
// general url: https://online.carrefour.com.tw/search
// requset parameter: map[string]string{"key": "澳洲梅花牛排", "orderBy": "0", "pageSize": "2", "pageIndex": "1", "minPrice": "0", "maxPrice": "1000"}
