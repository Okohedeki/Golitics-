package main

import (
	"os"
	"strconv"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type RequestLogParams struct {
	URL    string
	Method string
	Parser string
}

func main() {
	crawl()
}

func crawl() {
	var baseurl = "https://www.congress.gov/members?q=%7B%22congress%22%3A%22all%22%7D&pageSize=250&page=1"

	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("Request.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Fatalf("Error opening file: %v", err)
	}

	c := colly.NewCollector(
		/*colly.AllowedDomains("https://www.congress.gov/members"),*/
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"),
		/*colly.MaxDepth(5),*/
		/*colly.IgnoreRobotsTxt(),*/
	)

	/*websiteParser := c.Clone()*/

	c.OnRequest(func(r *colly.Request) {
		/*fmt.Println("Visiting: ", r.URL.String())*/
		log.WithFields(
			log.Fields{
				"parser": "Representative",
				"url":    r.URL.String(),
				/*"header": r.Headers,*/
				"paylod": r.Body,
				"type":   "Request",
			},
		).Info("GET Request")
	})

	c.OnResponse(func(r *colly.Response) {
		status_code := strconv.Itoa(r.StatusCode)
		log.WithFields(
			log.Fields{
				"parser":     "Representative",
				"url":        r.Request.URL.String(),
				"statusCode": status_code,
				/*"header":     r.Headers,*/
				"type": "Response",
			},
		).Info("GET Response")
	})

	c.Visit(baseurl)

	defer file.Close()

}
