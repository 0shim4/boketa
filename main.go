package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
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
	dir := fmt.Sprintf("./resource/%s", time.Now().Format(DATE_FORMAT))
	os.Mkdir(dir, 0777)
	for i := 1; i <= 25; i++ {
		url, title := scrape()
		fmt.Println(title)

		go extendImageBottom(download(url, fmt.Sprintf("%03d", i), dir))
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

func download(url string, number string, dir string) (path string) {
	res, _ := http.Get(url)
	defer res.Body.Close()

	path = fmt.Sprintf("%s/%s.png", dir, number)

	file, _ := os.Create(path)
	defer file.Close()

	io.Copy(file, res.Body)

	return path
}

func extendImageBottom(path string) {
	inputFile, _ := os.Open(path)
	defer inputFile.Close()

	img, _, _ := image.Decode(inputFile)

	outputFile, _ := os.Create(path)
	defer outputFile.Close()

	m := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()+100))
	c := color.RGBA{0, 0, 255, 255}

	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

	rct := image.Rectangle{image.Point{0, 0}, m.Bounds().Size()}

	draw.Draw(m, rct, img, image.Point{0, 0}, draw.Src)

	jpeg.Encode(outputFile, m, &jpeg.Options{Quality: 100})
}
