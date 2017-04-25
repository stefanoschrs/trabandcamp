package main

import (
	//"github.com/parnurzeal/gorequest"
	. "github.com/tj/go-debug"

	//"golang.org/x/net/html"

	"encoding/json"
	"io/ioutil"
	//"net/http"
	//"strings"
	//"regexp"
	"path"
	//"sync"
	"flag"
	"fmt"
	//"io"
	"os"
)


// Flags - Command line arguments
type Flags struct {
	ConfigLocation string
	IgnoreWarnings bool
}

// Configuration object
type Configuration struct {
	OutputDirectory string
	AlbumsLimit int
	Concurrency ConcurrencyConfiguration
}
// ConcurrencyConfiguration object
type ConcurrencyConfiguration struct{
	Track int8
}

// Band class
type Band struct {
	Name string
	Albums []Album
}

// Album class
type Album struct {
	Name string
	Tracks []Track
}

// Track class
type Track struct {
	Title   string `json:"title"`
	FileURL File   `json:"file"`
}

// File - Track URL helper class
type File struct {
	URL string `json:"mp3-128"`
}


var debug = Debug("app")
var config = Configuration{
	OutputDirectory: path.Dir(os.Args[0]) + "/data",
	AlbumsLimit: 10,
	Concurrency: ConcurrencyConfiguration{
		Track: 4,
	},
}
var flags Flags
var trackThrottle chan int


func _contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//func downloadTrack(path string, track Track, wg *sync.WaitGroup, throttle chan int) {
//	defer wg.Done()
//
//	fmt.Printf("[DEBUG] Downloading %s (%s)\n", track.Title, track.Album)
//
//	var fileName = fmt.Sprintf("%s/%s/%s.mp3", path, track.Album, track.Title)
//	output, err := os.Create(fileName)
//	if err != nil {
//		fmt.Println("[ERROR] Failed creating", fileName, "-", err)
//		<-throttle
//		return
//	}
//	defer output.Close()
//
//	var url = "https:" + track.FileURL.URL
//	response, err := http.Get(url)
//	if err != nil {
//		fmt.Println("[ERROR] Failed downloading", url, "-", err)
//		<-throttle
//		return
//	}
//	defer response.Body.Close()
//
//	_, err = io.Copy(output, response.Body)
//	if err != nil {
//		fmt.Println("[ERROR] Failed downloading", url, "-", err)
//		<-throttle
//		return
//	}
//
//	fmt.Printf("[DEBUG] Successfully Downloaded %s (%s)\n", track.Title, track.Album)
//
//	<-throttle
//}
//
//func fetchAlbumTracks(band string, album string) (tracks []Track) {
//	var url = "https://" + band + ".bandcamp.com" + album
//	_, body, errs := gorequest.New().Get(url).End()
//
//	if errs != nil {
//		fmt.Println("[ERROR] Failed to crawl \"" + url + "\"")
//		return
//	}
//
//	pattern, _ := regexp.Compile(`trackinfo.+}]`)
//	result := pattern.FindString(body)
//
//	json.Unmarshal([]byte(result[strings.Index(result, "[{"):]), &tracks)
//
//	return
//}
//
//func fetchAlbums(band string) (albums []string) {
//	var url = "https://" + band + ".bandcamp.com/music"
//	resp, err := http.Get(url)
//
//	if err != nil {
//		fmt.Println("[ERROR] Failed to crawl \"" + url + "\"")
//		return
//	}
//
//	b := resp.Body
//	defer b.Close()
//
//	z := html.NewTokenizer(b)
//
//	for {
//		tt := z.Next()
//
//		switch {
//		case tt == html.ErrorToken:
//			return
//		case tt == html.StartTagToken:
//			t := z.Token()
//
//			isAnchor := t.Data == "a"
//			if !isAnchor {
//				continue
//			}
//
//			for _, a := range t.Attr {
//				if a.Key == "href" {
//					href := a.Val
//					if strings.Index(href, "/album") == 0 && !_contains(albums, href) {
//						albums = append(albums, href)
//					}
//				}
//			}
//		}
//	}
//}
//
//func checkBandExistence(band string) bool {
//	debug("Checking %s existence", band)
//
//	var url = "https://" + band + ".bandcamp.com/music"
//	resp, err := http.Get(url)
//
//	if resp.StatusCode != 200 || err != nil {
//		return false
//	}
//
//	return true
//}

