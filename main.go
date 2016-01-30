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

	doc, err := goquery.NewDocument(popular)
	if err != nil {
		panic(err)
	}

	var arr []string
	doc.Find("#a-index article").Each(func(_ int, s *goquery.Selection) {
		a, _ := s.Attr("data-file-url")
		arr = append(arr, a)
	})

	httpClient := http.DefaultClient
	svc := dropbox4go.New(httpClient, os.Getenv("DB_ACCESS_TOKEN"))
	now := time.Now().UTC().In(time.FixedZone("Asia/Tokyo", 9*60*60))

	for _, img := range arr {
		resp, err := httpClient.Get(fqdn + img)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		pos := strings.LastIndex(img, "/")
		fileName := img[pos+1:]

		req := dropbox4go.Request{
			File: resp.Body,
			Parameters: dropbox4go.Parameters{
				Path: fmt.Sprintf(
					"/danbooru/popular/%s/%s",
					now.Format("2006-01-02"),
					fileName,
				),
				Mode:           "overwrite",
				AutoRename:     false,
				ClientModified: now.Format(time.RFC3339),
				Mute:           true,
			},
		}
		a, err := svc.Upload(req)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(fmt.Sprintf("file %s has uploaded.", a.Name))
		}
	}
}
