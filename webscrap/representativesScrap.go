//http://go-colly.org/docs/examples/coursera_courses/
package main

import (
	helper "GOLITICS/helper"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Representative struct {
	name        string `json:"name"`
	url         string `json:"url"`
	yearsServed string `json:"yearsServed"`
	state       string `json:"state"`
	party       string `json:"party"`
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
	var delim = '|'
	space := regexp.MustCompile(`([^:[A-Za-z])\s+`)
	politicianInfo := Representative{}
	//repInfo := make([]Representative, 0, 200)
	innerDataInfo := make([]Representative, 0, 200)

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
		log.WithFields(
			log.Fields{
				"parser": "Representative",
				"url":    e.Request.URL.String(),
			},
		).Info("Found Rep Table")

		e.ForEach("li.expanded", func(_ int, el *colly.HTMLElement) {

			politicianInfo.name = el.ChildText("span.result-heading")
			politicianInfo.url = "congress.gov" + el.ChildAttrs("a", "href")[0]

			data := el.ChildText("span.result-item")
			dataSpaceRemove := space.ReplaceAllString(data, "|")
			splitPoliticianInfo := helper.DelSplit(dataSpaceRemove, delim)

			for i, word := range splitPoliticianInfo {
				switch word {
				case "State:":
					politicianInfo.state = splitPoliticianInfo[i+1]
				case "Party:":
					politicianInfo.party = splitPoliticianInfo[i+1]
				case "Served:":
					politicianInfo.yearsServed = splitPoliticianInfo[i+1:][0]
				}
			}

		})
		innerDataInfo = append(innerDataInfo, politicianInfo)

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

	politicianInfoJson, err := json.Marshal(innerDataInfo)
	if err != nil {
		fmt.Println(politicianInfoJson)
	} else {
		log.WithFields(
			log.Fields{
				"parser": "Representative",
				"step":   "jsonConversion",
			},
		).Warn("Error converting struct to json")
	}

	defer file.Close()

	fmt.Println(politicianInfo)
	res, _ := json.MarshalIndent(politicianInfo, "", "\t")
	fmt.Println(string(res))

}
