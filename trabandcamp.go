package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/parnurzeal/gorequest"
)

// Configuration object
type Configuration struct {
	Directory string `json:"directory"`
}

// Track class
type Track struct {
	Title   string `json:"title"`
	FileURL File   `json:"file"`
	Album   string
}

// File - Track URL helper class
type File struct {
	URL string `json:"mp3-128"`
}

const CONCURRENCY = 4

var throttle = make(chan int, CONCURRENCY)

func _contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func downloadTrack(path string, track Track, wg *sync.WaitGroup, throttle chan int) {
	defer wg.Done()

	fmt.Printf("[DEBUG] Downloading %s (%s)\n", track.Title, track.Album)

	var fileName = fmt.Sprintf("%s/%s/%s.mp3", path, track.Album, track.Title)
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] Failed creating", fileName, "-", err)
		<-throttle
		return
	}
	defer output.Close()

	var url = "https:" + track.FileURL.URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("[ERROR] Failed downloading", url, "-", err)
		<-throttle
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("[ERROR] Failed downloading", url, "-", err)
		<-throttle
		return
	}

	fmt.Printf("[DEBUG] Successfully Downloaded %s (%s)\n", track.Title, track.Album)

	<-throttle
}

func fetchAlbumTracks(band string, album string) (tracks []Track) {
	var url = "https://" + band + ".bandcamp.com" + album
	_, body, errs := gorequest.New().Get(url).End()

	if errs != nil {
		fmt.Println("[ERROR] Failed to crawl \"" + url + "\"")
		return
	}

	pattern, _ := regexp.Compile(`trackinfo.+}]`)
	result := pattern.FindString(body)

	json.Unmarshal([]byte(result[strings.Index(result, "[{"):]), &tracks)

	return
}

func fetchAlbums(band string) (albums []string) {
	var url = "https://" + band + ".bandcamp.com/music"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("[ERROR] Failed to crawl \"" + url + "\"")
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
					if strings.Index(href, "/album") == 0 && !_contains(albums, href) {
						albums = append(albums, href)
					}
				}
			}
		}
	}
}

func getMusicPathRoot() string {
	root := path.Dir(os.Args[0]) + "/data"

	file, err := os.Open(path.Dir(os.Args[0]) + "/.trabandcamprc")
	if err == nil {
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err = decoder.Decode(&configuration)
		if err != nil {
			fmt.Println("[ERROR] Cannot parse configuration file.", err)
			os.Exit(1)
		}

		root = configuration.Directory
	}

	return root
}

func checkBandExistence(band string) bool {
	var url = "https://" + band + ".bandcamp.com/music"
	resp, err := http.Get(url)

	if resp.StatusCode != 200 || err != nil {
		return false
	}

	return true
}

func main() {
	fmt.Println("                          _                     _")
	fmt.Println("___________              | |                   | |")
	fmt.Println("\\__    ___/___________   | |__   __ _ _ __   __| | ___ __ _ _ __ ___  _ __")
	fmt.Println("  |    |  \\_  __ \\__  \\  | '_ \\ / _` | '_ \\ / _` |/ __/ _` | '_ ` _ \\| '_ \\")
	fmt.Println("  |    |   |  | \\// __ \\_| |_) | (_| | | | | (_| | (_| (_| | | | | | | |_) |")
	fmt.Println("  |____|   |__|  (____  /|_.__/ \\__,_|_| |_|\\__,_|\\___\\__,_|_| |_| |_| .__/")
	fmt.Println("                                                                     | |")
	fmt.Println("                                                                     |_| v.0.0.2")

	if len(os.Args) != 2 {
		fmt.Println("[ERROR] Missing Band Name")
		os.Exit(1)
	}

	band := os.Args[1]
	if !checkBandExistence(band) {
		fmt.Printf("[ERROR] Band `%s` doesn't exist\n", band)
		os.Exit(1)
	}

	var musicPath = getMusicPathRoot()
	musicPath = musicPath + "/" + band
	os.MkdirAll(musicPath, 0777)

	fmt.Println("[INFO] Analyzing " + band)

	var albums = fetchAlbums(band)
	fmt.Printf("[INFO] Found %d Albums\n", len(albums))
	fmt.Printf("[DEBUG] %q\n", albums)

	var tracks []Track
	for _, album := range albums {
		fmt.Printf("[DEBUG] Fetching album %s tracks\n", album[7:])
		var albumPath = musicPath + "/" + album[7:]
		os.MkdirAll(albumPath, 0777)
		tmpTracks := fetchAlbumTracks(band, album)

		for index := range tmpTracks {
			tmpTracks[index].Album = album[7:]
		}

		tracks = append(tracks, tmpTracks...)
	}
	fmt.Printf("[INFO] Found %d Tracks\n", len(tracks))

	var wg sync.WaitGroup
	for _, track := range tracks {
		throttle <- 1
		wg.Add(1)

		go downloadTrack(musicPath, track, &wg, throttle)
	}

	wg.Wait()
}
