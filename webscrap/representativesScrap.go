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
	url         string
	state       string
	party       string
	yearsServed string
}

func main() {
	crawl()
}

func crawl() {
	//var maxRepr string
	var baseurl = "https://www.congress.gov/members?q=%7B%22congress%22%3A%22all%22%7D&pageSize=250&page=1"

	repInfo := make([]Representative, 0, 200)

	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("Request.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Fatalf("Error opening file: %v", err)
	}

	c := colly.NewCollector(
	//colly.AllowedDomains("https://www.congress.gov/members"),
	// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"),
	// /*colly.MaxDepth(5),*/
	// /*colly.IgnoreRobotsTxt(),*/
	)

	c.Limit((&colly.LimitRule{
		Delay:       2 * time.Second,
		RandomDelay: 2 * time.Second,
	}))

	//infoCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
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

	c.OnResponse(func(e *colly.Response) {

		status_code := strconv.Itoa(e.StatusCode)
		if e.StatusCode == 200 || e.StatusCode == 203 {

			absoluteUrl := e.Request.AbsoluteURL((e.Request.URL.String()))
			fmt.Println(absoluteUrl)
			//infoCollector.Visit(absoluteUrl)
			log.WithFields(
				log.Fields{
					"parser":     "Representative",
					"url":        e.Request.URL.String(),
					"statusCode": status_code,
					/*"header":     r.Headers,*/
					"type": "Response",
				},
			).Info("GET Response")

		} else {
			log.WithFields(
				log.Fields{
					"parser": "Representative",
					"url":    e.Request.URL.String(),
				},
			).Warn("Could not find scraping data")

		}

	})

	c.OnHTML(".search-column-main.basic-search-results.nav-on", func(e *colly.HTMLElement) {
		fmt.Println("Test")
		log.WithFields(
			log.Fields{
				"parser": "Representative",
				"url":    e.Request.URL.String(),
			},
		).Info("Found Rep Table")

		politicianInfo := Representative{
			name: e.ChildText("span.result-heading"),
			url:  e.Request.AbsoluteURL(e.Attr("href")),
		}

		e.ForEach("li.expanded", func(_ int, el *colly.HTMLElement) {

			switch el.ChildText("span.result-item") {
			case "State:":
				politicianInfo.state = el.ChildText("span.result-item")
			case "Party:":
				politicianInfo.party = el.ChildText("span.result-item")
			case "Served:":
				politicianInfo.yearsServed = el.ChildText("span.result-item")
			}

		})
		repInfo = append(repInfo, politicianInfo)
		fmt.Println(repInfo)

	})

	// c.OnHTML("a.next", func(e *colly.HTMLElement) {
	// 	nextPage := e.Request.AbsoluteURL((e.Attr("href")))
	// 	c.Visit((nextPage))

	// },
	// )

	c.OnError(func(r *colly.Response, err error) {
		log.WithFields(
			log.Fields{
				"parser":     "Representative",
				"url":        r.Request.URL.String(),
				"statusCode": r.StatusCode,
				"type":       "Response",
			},
		).Warn("Error Parsing Webpage")
	})

	c.Visit(baseurl)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	enc.Encode(repInfo)

	defer file.Close()

}
