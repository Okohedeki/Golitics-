//http://go-colly.org/docs/examples/coursera_courses/
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Representative struct {
	name        string 
	state       string 
	party       string 
	yearsServed string 
	url         string 

func main() {
	crawl()
}

func crawl() {
	//var maxRepr string
	var baseurl = "https://www.congress.gov/members?q=%7B%22congress%22%3A%22all%22%7D&pageSize=250&page=1"

	repInfo := make([]Representative)

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

	c.Limit((&colly.LimitRule{
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	}))

	infoCollector := c.Clone()

	//stopRepCrawl := false

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

	infoCollector.OnHTML("ol", func(e *colly.HTMLElement) {
	
		log.WithFields(
			log.Fields{
				"parser": "Representative",
				"url":    e.URL.String(),
			},
		).Info("Found Rep Table")

		e.ForEach("li.expanded", func(_ int, el *colly.HTMLElement) {
			fmt.Printf("test")
			repInfo.name = e.ChildText("span.result-heading")
			repInfo.Url = e.Attr("href")
			switch el.ChildText("span.result-item") {
			case "State:":
				repInfo.state = el.ChildText("span.result-item")
			case "Party:":
				repInfo.party = el.ChildText("span.result-item")
			case "Served:"
				repInfo.yearsServed = el.ChildText("span.result-item")
			}


		})

		result, _ := json.MarshalIndent(repInfo, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(string(result))

	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL((e.Attr("href")))
		c.Visit((nextPage))

	},
	)

	c.Visit(baseurl)

	defer file.Close()

}
