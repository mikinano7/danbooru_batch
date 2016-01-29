package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/mikinano7/dropbox4go"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	fqdn    = "https://danbooru.donmai.us"
	popular = fqdn + "/explore/posts/popular"
)

func main() {
	godotenv.Load("go.env")

	doc, _ := goquery.NewDocument(popular)
	var arr []string
	doc.Find("#a-index article").Each(func(_ int, s *goquery.Selection) {
		a, _ := s.Attr("data-file-url")
		arr = append(arr, a)
	})

	httpClient := http.DefaultClient
	svc := dropbox4go.New(httpClient, os.Getenv("DB_ACCESS_TOKEN"))

	for _, img := range arr {
		resp, _ := httpClient.Get(fqdn + img)
		defer resp.Body.Close()

		pos := strings.LastIndex(img, "/")
		fileName := img[pos+1:]

		req := dropbox4go.Request{
			File: resp.Body,
			Parameters: dropbox4go.Parameters{
				Path:           fmt.Sprintf("/home/%s", fileName),
				Mode:           "overwrite",
				AutoRename:     false,
				ClientModified: time.Now().UTC().Format(time.RFC3339),
				Mute:           true,
			},
		}
		svc.Upload(req)
	}
}
