package goscraping

import (
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
)

func InitRxs(patterns []string) []*regexp.Regexp {
	rxs := []*regexp.Regexp{}
	// fmt.Println("patterns :", patterns)
	for _, pattern := range patterns {
		// fmt.Println("pattern :", pattern)
		rx := regexp.MustCompile(pattern)
		rxs = append(rxs, rx)
	}
	// fmt.Println("rxs :", rxs)
	return rxs
}

type Item struct{}

type Spider struct {
	AllowPatterns []string
	DenyPatterns  []string
	ItemPatterns  []string

	ParseItemFn   func(url string, res *http.Response, doc *goquery.Document)
	ProcessItemFn func(*Item)

	CrawlDelay   time.Duration
	SameHostOnly bool

	rxAllows []*regexp.Regexp
	rxDenies []*regexp.Regexp
	rxItems  []*regexp.Regexp
}

func NewSpider() *Spider {
	spider := &Spider{}
	return spider
}

func (spider *Spider) Run(startURLs []string) {
	spider.rxAllows = InitRxs(spider.AllowPatterns)
	// fmt.Println("Run rxAllows :", spider.rxAllows)

	spider.rxDenies = InitRxs(spider.DenyPatterns)
	spider.rxItems = InitRxs(spider.ItemPatterns)

	extender := new(SpiderExtender)
	extender.Spider = spider
	opts := gocrawl.NewOptions(extender)
	opts.CrawlDelay = spider.CrawlDelay * time.Second
	opts.LogFlags = gocrawl.LogError
	// opts.MaxVisits = 10
	opts.SameHostOnly = spider.SameHostOnly
	//opts.EnqueueChanBuffer = 1000
	//opts.EnqueueChanBuffer = 1
	//opts.WorkerIdleTTL = 0.0 * time.Second
	opts.URLNormalizationFlags = opts.URLNormalizationFlags ^ purell.FlagForceHTTP
	// opts.URLNormalizationFlags
	// fmt.Println(opts.URLNormalizationFlags)
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(startURLs)

}
