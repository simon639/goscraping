package goscraping

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type SpiderExtender struct {
	Spider *Spider
	gocrawl.DefaultExtender
}

func (extender *SpiderExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	url := ctx.NormalizedURL().String()
	spider := extender.Spider
	// fmt.Println("filter url :", url)
	if !isVisited {
		for _, rxDeny := range spider.rxDenies {
			if rxDeny.MatchString(url) {
				// fmt.Println("match deny !")
				return false
			}
		}

		// fmt.Println("rxAllows :", s.rxAllows)
		for _, rxAllow := range spider.rxAllows {
			// fmt.Println(rxAllow)
			if rxAllow.MatchString(url) {
				// fmt.Println("matched allow.")
				return true
			}
		}
		//fmt.Println("dont filter")
	}
	// fmt.Println("no match, filter url !")
	//fmt.Println("filtered")
	return false
}

func (extender *SpiderExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	// Use the goquery document or res.Body to manipulate the data
	// ...

	// Return nil and true - let gocrawl find the links
	url := ctx.NormalizedURL().String()
	fmt.Println("visit url :", url)
    spider := extender.Spider
	for _, rxItem := range spider.rxItems {
		if rxItem.MatchString(url) {
			fmt.Println("--------------------------------------------------------------")
			fmt.Println("item url :", url)
			//title := doc.Find("div.title h1").Text()
			//fmt.Println("title :", title)

			// TODO pass to ParseItem(item interface{})
			spider.ParseItemFn(url, res, doc)
			// if s.ProcessItemFn != nil {
			// 	s.ProcessItemFn(item)
			// }
			break
		}
	}

	return nil, true
}
