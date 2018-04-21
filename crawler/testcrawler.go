package crawler

import (
	"encoding/json"
	"log"
	"sync"
)

type queryData struct {
	Code  int
	Store string
	Data  []ParseData
}

func TestCrawler() {
	var qProducts []ParseData
	var craws = make([]BaseCrawler, 2)
	crf := CfrCrawler{}
	rt := RtCrawler{}
	craws[0] = crf
	craws[1] = rt
	var wg sync.WaitGroup
	queryChan := make(chan queryData, len(craws))

	msearch := map[string]interface{}{"cmd": "search", "key": "樂事"}
	for _, vv := range craws {
		wg.Add(1)
		go func(cr BaseCrawler) {
			defer wg.Done()
			log.Println(cr.GetStoreName())
			qd := queryData{
				Code:  -1,
				Store: cr.GetStoreName(),
			}
			rdata := cr.StartCrawlering(cr.MakeCrawCmd(msearch))
			if rdata != nil {
				pdata := cr.GetCrawlingData(rdata)
				if pdata != nil {
					qd.Code = 0
					qd.Data = pdata
				}
			}
			queryChan <- qd
		}(vv)
	}
	go func() {
		wg.Wait()
		close(queryChan)
	}()

	for tmp := range queryChan {
		if tmp.Code == 0 {
			qProducts = append(qProducts, tmp.Data...)
		}
	}

	jj, _ := json.Marshal(SortPriceParseData(qProducts)[:20])
	log.Println(string(jj))
}
