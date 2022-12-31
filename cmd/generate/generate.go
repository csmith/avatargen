package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

const (
	defaultSamplerName = "DPM2 a"
	defaultSteps       = 20
	defaultConfigScale = 7
	defaultWidth       = 512
	defaultHeight      = 512
)

type Request struct {
	Prompt      string `json:"prompt"`
	SamplerName string `json:"sampler_name"`
	Steps       int    `json:"steps"`
	ConfigScale int    `json:"cfg_scale"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

func RequestForPrompt(prompt string) Request {
	return Request{
		Prompt:      prompt,
		SamplerName: defaultSamplerName,
		Steps:       defaultSteps,
		ConfigScale: defaultConfigScale,
		Width:       defaultWidth,
		Height:      defaultHeight,
	}
}

type Response struct {
	Images []string `json:"images"`
}

func Generate(textToImageUrl, prompt string) ([]byte, error) {
	req := RequestForPrompt(prompt)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(textToImageUrl, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if len(response.Images) != 1 {
		return nil, fmt.Errorf("invalid number of images retured: %v", len(response.Images))
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(response.Images[0])))
	image, err := png.Decode(decoder)
	if err != nil {
		return nil, err
	}

	resized := resize.Resize(256, 256, image, resize.Bicubic)

	imageBytes, err := webp.EncodeRGB(resized, 60)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
