package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
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

		go makeImageForMovie(download(url, fmt.Sprintf("%03d", i), dir))
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

func makeImageForMovie(path string) {
	backgroundFile, _ := os.Open("./resource/background.png") // TODO: 毎回ファイルを開かず1回で済ませる
	inputFile, _ := os.Open(path)
	defer backgroundFile.Close()
	defer inputFile.Close()

	backgroundImg, _, _ := image.Decode(backgroundFile)
	png.Encode(new(bytes.Buffer), backgroundImg) // NOTE: pngとして識別されないためエンコード
	inputImg, _, _ := image.Decode(inputFile)

	startPointLogo := image.Point{(backgroundImg.Bounds().Dx() - inputImg.Bounds().Dx()) / 2, 0}

	logoRectangle := image.Rectangle{startPointLogo, startPointLogo.Add(inputImg.Bounds().Size())}
	originRectangle := image.Rectangle{image.Point{0, 0}, backgroundImg.Bounds().Size()}

	rgba := image.NewRGBA(originRectangle)
	draw.Draw(rgba, originRectangle, backgroundImg, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, logoRectangle, inputImg, image.Point{0, 0}, draw.Over)

	outputFile, _ := os.Create(path)
	defer outputFile.Close()
	jpeg.Encode(outputFile, rgba, nil)
}
