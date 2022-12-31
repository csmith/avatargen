package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.paperspace.io/machines/%s/start", *paperspaceMachine), nil)
	req.Header.Add("X-Api-Key", *paperspaceKey)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	waitForMachine()
}

func stopMachine() {
	log.Printf("Stopping paperspace machine %s", *paperspaceMachine)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.paperspace.io/machines/%s/stop", *paperspaceMachine), nil)
	req.Header.Add("X-Api-Key", *paperspaceKey)
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
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.paperspace.io/machines/getMachinePublic?machineId=%s", *paperspaceMachine), nil)
		req.Header.Add("X-Api-Key", *paperspaceKey)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		status := info{}
		if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
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
