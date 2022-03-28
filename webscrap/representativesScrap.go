//http://go-colly.org/docs/examples/coursera_courses/
package main

import (
	helper "GOLITICS/helper"
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Representative struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	YearsServed string `json:"yearsServed"`
	State       string `json:"state"`
	Party       string `json:"party"`
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

	tsql := `
	INSERT INTO Goliticians.dbo.Politicians (Name, URL, YearsServedRaw, State, Party) VALUES (@Name, @URL, @YearsServedRaw, @State, @Party);
	select isNull(SCOPE_IDENTITY(), -1);`

	space := regexp.MustCompile(`[^:a-zA-Z0-9]\s+`)
	//repInfo := make([]Representative, 0, 200)
	innerDataInfo := make([]Representative, 0, 200)

	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("Request.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Fatalf("Error opening file: %v", err)
	}

	c := colly.NewCollector()

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
		politicianInfo := Representative{}
		e.ForEach("li.expanded", func(_ int, el *colly.HTMLElement) {

			politicianInfo.Name = el.ChildText("span.result-heading")
			politicianInfo.URL = "congress.gov" + el.ChildAttrs("a", "href")[0]

			data := el.ChildText("span.result-item")
			dataSpaceRemove := space.ReplaceAllString(data, "|")

			fmt.Println(dataSpaceRemove)

			splitPoliticianInfo := helper.DelSplit(dataSpaceRemove, delim)

			for i, word := range splitPoliticianInfo {
				switch word {
				case "State:":
					politicianInfo.State = splitPoliticianInfo[i+1]
				case "Party:":
					politicianInfo.Party = splitPoliticianInfo[i+1]
				case "Served:":
					politicianInfo.YearsServed = splitPoliticianInfo[i+1:][0]
				}
			}

			innerDataInfo = append(innerDataInfo, politicianInfo)

		})

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
	db, ctx := helper.ConnectDB()
	//Need to log and go to next iteration if loop error
	for _, info := range innerDataInfo {
		insertRep(tsql, db, ctx, info)
	}

	defer file.Close()
	//fmt.Println(innerDataInfo)

}

//return newID, nil
func insertRep(tsql string, db *sql.DB, ctx context.Context, info Representative) {
	stmt, err := db.Prepare(tsql)
	if err != nil {
		//return -1, err
		fmt.Printf("Error")
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", info.Name),
		sql.Named("URL", info.URL),
		sql.Named("YearsServedRaw", info.YearsServed),
		sql.Named("State", info.State),
		sql.Named("Party", info.Party),
	)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		//return -1, err
		fmt.Printf("Error")
	}
}
