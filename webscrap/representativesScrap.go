package representatives_scrap

import (
	"github.com/gocolly/colly"
	"os"
	log "github.com/sirupsen/logrus"
)

type RequestLogParams struct {
	URL string 
	Method string 
	Parser string 

}

func main() {
	crawl()
}

func crawl() {
	log.setFormatter(&log.JSONFormatter{})

	c := colly.NewCollector(
		colly.AllowedDomains("https://www.congress.gov/members"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
	)

	websiteParser := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		
	}

}
