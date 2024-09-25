package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"golang.org/x/exp/rand"
	"log"
)

type Posts struct {
	Images []Image `json:"images"`
	Links []Link `json:"links"`
}

type Image string

type Link string

func advertiseRelay() string {
	url:= "https://raw.githubusercontent.com/pvorangecrush/nostr-pages/refs/heads/main/poster/posts.json"

	postClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "nostr-ad-bot")

	res, getErr := postClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if  readErr != nil {
		log.Fatal(readErr)
	}

	posts1 := Posts{}
	jsonErr := json.Unmarshal(body, &posts1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	randomImageIndex := rand.Intn(len(posts1.Images))
	randomLinkIndex := rand.Intn(len(posts1.Links))
	pick_image := posts1.Images[randomImageIndex]
	pick_link := posts1.Links[randomLinkIndex]
	new_string := string(pick_image) +  "\n" + string(pick_link)
	return new_string
}
