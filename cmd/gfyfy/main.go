package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/f2prateek/gfyfy/Godeps/_workspace/src/github.com/tj/docopt"
)

const (
	usage = `Gfyfy.

Usage:
  gfyfy [--addr=<a>]
  gfyfy -h | --help
  gfyfy --version

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

func serve(w http.ResponseWriter, r *http.Request) {
	gifURL := r.URL.Path[1:]
	if gifURL == "" {
		w.WriteHeader(200)
		fmt.Fprintln(w, http.StatusText(200))
		return
	}
	if gifURL == "favicon.ico" {
		http.NotFound(w, r)
		return
	}
	if _, err := url.Parse(gifURL); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	logger := log.New(os.Stderr, gifURL+" ", log.LstdFlags)
	gfycatURL := gfycatURL(gifURL)
	logger.Println("gfycat url", gfycatURL)
	resp, err := http.Get(gfycatURL)
	if err != nil {
		logger.Println("error fetching gfycat response", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	var gyfycatResponse GyfycatResponse
	if err := json.NewDecoder(resp.Body).Decode(&gyfycatResponse); err != nil {
		logger.Println("error decoding gfycat response", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	if gyfycatResponse.WebmURL == "" {
		logger.Println("no webmUrl returned")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	logger.Println("redirecting to", gyfycatResponse.WebmURL)
	http.Redirect(w, r, gyfycatResponse.WebmURL, http.StatusMovedPermanently)
}

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatal(err)
	}

	addr := args["--addr"].(string)
	if addr == ":8080" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			log.Println("using $PORT", envPort)
			addr = ":" + envPort
		}
	}

	log.Println("starting gfyfy server on", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(serve)); err != nil {
		log.Fatal(err)
	}
}
