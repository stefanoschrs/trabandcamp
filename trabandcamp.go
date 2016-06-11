package main

import (
	"io"
	"os"
	"fmt"
	"path"
	"time"
	"regexp"
	"strings"
	"net/http"
	"encoding/json"

	"golang.org/x/net/html"

	"github.com/parnurzeal/gorequest"
)

// Configuration object
type Configuration struct {
	Directory string `json:directory`
}

// Track class
type Track struct {
	Title 	string 	`json:"title"`
	FileURL File 	`json:"file"`
}

// File - Track URL helper class
type File struct{
	URL string `json:"mp3-128"`
}

func _contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func fetchAlbums(band string) (albums []string) {
	var url = "https://" + band + ".bandcamp.com/music"
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
					if strings.Index(href, "/album") == 0 && !_contains(albums, href){
						albums = append(albums, href)
					}
				}
			}
		}
	}
}

func fetchAlbum(band string, album string) (tracks []Track) {
	var url = "https://" + band + ".bandcamp.com" + album
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

	var fileName = path + "/" + track.Title + ".mp3"
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("ERROR: Failed creating", fileName, "-", err)
		return
	}
	defer output.Close()

	var url = "http:" + track.FileURL.URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("ERROR: Failed downloading", url, "-", err)
		return
	}

	fmt.Printf("Successfully Downloaded %s (%s)\n", track.Title, album)
}

func checkBandExistance(band string) bool{
	var url = "https://" + band + ".bandcamp.com/music"
	resp, err := http.Get(url)

	if resp.StatusCode != 200 || err != nil {
		return false
	}

	return true
}

func main(){
	fmt.Println("                          _                     _")
	fmt.Println("___________              | |                   | |")
	fmt.Println("\\__    ___/___________   | |__   __ _ _ __   __| | ___ __ _ _ __ ___  _ __")
	fmt.Println("  |    |  \\_  __ \\__  \\  | '_ \\ / _` | '_ \\ / _` |/ __/ _` | '_ ` _ \\| '_ \\")
	fmt.Println("  |    |   |  | \\// __ \\_| |_) | (_| | | | | (_| | (_| (_| | | | | | | |_) |")
	fmt.Println("  |____|   |__|  (____  /|_.__/ \\__,_|_| |_|\\__,_|\\___\\__,_|_| |_| |_| .__/")
	fmt.Println("                                                                     | |")
	fmt.Println("                                                                     |_| v.0.0.2")

	if len(os.Args) != 2 {
		fmt.Println("ERROR: Missing Band Name")
		os.Exit(1)
	}

	band := os.Args[1]
	if !checkBandExistance(band){
		fmt.Printf("ERROR: Band `%s` doesn't exist\n", band)
		os.Exit(1)
	}

	var musicPath = path.Dir(os.Args[0]) + "/data"
	file, err := os.Open(path.Dir(os.Args[0]) + "/.trabandcamprc")
	if err == nil {
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			fmt.Println("ERROR: Cannot parse configuration file.", err)
			os.Exit(1)
		}

		musicPath = configuration.Directory
	}
	musicPath = musicPath + "/" + band
	os.MkdirAll(musicPath, 0777)

	fmt.Println("Analyzing " + band)

	var albums = fetchAlbums(band)
	fmt.Printf("Albums: %q\n", albums)

	for _, album := range albums{
		var tracks = fetchAlbum(band, album)
		var albumPath = musicPath + "/" + album[6:]
		os.MkdirAll(albumPath, 0777)

		for _, v := range tracks{
			go download(albumPath, v, album[7:])
		}
	}

	time.Sleep(10 * time.Hour)
}
