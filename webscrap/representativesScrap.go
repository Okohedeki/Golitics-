//http://go-colly.org/docs/examples/coursera_courses/
package main

import (
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
	houseSenate string
	state       string
	party       string
	yearsServed string
}

type InnerData struct {
	InnerData string
}

func main() {
	crawl()
}

func crawl() {
	//var maxRepr string
	var baseurl = "https://www.congress.gov/members?q=%7B%22congress%22%3A%22all%22%7D&pageSize=250&page=1"

	//repInfo := make([]Representative, 0, 200)
	innerDataInfo := make([]InnerData, 0, 200)

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
			//infoCollector.Visit(absoluteUrl)
			log.WithFields(
				log.Fields{
					"parser":     "Representative",
					"url":        absoluteUrl,
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

		politicianInfo := Representative{}
		innerDataStruct := InnerData{}

		e.ForEach("li.expanded", func(_ int, el *colly.HTMLElement) {

			politicianInfo.name = el.ChildText("span.result-heading")
			politicianInfo.url = "congress.gov" + el.ChildAttrs("a", "href")[0]

			data := el.ChildText("span.result-item")
			innerDataStruct.InnerData = data
			// switch el.ChildText("span", "result-item") {
			// case "State:":
			// 	fmt.Println(el.ChildText("span.result-item"))

			// case "Party:":
			// 	politicianInfo.party = el.ChildText("div.member-profile.member-image-exists > span.result-item > strong")
			// case "Served:":
			// 	politicianInfo.yearsServed = el.ChildText("div.member-profile.member-image-exists > span.result-item > strong")
			// }

		})
		innerDataInfo = append(innerDataInfo, innerDataStruct)
		//repInfo = append(repInfo, politicianInfo)

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
	defer file.Close()

	fmt.Println(innerDataInfo)

}
