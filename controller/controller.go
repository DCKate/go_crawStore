package controller

import (
	"go_findstuff/crawler"
	"log"
	"sort"
	"sync"
)

// RunSearchCrawler : use goroutine to start crawling website
//		craws : all of the crawler instance
//		msearch : the http parameter for the url used by crawler
func RunSearchCrawler(craws []crawler.BaseCrawler, msearch map[string]interface{}) []crawler.ParseData {
	var wg sync.WaitGroup
	var qProducts []crawler.ParseData
	queryChan := make(chan crawler.QueryData, len(craws))
	for _, vv := range craws {
		wg.Add(1)
		go func(cr crawler.BaseCrawler) {
			defer wg.Done()
			log.Println(cr.GetStoreName())
			qd := crawler.QueryData{
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
	return qProducts
}

// ProductsSearchController : accoding the input parameter to generate instance of crawler
//		para : input parameter, only "key", "store", "limit" are used
func ProductsSearchController(para map[string]interface{}) (int, []crawler.ParseData) {
	var craws []crawler.BaseCrawler
	if vv, ok := para["store"]; ok {
		stores := vv.([]interface{})
		for _, name := range stores {
			switch name.(string) {
			case crawler.RtStoreName:
				cr := crawler.RtCrawler{}
				craws = append(craws, cr)
			case crawler.CrfStoreName:
				cr := crawler.CfrCrawler{}
				craws = append(craws, cr)
			}
		}
	} else {
		crf := crawler.CfrCrawler{}
		rt := crawler.RtCrawler{}
		craws = append(craws, crf, rt)
	}
	skey := para["key"].(string)
	msearch := map[string]interface{}{"cmd": "search", "key": skey}
	ret := RunSearchCrawler(craws, msearch)
	disNum := len(ret)
	if vv, ok := para["limit"]; ok {
		limitNum := vv.(int)
		if disNum > limitNum {
			disNum = limitNum
		}
	}
	if len(ret) > 0 {
		sort.Sort(crawler.ParseDataGroup(ret))
		return 0, ret[:disNum]
	}
	return -1, nil
}
