package crawler

import (
	"encoding/json"
	"log"
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
)

type rturl struct {
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
	ID         string `json:"id"`
	DetailURL  string `json:"detailUrl"`
	PictureURL string `json:"pictureUrl"`
	Price      string `json:"price"`
	Name       string `json:"name"`
	Category   string `json:"category"`
}

type RtCrawler struct {
	QueryProducts []rtProduct
}

// func (*RtCrawler) StartCrawlering(key map[string]string) {

// }

// func (*CfrCrawler) GetCrawlingData(cond map[string]string) []ParseData {

// }
//https://www.rt-mart.com.tw/direct/index.php?action=product_detail&prod_no=P0000200716545
//http://www.rt-mart.com.tw/direct/index.php?action=product_search

func parseRtHtml(domDoc *html.Tokenizer) (map[int]rturl, []rtData) {
	var jsitems rtSearchResp
	donitems := make(map[int]rturl)

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
				if token.Data == "div" {
					if nst == nodeProduct {
						for _, a := range token.Attr {
							if a.Key == "class" && a.Val == "for_imgbox" {
								iDetailNode++
								nst = nodeDetail
							}
						}
					} else {
						for _, a := range token.Attr {
							if a.Key == "class" && a.Val == "indexProList" {
								nst = nodeProduct
							}
						}
					}

				}
				if token.Data == "a" {
					for _, a := range token.Attr {
						if a.Key == "href" {
							if vv, ok := donitems[iDetailNode]; ok {
								vv.Detail = a.Val
								donitems[iDetailNode] = vv
							} else {
								donitems[iDetailNode] = rturl{
									Detail: a.Val,
									Img:    "",
								}
							}
						}
					}
				}
			}
		case html.SelfClosingTagToken:
			token = domDoc.Token()
			if nst >= nodeProduct {
				if token.Data == "img" {
					for _, a := range token.Attr {
						if a.Key == "src" {
							if vv, ok := donitems[iDetailNode]; ok {
								vv.Img = a.Val
								donitems[iDetailNode] = vv
							} else {
								donitems[iDetailNode] = rturl{
									Detail: "",
									Img:    a.Val,
								}
							}
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

func (crw *RtCrawler) makeQueryProduct(purls map[int]rturl, pitems []rtData) {
	for _, vv := range pitems {
		tmp := rtProduct{
			ID:         vv.ID,
			DetailURL:  "",
			PictureURL: "",
			Price:      vv.Price,
			Name:       vv.Name,
			Category:   vv.Category,
		}
		ind, err := strconv.Atoi(vv.ItemPosition)
		if err == nil {
			if mm, ok := purls[ind]; ok {
				tmp.PictureURL = mm.Img
				tmp.DetailURL = mm.Detail
				// pitems[ii] = vv
			}
		}
		crw.QueryProducts = append(crw.QueryProducts, tmp)
	}
}
