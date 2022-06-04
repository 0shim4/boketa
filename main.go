package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func main() {
	for i := 0; i < 25; i++ {
		fmt.Println(scrape())
		time.Sleep(time.Second * 1)
	}
}

func scrape() (string, string) {
	url := "https://bokete.jp/"

	res, _ := http.Get(url)
	defer res.Body.Close()

	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)

	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

	doc, _ := goquery.NewDocumentFromReader(reader)

	src, _ := doc.Find("div.photo-content > a").Find("img").Attr("src")
	text := doc.Find("a.boke-text > div").Text()

	return "https:" + src, text[5:]
}
