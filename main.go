package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tj/docopt"
)

const (
	usage = `Gfycat Server.

Usage:
  gfycat [--addr=<a>]
  gfycat -h | --help
  gfycat --version

Options:
  -h --help      Show this screen.
  --version      Show version.
  --addr=<a>     Bind address [default: :8080].`
	version = "0.1.0"
)

type GyfycatResponse struct {
	FrameRate int    `json:"frameRate"`
	GfyName   string `json:"gfyName"`
	Gfyname   string `json:"gfyname"`
	Gfysize   int    `json:"gfysize"`
	GifSize   int    `json:"gifSize"`
	GifURL    string `json:"gifUrl"`
	GifWidth  int    `json:"gifWidth"`
	Mp4Url    string `json:"mp4Url"`
	WebmURL   string `json:"webmUrl"`
}

func gfycatURL(url string) string {
	return fmt.Sprintf("http://upload.gfycat.com/transcode?fetchUrl=%s", url)
}

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path[1:]
		log.Println("url", url)
		gfycatURL := gfycatURL(url)
		log.Println("gfycat url", gfycatURL)
		resp, err := http.Get(gfycatURL)
		if err != nil {
			log.Println("error fetching gfycat response", err)
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		var gyfycatResponse GyfycatResponse
		if err := json.NewDecoder(resp.Body).Decode(&gyfycatResponse); err != nil {
			log.Println("error decoding gfycat response", err)
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		if gyfycatResponse.WebmURL == "" {
			log.Println("no webmUrl returned")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		log.Println("redirecting to", gyfycatResponse.WebmURL)
		http.Redirect(w, r, gyfycatResponse.WebmURL, http.StatusMovedPermanently)
	})

	addr := args["--addr"].(string)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
