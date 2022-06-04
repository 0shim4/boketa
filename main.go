package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

const DATE_FORMAT string = "20060102150405"

func main() {
	executeTime := time.Now().Format(DATE_FORMAT)
	for i := 1; i <= 25; i++ {
		url, title := scrape()
		fmt.Println(title)

		download(url, fmt.Sprintf("%03d", i), executeTime)
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

func download(url string, number string, executeTime string) {

	res, _ := http.Get(url)
	defer res.Body.Close()

	dir := fmt.Sprintf("./resource/%s", executeTime)
	os.Mkdir(dir, 0777)
	file, _ := os.Create(fmt.Sprintf("%s/%s.jpg", dir, number))
	defer file.Close()

	io.Copy(file, res.Body)
}
