package main

import (
	"encoding/json"
	"io/ioutil"
	"golang.org/x/exp/rand"
	"log"
	"os"
)

type Posts struct {
	Images []Image `json:"images"`
	Links []Link `json:"links"`
}

type Image string

type Link string

func advertiseRelay(file string) string {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	body, readErr := ioutil.ReadAll(jsonFile)
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
