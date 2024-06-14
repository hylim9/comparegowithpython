package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scrapeJob struct {
	title string
	company_name string
	position string
	region string
	url string
}

func Scrape(searchWord string) {
	start := time.Now()
	var url string = "https://remoteok.com/remote-" + searchWord + "-jobs"
	var jobs []scrapeJob
	c := make(chan scrapeJob)
	fmt.Println("Requesting: ", url)
	// res, err := http.Get(url)
	// checkErr(err)
	// checkStatusCode(res)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	
	client := &http.Client{}
	res, err := client.Do(req)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find("table#jobsboard>tbody>tr.job")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	fmt.Println("***")
	fmt.Println(jobs)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	
}

// func getPage(url string, mainC chan<- []scrapeJob) {
	
// }

func extractJob(card *goquery.Selection, c chan<- scrapeJob) {
	company_info := card.Find("td.company")
	url, _ := company_info.Find("a").Attr("href")
	title := CleanString(company_info.Find("a").Text())
	company_name := CleanString(company_info.Find("span.companyLink>h3").Text())
	region := CleanString(company_info.Find("div.location").Text())
	tag := CleanString(card.Find("td.tags>a").First().Find("div>h3").Text())
	// fmt.Println(title, company_name, region, url, tag)
	c <- scrapeJob{
		title: title, 
		company_name: company_name,
		position: tag,
		region: region,
		url: url,
	}
}

// CleanString 문자열 정리
func CleanString(str string) string {
	// Fields 문자열을 분리 , TrimSpace 양쪽 끝에 공백을 제거, Join 배열을 separater 기준 join
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}