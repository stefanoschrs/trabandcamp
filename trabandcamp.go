package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"regexp"
	"strings"
	"net/http"
	"encoding/json"
	// "strconv"

    "github.com/parnurzeal/gorequest"
	"golang.org/x/net/html"
)

type File struct{
	Url string `json:"mp3-128"`
}

type Track struct {
	Title 	string 	`json:"title"`
	FileUrl File 	`json:"file"`
}

func fetchAlbums(band string) (albums []string) {
	var url string = "https://" + band + ".bandcamp.com/music"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
			case tt == html.ErrorToken:
				return
			case tt == html.StartTagToken:
				t := z.Token()

				isAnchor := t.Data == "a"
				if !isAnchor {
					continue
				}

				for _, a := range t.Attr {
					if a.Key == "href" {
						href := a.Val
						if strings.Index(href, "/album") == 0 {
							albums = append(albums, href)
						}
					}
				}
		}
	}
}

func fetchAlbum(band string, album string) (tracks []Track) {
	var url string = "https://" + band + ".bandcamp.com" + album
	_, body, errs := gorequest.New().Get(url).End()

	if errs != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	pattern, _ := regexp.Compile(`trackinfo.+}]`)
	result := pattern.FindString(body)

	json.Unmarshal([]byte(result[strings.Index(result, "[{"):]), &tracks)

	return
}

func download(path string, track Track, album string){
	fmt.Printf("Downloading %s (%s)\n", track.Title, album)

	var fileName string = path + "/" + track.Title + ".mp3"
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	var url string = "http:" + track.FileUrl.Url
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Printf("Successfully Downloaded %s (%s)\n", track.Title, album)
}

func main(){
	if len(os.Args) != 2 {
		panic("Missing Band Name")
	}

	var path string = "./data"
	os.MkdirAll(path, 0777)

	band := os.Args[1]
	fmt.Println("Analyzing " + band)
	path = path + "/" + band
	os.MkdirAll(path, 0777)

	var albums []string = fetchAlbums(band)
	fmt.Println("Albums:")
	fmt.Println(albums)

	for _, album := range albums{
		var tracks []Track = fetchAlbum(band, album)
		var albumPath string = path + "/" + album[6:]
		os.MkdirAll(albumPath, 0777)

		for _, v := range tracks{
			go download(albumPath, v, album[7:])
		}
	}

	time.Sleep(1 * time.Hour)
}
