package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/csmith/envflag"
)

var (
	fallbackDir = flag.String("fallback-dir", "output", "path to directory containing backup images")
	uniqueDir   = flag.String("unique-dir", "unique", "path to directory containing unique images")
)

const (
	ipAddressLimit      = time.Hour * 24
	globalLimit         = time.Minute * 10
	fallbackDeleteRate  = time.Hour
	maxGlobalLimitDrift = time.Hour * 24
)

//go:embed frame.png
var frame []byte

var frameImage image.Image

var uniquesByIp = make(map[string]time.Time)
var uniquesMutex sync.Mutex
var lastUnique = time.Now().Add(-globalLimit)
var lastFallback = time.Time{}

func main() {
	envflag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	var err error
	frameImage, err = png.Decode(bytes.NewReader(frame))
	if err != nil {
		panic(err)
	}

	go prune()

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveImage)

	http.ListenAndServe(":8080", mux)
}

func prune() {
	for {
		time.Sleep(ipAddressLimit)
		uniquesMutex.Lock()
		for i := range uniquesByIp {
			if uniquesByIp[i].Add(ipAddressLimit).Before(time.Now()) {
				delete(uniquesByIp, i)
			}
		}
		uniquesMutex.Unlock()
	}
}

func getsUnique(ip string) bool {
	uniquesMutex.Lock()
	defer uniquesMutex.Unlock()

	if lastUnique.Add(globalLimit).After(time.Now()) {
		// Global limit reached
		return false
	}

	if last, ok := uniquesByIp[ip]; ok && last.Add(ipAddressLimit).After(time.Now()) {
		// Individual limit reached
		return false
	}

	// We're go for a unique!

	// Don't allow the lastUnique time to trail _too_ far behind us - keep it capped within the maxGlobalLimitDrift
	if lastUnique.Add(maxGlobalLimitDrift).Before(time.Now()) {
		lastUnique = time.Now().Add(globalLimit).Add(-maxGlobalLimitDrift)
	} else {
		lastUnique = lastUnique.Add(globalLimit)
	}
	
	uniquesByIp[ip] = time.Now()
	return true
}

func serveImage(writer http.ResponseWriter, request *http.Request) {
	// Short-circuit any request with an If-Modified-Since that's within the limit
	if h := request.Header.Get("If-Modified-Since"); h != "" {
		d, err := time.Parse(time.RFC1123, h)
		if err == nil && d.Add(ipAddressLimit).After(time.Now()) {
			writer.WriteHeader(http.StatusNotModified)
			return
		}
	}

	unique := getsUnique(request.Header.Get("X-Forwarded-For"))

	var dir string
	if unique {
		dir = *uniqueDir
	} else {
		dir = *fallbackDir
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(writer, "unable to find avatars", http.StatusInternalServerError)
		log.Printf("Failed to read directory: %v", err)
	}

	if unique && len(files) == 0 {
		// Whelp! Fall back to the fallback.
		unique = false
		dir = *fallbackDir

		files, err = os.ReadDir(dir)
		if err != nil {
			http.Error(writer, "unable to find avatars", http.StatusInternalServerError)
			log.Printf("Failed to read directory: %v", err)
		}
	}

	f := files[rand.Intn(len(files))]

	fh, err := os.Open(filepath.Join(dir, f.Name()))
	if err != nil {
		http.Error(writer, "unable to open image", http.StatusInternalServerError)
		log.Printf("Failed to open file: %v", err)
	}
	defer fh.Close()

	var b []byte
	if unique {
		out := image.NewRGBA(frameImage.Bounds())

		im, err := webp.Decode(fh)
		if err != nil {
			http.Error(writer, "unable to decode image", http.StatusInternalServerError)
			log.Printf("Failed to decode file: %v", err)
		}

		draw.Draw(out, im.Bounds(), im, image.Pt(0, 0), draw.Over)
		draw.Draw(out, im.Bounds(), frameImage, image.Pt(0, 0), draw.Over)

		b, err = webp.EncodeRGB(out, 60)
		if err != nil {
			http.Error(writer, "unable to encode image", http.StatusInternalServerError)
			log.Printf("Failed to encode unique image: %v", err)
		}

		defer os.Remove(filepath.Join(dir, f.Name()))
	} else {
		b, err = io.ReadAll(fh)
		if err != nil {
			http.Error(writer, "unable to read image", http.StatusInternalServerError)
			log.Printf("Failed to read file: %v", err)
		}

		if lastFallback.Add(fallbackDeleteRate).Before(time.Now()) {
			defer os.Remove(filepath.Join(dir, f.Name()))
		}
	}

	writer.Header().Add("Cache-Control", fmt.Sprintf("private, maxage=%d", int(ipAddressLimit.Seconds())))
	writer.Header().Add("Last-Modified", time.Now().Format(time.RFC1123))
	writer.Header().Add("Content-Type", "image/webp")
	writer.Header().Add("X-Unique", fmt.Sprintf("%t", unique))
	writer.Header().Add("X-Source", f.Name())
	writer.WriteHeader(http.StatusOK)
	writer.Write(b)
}
