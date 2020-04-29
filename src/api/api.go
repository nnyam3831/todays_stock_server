package api

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"golang.org/x/net/html/charset"
)

// Stock ...
type Stock struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Price   string `json:"price"`
	Percent string `json:"percent"`
}

// euc-kr로 인코딩
var baseURL string = "https://finance.naver.com/sise/"
var goldenCross string = baseURL + "item_gold.nhn"
var riseURL string = baseURL + "sise_rise.nhn"
var lastSearch string = baseURL + "lastsearch2.nhn"
var gcross []Stock
var rise []Stock
var search []Stock
var KOSPI Stock
var KOSDAQ Stock

// Home API
func Home(c echo.Context) error {

	defer c.Request().Body.Close()
	done1 := make(chan bool)
	done2 := make(chan bool)
	done3 := make(chan bool)
	done4 := make(chan bool)
	go getGoldenCross(done1)
	go getKOS(done2)
	go getLastSearch(done3)
	go getRise(done4)

	<-done1
	<-done2
	<-done3
	<-done4
	return c.JSON(http.StatusOK, "scrapping completed")

}

// GetGQ GoldenCross Data
func GetGQ(c echo.Context) error {
	jsonData, err := json.Marshal(gcross)
	if err != nil {
		log.Println(err)
	}
	return c.JSONBlob(http.StatusOK, jsonData)
}

// GetRise ...
func GetRise(c echo.Context) error {
	jsonData, err := json.Marshal(rise)
	if err != nil {
		log.Println(err)
	}
	return c.JSONBlob(http.StatusOK, jsonData)
}

// GetSearch ...
func GetSearch(c echo.Context) error {
	jsonData, err := json.Marshal(search)
	if err != nil {
		log.Println(err)
	}
	return c.JSONBlob(http.StatusOK, jsonData)
}

// GETKOS ...
func GetKOS(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"kospi_now":     KOSPI.Price,
		"kospi_change":  KOSPI.Percent,
		"kosdaq_now":    KOSDAQ.Price,
		"kosdaq_change": KOSDAQ.Percent,
	})
}
func getGoldenCross(done chan bool) {
	res, err := http.Get(goldenCross)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)
	box := doc.Find(".box_type_l")
	box.Find("tr").Each(func(i int, s *goquery.Selection) {

		title := s.Find("td:nth-child(2)").Text()
		href, _ := s.Find("td:nth-child(2)").Find("a").Attr("href")
		price := s.Find("td:nth-child(3)").Text()
		percent := s.Find("td:nth-child(5) > span").Text()
		tit, _ := iconv.ConvertString(string(title), "euc-kr", "utf-8")
		link, _ := iconv.ConvertString(string(href), "euc-kr", "utf-8")
		pr, _ := iconv.ConvertString(string(price), "euc-kr", "utf-8")
		per, _ := iconv.ConvertString(string(percent), "euc-kr", "utf-8")

		if tit != "" && pr != "" && per != "" && link != "" {
			newStock := Stock{
				Title:   cleanString(tit),
				Link:    cleanString(link),
				Price:   cleanString(pr),
				Percent: cleanString(per),
			}
			gcross = append(gcross, newStock)
		}
	})
	done <- true
}

func getRise(done chan bool) {

	res, err := http.Get(riseURL)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)
	box := doc.Find(".box_type_l")
	box.Find("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {

		title := s.Find("td:nth-child(2)").Text()
		href, _ := s.Find("td:nth-child(2)").Find("a").Attr("href")
		price := s.Find("td:nth-child(3)").Text()
		percent := s.Find("td:nth-child(5) > span").Text()
		tit, _ := iconv.ConvertString(string(title), "euc-kr", "utf-8")
		link, _ := iconv.ConvertString(string(href), "euc-kr", "utf-8")
		pr, _ := iconv.ConvertString(string(price), "euc-kr", "utf-8")
		per, _ := iconv.ConvertString(string(percent), "euc-kr", "utf-8")

		if tit != "" && pr != "" && per != "" && link != "" {
			newStock := Stock{
				Title:   cleanString(tit),
				Link:    cleanString(link),
				Price:   cleanString(pr),
				Percent: cleanString(per),
			}
			rise = append(rise, newStock)
			if len(rise) == 100 {
				return false
			}
		}
		return true
	})
	done <- true
}

func getKOS(done chan bool) {
	res, err := http.Get(baseURL)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)
	kpNow := doc.Find("#KOSPI_now").Text()
	kpChange := doc.Find("#KOSPI_change").Text()
	kospiPrice, _ := iconv.ConvertString(string(kpNow), "euc-kr", "utf-8")
	kospiPer, _ := iconv.ConvertString(string(kpChange), "euc-kr", "utf-8")

	kdNow := doc.Find("span#KOSDAQ_now").Text()
	kdChange := doc.Find("span#KOSDAQ_change").Text()
	kosdaqPrice, _ := iconv.ConvertString(string(kdNow), "euc-kr", "utf-8")
	kosdaqPer, _ := iconv.ConvertString(string(kdChange), "euc-kr", "utf-8")

	KOSPI.Title = "코스피"
	KOSPI.Price = cleanString(kospiPrice)
	KOSPI.Percent = cleanString(kospiPer)
	KOSDAQ.Title = "코스닥"
	KOSDAQ.Price = cleanString(kosdaqPrice)
	KOSDAQ.Percent = cleanString(kosdaqPer)

	done <- true
}
func getLastSearch(done chan bool) {
	res, err := http.Get(lastSearch)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)
	box := doc.Find(".box_type_l")
	box.Find("tr").Each(func(i int, s *goquery.Selection) {

		title := s.Find("td:nth-child(2)").Text()
		href, _ := s.Find("td:nth-child(2)").Find("a").Attr("href")
		price := s.Find("td:nth-child(3)").Text()
		percent := s.Find("td:nth-child(5) > span").Text()
		tit, _ := iconv.ConvertString(string(title), "euc-kr", "utf-8")
		link, _ := iconv.ConvertString(string(href), "euc-kr", "utf-8")
		pr, _ := iconv.ConvertString(string(price), "euc-kr", "utf-8")
		per, _ := iconv.ConvertString(string(percent), "euc-kr", "utf-8")

		if tit != "" && pr != "" && per != "" && link != "" {
			newStock := Stock{
				Title:   cleanString(tit),
				Link:    cleanString(link),
				Price:   cleanString(pr),
				Percent: cleanString(per),
			}
			search = append(search, newStock)
		}
	})
	done <- true
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request Failed: ", res.StatusCode)
	}
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func cleanString(str string) string {
	return strings.TrimSpace(str)
}
