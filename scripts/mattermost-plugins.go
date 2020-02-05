package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const URL = "https://raw.githubusercontent.com/mattermost/mattermost-marketplace/master/plugins.json"

func main() {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Invalid status: " + resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var p []Plugin
	if err := json.Unmarshal(b, &p); err != nil {
		log.Fatal(err)
	}
	var url = make(map[string]interface{})
	for _, v := range p {
		url[v.URL] = nil
	}
	for k, _ := range url {
		fmt.Println(k)
	}
}

type Plugin struct {
	URL string `json:"homepage_url"`
}