func loadConfig() {
	file, err := ioutil.ReadFile(flags.ConfigLocation)
	if err != nil {
		return
	}

	var tmpConf Configuration
	if err = json.Unmarshal(file, &tmpConf); err != nil {
		fmt.Println("parsing config file", err)
		os.Exit(1)
	}
	config = tmpConf
}

func parseFlags() {
	configLocation := flag.String("config", path.Dir(os.Args[0]), "Configuration file location")
	ignoreWarnings := flag.Bool("y", false, "Ignore Warnings")

	flag.Parse()

	flags.IgnoreWarnings = *ignoreWarnings
	flags.ConfigLocation = *configLocation

	if len(flag.Args()) == 0 {
		fmt.Println("[ERROR] Missing Band Name/s")
		os.Exit(1)
	}
}

func main() {
	fmt.Println("                          _                     _")
	fmt.Println("___________              | |                   | |")
	fmt.Println("\\__    ___/___________   | |__   __ _ _ __   __| | ___ __ _ _ __ ___  _ __")
	fmt.Println("  |    |  \\_  __ \\__  \\  | '_ \\ / _` | '_ \\ / _` |/ __/ _` | '_ ` _ \\| '_ \\")
	fmt.Println("  |    |   |  | \\// __ \\_| |_) | (_| | | | | (_| | (_| (_| | | | | | | |_) |")
	fmt.Println("  |____|   |__|  (____  /|_.__/ \\__,_|_| |_|\\__,_|\\___\\__,_|_| |_| |_| .__/")
	fmt.Println("                                                                     | |")
	fmt.Println("                                                                     |_| v.0.0.3")

	if len(os.Args) < 2 {
		fmt.Println("[ERROR] Missing Band Name/s")
		os.Exit(1)
	}

	parseFlags()
	loadConfig()

	debug("Flags %v", flags)
	debug("Config %v", config)
	debug("Bands %v", flag.Args())

	var bands []Band
	for _, bandName := range flag.Args() {
		//if !checkBandExistence(band) {
		//	fmt.Printf("[ERROR] Band `%s` doesn't exist\n", band)
		//	os.Exit(1)
		//}

		var bandObject Band
		bandObject.Name = bandName
		bands = append(bands, bandObject)
	}

	albumCount := 0
	// TODO: Run the analysis in parallel
	for index, band := range bands {
		fmt.Printf("Analyzing %s\n", band.Name)

		//var bandAlbums = fetchAlbums(band)
		var bandAlbums = []string{"aaaaaaaa","bbbbbbbb"}
		debug("Band: %s, Albums: %v", band.Name, bandAlbums)
		fmt.Printf("Found %d Albums\n", len(bandAlbums))

		albumCount = albumCount + len(bandAlbums)
		for _, albumName := range bandAlbums {
			var albumObject Album
			albumObject.Name = albumName
			bands[index].Albums = append(bands[index].Albums, albumObject)
		}
	}

	debug("%v", bands)

	debug("Total albums found: %d", albumCount)
	if flags.IgnoreWarnings == false && albumCount > config.AlbumsLimit {
		for {
			var res string
			fmt.Printf("Are you sure you want to download %d albums? (Y/n) ", albumCount)
			fmt.Scanln(&res)

			if res == "n" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}

			if res == "" || res == "y" {
				break
			}

			fmt.Println("Unrecognized command..")
		}
	}


	for bandIndex, band := range bands {
		for albumIndex, album := range band.Albums {
			debug("Fetching album %s tracks", album.Name[7:])

			//var albumTracks = fetchAlbumTracks(band, album)
			var albumTracks = []Track{Track{"a", File{"a"}}}
			for _, track := range albumTracks {
				bands[bandIndex].
					Albums[albumIndex].
					Tracks = append(bands[bandIndex].Albums[albumIndex].Tracks, track)
			}
		}
	}
	debug("%v", bands)
	//fmt.Printf("[INFO] Found %d Tracks\n", len(tracks))
	os.Exit(0)




	//var band string
	//
	//trackThrottle = make(chan int, config.Concurrency.Track)
	//
	//
	//var musicPath = config.OutputDirectory
	//musicPath = musicPath + "/" + band
	//os.MkdirAll(musicPath, 0777)



	//
	//var wg sync.WaitGroup
	//for _, track := range tracks {
	//	trackThrottle <- 1
	//	wg.Add(1)
	//
	//	go downloadTrack(musicPath, track, &wg, trackThrottle)
	//}
	//
	//wg.Wait()
}
