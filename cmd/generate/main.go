package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/csmith/aca"
	"github.com/csmith/envflag"
)

var (
	sdUrl = flag.String("sd-url", "http://example.com:7860/sdapi/v1/txt2img", "URL to the Static Diffusion Web UI's txt2img endpoint")
	dirs  = flag.String("dirs", "output:1000-2500,unique:1000-2000", "Comma-separated list of directories and their target picture amounts")
)

type Config struct {
	dir string
	min int
	max int
}

func main() {
	envflag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	generator, err := aca.NewGenerator("-", rand.New(rand.NewSource(rand.Int63())))
	if err != nil {
		panic(err)
	}

	targets := strings.Split(*dirs, ",")
	configs := make([]Config, len(targets))
	for i := range targets {
		parts := strings.Split(targets[i], ":")
		limits := strings.Split(parts[1], "-")
		min, _ := strconv.Atoi(limits[0])
		max, _ := strconv.Atoi(limits[1])

		configs[i] = Config{
			dir: parts[0],
			min: min,
			max: max,
		}
	}

	for {
		log.Printf("Checking directories...")

		var needed []string
		for i := range configs {
			files, _ := os.ReadDir(configs[i].dir)
			if len(files) < configs[i].min {
				log.Printf("Directory %s is below minimum %d: %d", configs[i].dir, configs[i].min, len(files))
				for j := len(files); j < configs[i].max; j++ {
					needed = append(needed, filepath.Join(configs[i].dir, fmt.Sprintf("%s.webp", generator.Generate())))
				}
			}
		}

		if len(needed) > 0 {
			generateBatch(needed)
		}

		time.Sleep(time.Minute * 10)
	}
}

func generateBatch(targets []string) {
	launchMachine()
	defer stopMachine()

	guard := make(chan struct{}, 4)
	for t := range targets {
		guard <- struct{}{}
		go func(t int) {
			prompt := Prompt()
			image, err := Generate(*sdUrl, prompt)
			if err != nil {
				for j := 0; j < 10 && err != nil; j++ {
					log.Printf("Request failed: %v, sleeping and trying again", err)
					time.Sleep(time.Second)
					image, err = Generate(*sdUrl, prompt)
				}
				if err != nil {
					panic(err)
				}
			}

			err = os.WriteFile(targets[t], image, os.FileMode(0755))
			if err != nil {
				panic(err)
			}

			log.Printf("Generated %s with prompt: %s", targets[t], prompt)
			<-guard
		}(t)
	}
}
