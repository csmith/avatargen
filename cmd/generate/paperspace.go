package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	paperspaceMachine = flag.String("paperspace-machine", "", "ID of the paperspace machine to start")
	paperspaceKey     = flag.String("paperspace-key", "", "API key for starting and stopping paperspace instances")
)

func launchMachine() {
	log.Printf("Launching paperspace machine %s", *paperspaceMachine)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("https://api.paperspace.com/v1/machines/%s/start", *paperspaceMachine), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *paperspaceKey))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Response from launching: %s", b)

	waitForMachine()
}

func stopMachine() {
	log.Printf("Stopping paperspace machine %s", *paperspaceMachine)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("https://api.paperspace.com/v1/machines/%s/stop", *paperspaceMachine), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *paperspaceKey))
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
}

func waitForMachine() {
	type info struct {
		State string `json:"state"`
	}

	for i := 0; i < 30; i++ {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.paperspace.com/v1/machines/%s", *paperspaceMachine), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *paperspaceKey))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		log.Printf("Response from get machine: %s", b)

		status := info{}
		if err := json.Unmarshal(b, &status); err != nil {
			panic(err)
		}

		log.Printf("Machine state is: %s", status.State)
		if strings.Contains(status.State, "ready") {
			return
		}

		time.Sleep(time.Second * 10)
	}

	panic("Machine not ready after 5 minutes")
}
