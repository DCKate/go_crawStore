package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type nodeState int

const (
	nodeNone    = 0
	nodeList    = 1
	nodeProduct = 2
	nodeDetail  = 3

	RtStoreName     = "rt-mart"
	rtCrawlerDomain = "https://www.rt-mart.com.tw/direct/index.php"
	//RtSearch used to call the action of search on rt-mart
	RtSearch OperationCmd = "product_search"
	RtHtml   OperationCmd = "parse_rtmart_html"
)

type rtUrl struct {
	Detail string
	Img    string
}

type rtData struct {
	ID           string `json:"id"`
	PictureURL   string `json:"pictureUrl"`
	Price        string `json:"price"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	ItemPosition string `json:"position"`
	Dimension    string `json:"dimension2"`
	// Data                 json.RawMessage
}

type rtSearchContent struct {
	Currency string   `json:"currencyCode"`
	Products []rtData `json:"impressions"`
	// Data       json.RawMessage
}

type rtSearchResp struct {
	Ecommerce rtSearchContent `json:"ecommerce"`
	// Data     json.RawMessage
}

type rtProduct struct {
	Store      string `json:"store"`
	ID         string `json:"id"`
	SeName     string `json:"sename"`
	PictureURL string `json:"pictureUrl"`
	Price      string `json:"price"`
	Name       string `json:"name"`
	Category   string `json:"profile"`
}

//RtCrawler implement Crawler's interface and cache the last query data(QueryProducts)
type RtCrawler struct {
}

//GetStoreName get the store name of this crawler used to
func (cr RtCrawler) GetStoreName() string {
	return RtStoreName
}

func (cr RtCrawler) MakeCrawCmd(para map[string]interface{}) CrawleCmd {
	cmd := CrawleCmd{
		Cmd:       "",
		Parameter: make(map[string]string),
	}
	if vv, ok := para["cmd"]; ok {
		switch vv {
		case "search":
			cmd.Cmd = RtSearch
			cmd.Parameter["prod_keyword"] = para["key"].(string)
		}
	}
	if vv, ok := para["price_range"]; ok {
		cmd.Parameter["price_range"] = strconv.Itoa(vv.(int))
	}
	return cmd
}

func (cr RtCrawler) StartCrawlering(data CrawleCmd) interface{} {
	switch data.Cmd {
	case RtSearch:
		domain := fmt.Sprintf("%s?action=%s", rtCrawlerDomain, RtSearch)
		st, urls, rtd := parseRespHtml("POST", domain, data.Parameter, ParseRtHtml)
		if st == 0 {
			return cr.makeQueryProduct(urls.(map[int]rtUrl), rtd.([]rtData))
		}
	case RtHtml:
		urls, rtd := ParseHtmlDoc(data.Parameter["file"], ParseRtHtml)
		return cr.makeQueryProduct(urls.(map[int]rtUrl), rtd.([]rtData))
	}
	return nil
}

func (cr RtCrawler) GetCrawlingData(data interface{}) []ParseData {
	qpros := data.([]rtProduct)
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
				Profile: vv.Category,
				SeName:  vv.SeName,
			}
			tda[ii] = tmp
		}
	}
	return tda
}

func (cr RtCrawler) GetProductDetail(pro ParseData) interface{} {
	return fmt.Sprintf("%s?action=product_detail&prod_no=%s", rtCrawlerDomain, pro.SeName)
}

func (crw RtCrawler) makeQueryProduct(purls map[int]rtUrl, pitems []rtData) []rtProduct {
	qpros := make([]rtProduct, len(pitems))
	for ii, vv := range pitems {
		tmp := rtProduct{
			Store:      RtStoreName,
			ID:         vv.ID,
			SeName:     "",
			PictureURL: "",
			Price:      vv.Price,
			Name:       vv.Name,
			Category:   vv.Category,
		}
		ind, err := strconv.Atoi(vv.ItemPosition)
		if err == nil {
			if mm, ok := purls[ind]; ok {
				tmp.PictureURL = mm.Img
				re := regexp.MustCompile("P[0-9]+")
				serial := re.FindString(mm.Detail)
				tmp.SeName = serial
				// pitems[ii] = vv
			}
		}
		qpros[ii] = tmp
	}
	return qpros
}

func tokenGetAttr(token html.Token, tkdata string, attkey string, attval string) (bool, string) {
	if token.Data == tkdata {
		for _, a := range token.Attr {
			if a.Key == attkey && (len(attval) == 0 || a.Val == attval) {
				return true, a.Val
			}
		}
	}
	return false, ""
}

// ParseRtHtml : use to parse the html response return by rt-mart
func ParseRtHtml(domDoc *html.Tokenizer) (interface{}, interface{}) {
	var jsitems rtSearchResp
	donitems := make(map[int]rtUrl)

	token := domDoc.Token()
	var nst nodeState = nodeNone
	iDetailNode := -1
	for {
		tt := domDoc.Next()
		switch tt {
		case html.ErrorToken:
			return donitems, jsitems.Ecommerce.Products
		case html.StartTagToken:
			token = domDoc.Token()
			if token.Data == "li" {
				nst = nodeList
			}
			if nst >= nodeList {
				if nst == nodeProduct {
					if ok, _ := tokenGetAttr(token, "div", "class", "for_imgbox"); ok {
						iDetailNode++
						nst = nodeDetail
					}
				} else {
					if ok, _ := tokenGetAttr(token, "div", "class", "indexProList"); ok {
						nst = nodeProduct
					}
				}

				if ok, href := tokenGetAttr(token, "a", "href", ""); ok {
					if vv, yy := donitems[iDetailNode]; yy {
						vv.Detail = href
						donitems[iDetailNode] = vv
					} else {
						donitems[iDetailNode] = rtUrl{
							Detail: href,
							Img:    "",
						}
					}
				}
			}
		case html.SelfClosingTagToken:
			token = domDoc.Token()
			if nst >= nodeProduct {
				if ok, isrc := tokenGetAttr(token, "img", "src", ""); ok {
					if vv, ok := donitems[iDetailNode]; ok {
						vv.Img = isrc
						donitems[iDetailNode] = vv
					} else {
						donitems[iDetailNode] = rtUrl{
							Detail: "",
							Img:    isrc,
						}
					}
				}
			}

		case html.EndTagToken:
			if token.Data == "li" {
				nst = nodeNone
			}

		case html.TextToken:
			if token.Data != "script" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDoc.Text())))

			if len(TxtContent) > 0 && strings.Contains(TxtContent, "dataLayer.push(") {
				TxtContent = strings.TrimPrefix(TxtContent, "dataLayer.push(")
				TxtContent = strings.TrimRight(TxtContent, ");")
				TxtContent = strings.Replace(TxtContent, "\\'", "'", -1)
				if err := json.Unmarshal([]byte(TxtContent), &jsitems); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
